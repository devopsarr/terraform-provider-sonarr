package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ provider.DataSourceType = dataSeriesType{}
	_ datasource.DataSource   = dataSeries{}
)

type dataSeriesType struct{}

type dataSeries struct {
	provider sonarrProvider
}

func (t dataSeriesType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Single [Series](../resources/series).",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Series ID.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"title": {
				MarkdownDescription: "Series Title.",
				Required:            true,
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
				Optional:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
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
	var data Series
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

	series, err := findSeries(data.Title.Value, response)
	if err != nil {
		resp.Diagnostics.AddError("Data Source Error", fmt.Sprintf("Unable to find series, got error: %s", err))

		return
	}

	result := writeSeries(ctx, series)
	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
}

func findSeries(title string, series []*sonarr.Series) (*sonarr.Series, error) {
	for _, s := range series {
		if s.Title == title {
			return s, nil
		}
	}

	return nil, fmt.Errorf("no series with title %s", title)
}
