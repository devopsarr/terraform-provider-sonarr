package provider

import (
	"context"
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

const customFormatResourceName = "custom_format"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &CustomFormatResource{}
	_ resource.ResourceWithImportState = &CustomFormatResource{}
)

func NewCustomFormatResource() resource.Resource {
	return &CustomFormatResource{}
}

// CustomFormatResource defines the custom format implementation.
type CustomFormatResource struct {
	client *sonarr.APIClient
}

// CustomFormat describes the custom format data model.
type CustomFormat struct {
	Specifications                  types.Set    `tfsdk:"specifications"`
	Name                            types.String `tfsdk:"name"`
	ID                              types.Int64  `tfsdk:"id"`
	IncludeCustomFormatWhenRenaming types.Bool   `tfsdk:"include_custom_format_when_renaming"`
}

func (c CustomFormat) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"include_custom_format_when_renaming": types.BoolType,
			"id":                                  types.Int64Type,
			"name":                                types.StringType,
			"specifications":                      types.SetType{}.WithElementType(CustomFormatCondition{}.getType()),
		})
}

func (r *CustomFormatResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + customFormatResourceName
}

func (r *CustomFormatResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Profiles -->Custom Format resource.\nFor more information refer to [Custom Format](https://wiki.servarr.com/sonarr/settings#custom-formats).",
		Attributes: map[string]schema.Attribute{
			"include_custom_format_when_renaming": schema.BoolAttribute{
				MarkdownDescription: "Include custom format when renaming flag.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Custom Format name.",
				Required:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Custom Format ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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

func (r CustomFormatResource) getSpecificationSchema() schema.Schema {
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
			"min": schema.Int64Attribute{
				MarkdownDescription: "Min.",
				Optional:            true,
				Computed:            true,
			},
			"max": schema.Int64Attribute{
				MarkdownDescription: "Max.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *CustomFormatResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *CustomFormatResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *CustomFormat

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new CustomFormat
	request := client.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.CustomFormatApi.CreateCustomFormat(ctx).CustomFormatResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, customFormatResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+customFormatResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state CustomFormat

	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *CustomFormatResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client CustomFormat

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get CustomFormat current value
	response, _, err := r.client.CustomFormatApi.GetCustomFormatById(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, customFormatResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+customFormatResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	// this is needed because of many empty fields are unknown in both plan and read
	var state CustomFormat

	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *CustomFormatResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *CustomFormat

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update CustomFormat
	request := client.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.CustomFormatApi.UpdateCustomFormat(ctx, strconv.Itoa(int(request.GetId()))).CustomFormatResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, customFormatResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+customFormatResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state CustomFormat

	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *CustomFormatResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete CustomFormat current value
	_, err := r.client.CustomFormatApi.DeleteCustomFormat(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, customFormatResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+customFormatResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *CustomFormatResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+customFormatResourceName+": "+req.ID)
}

func (c *CustomFormat) write(ctx context.Context, customFormat *sonarr.CustomFormatResource, diags *diag.Diagnostics) {
	var tempDiag diag.Diagnostics

	c.ID = types.Int64Value(int64(customFormat.GetId()))
	c.Name = types.StringValue(customFormat.GetName())
	c.IncludeCustomFormatWhenRenaming = types.BoolValue(customFormat.GetIncludeCustomFormatWhenRenaming())

	specs := make([]CustomFormatCondition, len(customFormat.Specifications))
	for n, c := range customFormat.Specifications {
		specs[n].write(ctx, c)
	}

	c.Specifications, tempDiag = types.SetValueFrom(ctx, CustomFormatResource{}.getSpecificationSchema().Type(), specs)
	diags.Append(tempDiag...)
}

func (c *CustomFormat) read(ctx context.Context, diags *diag.Diagnostics) *sonarr.CustomFormatResource {
	specifications := make([]CustomFormatCondition, len(c.Specifications.Elements()))
	diags.Append(c.Specifications.ElementsAs(ctx, &specifications, false)...)
	specs := make([]*sonarr.CustomFormatSpecificationSchema, len(specifications))

	for n, d := range specifications {
		specs[n] = d.read(ctx)
	}

	format := sonarr.NewCustomFormatResource()
	format.SetId(int32(c.ID.ValueInt64()))
	format.SetName(c.Name.ValueString())
	format.SetIncludeCustomFormatWhenRenaming(c.IncludeCustomFormatWhenRenaming.ValueBool())
	format.SetSpecifications(specs)

	return format
}
