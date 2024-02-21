package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const tagResourceName = "tag"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &TagResource{}
	_ resource.ResourceWithImportState = &TagResource{}
)

func NewTagResource() resource.Resource {
	return &TagResource{}
}

// TagResource defines the tag implementation.
type TagResource struct {
	client *sonarr.APIClient
	auth   context.Context
}

// Tag describes the tag data model.
type Tag struct {
	Label types.String `tfsdk:"label"`
	ID    types.Int64  `tfsdk:"id"`
}

func (t Tag) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"id":    types.Int64Type,
			"label": types.StringType,
		})
}

func (r *TagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + tagResourceName
}

func (r *TagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Tags -->\nTag resource.\nFor more information refer to [Tags](https://wiki.servarr.com/sonarr/settings#tags) documentation.",
		Attributes: map[string]schema.Attribute{
			"label": schema.StringAttribute{
				MarkdownDescription: "Tag label. It must be lowercase.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^.*[^A-Z]+.*$`),
						"String cannot contains uppercase values",
					),
				},
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Tag ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *TagResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *TagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var tag *Tag

	resp.Diagnostics.Append(req.Plan.Get(ctx, &tag)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Tag
	request := *sonarr.NewTagResource()
	request.SetLabel(tag.Label.ValueString())

	response, _, err := r.client.TagAPI.CreateTag(r.auth).TagResource(request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, tagResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+tagResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	tag.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &tag)...)
}

func (r *TagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var tag *Tag

	resp.Diagnostics.Append(req.State.Get(ctx, &tag)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get tag current value
	response, _, err := r.client.TagAPI.GetTagById(r.auth, int32(tag.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, tagResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+tagResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	tag.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &tag)...)
}

func (r *TagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var tag *Tag

	resp.Diagnostics.Append(req.Plan.Get(ctx, &tag)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Tag
	tagResource := *sonarr.NewTagResource()
	tagResource.SetLabel(tag.Label.ValueString())
	tagResource.SetId(int32(tag.ID.ValueInt64()))

	response, _, err := r.client.TagAPI.UpdateTag(r.auth, fmt.Sprint(tagResource.GetId())).TagResource(tagResource).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, tagResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+tagResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	tag.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &tag)...)
}

func (r *TagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete tag current value
	_, err := r.client.TagAPI.DeleteTag(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, tagResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+tagResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *TagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+tagResourceName+": "+req.ID)
}

func (t *Tag) write(tag *sonarr.TagResource) {
	t.ID = types.Int64Value(int64(tag.GetId()))
	t.Label = types.StringValue(tag.GetLabel())
}
