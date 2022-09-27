package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const seriesDataSourceName = "series"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SeriesDataSource{}

func NewSeriesDataSource() datasource.DataSource {
	return &SeriesDataSource{}
}

// SeriesDataSource defines the tags implementation.
type SeriesDataSource struct {
	client *sonarr.Sonarr
}

func (d *SeriesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + seriesDataSourceName
}

func (d *SeriesDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "[subcategory:Series]: #\nSingle [Series](../resources/series).",
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

func (d *SeriesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SeriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data Series

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get series current value
	response, err := d.client.GetAllSeriesContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", seriesDataSourceName, err))

		return
	}

	series, err := findSeries(data.Title.Value, response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", seriesDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+seriesDataSourceName)
	result := writeSeries(ctx, series)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func findSeries(title string, series []*sonarr.Series) (*sonarr.Series, error) {
	for _, s := range series {
		if s.Title == title {
			return s, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(seriesDataSourceName, "title", title)
}
