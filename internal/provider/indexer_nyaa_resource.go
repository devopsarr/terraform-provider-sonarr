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
	indexerNyaaResourceName   = "indexer_nyaa"
	IndexerNyaaImplementation = "Nyaa"
	IndexerNyaaConfigContrat  = "NyaaSettings"
	IndexerNyaaProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IndexerNyaaResource{}
var _ resource.ResourceWithImportState = &IndexerNyaaResource{}

func NewIndexerNyaaResource() resource.Resource {
	return &IndexerNyaaResource{}
}

// IndexerNyaaResource defines the Nyaa indexer implementation.
type IndexerNyaaResource struct {
	client *sonarr.Sonarr
}

// IndexerNyaa describes the Nyaa indexer data model.
type IndexerNyaa struct {
	Tags                      types.Set     `tfsdk:"tags"`
	Name                      types.String  `tfsdk:"name"`
	BaseURL                   types.String  `tfsdk:"base_url"`
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

func (i IndexerNyaa) toIndexer() *Indexer {
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
		Tags:                      i.Tags,
	}
}

func (i *IndexerNyaa) fromIndexer(indexer *Indexer) {
	i.EnableAutomaticSearch = indexer.EnableAutomaticSearch
	i.EnableInteractiveSearch = indexer.EnableInteractiveSearch
	i.EnableRss = indexer.EnableRss
	i.AnimeStandardFormatSearch = indexer.AnimeStandardFormatSearch
	i.Priority = indexer.Priority
	i.DownloadClientID = indexer.DownloadClientID
	i.ID = indexer.ID
	i.Name = indexer.Name
	i.AdditionalParameters = indexer.AdditionalParameters
	i.MinimumSeeders = indexer.MinimumSeeders
	i.SeasonPackSeedTime = indexer.SeasonPackSeedTime
	i.SeedTime = indexer.SeedTime
	i.SeedRatio = indexer.SeedRatio
	i.BaseURL = indexer.BaseURL
	i.Tags = indexer.Tags
}

func (r *IndexerNyaaResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerNyaaResourceName
}

func (r *IndexerNyaaResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "<!-- subcategory:Indexers -->Indexer Nyaa resource.\nFor more information refer to [Indexer](https://wiki.servarr.com/sonarr/settings#indexers) and [Nyaa](https://wiki.servarr.com/sonarr/supported#nyaa).",
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
				MarkdownDescription: "IndexerNyaa name.",
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
				MarkdownDescription: "IndexerNyaa ID.",
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
		},
	}, nil
}

func (r *IndexerNyaaResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IndexerNyaaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var indexer *IndexerNyaa

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerNyaa
	request := indexer.read(ctx)

	response, err := r.client.AddIndexerContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", indexerNyaaResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerNyaaResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerNyaaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var indexer *IndexerNyaa

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IndexerNyaa current value
	response, err := r.client.GetIndexerContext(ctx, indexer.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", indexerNyaaResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerNyaaResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerNyaaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var indexer *IndexerNyaa

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update IndexerNyaa
	request := indexer.read(ctx)

	response, err := r.client.UpdateIndexerContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update "+indexerNyaaResourceName+", got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerNyaaResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerNyaaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var indexer *IndexerNyaa

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerNyaa current value
	err := r.client.DeleteIndexerContext(ctx, indexer.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", indexerNyaaResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerNyaaResourceName+": "+strconv.Itoa(int(indexer.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerNyaaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+indexerNyaaResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (i *IndexerNyaa) write(ctx context.Context, indexer *sonarr.IndexerOutput) {
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

func (i *IndexerNyaa) read(ctx context.Context) *sonarr.IndexerInput {
	var tags []int

	tfsdk.ValueAs(ctx, i.Tags, &tags)

	return &sonarr.IndexerInput{
		EnableAutomaticSearch:   i.EnableAutomaticSearch.ValueBool(),
		EnableInteractiveSearch: i.EnableInteractiveSearch.ValueBool(),
		EnableRss:               i.EnableRss.ValueBool(),
		Priority:                i.Priority.ValueInt64(),
		DownloadClientID:        i.DownloadClientID.ValueInt64(),
		ID:                      i.ID.ValueInt64(),
		ConfigContract:          IndexerNyaaConfigContrat,
		Implementation:          IndexerNyaaImplementation,
		Name:                    i.Name.ValueString(),
		Protocol:                IndexerNyaaProtocol,
		Tags:                    tags,
		Fields:                  i.toIndexer().readFields(ctx),
	}
}
