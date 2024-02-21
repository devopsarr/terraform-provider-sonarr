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

const allSeriesDataSourceName = "all_series"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AllSeriessDataSource{}

func NewAllSeriessDataSource() datasource.DataSource {
	return &AllSeriessDataSource{}
}

// AllSeriessDataSource defines the tags implementation.
type AllSeriessDataSource struct {
	client *sonarr.APIClient
	auth   context.Context
}

// AllSeriess describes the series(es) data model.
type SeriesList struct {
	Series types.Set    `tfsdk:"series"`
	ID     types.String `tfsdk:"id"`
}

func (d *AllSeriessDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + allSeriesDataSourceName
}

func (d *AllSeriessDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Series -->List all available [Series](../resources/series).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"series": schema.SetNestedAttribute{
				MarkdownDescription: "Series list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
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
				},
			},
		},
	}
}

func (d *AllSeriessDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *AllSeriessDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get series current value
	response, _, err := d.client.SeriesAPI.ListSeries(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, allSeriesDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+allSeriesDataSourceName)
	// Map response body to resource schema attribute
	series := make([]Series, len(response))
	for i, t := range response {
		series[i].write(ctx, &t, &resp.Diagnostics)
	}

	seriesList, diags := types.SetValueFrom(ctx, Series{}.getType(), series)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, SeriesList{Series: seriesList, ID: types.StringValue(strconv.Itoa(len(response)))})...)
}
