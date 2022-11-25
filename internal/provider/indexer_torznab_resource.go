package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const (
	indexerTorznabResourceName   = "indexer_torznab"
	IndexerTorznabImplementation = "Torznab"
	IndexerTorznabConfigContrat  = "TorznabSettings"
	IndexerTorznabProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IndexerTorznabResource{}
var _ resource.ResourceWithImportState = &IndexerTorznabResource{}

func NewIndexerTorznabResource() resource.Resource {
	return &IndexerTorznabResource{}
}

// IndexerTorznabResource defines the Torznab indexer implementation.
type IndexerTorznabResource struct {
	client *sonarr.Sonarr
}

// IndexerTorznab describes the Torznab indexer data model.
type IndexerTorznab struct {
	Tags                      types.Set     `tfsdk:"tags"`
	Categories                types.Set     `tfsdk:"categories"`
	AnimeCategories           types.Set     `tfsdk:"anime_categories"`
	Name                      types.String  `tfsdk:"name"`
	BaseURL                   types.String  `tfsdk:"base_url"`
	APIPath                   types.String  `tfsdk:"api_path"`
	APIKey                    types.String  `tfsdk:"api_key"`
	AdditionalParameters      types.String  `tfsdk:"additional_parameters"`
	Priority                  types.Int64   `tfsdk:"priority"`
	ID                        types.Int64   `tfsdk:"id"`
	DownloadClientID          types.Int64   `tfsdk:"download_client_id"`
	MinimumSeeders            types.Int64   `tfsdk:"minimum_seeders"`
	SeasonPackSeedTime        types.Int64   `tfsdk:"season_pack_seed_time"`
	SeedTime                  types.Int64   `tfsdk:"seed_time"`
	SeedRatio                 types.Float64 `tfsdk:"seed_ratio"`
	AnimeStandardFormatSearch types.Bool    `tfsdk:"anime_standard_format_search"`
	EnableAutomaticSearch     types.Bool    `tfsdk:"enable_automatic_search"`
	EnableRss                 types.Bool    `tfsdk:"enable_rss"`
	EnableInteractiveSearch   types.Bool    `tfsdk:"enable_interactive_search"`
}

func (i IndexerTorznab) toIndexer() *Indexer {
	return &Indexer{
		EnableAutomaticSearch:     i.EnableAutomaticSearch,
		EnableInteractiveSearch:   i.EnableInteractiveSearch,
		EnableRss:                 i.EnableRss,
		AnimeStandardFormatSearch: i.AnimeStandardFormatSearch,
		Priority:                  i.Priority,
		DownloadClientID:          i.DownloadClientID,
		ID:                        i.ID,
		Name:                      i.Name,
		AdditionalParameters:      i.AdditionalParameters,
		MinimumSeeders:            i.MinimumSeeders,
		SeasonPackSeedTime:        i.SeasonPackSeedTime,
		SeedTime:                  i.SeedTime,
		SeedRatio:                 i.SeedRatio,
		BaseURL:                   i.BaseURL,
		APIPath:                   i.APIPath,
		APIKey:                    i.APIKey,
		Tags:                      i.Tags,
		Categories:                i.Categories,
		AnimeCategories:           i.AnimeCategories,
	}
}

func (i *IndexerTorznab) fromIndexer(indexer *Indexer) {
	i.EnableAutomaticSearch = indexer.EnableAutomaticSearch
	i.EnableInteractiveSearch = indexer.EnableInteractiveSearch
	i.EnableRss = indexer.EnableRss
	i.AnimeStandardFormatSearch = indexer.AnimeStandardFormatSearch
	i.Priority = indexer.Priority
	i.DownloadClientID = indexer.DownloadClientID
	i.ID = indexer.ID
	i.Name = indexer.Name
	i.APIPath = indexer.APIPath
	i.APIKey = indexer.APIKey
	i.AdditionalParameters = indexer.AdditionalParameters
	i.MinimumSeeders = indexer.MinimumSeeders
	i.SeasonPackSeedTime = indexer.SeasonPackSeedTime
	i.SeedTime = indexer.SeedTime
	i.SeedRatio = indexer.SeedRatio
	i.BaseURL = indexer.BaseURL
	i.Tags = indexer.Tags
	i.Categories = indexer.Categories
	i.AnimeCategories = indexer.AnimeCategories
}

func (r *IndexerTorznabResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerTorznabResourceName
}

func (r *IndexerTorznabResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "<!-- subcategory:Indexers -->Indexer Torznab resource.\nFor more information refer to [Indexer](https://wiki.servarr.com/sonarr/settings#indexers) and [Torznab](https://wiki.servarr.com/sonarr/supported#torznab).",
		Attributes: map[string]tfsdk.Attribute{
			"enable_automatic_search": {
				MarkdownDescription: "Enable automatic search flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"enable_interactive_search": {
				MarkdownDescription: "Enable interactive search flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"enable_rss": {
				MarkdownDescription: "Enable RSS flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"priority": {
				MarkdownDescription: "Priority.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"download_client_id": {
				MarkdownDescription: "Download client ID.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"name": {
				MarkdownDescription: "IndexerTorznab name.",
				Required:            true,
				Type:                types.StringType,
			},
			"tags": {
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"id": {
				MarkdownDescription: "IndexerTorznab ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			// Field values
			"anime_standard_format_search": {
				MarkdownDescription: "Search anime in standard format.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"minimum_seeders": {
				MarkdownDescription: "Minimum seeders.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"season_pack_seed_time": {
				MarkdownDescription: "Season seed time.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"seed_time": {
				MarkdownDescription: "Seed time.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"seed_ratio": {
				MarkdownDescription: "Seed ratio.",
				Optional:            true,
				Computed:            true,
				Type:                types.Float64Type,
			},
			"additional_parameters": {
				MarkdownDescription: "Additional parameters.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"base_url": {
				MarkdownDescription: "Base URL.",
				Required:            true,
				Type:                types.StringType,
			},
			"api_key": {
				MarkdownDescription: "API key.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				Type:                types.StringType,
			},
			"api_path": {
				MarkdownDescription: "API path.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"categories": {
				MarkdownDescription: "Categories list.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"anime_categories": {
				MarkdownDescription: "Anime categories list.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
		},
	}, nil
}

func (r *IndexerTorznabResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *IndexerTorznabResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var indexer *IndexerTorznab

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerTorznab
	request := indexer.read(ctx)

	response, err := r.client.AddIndexerContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", indexerTorznabResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerTorznabResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerTorznabResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var indexer *IndexerTorznab

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IndexerTorznab current value
	response, err := r.client.GetIndexerContext(ctx, indexer.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", indexerTorznabResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerTorznabResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerTorznabResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var indexer *IndexerTorznab

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update IndexerTorznab
	request := indexer.read(ctx)

	response, err := r.client.UpdateIndexerContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update "+indexerTorznabResourceName+", got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerTorznabResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerTorznabResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var indexer *IndexerTorznab

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerTorznab current value
	err := r.client.DeleteIndexerContext(ctx, indexer.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", indexerTorznabResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerTorznabResourceName+": "+strconv.Itoa(int(indexer.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerTorznabResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+indexerTorznabResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (i *IndexerTorznab) write(ctx context.Context, indexer *sonarr.IndexerOutput) {
	genericIndexer := Indexer{
		EnableAutomaticSearch:   types.BoolValue(indexer.EnableAutomaticSearch),
		EnableInteractiveSearch: types.BoolValue(indexer.EnableInteractiveSearch),
		EnableRss:               types.BoolValue(indexer.EnableRss),
		Priority:                types.Int64Value(indexer.Priority),
		DownloadClientID:        types.Int64Value(indexer.DownloadClientID),
		ID:                      types.Int64Value(indexer.ID),
		Name:                    types.StringValue(indexer.Name),
		Tags:                    types.SetValueMust(types.Int64Type, nil),
	}
	tfsdk.ValueFrom(ctx, indexer.Tags, genericIndexer.Tags.Type(ctx), &genericIndexer.Tags)
	genericIndexer.writeFields(ctx, indexer.Fields)
	i.fromIndexer(&genericIndexer)
}

func (i *IndexerTorznab) read(ctx context.Context) *sonarr.IndexerInput {
	var tags []int

	tfsdk.ValueAs(ctx, i.Tags, &tags)

	return &sonarr.IndexerInput{
		EnableAutomaticSearch:   i.EnableAutomaticSearch.ValueBool(),
		EnableInteractiveSearch: i.EnableInteractiveSearch.ValueBool(),
		EnableRss:               i.EnableRss.ValueBool(),
		Priority:                i.Priority.ValueInt64(),
		DownloadClientID:        i.DownloadClientID.ValueInt64(),
		ID:                      i.ID.ValueInt64(),
		ConfigContract:          IndexerTorznabConfigContrat,
		Implementation:          IndexerTorznabImplementation,
		Name:                    i.Name.ValueString(),
		Protocol:                IndexerTorznabProtocol,
		Tags:                    tags,
		Fields:                  i.toIndexer().readFields(ctx),
	}
}
