package provider

import (
	"context"
	"strings"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mitchellh/hashstructure/v2"
)

const autoTagConditionDataSourceName = "auto_tag_condition"

var autoTagFields = helpers.Fields{
	Strings: []string{"value"},
}

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AutoTagConditionDataSource{}

func NewAutoTagConditionDataSource() datasource.DataSource {
	return &AutoTagConditionDataSource{}
}

// AutoTagConditionDataSource defines the auto tag condition implementation.
type AutoTagConditionDataSource struct {
	client *sonarr.APIClient
}

// AutoTagCondition describes the auto tag condition data model.
type AutoTagCondition struct {
	Name           types.String `tfsdk:"name"`
	Implementation types.String `tfsdk:"implementation"`
	Value          types.String `tfsdk:"value"`
	Negate         types.Bool   `tfsdk:"negate"`
	Required       types.Bool   `tfsdk:"required"`
}

func (c AutoTagCondition) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"name":           types.StringType,
			"implementation": types.StringType,
			"value":          types.StringType,
			"negate":         types.BoolType,
			"required":       types.BoolType,
		})
}

func (d *AutoTagConditionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + autoTagConditionDataSourceName
}

func (d *AutoTagConditionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Tags --> Generic Auto Tag Condition data source. When possible use a specific data source instead.\nFor more information refer to [ Format Conditions](https://wiki.servarr.com/sonarr/settings#conditions).\n To be used in conjunction with [ Format](../resources/auto_tag).",
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
				MarkdownDescription: " format condition ID.",
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

func (d *AutoTagConditionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *AutoTagConditionDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *AutoTagCondition

	hash, err := hashstructure.Hash(&data, hashstructure.FormatV2, nil)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, helpers.ParseClientError(helpers.Create, autoTagConditionDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+autoTagConditionDataSourceName)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), int64(hash))...)
}

func (c *AutoTagCondition) write(ctx context.Context, spec *sonarr.AutoTaggingSpecificationSchema) {
	c.Implementation = types.StringValue(spec.GetImplementation())
	c.Name = types.StringValue(spec.GetName())
	c.Negate = types.BoolValue(spec.GetNegate())
	c.Required = types.BoolValue(spec.GetRequired())
	helpers.WriteFields(ctx, c, spec.GetFields(), autoTagFields)

	// workaround to manage list.
	if spec.GetFields()[0].GetType() == "tag" {
		c.Value = types.StringValue(strings.Trim(c.Value.ValueString(), "[]"))
	}
}

func (c *AutoTagCondition) read(ctx context.Context) *sonarr.AutoTaggingSpecificationSchema {
	spec := sonarr.NewAutoTaggingSpecificationSchema()
	spec.SetName(c.Name.ValueString())

	spec.SetImplementation(c.Implementation.ValueString())
	spec.SetNegate(c.Negate.ValueBool())
	spec.SetRequired(c.Required.ValueBool())
	spec.SetFields(helpers.ReadFields(ctx, c, autoTagFields))

	return spec
}
