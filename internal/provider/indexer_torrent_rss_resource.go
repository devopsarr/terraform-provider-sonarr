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
	indexerTorrentRssResourceName   = "indexer_torrent_rss"
	indexerTorrentRssImplementation = "TorrentRssIndexer"
	indexerTorrentRssConfigContract = "TorrentRssIndexerSettings"
	indexerTorrentRssProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerTorrentRssResource{}
	_ resource.ResourceWithImportState = &IndexerTorrentRssResource{}
)

func NewIndexerTorrentRssResource() resource.Resource {
	return &IndexerTorrentRssResource{}
}

// IndexerTorrentRssResource defines the TorrentRss indexer implementation.
type IndexerTorrentRssResource struct {
	client *sonarr.APIClient
}

// IndexerTorrentRss describes the TorrentRss indexer data model.
type IndexerTorrentRss struct {
	SeedRatio               types.Float64 `tfsdk:"seed_ratio"`
	Tags                    types.Set     `tfsdk:"tags"`
	Name                    types.String  `tfsdk:"name"`
	BaseURL                 types.String  `tfsdk:"base_url"`
	Cookie                  types.String  `tfsdk:"cookie"`
	DownloadClientID        types.Int64   `tfsdk:"download_client_id"`
	ID                      types.Int64   `tfsdk:"id"`
	MinimumSeeders          types.Int64   `tfsdk:"minimum_seeders"`
	SeasonPackSeedTime      types.Int64   `tfsdk:"season_pack_seed_time"`
	SeedTime                types.Int64   `tfsdk:"seed_time"`
	Priority                types.Int64   `tfsdk:"priority"`
	AllowZeroSize           types.Bool    `tfsdk:"allow_zero_size"`
	EnableAutomaticSearch   types.Bool    `tfsdk:"enable_automatic_search"`
	EnableRss               types.Bool    `tfsdk:"enable_rss"`
	EnableInteractiveSearch types.Bool    `tfsdk:"enable_interactive_search"`
}

func (i IndexerTorrentRss) toIndexer() *Indexer {
	return &Indexer{
		EnableAutomaticSearch:   i.EnableAutomaticSearch,
		EnableInteractiveSearch: i.EnableInteractiveSearch,
		EnableRss:               i.EnableRss,
		AllowZeroSize:           i.AllowZeroSize,
		Priority:                i.Priority,
		DownloadClientID:        i.DownloadClientID,
		ID:                      i.ID,
		Name:                    i.Name,
		Cookie:                  i.Cookie,
		MinimumSeeders:          i.MinimumSeeders,
		SeasonPackSeedTime:      i.SeasonPackSeedTime,
		SeedTime:                i.SeedTime,
		SeedRatio:               i.SeedRatio,
		BaseURL:                 i.BaseURL,
		Tags:                    i.Tags,
		ConfigContract:          types.StringValue(indexerTorrentRssConfigContract),
		Implementation:          types.StringValue(indexerTorrentRssImplementation),
		Protocol:                types.StringValue(indexerTorrentRssProtocol),
	}
}

func (i *IndexerTorrentRss) fromIndexer(indexer *Indexer) {
	i.EnableAutomaticSearch = indexer.EnableAutomaticSearch
	i.EnableInteractiveSearch = indexer.EnableInteractiveSearch
	i.EnableRss = indexer.EnableRss
	i.AllowZeroSize = indexer.AllowZeroSize
	i.Priority = indexer.Priority
	i.DownloadClientID = indexer.DownloadClientID
	i.ID = indexer.ID
	i.Name = indexer.Name
	i.Cookie = indexer.Cookie
	i.MinimumSeeders = indexer.MinimumSeeders
	i.SeasonPackSeedTime = indexer.SeasonPackSeedTime
	i.SeedTime = indexer.SeedTime
	i.SeedRatio = indexer.SeedRatio
	i.BaseURL = indexer.BaseURL
	i.Tags = indexer.Tags
}

func (r *IndexerTorrentRssResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerTorrentRssResourceName
}

func (r *IndexerTorrentRssResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexers -->Indexer Torrent RSS resource.\nFor more information refer to [Indexer](https://wiki.servarr.com/sonarr/settings#indexers) and [Torrent RSS](https://wiki.servarr.com/sonarr/supported#torrentrssindexer).",
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
				MarkdownDescription: "IndexerTorrentRss name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "IndexerTorrentRss ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"allow_zero_size": schema.BoolAttribute{
				MarkdownDescription: "Allow zero size files.",
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
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Required:            true,
			},
			"cookie": schema.StringAttribute{
				MarkdownDescription: "Cookie.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *IndexerTorrentRssResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *IndexerTorrentRssResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var indexer *IndexerTorrentRss

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerTorrentRss
	request := indexer.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.IndexerApi.CreateIndexer(ctx).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, indexerTorrentRssResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerTorrentRssResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerTorrentRssResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var indexer *IndexerTorrentRss

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IndexerTorrentRss current value
	response, _, err := r.client.IndexerApi.GetIndexerById(ctx, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerTorrentRssResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerTorrentRssResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	indexer.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerTorrentRssResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var indexer *IndexerTorrentRss

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update IndexerTorrentRss
	request := indexer.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.IndexerApi.UpdateIndexer(ctx, strconv.Itoa(int(request.GetId()))).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, indexerTorrentRssResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerTorrentRssResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerTorrentRssResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerTorrentRss current value
	_, err := r.client.IndexerApi.DeleteIndexer(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, indexerTorrentRssResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerTorrentRssResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerTorrentRssResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+indexerTorrentRssResourceName+": "+req.ID)
}

func (i *IndexerTorrentRss) write(ctx context.Context, indexer *sonarr.IndexerResource, diags *diag.Diagnostics) {
	genericIndexer := i.toIndexer()
	genericIndexer.write(ctx, indexer, diags)
	i.fromIndexer(genericIndexer)
}

func (i *IndexerTorrentRss) read(ctx context.Context, diags *diag.Diagnostics) *sonarr.IndexerResource {
	return i.toIndexer().read(ctx, diags)
}
