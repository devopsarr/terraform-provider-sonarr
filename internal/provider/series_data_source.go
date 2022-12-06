package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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

func (d *SeriesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "<!-- subcategory:Series -->Single [Series](../resources/series).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Series ID.",
				Computed:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Series Title.",
				Required:            true,
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
			"language_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Language Profile ID .",
				Computed:            true,
			},
			"quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Quality Profile ID.",
				Computed:            true,
			},
			"tvdb_id": schema.Int64Attribute{
				MarkdownDescription: "TVDB ID.",
				Computed:            true,
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

func (d *SeriesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *SeriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Series

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get series current value
	response, err := d.client.GetAllSeriesContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", seriesDataSourceName, err))

		return
	}

	series, err := findSeries(data.Title.ValueString(), response)
	if err != nil {
		resp.Diagnostics.AddError(tools.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", seriesDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+seriesDataSourceName)
	data.write(ctx, series)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findSeries(title string, series []*sonarr.Series) (*sonarr.Series, error) {
	for _, s := range series {
		if s.Title == title {
			return s, nil
		}
	}

	return nil, tools.ErrDataNotFoundError(seriesDataSourceName, "title", title)
}
