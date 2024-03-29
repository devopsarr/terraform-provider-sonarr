package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const searchSearchSeriesDataSourceName = "search_series"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SearchSeriesDataSource{}

func NewSearchSeriesDataSource() datasource.DataSource {
	return &SearchSeriesDataSource{}
}

// SearchSeriesDataSource defines the tags implementation.
type SearchSeriesDataSource struct {
	client *sonarr.APIClient
	auth   context.Context
}

func (d *SearchSeriesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + searchSearchSeriesDataSourceName
}

func (d *SearchSeriesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Series -->\nSearch a Single [Series](../resources/series) via tvdb_id.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Series ID.",
				Computed:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Series Title.",
				Computed:            true,
			},
			"title_slug": schema.StringAttribute{
				MarkdownDescription: "Series Title in kebab format.",
				Computed:            true,
			},
			"monitored": schema.BoolAttribute{
				MarkdownDescription: "Monitored flag.",
				Computed:            true,
			},
			"season_folder": schema.BoolAttribute{
				MarkdownDescription: "Season Folder flag.",
				Computed:            true,
			},
			"use_scene_numbering": schema.BoolAttribute{
				MarkdownDescription: "Scene numbering flag.",
				Computed:            true,
			},
			"quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Quality Profile ID.",
				Computed:            true,
			},
			"tvdb_id": schema.Int64Attribute{
				MarkdownDescription: "TVDB ID.",
				Required:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Series Path.",
				Computed:            true,
			},
			"root_folder_path": schema.StringAttribute{
				MarkdownDescription: "Series Root Folder.",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (d *SearchSeriesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *SearchSeriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Series

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get series current value
	response, _, err := d.client.SeriesLookupAPI.ListSeriesLookup(d.auth).Term(strconv.Itoa(int(data.TvdbID.ValueInt64()))).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, searchSearchSeriesDataSourceName, err))

		return
	}

	if !(int64(response[0].GetTvdbId()) == data.TvdbID.ValueInt64()) {
		resp.Diagnostics.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(searchSearchSeriesDataSourceName, "TVDBID", strconv.Itoa(int(data.TvdbID.ValueInt64()))))

		return
	}

	tflog.Trace(ctx, "read "+searchSearchSeriesDataSourceName)
	data.write(ctx, &response[0], &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
