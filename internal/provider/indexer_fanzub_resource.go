package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	indexerFanzubResourceName   = "indexer_fanzub"
	indexerFanzubImplementation = "Fanzub"
	indexerFanzubConfigContract = "FanzubSettings"
	indexerFanzubProtocol       = "usenet"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerFanzubResource{}
	_ resource.ResourceWithImportState = &IndexerFanzubResource{}
)

func NewIndexerFanzubResource() resource.Resource {
	return &IndexerFanzubResource{}
}

// IndexerFanzubResource defines the Fanzub indexer implementation.
type IndexerFanzubResource struct {
	client *sonarr.APIClient
}

// IndexerFanzub describes the Fanzub indexer data model.
type IndexerFanzub struct {
	Tags                      types.Set    `tfsdk:"tags"`
	BaseURL                   types.String `tfsdk:"base_url"`
	Name                      types.String `tfsdk:"name"`
	ID                        types.Int64  `tfsdk:"id"`
	DownloadClientID          types.Int64  `tfsdk:"download_client_id"`
	Priority                  types.Int64  `tfsdk:"priority"`
	AnimeStandardFormatSearch types.Bool   `tfsdk:"anime_standard_format_search"`
	EnableRss                 types.Bool   `tfsdk:"enable_rss"`
	EnableInteractiveSearch   types.Bool   `tfsdk:"enable_interactive_search"`
	EnableAutomaticSearch     types.Bool   `tfsdk:"enable_automatic_search"`
}

func (i IndexerFanzub) toIndexer() *Indexer {
	return &Indexer{
		AnimeStandardFormatSearch: i.AnimeStandardFormatSearch,
		EnableAutomaticSearch:     i.EnableAutomaticSearch,
		EnableInteractiveSearch:   i.EnableInteractiveSearch,
		EnableRss:                 i.EnableRss,
		Priority:                  i.Priority,
		DownloadClientID:          i.DownloadClientID,
		ID:                        i.ID,
		Name:                      i.Name,
		BaseURL:                   i.BaseURL,
		Tags:                      i.Tags,
	}
}

func (i *IndexerFanzub) fromIndexer(indexer *Indexer) {
	i.AnimeStandardFormatSearch = indexer.AnimeStandardFormatSearch
	i.EnableAutomaticSearch = indexer.EnableAutomaticSearch
	i.EnableInteractiveSearch = indexer.EnableInteractiveSearch
	i.EnableRss = indexer.EnableRss
	i.Priority = indexer.Priority
	i.DownloadClientID = indexer.DownloadClientID
	i.ID = indexer.ID
	i.Name = indexer.Name
	i.BaseURL = indexer.BaseURL
	i.Tags = indexer.Tags
}

func (r *IndexerFanzubResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerFanzubResourceName
}

func (r *IndexerFanzubResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexers -->Indexer Fanzub resource.\nFor more information refer to [Indexer](https://wiki.servarr.com/sonarr/settings#indexers) and [Fanzub](https://wiki.servarr.com/sonarr/supported#fanzub).",
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
				MarkdownDescription: "IndexerFanzub name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "IndexerFanzub ID.",
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
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *IndexerFanzubResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *IndexerFanzubResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var indexer *IndexerFanzub

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerFanzub
	request := indexer.read(ctx)

	response, _, err := r.client.IndexerApi.CreateIndexer(ctx).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, indexerFanzubResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerFanzubResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerFanzubResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var indexer *IndexerFanzub

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IndexerFanzub current value
	response, _, err := r.client.IndexerApi.GetIndexerById(ctx, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerFanzubResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerFanzubResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerFanzubResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var indexer *IndexerFanzub

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update IndexerFanzub
	request := indexer.read(ctx)

	response, _, err := r.client.IndexerApi.UpdateIndexer(ctx, strconv.Itoa(int(request.GetId()))).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, indexerFanzubResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerFanzubResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerFanzubResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var indexer IndexerFanzub

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerFanzub current value
	_, err := r.client.IndexerApi.DeleteIndexer(ctx, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerFanzubResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerFanzubResourceName+": "+strconv.Itoa(int(indexer.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerFanzubResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+indexerFanzubResourceName+": "+req.ID)
}

func (i *IndexerFanzub) write(ctx context.Context, indexer *sonarr.IndexerResource) {
	genericIndexer := Indexer{
		EnableAutomaticSearch:   types.BoolValue(indexer.GetEnableAutomaticSearch()),
		EnableInteractiveSearch: types.BoolValue(indexer.GetEnableInteractiveSearch()),
		EnableRss:               types.BoolValue(indexer.GetEnableRss()),
		Priority:                types.Int64Value(int64(indexer.GetPriority())),
		DownloadClientID:        types.Int64Value(int64(indexer.GetDownloadClientId())),
		ID:                      types.Int64Value(int64(indexer.GetId())),
		Name:                    types.StringValue(indexer.GetName()),
		AnimeCategories:         types.SetValueMust(types.Int64Type, nil),
		Categories:              types.SetValueMust(types.Int64Type, nil),
	}
	genericIndexer.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, indexer.Tags)
	genericIndexer.writeFields(ctx, indexer.GetFields())
	i.fromIndexer(&genericIndexer)
}

func (i *IndexerFanzub) read(ctx context.Context) *sonarr.IndexerResource {
	var tags []*int32

	tfsdk.ValueAs(ctx, i.Tags, &tags)

	indexer := sonarr.NewIndexerResource()
	indexer.SetEnableAutomaticSearch(i.EnableAutomaticSearch.ValueBool())
	indexer.SetEnableInteractiveSearch(i.EnableInteractiveSearch.ValueBool())
	indexer.SetEnableRss(i.EnableRss.ValueBool())
	indexer.SetPriority(int32(i.Priority.ValueInt64()))
	indexer.SetDownloadClientId(int32(i.DownloadClientID.ValueInt64()))
	indexer.SetId(int32(i.ID.ValueInt64()))
	indexer.SetConfigContract(indexerFanzubConfigContract)
	indexer.SetImplementation(indexerFanzubImplementation)
	indexer.SetName(i.Name.ValueString())
	indexer.SetProtocol(indexerFanzubProtocol)
	indexer.SetTags(tags)
	indexer.SetFields(i.toIndexer().readFields(ctx))

	return indexer
}
