package provider

import (
	"context"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mitchellh/hashstructure/v2"
	"golang.org/x/exp/slices"
)

const customFormatConditionDataSourceName = "custom_format_condition"

var (
	customFormatStringFields = []string{"value"}
	customFormatIntFields    = []string{"min", "max"}
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CustomFormatConditionDataSource{}

func NewCustomFormatConditionDataSource() datasource.DataSource {
	return &CustomFormatConditionDataSource{}
}

// CustomFormatConditionDataSource defines the custom format condition implementation.
type CustomFormatConditionDataSource struct {
	client *sonarr.APIClient
}

// CustomFormatCondition describes the custom format condition data model.
type CustomFormatCondition struct {
	Name           types.String `tfsdk:"name"`
	Implementation types.String `tfsdk:"implementation"`
	Value          types.String `tfsdk:"value"`
	Min            types.Int64  `tfsdk:"min"`
	Max            types.Int64  `tfsdk:"max"`
	Negate         types.Bool   `tfsdk:"negate"`
	Required       types.Bool   `tfsdk:"required"`
}

// CustomFormatValue describes the custom format value data model.
type CustomFormatConditionValue struct {
	Name     types.String `tfsdk:"name"`
	Value    types.String `tfsdk:"value"`
	Negate   types.Bool   `tfsdk:"negate"`
	Required types.Bool   `tfsdk:"required"`
}

// CustomFormatMinMax describes the custom format min max data model.
type CustomFormatConditionMinMax struct {
	Name     types.String `tfsdk:"name"`
	Min      types.Int64  `tfsdk:"min"`
	Max      types.Int64  `tfsdk:"max"`
	Negate   types.Bool   `tfsdk:"negate"`
	Required types.Bool   `tfsdk:"required"`
}

func (d *CustomFormatConditionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + customFormatConditionDataSourceName
}

func (d *CustomFormatConditionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Profiles --> Generic Custom Format Condition data source. When possible use a specific data source instead.\nFor more information refer to [Custom Format Conditions](https://wiki.servarr.com/sonarr/settings#conditions).\n To be used in conjunction with [Custom Format](../resources/custom_format).",
		Attributes: map[string]schema.Attribute{
			"negate": schema.BoolAttribute{
				MarkdownDescription: "Negate flag.",
				Required:            true,
			},
			"required": schema.BoolAttribute{
				MarkdownDescription: "Computed flag.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Specification name.",
				Required:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "Implementation.",
				Required:            true,
			},
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.Int64Attribute{
				MarkdownDescription: "Custom format condition ID.",
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

func (d *CustomFormatConditionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *CustomFormatConditionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *CustomFormatCondition

	hash, err := hashstructure.Hash(&data, hashstructure.FormatV2, nil)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, helpers.ParseClientError(helpers.Create, customFormatConditionDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+customFormatConditionDataSourceName)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), int64(hash))...)
}

func (s *CustomFormatCondition) write(spec *sonarr.CustomFormatSpecificationSchema) {
	s.Implementation = types.StringValue(spec.GetImplementation())
	s.Name = types.StringValue(spec.GetName())
	s.Negate = types.BoolValue(spec.GetNegate())
	s.Required = types.BoolValue(spec.GetRequired())
	s.writeFields(spec.GetFields())
}

func (s *CustomFormatCondition) writeFields(fields []*sonarr.Field) {
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

func (s *CustomFormatCondition) read() *sonarr.CustomFormatSpecificationSchema {
	spec := sonarr.NewCustomFormatSpecificationSchema()
	spec.SetName(s.Name.ValueString())

	spec.SetImplementation(s.Implementation.ValueString())
	spec.SetNegate(s.Negate.ValueBool())
	spec.SetRequired(s.Required.ValueBool())
	spec.SetFields(s.readFields())

	return spec
}

func (s *CustomFormatCondition) readFields() []*sonarr.Field {
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
