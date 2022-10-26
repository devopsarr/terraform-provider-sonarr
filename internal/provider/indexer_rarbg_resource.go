package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const (
	indexerRarbgResourceName   = "indexer_rarbg"
	IndexerRarbgImplementation = "Rarbg"
	IndexerRarbgConfigContrat  = "RarbgSettings"
	IndexerRarbgProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IndexerRarbgResource{}
var _ resource.ResourceWithImportState = &IndexerRarbgResource{}

func NewIndexerRarbgResource() resource.Resource {
	return &IndexerRarbgResource{}
}

// IndexerRarbgResource defines the Rarbg indexer implementation.
type IndexerRarbgResource struct {
	client *sonarr.Sonarr
}

// IndexerRarbg describes the Rarbg indexer data model.
type IndexerRarbg struct {
	Tags                    types.Set     `tfsdk:"tags"`
	Name                    types.String  `tfsdk:"name"`
	CaptchaToken            types.String  `tfsdk:"captcha_token"`
	BaseURL                 types.String  `tfsdk:"base_url"`
	Priority                types.Int64   `tfsdk:"priority"`
	ID                      types.Int64   `tfsdk:"id"`
	DownloadClientID        types.Int64   `tfsdk:"download_client_id"`
	MinimumSeeders          types.Int64   `tfsdk:"minimum_seeders"`
	SeasonPackSeedTime      types.Int64   `tfsdk:"season_pack_seed_time"`
	SeedTime                types.Int64   `tfsdk:"seed_time"`
	SeedRatio               types.Float64 `tfsdk:"seed_ratio"`
	EnableAutomaticSearch   types.Bool    `tfsdk:"enable_automatic_search"`
	RankedOnly              types.Bool    `tfsdk:"ranked_only"`
	EnableRss               types.Bool    `tfsdk:"enable_rss"`
	EnableInteractiveSearch types.Bool    `tfsdk:"enable_interactive_search"`
}

func (i IndexerRarbg) toIndexer() *Indexer {
	return &Indexer{
		EnableAutomaticSearch:   i.EnableAutomaticSearch,
		EnableInteractiveSearch: i.EnableInteractiveSearch,
		EnableRss:               i.EnableRss,
		Priority:                i.Priority,
		DownloadClientID:        i.DownloadClientID,
		ID:                      i.ID,
		Name:                    i.Name,
		RankedOnly:              i.RankedOnly,
		MinimumSeeders:          i.MinimumSeeders,
		SeasonPackSeedTime:      i.SeasonPackSeedTime,
		SeedTime:                i.SeedTime,
		SeedRatio:               i.SeedRatio,
		CaptchaToken:            i.CaptchaToken,
		BaseURL:                 i.BaseURL,
		Tags:                    i.Tags,
	}
}

func (i *IndexerRarbg) fromIndexer(indexer *Indexer) {
	i.EnableAutomaticSearch = indexer.EnableAutomaticSearch
	i.EnableInteractiveSearch = indexer.EnableInteractiveSearch
	i.EnableRss = indexer.EnableRss
	i.Priority = indexer.Priority
	i.DownloadClientID = indexer.DownloadClientID
	i.ID = indexer.ID
	i.Name = indexer.Name
	i.RankedOnly = indexer.RankedOnly
	i.MinimumSeeders = indexer.MinimumSeeders
	i.SeasonPackSeedTime = indexer.SeasonPackSeedTime
	i.SeedTime = indexer.SeedTime
	i.SeedRatio = indexer.SeedRatio
	i.CaptchaToken = indexer.CaptchaToken
	i.BaseURL = indexer.BaseURL
	i.Tags = indexer.Tags
}

func (r *IndexerRarbgResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerRarbgResourceName
}

func (r *IndexerRarbgResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "[subcategory:Indexers]: #\nIndexerRarbg resource.\nFor more information refer to [Indexer](https://wiki.servarr.com/sonarr/settings#indexers) documentation.",
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
				MarkdownDescription: "IndexerRarbg name.",
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
				MarkdownDescription: "IndexerRarbg ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			// Field values
			"ranked_only": {
				MarkdownDescription: "Allow ranked only.",
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
			"base_url": {
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"captcha_token": {
				MarkdownDescription: "Captcha token.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (r *IndexerRarbgResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			helpers.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *IndexerRarbgResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan IndexerRarbg

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerRarbg
	request := readIndexerRarbg(ctx, &plan)

	response, err := r.client.AddIndexerContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", indexerRarbgResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerRarbgResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeIndexerRarbg(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *IndexerRarbgResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state IndexerRarbg

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IndexerRarbg current value
	response, err := r.client.GetIndexerContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", indexerRarbgResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerRarbgResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	result := writeIndexerRarbg(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *IndexerRarbgResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan IndexerRarbg

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update IndexerRarbg
	request := readIndexerRarbg(ctx, &plan)

	response, err := r.client.UpdateIndexerContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to update "+indexerRarbgResourceName+", got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerRarbgResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeIndexerRarbg(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *IndexerRarbgResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state IndexerRarbg

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerRarbg current value
	err := r.client.DeleteIndexerContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", indexerRarbgResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerRarbgResourceName+": "+strconv.Itoa(int(state.ID.Value)))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerRarbgResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			helpers.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+indexerRarbgResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func writeIndexerRarbg(ctx context.Context, indexer *sonarr.IndexerOutput) *IndexerRarbg {
	var output IndexerRarbg

	genericIndexer := Indexer{
		EnableAutomaticSearch:   types.Bool{Value: indexer.EnableAutomaticSearch},
		EnableInteractiveSearch: types.Bool{Value: indexer.EnableInteractiveSearch},
		EnableRss:               types.Bool{Value: indexer.EnableRss},
		Priority:                types.Int64{Value: indexer.Priority},
		DownloadClientID:        types.Int64{Value: indexer.DownloadClientID},
		ID:                      types.Int64{Value: indexer.ID},
		Name:                    types.String{Value: indexer.Name},
		Tags:                    types.Set{ElemType: types.Int64Type},
	}
	tfsdk.ValueFrom(ctx, indexer.Tags, genericIndexer.Tags.Type(ctx), &genericIndexer.Tags)
	genericIndexer.writeIndexerFields(ctx, indexer.Fields)
	output.fromIndexer(&genericIndexer)

	return &output
}

func readIndexerRarbg(ctx context.Context, indexer *IndexerRarbg) *sonarr.IndexerInput {
	var tags []int

	tfsdk.ValueAs(ctx, indexer.Tags, &tags)

	return &sonarr.IndexerInput{
		EnableAutomaticSearch:   indexer.EnableAutomaticSearch.Value,
		EnableInteractiveSearch: indexer.EnableInteractiveSearch.Value,
		EnableRss:               indexer.EnableRss.Value,
		Priority:                indexer.Priority.Value,
		DownloadClientID:        indexer.DownloadClientID.Value,
		ID:                      indexer.ID.Value,
		ConfigContract:          IndexerRarbgConfigContrat,
		Implementation:          IndexerRarbgImplementation,
		Name:                    indexer.Name.Value,
		Protocol:                IndexerRarbgProtocol,
		Tags:                    tags,
		Fields:                  readIndexerFields(ctx, indexer.toIndexer()),
	}
}
