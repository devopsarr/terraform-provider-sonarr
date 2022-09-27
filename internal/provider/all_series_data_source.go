package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const allSeriesDataSourceName = "all_series"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AllSeriessDataSource{}

func NewAllSeriessDataSource() datasource.DataSource {
	return &AllSeriessDataSource{}
}

// AllSeriessDataSource defines the tags implementation.
type AllSeriessDataSource struct {
	client *sonarr.Sonarr
}

// AllSeriess describes the series(es) data model.
type SeriesList struct {
	Series types.Set    `tfsdk:"series"`
	ID     types.String `tfsdk:"id"`
}

func (d *AllSeriessDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + allSeriesDataSourceName
}

func (d *AllSeriessDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "[subcategory:Series]: #\nList all available [Series](../resources/series).",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"series": {
				MarkdownDescription: "Series list.",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						MarkdownDescription: "Series ID.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"title": {
						MarkdownDescription: "Series Title.",
						Computed:            true,
						Type:                types.StringType,
					},
					"title_slug": {
						MarkdownDescription: "Series Title in kebab format.",
						Computed:            true,
						Type:                types.StringType,
					},
					"monitored": {
						MarkdownDescription: "Monitored flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"season_folder": {
						MarkdownDescription: "Season Folder flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"use_scene_numbering": {
						MarkdownDescription: "Scene numbering flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"language_profile_id": {
						MarkdownDescription: "Language Profile ID .",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"quality_profile_id": {
						MarkdownDescription: "Quality Profile ID.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"tvdb_id": {
						MarkdownDescription: "TVDB ID.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"path": {
						MarkdownDescription: "Series Path.",
						Computed:            true,
						Type:                types.StringType,
					},
					"root_folder_path": {
						MarkdownDescription: "Series Root Folder.",
						Computed:            true,
						Type:                types.StringType,
					},
					"tags": {
						MarkdownDescription: "Tags.",
						Computed:            true,
						Type: types.SetType{
							ElemType: types.Int64Type,
						},
					},
				}),
			},
		},
	}, nil
}

func (d *AllSeriessDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			helpers.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *AllSeriessDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SeriesList

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get series current value
	response, err := d.client.GetAllSeriesContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", allSeriesDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+allSeriesDataSourceName)
	// Map response body to resource schema attribute
	series := *writeSeriesList(ctx, response)
	tfsdk.ValueFrom(ctx, series, data.Series.Type(context.Background()), &data.Series)

	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.String{Value: strconv.Itoa(len(response))}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func writeSeriesList(ctx context.Context, series []*sonarr.Series) *[]Series {
	output := make([]Series, len(series))
	for i, t := range series {
		output[i] = *writeSeries(ctx, t)
	}

	return &output
}
