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
	autoTagConditionSeriesTypeDataSourceName = "auto_tag_condition_series_type"
	autoTagConditionSeriesTypeImplementation = "SeriesTypeSpecification"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AutoTagConditionSeriesTypeDataSource{}

func NewAutoTagConditionSeriesTypeDataSource() datasource.DataSource {
	return &AutoTagConditionSeriesTypeDataSource{}
}

// AutoTagConditionSeriesTypeDataSource defines the auto_tag_condition_series type implementation.
type AutoTagConditionSeriesTypeDataSource struct {
	client *sonarr.APIClient
}

func (d *AutoTagConditionSeriesTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + autoTagConditionSeriesTypeDataSourceName
}

func (d *AutoTagConditionSeriesTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Profiles --> Auto Tag Condition Series Type data source.\nFor more intagion refer to [Auto Tag Conditions](https://wiki.servarr.com/sonarr/settings#conditions).",
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
				MarkdownDescription: "Auto tag condition series type ID.",
				Computed:            true,
			},
			// Field values
			"value": schema.StringAttribute{
				MarkdownDescription: "Series type. `0` Standard, `1` Daily, `2` Anime.",
				Required:            true,
			},
		},
	}
}

func (d *AutoTagConditionSeriesTypeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *AutoTagConditionSeriesTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *AutoTagCondition

	hash, err := hashstructure.Hash(&data, hashstructure.FormatV2, nil)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, helpers.ParseClientError(helpers.Create, autoTagConditionSeriesTypeDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+autoTagConditionSeriesTypeDataSourceName)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("implementation"), autoTagConditionSeriesTypeImplementation)...)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), int64(hash))...)
}
