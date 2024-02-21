package provider

import (
	"context"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mitchellh/hashstructure/v2"
)

const (
	customFormatConditionSizeDataSourceName = "custom_format_condition_size"
	customFormatConditionSizeImplementation = "SizeSpecification"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CustomFormatConditionSizeDataSource{}

func NewCustomFormatConditionSizeDataSource() datasource.DataSource {
	return &CustomFormatConditionSizeDataSource{}
}

// CustomFormatConditionSizeDataSource defines the custom_format_condition_size implementation.
type CustomFormatConditionSizeDataSource struct {
	client *sonarr.APIClient
	auth   context.Context
}

func (d *CustomFormatConditionSizeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + customFormatConditionSizeDataSourceName
}

func (d *CustomFormatConditionSizeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Profiles --> Custom Format Condition Size data source.\nFor more information refer to [Custom Format Conditions](https://wiki.servarr.com/sonarr/settings#conditions).",
		Attributes: map[string]schema.Attribute{
			"negate": schema.BoolAttribute{
				MarkdownDescription: "Negate flag.",
				Required:            true,
			},
			"required": schema.BoolAttribute{
				MarkdownDescription: "Computed flag.",
				Required:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "Implementation.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Specification name.",
				Required:            true,
			},
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.Int64Attribute{
				MarkdownDescription: "Custom format condition size ID.",
				Computed:            true,
			},
			// Field values
			"min": schema.Int64Attribute{
				MarkdownDescription: "Min size in GB.",
				Required:            true,
			},
			"max": schema.Int64Attribute{
				MarkdownDescription: "Max size in GB.",
				Required:            true,
			},
		},
	}
}

func (d *CustomFormatConditionSizeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *CustomFormatConditionSizeDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *CustomFormatConditionMinMax

	hash, err := hashstructure.Hash(&data, hashstructure.FormatV2, nil)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, helpers.ParseClientError(helpers.Create, customFormatConditionSizeDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+customFormatConditionSizeDataSourceName)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("implementation"), customFormatConditionSizeImplementation)...)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), int64(hash))...)
}
