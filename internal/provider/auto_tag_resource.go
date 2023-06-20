package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const autoTagResourceName = "auto_tag"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &AutoTagResource{}
	_ resource.ResourceWithImportState = &AutoTagResource{}
)

func NewAutoTagResource() resource.Resource {
	return &AutoTagResource{}
}

// AutoTagResource defines the tag implementation.
type AutoTagResource struct {
	client *sonarr.APIClient
}

// AutoTag describes the tag data model.
type AutoTag struct {
	Specifications          types.Set    `tfsdk:"specifications"`
	Tags                    types.Set    `tfsdk:"tags"`
	Name                    types.String `tfsdk:"name"`
	ID                      types.Int64  `tfsdk:"id"`
	RemoveTagsAutomatically types.Bool   `tfsdk:"remove_tags_automatically"`
}

func (r *AutoTagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + autoTagResourceName
}

func (r *AutoTagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Tags -->Auto Tag resource.\nFor more information refer to [Tags](https://wiki.servarr.com/sonarr/settings#tags) documentation.",
		Attributes: map[string]schema.Attribute{
			"remove_tags_automatically": schema.BoolAttribute{
				MarkdownDescription: "Remove tags automatically flag.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Auto Tag name.",
				Required:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Auto Tag ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"specifications": schema.SetNestedAttribute{
				MarkdownDescription: "Specifications.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: r.getSpecificationSchema().Attributes,
				},
			},
		},
	}
}

func (r AutoTagResource) getSpecificationSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"negate": schema.BoolAttribute{
				MarkdownDescription: "Negate flag.",
				Optional:            true,
				Computed:            true,
			},
			"required": schema.BoolAttribute{
				MarkdownDescription: "Required flag.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Specification name.",
				Optional:            true,
				Computed:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "Implementation.",
				Optional:            true,
				Computed:            true,
			},
			// Field values
			"value": schema.StringAttribute{
				MarkdownDescription: "Value.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *AutoTagResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *AutoTagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var autoTag *AutoTag

	resp.Diagnostics.Append(req.Plan.Get(ctx, &autoTag)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new auto tag
	request := autoTag.read(ctx)

	response, _, err := r.client.AutoTaggingApi.CreateAutoTagging(ctx).AutoTaggingResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, autoTagResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+autoTagResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	autoTag.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &autoTag)...)
}

func (r *AutoTagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var autoTag *AutoTag

	resp.Diagnostics.Append(req.State.Get(ctx, &autoTag)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get auto tag current value
	response, _, err := r.client.AutoTaggingApi.GetAutoTaggingById(ctx, int32(autoTag.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, autoTagResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+autoTagResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	autoTag.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &autoTag)...)
}

func (r *AutoTagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var autoTag *AutoTag

	resp.Diagnostics.Append(req.Plan.Get(ctx, &autoTag)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update auto tag
	request := autoTag.read(ctx)

	response, _, err := r.client.AutoTaggingApi.UpdateAutoTagging(ctx, fmt.Sprint(request.GetId())).AutoTaggingResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, autoTagResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+autoTagResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	autoTag.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &autoTag)...)
}

func (r *AutoTagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var autoTag *AutoTag

	resp.Diagnostics.Append(req.State.Get(ctx, &autoTag)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete auto tag current value
	_, err := r.client.AutoTaggingApi.DeleteAutoTagging(ctx, int32(autoTag.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, autoTagResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+autoTagResourceName+": "+strconv.Itoa(int(autoTag.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *AutoTagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+autoTagResourceName+": "+req.ID)
}

func (t *AutoTag) write(ctx context.Context, autoTag *sonarr.AutoTaggingResource) {
	t.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, autoTag.GetTags())
	t.ID = types.Int64Value(int64(autoTag.GetId()))
	t.Name = types.StringValue(autoTag.GetName())
	t.RemoveTagsAutomatically = types.BoolValue(autoTag.GetRemoveTagsAutomatically())
	t.Specifications = types.SetValueMust(AutoTagResource{}.getSpecificationSchema().Type(), nil)

	specs := make([]AutoTagCondition, len(autoTag.Specifications))
	for n, s := range autoTag.Specifications {
		specs[n].write(ctx, s)
	}

	tfsdk.ValueFrom(ctx, specs, t.Specifications.Type(ctx), &t.Specifications)
}

func (t *AutoTag) read(ctx context.Context) *sonarr.AutoTaggingResource {
	tags := make([]*int32, len(t.Tags.Elements()))
	tfsdk.ValueAs(ctx, t.Tags, &tags)

	specifications := make([]AutoTagCondition, len(t.Specifications.Elements()))
	tfsdk.ValueAs(ctx, t.Specifications, &specifications)
	specs := make([]*sonarr.AutoTaggingSpecificationSchema, len(specifications))

	for n, s := range specifications {
		specs[n] = s.read(ctx)
	}

	autoTag := sonarr.NewAutoTaggingResource()
	autoTag.SetId(int32(t.ID.ValueInt64()))
	autoTag.SetName(t.Name.ValueString())
	autoTag.SetRemoveTagsAutomatically(t.RemoveTagsAutomatically.ValueBool())
	autoTag.SetSpecifications(specs)
	autoTag.SetTags(tags)

	return autoTag
}
