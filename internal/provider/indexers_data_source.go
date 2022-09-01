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
var (
	_ provider.DataSourceType = dataIndexersType{}
	_ datasource.DataSource   = dataIndexers{}
)

type dataIndexersType struct{}

type dataIndexers struct {
	provider sonarrProvider
}

// Indexers is a list of Indexer.
type Indexers struct {
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	ID       types.String `tfsdk:"id"`
	Indexers types.Set    `tfsdk:"indexers"`
}

func (t dataIndexersType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "List all available [Indexers](../resources/indexer).",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"indexers": {
				MarkdownDescription: "Indexer list.",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"enable_automatic_search": {
						MarkdownDescription: "Enable automatic search flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"enable_interactive_search": {
						MarkdownDescription: "Enable interactive search flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"enable_rss": {
						MarkdownDescription: "Enable RSS flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"priority": {
						MarkdownDescription: "Priority.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"download_client_id": {
						MarkdownDescription: "Download client ID.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"config_contract": {
						MarkdownDescription: "Indexer configuration template.",
						Computed:            true,
						Type:                types.StringType,
					},
					"implementation": {
						MarkdownDescription: "Indexer implementation name.",
						Computed:            true,
						Type:                types.StringType,
					},
					"name": {
						MarkdownDescription: "Indexer name.",
						Computed:            true,
						Type:                types.StringType,
					},
					"protocol": {
						MarkdownDescription: "Protocol. Valid values are 'usenet' and 'torrent'.",
						Computed:            true,
						Type:                types.StringType,
					},
					"tags": {
						MarkdownDescription: "List of associated tags.",
						Computed:            true,
						Type: types.SetType{
							ElemType: types.Int64Type,
						},
					},
					"id": {
						MarkdownDescription: "Indexer ID.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					// Field values
					"allow_zero_size": {
						MarkdownDescription: "Allow zero size files.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"anime_standard_format_search": {
						MarkdownDescription: "Search anime in standard format.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"ranked_only": {
						MarkdownDescription: "Allow ranked only.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"delay": {
						MarkdownDescription: "Delay before grabbing.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"minimum_seeders": {
						MarkdownDescription: "Minimum seeders.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"season_pack_seed_time": {
						MarkdownDescription: "Season seed time.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"seed_time": {
						MarkdownDescription: "Seed time.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"seed_ratio": {
						MarkdownDescription: "Seed ratio.",
						Computed:            true,
						Type:                types.Float64Type,
					},
					"additional_parameters": {
						MarkdownDescription: "Additional parameters.",
						Computed:            true,
						Type:                types.StringType,
					},
					"api_key": {
						MarkdownDescription: "API key.",
						Computed:            true,
						Type:                types.StringType,
					},
					"api_path": {
						MarkdownDescription: "API path.",
						Computed:            true,
						Type:                types.StringType,
					},
					"base_url": {
						MarkdownDescription: "Base URL.",
						Computed:            true,
						Type:                types.StringType,
					},
					"captcha_token": {
						MarkdownDescription: "Captcha token.",
						Computed:            true,
						Type:                types.StringType,
					},
					"cookie": {
						MarkdownDescription: "Cookie.",
						Computed:            true,
						Type:                types.StringType,
					},
					"passkey": {
						MarkdownDescription: "Passkey.",
						Computed:            true,
						Type:                types.StringType,
					},
					"username": {
						MarkdownDescription: "Username.",
						Computed:            true,
						Type:                types.StringType,
					},
					"categories": {
						MarkdownDescription: "Series list.",
						Computed:            true,
						Type: types.SetType{
							ElemType: types.Int64Type,
						},
					},
					"anime_categories": {
						MarkdownDescription: "Anime list.",
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

func (t dataIndexersType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return dataIndexers{
		provider: provider,
	}, diags
}

func (d dataIndexers) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data Indexers
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get indexers current value
	response, err := d.provider.client.GetIndexersContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read indexers, got error: %s", err))

		return
	}
	// Map response body to resource schema attribute
	profiles := *writeIndexers(ctx, response)
	tfsdk.ValueFrom(ctx, profiles, data.Indexers.Type(context.Background()), &data.Indexers)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.String{Value: strconv.Itoa(len(response))}
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func writeIndexers(ctx context.Context, delays []*sonarr.IndexerOutput) *[]Indexer {
	output := make([]Indexer, len(delays))
	for i, p := range delays {
		output[i] = *writeIndexer(ctx, p)
	}

	return &output
}
