package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
	auth   context.Context
}

// AutoTag describes the tag data model.
type AutoTag struct {
	Specifications          types.Set    `tfsdk:"specifications"`
	Tags                    types.Set    `tfsdk:"tags"`
	Name                    types.String `tfsdk:"name"`
	ID                      types.Int64  `tfsdk:"id"`
	RemoveTagsAutomatically types.Bool   `tfsdk:"remove_tags_automatically"`
}

func (t AutoTag) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"specifications":            types.SetType{}.WithElementType(AutoTagCondition{}.getType()),
			"tags":                      types.SetType{}.WithElementType(types.Int64Type),
			"name":                      types.StringType,
			"id":                        types.Int64Type,
			"remove_tags_automatically": types.BoolType,
		})
}

func (r *AutoTagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + autoTagResourceName
}

func (r *AutoTagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
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
	request := autoTag.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.AutoTaggingAPI.CreateAutoTagging(r.auth).AutoTaggingResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, autoTagResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+autoTagResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	autoTag.write(ctx, response, &resp.Diagnostics)
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
	response, _, err := r.client.AutoTaggingAPI.GetAutoTaggingById(r.auth, int32(autoTag.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, autoTagResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+autoTagResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	autoTag.write(ctx, response, &resp.Diagnostics)
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
	request := autoTag.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.AutoTaggingAPI.UpdateAutoTagging(r.auth, fmt.Sprint(request.GetId())).AutoTaggingResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, autoTagResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+autoTagResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	autoTag.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &autoTag)...)
}

func (r *AutoTagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete auto tag current value
	_, err := r.client.AutoTaggingAPI.DeleteAutoTagging(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, autoTagResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+autoTagResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *AutoTagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+autoTagResourceName+": "+req.ID)
}

func (t *AutoTag) write(ctx context.Context, autoTag *sonarr.AutoTaggingResource, diags *diag.Diagnostics) {
	var tempDiag diag.Diagnostics

	t.ID = types.Int64Value(int64(autoTag.GetId()))
	t.Name = types.StringValue(autoTag.GetName())
	t.RemoveTagsAutomatically = types.BoolValue(autoTag.GetRemoveTagsAutomatically())

	specs := make([]AutoTagCondition, len(autoTag.Specifications))
	for n, s := range autoTag.Specifications {
		specs[n].write(ctx, &s)
	}

	t.Tags, tempDiag = types.SetValueFrom(ctx, types.Int64Type, autoTag.GetTags())
	diags.Append(tempDiag...)
	t.Specifications, tempDiag = types.SetValueFrom(ctx, AutoTagCondition{}.getType(), specs)
	diags.Append(tempDiag...)
}

func (t *AutoTag) read(ctx context.Context, diags *diag.Diagnostics) *sonarr.AutoTaggingResource {
	specifications := make([]AutoTagCondition, len(t.Specifications.Elements()))
	diags.Append(t.Specifications.ElementsAs(ctx, &specifications, false)...)
	specs := make([]sonarr.AutoTaggingSpecificationSchema, len(specifications))

	for n, s := range specifications {
		specs[n] = *s.read(ctx)
	}

	autoTag := sonarr.NewAutoTaggingResource()
	autoTag.SetId(int32(t.ID.ValueInt64()))
	autoTag.SetName(t.Name.ValueString())
	autoTag.SetRemoveTagsAutomatically(t.RemoveTagsAutomatically.ValueBool())
	autoTag.SetSpecifications(specs)
	diags.Append(t.Tags.ElementsAs(ctx, &autoTag.Tags, true)...)

	return autoTag
}
