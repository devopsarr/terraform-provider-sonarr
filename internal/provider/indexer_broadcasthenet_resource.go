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
	indexerBroadcastheNetResourceName   = "indexer_broadcasthenet"
	indexerBroadcastheNetImplementation = "BroadcastheNet"
	indexerBroadcastheNetConfigContract = "BroadcastheNetSettings"
	indexerBroadcastheNetProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerBroadcastheNetResource{}
	_ resource.ResourceWithImportState = &IndexerBroadcastheNetResource{}
)

func NewIndexerBroadcastheNetResource() resource.Resource {
	return &IndexerBroadcastheNetResource{}
}

// IndexerBroadcastheNetResource defines the BroadcastheNet indexer implementation.
type IndexerBroadcastheNetResource struct {
	client *sonarr.APIClient
	auth   context.Context
}

// IndexerBroadcastheNet describes the BroadcastheNet indexer data model.
type IndexerBroadcastheNet struct {
	SeedRatio               types.Float64 `tfsdk:"seed_ratio"`
	Tags                    types.Set     `tfsdk:"tags"`
	Name                    types.String  `tfsdk:"name"`
	APIKey                  types.String  `tfsdk:"api_key"`
	BaseURL                 types.String  `tfsdk:"base_url"`
	ID                      types.Int64   `tfsdk:"id"`
	Priority                types.Int64   `tfsdk:"priority"`
	DownloadClientID        types.Int64   `tfsdk:"download_client_id"`
	MinimumSeeders          types.Int64   `tfsdk:"minimum_seeders"`
	SeasonPackSeedTime      types.Int64   `tfsdk:"season_pack_seed_time"`
	SeedTime                types.Int64   `tfsdk:"seed_time"`
	EnableAutomaticSearch   types.Bool    `tfsdk:"enable_automatic_search"`
	EnableRss               types.Bool    `tfsdk:"enable_rss"`
	EnableInteractiveSearch types.Bool    `tfsdk:"enable_interactive_search"`
}

func (i IndexerBroadcastheNet) toIndexer() *Indexer {
	return &Indexer{
		EnableAutomaticSearch:   i.EnableAutomaticSearch,
		EnableInteractiveSearch: i.EnableInteractiveSearch,
		EnableRss:               i.EnableRss,
		Priority:                i.Priority,
		DownloadClientID:        i.DownloadClientID,
		ID:                      i.ID,
		Name:                    i.Name,
		MinimumSeeders:          i.MinimumSeeders,
		SeasonPackSeedTime:      i.SeasonPackSeedTime,
		SeedTime:                i.SeedTime,
		SeedRatio:               i.SeedRatio,
		APIKey:                  i.APIKey,
		BaseURL:                 i.BaseURL,
		Tags:                    i.Tags,
		ConfigContract:          types.StringValue(indexerBroadcastheNetConfigContract),
		Implementation:          types.StringValue(indexerBroadcastheNetImplementation),
		Protocol:                types.StringValue(indexerBroadcastheNetProtocol),
	}
}

func (i *IndexerBroadcastheNet) fromIndexer(indexer *Indexer) {
	i.EnableAutomaticSearch = indexer.EnableAutomaticSearch
	i.EnableInteractiveSearch = indexer.EnableInteractiveSearch
	i.EnableRss = indexer.EnableRss
	i.Priority = indexer.Priority
	i.DownloadClientID = indexer.DownloadClientID
	i.ID = indexer.ID
	i.Name = indexer.Name
	i.MinimumSeeders = indexer.MinimumSeeders
	i.SeasonPackSeedTime = indexer.SeasonPackSeedTime
	i.SeedTime = indexer.SeedTime
	i.SeedRatio = indexer.SeedRatio
	i.APIKey = indexer.APIKey
	i.BaseURL = indexer.BaseURL
	i.Tags = indexer.Tags
}

func (r *IndexerBroadcastheNetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerBroadcastheNetResourceName
}

func (r *IndexerBroadcastheNetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexers -->\nIndexer BroadcastheNet resource.\nFor more information refer to [Indexer](https://wiki.servarr.com/sonarr/settings#indexers) and [BroadcastheNet](https://wiki.servarr.com/sonarr/supported#broadcasthenet).",
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
				MarkdownDescription: "IndexerBroadcastheNet name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "IndexerBroadcastheNet ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
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
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *IndexerBroadcastheNetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *IndexerBroadcastheNetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var indexer *IndexerBroadcastheNet

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerBroadcastheNet
	request := indexer.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.IndexerAPI.CreateIndexer(r.auth).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, indexerBroadcastheNetResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerBroadcastheNetResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerBroadcastheNetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var indexer *IndexerBroadcastheNet

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IndexerBroadcastheNet current value
	response, _, err := r.client.IndexerAPI.GetIndexerById(r.auth, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerBroadcastheNetResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerBroadcastheNetResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	indexer.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerBroadcastheNetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var indexer *IndexerBroadcastheNet

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update IndexerBroadcastheNet
	request := indexer.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.IndexerAPI.UpdateIndexer(r.auth, request.GetId()).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, indexerBroadcastheNetResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerBroadcastheNetResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerBroadcastheNetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerBroadcastheNet current value
	_, err := r.client.IndexerAPI.DeleteIndexer(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, indexerBroadcastheNetResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerBroadcastheNetResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerBroadcastheNetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+indexerBroadcastheNetResourceName+": "+req.ID)
}

func (i *IndexerBroadcastheNet) write(ctx context.Context, indexer *sonarr.IndexerResource, diags *diag.Diagnostics) {
	genericIndexer := i.toIndexer()
	genericIndexer.write(ctx, indexer, diags)
	i.fromIndexer(genericIndexer)
}

func (i *IndexerBroadcastheNet) read(ctx context.Context, diags *diag.Diagnostics) *sonarr.IndexerResource {
	return i.toIndexer().read(ctx, diags)
}
