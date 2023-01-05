package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const indexersDataSourceName = "indexers"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IndexersDataSource{}

func NewIndexersDataSource() datasource.DataSource {
	return &IndexersDataSource{}
}

// IndexersDataSource defines the indexers implementation.
type IndexersDataSource struct {
	client *sonarr.APIClient
}

// Indexers describes the indexers data model.
type Indexers struct {
	Indexers types.Set    `tfsdk:"indexers"`
	ID       types.String `tfsdk:"id"`
}

func (d *IndexersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexersDataSourceName
}

func (d *IndexersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Indexers -->List all available [Indexers](../resources/indexer).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"indexers": schema.SetNestedAttribute{
				MarkdownDescription: "Indexer list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"enable_automatic_search": schema.BoolAttribute{
							MarkdownDescription: "Enable automatic search flag.",
							Computed:            true,
						},
						"enable_interactive_search": schema.BoolAttribute{
							MarkdownDescription: "Enable interactive search flag.",
							Computed:            true,
						},
						"enable_rss": schema.BoolAttribute{
							MarkdownDescription: "Enable RSS flag.",
							Computed:            true,
						},
						"priority": schema.Int64Attribute{
							MarkdownDescription: "Priority.",
							Computed:            true,
						},
						"download_client_id": schema.Int64Attribute{
							MarkdownDescription: "Download client ID.",
							Computed:            true,
						},
						"config_contract": schema.StringAttribute{
							MarkdownDescription: "Indexer configuration template.",
							Computed:            true,
						},
						"implementation": schema.StringAttribute{
							MarkdownDescription: "Indexer implementation name.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Indexer name.",
							Computed:            true,
						},
						"protocol": schema.StringAttribute{
							MarkdownDescription: "Protocol. Valid values are 'usenet' and 'torrent'.",
							Computed:            true,
						},
						"tags": schema.SetAttribute{
							MarkdownDescription: "List of associated tags.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
						"id": schema.Int64Attribute{
							MarkdownDescription: "Indexer ID.",
							Computed:            true,
						},
						// Field values
						"allow_zero_size": schema.BoolAttribute{
							MarkdownDescription: "Allow zero size files.",
							Computed:            true,
						},
						"anime_standard_format_search": schema.BoolAttribute{
							MarkdownDescription: "Search anime in standard format.",
							Computed:            true,
						},
						"ranked_only": schema.BoolAttribute{
							MarkdownDescription: "Allow ranked only.",
							Computed:            true,
						},
						"delay": schema.Int64Attribute{
							MarkdownDescription: "Delay before grabbing.",
							Computed:            true,
						},
						"minimum_seeders": schema.Int64Attribute{
							MarkdownDescription: "Minimum seeders.",
							Computed:            true,
						},
						"season_pack_seed_time": schema.Int64Attribute{
							MarkdownDescription: "Season seed time.",
							Computed:            true,
						},
						"seed_time": schema.Int64Attribute{
							MarkdownDescription: "Seed time.",
							Computed:            true,
						},
						"seed_ratio": schema.Float64Attribute{
							MarkdownDescription: "Seed ratio.",
							Computed:            true,
						},
						"additional_parameters": schema.StringAttribute{
							MarkdownDescription: "Additional parameters.",
							Computed:            true,
						},
						"api_key": schema.StringAttribute{
							MarkdownDescription: "API key.",
							Computed:            true,
							Sensitive:           true,
						},
						"api_path": schema.StringAttribute{
							MarkdownDescription: "API path.",
							Computed:            true,
						},
						"base_url": schema.StringAttribute{
							MarkdownDescription: "Base URL.",
							Computed:            true,
						},
						"captcha_token": schema.StringAttribute{
							MarkdownDescription: "Captcha token.",
							Computed:            true,
						},
						"cookie": schema.StringAttribute{
							MarkdownDescription: "Cookie.",
							Computed:            true,
						},
						"passkey": schema.StringAttribute{
							MarkdownDescription: "Passkey.",
							Computed:            true,
							Sensitive:           true,
						},
						"username": schema.StringAttribute{
							MarkdownDescription: "Username.",
							Computed:            true,
						},
						"categories": schema.SetAttribute{
							MarkdownDescription: "Series list.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
						"anime_categories": schema.SetAttribute{
							MarkdownDescription: "Anime list.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
					},
				},
			},
		},
	}
}

func (d *IndexersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *sonarr.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *IndexersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Indexers

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get indexers current value
	response, _, err := d.client.IndexerApi.ListIndexer(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", indexersDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexersDataSourceName)
	// Map response body to resource schema attribute
	indexers := make([]Indexer, len(response))
	for j, i := range response {
		indexers[j].Tags = types.SetNull(types.Int64Type)
		indexers[j].write(ctx, i)
	}

	tfsdk.ValueFrom(ctx, indexers, data.Indexers.Type(ctx), &data.Indexers)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
