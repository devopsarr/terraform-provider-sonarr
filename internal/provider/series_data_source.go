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
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.DataSourceType = dataSeriesType{}
var _ datasource.DataSource = dataSeries{}

type dataSeriesType struct{}

type dataSeries struct {
	provider sonarrProvider
}

// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
// QualityProfiles is a list of QualityProfile.
type SeriesList struct {
	ID     types.String `tfsdk:"id"`
	Series []Series     `tfsdk:"series"`
}

func (t dataSeriesType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "List all available series",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"series": {
				MarkdownDescription: "List of series",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						MarkdownDescription: "ID of tag",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"title": {
						MarkdownDescription: "Series Title",
						Required:            true,
						Type:                types.StringType,
					},
					"title_slug": {
						MarkdownDescription: "Series Title in kebab format",
						Required:            true,
						Type:                types.StringType,
					},
					"monitored": {
						MarkdownDescription: "Monitored flag",
						Required:            true,
						Type:                types.BoolType,
					},
					"season_folder": {
						MarkdownDescription: "Season Folder flag",
						Required:            true,
						Type:                types.BoolType,
					},
					"use_scene_numbering": {
						MarkdownDescription: "Scene numbering flag",
						Required:            true,
						Type:                types.BoolType,
					},
					"language_profile_id": {
						MarkdownDescription: "Language Profile ID ",
						Required:            true,
						Type:                types.Int64Type,
					},
					"quality_profile_id": {
						MarkdownDescription: "Quality Profile ID",
						Required:            true,
						Type:                types.Int64Type,
					},
					"tvdb_id": {
						MarkdownDescription: "TVDB ID",
						Required:            true,
						Type:                types.Int64Type,
					},
					"path": {
						MarkdownDescription: "Series Path",
						Required:            true,
						Type:                types.StringType,
					},
					"root_folder_path": {
						MarkdownDescription: "Series Root Folder",
						Required:            true,
						Type:                types.StringType,
					},
					"tags": {
						MarkdownDescription: "Tags",
						Optional:            true,
						Type: types.SetType{
							ElemType: types.Int64Type,
						},
					},
				}),
			},
		},
	}, nil
}

func (t dataSeriesType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return dataSeries{
		provider: provider,
	}, diags
}

func (d dataSeries) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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
	for _, s := range response {
		data.Series = append(data.Series, *writeSeries(ctx, s))
	}

	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.String{Value: strconv.Itoa(len(response))}
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
