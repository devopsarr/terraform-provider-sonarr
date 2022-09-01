package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ provider.DataSourceType = dataAllSeriesType{}
var _ datasource.DataSource = dataAllSeries{}

type dataAllSeriesType struct{}

type dataAllSeries struct {
	provider sonarrProvider
}

// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
// QualityProfiles is a list of QualityProfile.
type SeriesList struct {
	ID     types.String `tfsdk:"id"`
	Series types.Set    `tfsdk:"series"`
}

func (t dataAllSeriesType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "List all available [Series](../resources/series).",
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

func (t dataAllSeriesType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return dataAllSeries{
		provider: provider,
	}, diags
}

func (d dataAllSeries) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SeriesList
	diags := resp.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get series current value
	response, err := d.provider.client.GetAllSeriesContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read series, got error: %s", err))

		return
	}
	// Map response body to resource schema attribute
	series := *writeSeriesList(ctx, response)
	tfsdk.ValueFrom(ctx, series, data.Series.Type(context.Background()), &data.Series)

	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.String{Value: strconv.Itoa(len(response))}
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func writeSeriesList(ctx context.Context, series []*sonarr.Series) *[]Series {
	output := make([]Series, len(series))
	for i, t := range series {
		output[i] = *writeSeries(ctx, t)
	}

	return &output
}
