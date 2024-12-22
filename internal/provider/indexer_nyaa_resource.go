package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	indexerNyaaResourceName   = "indexer_nyaa"
	indexerNyaaImplementation = "Nyaa"
	indexerNyaaConfigContract = "NyaaSettings"
	indexerNyaaProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerNyaaResource{}
	_ resource.ResourceWithImportState = &IndexerNyaaResource{}
)

func NewIndexerNyaaResource() resource.Resource {
	return &IndexerNyaaResource{}
}

// IndexerNyaaResource defines the Nyaa indexer implementation.
type IndexerNyaaResource struct {
	client *sonarr.APIClient
	auth   context.Context
}

// IndexerNyaa describes the Nyaa indexer data model.
type IndexerNyaa struct {
	SeedRatio                 types.Float64 `tfsdk:"seed_ratio"`
	Tags                      types.Set     `tfsdk:"tags"`
	Name                      types.String  `tfsdk:"name"`
	BaseURL                   types.String  `tfsdk:"base_url"`
	AdditionalParameters      types.String  `tfsdk:"additional_parameters"`
	DownloadClientID          types.Int64   `tfsdk:"download_client_id"`
	ID                        types.Int64   `tfsdk:"id"`
	MinimumSeeders            types.Int64   `tfsdk:"minimum_seeders"`
	SeasonPackSeedTime        types.Int64   `tfsdk:"season_pack_seed_time"`
	SeedTime                  types.Int64   `tfsdk:"seed_time"`
	Priority                  types.Int64   `tfsdk:"priority"`
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
		ConfigContract:            types.StringValue(indexerNyaaConfigContract),
		Implementation:            types.StringValue(indexerNyaaImplementation),
		Protocol:                  types.StringValue(indexerNyaaProtocol),
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

func (r *IndexerNyaaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerNyaaResourceName
}

func (r *IndexerNyaaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexers -->\nIndexer Nyaa resource.\nFor more information refer to [Indexer](https://wiki.servarr.com/sonarr/settings#indexers) and [Nyaa](https://wiki.servarr.com/sonarr/supported#nyaa).",
		Attributes: map[string]schema.Attribute{
			"enable_automatic_search": schema.BoolAttribute{
				MarkdownDescription: "Enable automatic search flag.",
				Optional:            true,
				Computed:            true,
			},
			"enable_interactive_search": schema.BoolAttribute{
				MarkdownDescription: "Enable interactive search flag.",
				Optional:            true,
				Computed:            true,
			},
			"enable_rss": schema.BoolAttribute{
				MarkdownDescription: "Enable RSS flag.",
				Optional:            true,
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.",
				Optional:            true,
				Computed:            true,
			},
			"download_client_id": schema.Int64Attribute{
				MarkdownDescription: "Download client ID.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "IndexerNyaa name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "IndexerNyaa ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"anime_standard_format_search": schema.BoolAttribute{
				MarkdownDescription: "Search anime in standard format.",
				Optional:            true,
				Computed:            true,
			},
			"minimum_seeders": schema.Int64Attribute{
				MarkdownDescription: "Minimum seeders.",
				Optional:            true,
				Computed:            true,
			},
			"season_pack_seed_time": schema.Int64Attribute{
				MarkdownDescription: "Season seed time.",
				Optional:            true,
				Computed:            true,
			},
			"seed_time": schema.Int64Attribute{
				MarkdownDescription: "Seed time.",
				Optional:            true,
				Computed:            true,
			},
			"seed_ratio": schema.Float64Attribute{
				MarkdownDescription: "Seed ratio.",
				Optional:            true,
				Computed:            true,
			},
			"additional_parameters": schema.StringAttribute{
				MarkdownDescription: "Additional parameters.",
				Optional:            true,
				Computed:            true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Required:            true,
			},
		},
	}
}

func (r *IndexerNyaaResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *IndexerNyaaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var indexer *IndexerNyaa

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerNyaa
	request := indexer.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.IndexerAPI.CreateIndexer(r.auth).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, indexerNyaaResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerNyaaResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response, &resp.Diagnostics)
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
	response, _, err := r.client.IndexerAPI.GetIndexerById(r.auth, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerNyaaResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerNyaaResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	indexer.write(ctx, response, &resp.Diagnostics)
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
	request := indexer.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.IndexerAPI.UpdateIndexer(r.auth, request.GetId()).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, indexerNyaaResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerNyaaResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerNyaaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerNyaa current value
	_, err := r.client.IndexerAPI.DeleteIndexer(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, indexerNyaaResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerNyaaResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerNyaaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+indexerNyaaResourceName+": "+req.ID)
}

func (i *IndexerNyaa) write(ctx context.Context, indexer *sonarr.IndexerResource, diags *diag.Diagnostics) {
	genericIndexer := i.toIndexer()
	genericIndexer.write(ctx, indexer, diags)
	i.fromIndexer(genericIndexer)
}

func (i *IndexerNyaa) read(ctx context.Context, diags *diag.Diagnostics) *sonarr.IndexerResource {
	return i.toIndexer().read(ctx, diags)
}
