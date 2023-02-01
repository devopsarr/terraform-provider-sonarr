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
	"golang.org/x/exp/slices"
)

const customFormatResourceName = "custom_format"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &CustomFormatResource{}
	_ resource.ResourceWithImportState = &CustomFormatResource{}
)

var (
	customFormatStringFields = []string{"value"}
	customFormatIntFields    = []string{"min", "max"}
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

// Specification is part of CustomFormat.
type Specification struct {
	Name           types.String `tfsdk:"name"`
	Implementation types.String `tfsdk:"implementation"`
	Value          types.String `tfsdk:"value"`
	Min            types.Int64  `tfsdk:"min"`
	Max            types.Int64  `tfsdk:"max"`
	Negate         types.Bool   `tfsdk:"negate"`
	Required       types.Bool   `tfsdk:"required"`
}

func (r *CustomFormatResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + customFormatResourceName
}

func (r *CustomFormatResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
	request := client.read(ctx)

	response, _, err := r.client.CustomFormatApi.CreateCustomFormat(ctx).CustomFormatResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", customFormatResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+customFormatResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state CustomFormat

	state.write(ctx, response)
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
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", customFormatResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+customFormatResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	// this is needed because of many empty fields are unknown in both plan and read
	var state CustomFormat

	state.write(ctx, response)
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
	request := client.read(ctx)

	response, _, err := r.client.CustomFormatApi.UpdateCustomFormat(ctx, strconv.Itoa(int(request.GetId()))).CustomFormatResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", customFormatResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+customFormatResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state CustomFormat

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *CustomFormatResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var client *CustomFormat

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete CustomFormat current value
	_, err := r.client.CustomFormatApi.DeleteCustomFormat(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", customFormatResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+customFormatResourceName+": "+strconv.Itoa(int(client.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *CustomFormatResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+customFormatResourceName+": "+req.ID)
}

func (c *CustomFormat) write(ctx context.Context, customFormat *sonarr.CustomFormatResource) {
	c.ID = types.Int64Value(int64(customFormat.GetId()))
	c.Name = types.StringValue(customFormat.GetName())
	c.IncludeCustomFormatWhenRenaming = types.BoolValue(customFormat.GetIncludeCustomFormatWhenRenaming())
	c.Specifications = types.SetValueMust(CustomFormatResource{}.getSpecificationSchema().Type(), nil)

	specs := make([]Specification, len(customFormat.Specifications))
	for n, c := range customFormat.Specifications {
		specs[n].write(c)
	}

	tfsdk.ValueFrom(ctx, specs, c.Specifications.Type(ctx), &c.Specifications)
}

func (s *Specification) write(spec *sonarr.CustomFormatSpecificationSchema) {
	s.Implementation = types.StringValue(spec.GetImplementation())
	s.Name = types.StringValue(spec.GetName())
	s.Negate = types.BoolValue(spec.GetNegate())
	s.Required = types.BoolValue(spec.GetRequired())
	s.writeFields(spec.GetFields())
}

func (s *Specification) writeFields(fields []*sonarr.Field) {
	for _, f := range fields {
		if f.Value == nil {
			continue
		}

		if slices.Contains(customFormatStringFields, f.GetName()) {
			helpers.WriteStringField(f, s)

			continue
		}

		if slices.Contains(customFormatIntFields, f.GetName()) {
			helpers.WriteIntField(f, s)

			continue
		}
	}
}

func (c *CustomFormat) read(ctx context.Context) *sonarr.CustomFormatResource {
	specifications := make([]Specification, len(c.Specifications.Elements()))
	tfsdk.ValueAs(ctx, c.Specifications, &specifications)
	specs := make([]*sonarr.CustomFormatSpecificationSchema, len(specifications))

	for n, d := range specifications {
		specs[n] = d.read()
	}

	format := sonarr.NewCustomFormatResource()
	format.SetId(int32(c.ID.ValueInt64()))
	format.SetName(c.Name.ValueString())
	format.SetIncludeCustomFormatWhenRenaming(c.IncludeCustomFormatWhenRenaming.ValueBool())
	format.SetSpecifications(specs)

	return format
}

func (s *Specification) read() *sonarr.CustomFormatSpecificationSchema {
	spec := sonarr.NewCustomFormatSpecificationSchema()
	spec.SetName(s.Name.ValueString())

	spec.SetImplementation(s.Implementation.ValueString())
	spec.SetNegate(s.Negate.ValueBool())
	spec.SetRequired(s.Required.ValueBool())
	spec.SetFields(s.readFields())

	return spec
}

func (s *Specification) readFields() []*sonarr.Field {
	var output []*sonarr.Field

	for _, i := range customFormatIntFields {
		if field := helpers.ReadIntField(i, s); field != nil {
			output = append(output, field)
		}
	}

	for _, str := range customFormatStringFields {
		if field := helpers.ReadStringField(str, s); field != nil {
			output = append(output, field)
		}
	}

	return output
}
