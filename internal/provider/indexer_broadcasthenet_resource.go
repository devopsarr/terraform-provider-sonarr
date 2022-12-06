package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const (
	indexerBroadcastheNetResourceName   = "indexer_broadcasthenet"
	IndexerBroadcastheNetImplementation = "BroadcastheNet"
	IndexerBroadcastheNetConfigContrat  = "BroadcastheNetSettings"
	IndexerBroadcastheNetProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IndexerBroadcastheNetResource{}
var _ resource.ResourceWithImportState = &IndexerBroadcastheNetResource{}

func NewIndexerBroadcastheNetResource() resource.Resource {
	return &IndexerBroadcastheNetResource{}
}

// IndexerBroadcastheNetResource defines the BroadcastheNet indexer implementation.
type IndexerBroadcastheNetResource struct {
	client *sonarr.Sonarr
}

// IndexerBroadcastheNet describes the BroadcastheNet indexer data model.
type IndexerBroadcastheNet struct {
	Tags                    types.Set     `tfsdk:"tags"`
	Name                    types.String  `tfsdk:"name"`
	APIKey                  types.String  `tfsdk:"api_key"`
	BaseURL                 types.String  `tfsdk:"base_url"`
	Priority                types.Int64   `tfsdk:"priority"`
	ID                      types.Int64   `tfsdk:"id"`
	DownloadClientID        types.Int64   `tfsdk:"download_client_id"`
	MinimumSeeders          types.Int64   `tfsdk:"minimum_seeders"`
	SeasonPackSeedTime      types.Int64   `tfsdk:"season_pack_seed_time"`
	SeedTime                types.Int64   `tfsdk:"seed_time"`
	SeedRatio               types.Float64 `tfsdk:"seed_ratio"`
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

func (r *IndexerBroadcastheNetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerBroadcastheNetResourceName
}

func (r *IndexerBroadcastheNetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexers -->Indexer BroadcastheNet resource.\nFor more information refer to [Indexer](https://wiki.servarr.com/sonarr/settings#indexers) and [BroadcastheNet](https://wiki.servarr.com/sonarr/supported#broadcasthenet).",
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

func (r *IndexerBroadcastheNetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var indexer *IndexerBroadcastheNet

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerBroadcastheNet
	request := indexer.read(ctx)

	response, err := r.client.AddIndexerContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", indexerBroadcastheNetResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerBroadcastheNetResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	indexer.write(ctx, response)
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
	response, err := r.client.GetIndexerContext(ctx, indexer.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", indexerBroadcastheNetResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerBroadcastheNetResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	indexer.write(ctx, response)
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
	request := indexer.read(ctx)

	response, err := r.client.UpdateIndexerContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update "+indexerBroadcastheNetResourceName+", got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerBroadcastheNetResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerBroadcastheNetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var indexer *IndexerBroadcastheNet

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerBroadcastheNet current value
	err := r.client.DeleteIndexerContext(ctx, indexer.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", indexerBroadcastheNetResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerBroadcastheNetResourceName+": "+strconv.Itoa(int(indexer.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerBroadcastheNetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+indexerBroadcastheNetResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (i *IndexerBroadcastheNet) write(ctx context.Context, indexer *sonarr.IndexerOutput) {
	genericIndexer := Indexer{
		EnableAutomaticSearch:   types.BoolValue(indexer.EnableAutomaticSearch),
		EnableInteractiveSearch: types.BoolValue(indexer.EnableInteractiveSearch),
		EnableRss:               types.BoolValue(indexer.EnableRss),
		Priority:                types.Int64Value(indexer.Priority),
		DownloadClientID:        types.Int64Value(indexer.DownloadClientID),
		ID:                      types.Int64Value(indexer.ID),
		Name:                    types.StringValue(indexer.Name),
		Tags:                    types.SetNull(types.Int64Type),
	}
	if !i.Tags.IsNull() || !(len(indexer.Tags) == 0) {
		genericIndexer.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, indexer.Tags)
	}

	genericIndexer.writeFields(ctx, indexer.Fields)
	i.fromIndexer(&genericIndexer)
}

func (i *IndexerBroadcastheNet) read(ctx context.Context) *sonarr.IndexerInput {
	var tags []int

	tfsdk.ValueAs(ctx, i.Tags, &tags)

	return &sonarr.IndexerInput{
		EnableAutomaticSearch:   i.EnableAutomaticSearch.ValueBool(),
		EnableInteractiveSearch: i.EnableInteractiveSearch.ValueBool(),
		EnableRss:               i.EnableRss.ValueBool(),
		Priority:                i.Priority.ValueInt64(),
		DownloadClientID:        i.DownloadClientID.ValueInt64(),
		ID:                      i.ID.ValueInt64(),
		ConfigContract:          IndexerBroadcastheNetConfigContrat,
		Implementation:          IndexerBroadcastheNetImplementation,
		Name:                    i.Name.ValueString(),
		Protocol:                IndexerBroadcastheNetProtocol,
		Tags:                    tags,
		Fields:                  i.toIndexer().readFields(ctx),
	}
}
