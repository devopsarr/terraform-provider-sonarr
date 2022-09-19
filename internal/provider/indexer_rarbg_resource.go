package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr"
	"golift.io/starr/sonarr"
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

const (
	IndexerRarbgImplementation = "Rarbg"
	IndexerRarbgConfigContrat  = "RarbgSettings"
	IndexerRarbgProtocol       = "torrent"
)

// IndexerRarbg describes the Rarbg indexer data model.
type IndexerRarbg struct {
	EnableAutomaticSearch   types.Bool   `tfsdk:"enable_automatic_search"`
	EnableInteractiveSearch types.Bool   `tfsdk:"enable_interactive_search"`
	EnableRss               types.Bool   `tfsdk:"enable_rss"`
	Priority                types.Int64  `tfsdk:"priority"`
	DownloadClientID        types.Int64  `tfsdk:"download_client_id"`
	ID                      types.Int64  `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Tags                    types.Set    `tfsdk:"tags"`
	// Fields values
	RankedOnly         types.Bool    `tfsdk:"ranked_only"`
	MinimumSeeders     types.Int64   `tfsdk:"minimum_seeders"`
	SeasonPackSeedTime types.Int64   `tfsdk:"season_pack_seed_time"`
	SeedTime           types.Int64   `tfsdk:"seed_time"`
	SeedRatio          types.Float64 `tfsdk:"seed_ratio"`
	BaseURL            types.String  `tfsdk:"base_url"`
	CaptchaToken       types.String  `tfsdk:"captcha_token"`
}

func (r *IndexerRarbgResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_indexer_rarbg"
}

func (r *IndexerRarbgResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Indexer resource.<br/>For more information refer to [Indexer](https://wiki.servarr.com/sonarr/settings#indexers) documentation.",
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
				MarkdownDescription: "Indexer name.",
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
				MarkdownDescription: "Indexer ID.",
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
			UnexpectedResourceConfigureType,
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
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to create IndexerRarbg, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "created indexer: "+strconv.Itoa(int(response.ID)))
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
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read IndexerRarbgs, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "read indexer: "+strconv.Itoa(int(response.ID)))
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
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to update IndexerRarbg, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated indexer: "+strconv.Itoa(int(response.ID)))
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
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read IndexerRarbgs, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "deleted indexer: "+strconv.Itoa(int(state.ID.Value)))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerRarbgResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported indexer: "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func writeIndexerRarbg(ctx context.Context, indexer *sonarr.IndexerOutput) *IndexerRarbg {
	output := IndexerRarbg{
		EnableAutomaticSearch:   types.Bool{Value: indexer.EnableAutomaticSearch},
		EnableInteractiveSearch: types.Bool{Value: indexer.EnableInteractiveSearch},
		EnableRss:               types.Bool{Value: indexer.EnableRss},
		Priority:                types.Int64{Value: indexer.Priority},
		DownloadClientID:        types.Int64{Value: indexer.DownloadClientID},
		ID:                      types.Int64{Value: indexer.ID},
		Name:                    types.String{Value: indexer.Name},
		Tags:                    types.Set{ElemType: types.Int64Type},
	}
	tfsdk.ValueFrom(ctx, indexer.Tags, output.Tags.Type(ctx), &output.Tags)

	for _, f := range indexer.Fields {
		if f.Value != nil {
			switch f.Name {
			case "rankedOnly":
				output.RankedOnly = types.Bool{Value: f.Value.(bool)}
			case "minimumSeeders":
				output.MinimumSeeders = types.Int64{Value: int64(f.Value.(float64))}
			case "seasonPackSeedTime":
				output.SeasonPackSeedTime = types.Int64{Value: int64(f.Value.(float64))}
			case "seedTime":
				output.SeedTime = types.Int64{Value: int64(f.Value.(float64))}
			case "seedRatio":
				output.SeedRatio = types.Float64{Value: f.Value.(float64)}
			case "baseUrl":
				output.BaseURL = types.String{Value: f.Value.(string)}
			case "captchaToken":
				output.CaptchaToken = types.String{Value: f.Value.(string)}
			// TODO: manage unknown values
			default:
			}
		}
	}

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
		Fields:                  readIndexerRarbgFields(indexer),
	}
}

func readIndexerRarbgFields(indexer *IndexerRarbg) []*starr.FieldInput {
	var output []*starr.FieldInput

	if !indexer.RankedOnly.IsNull() && !indexer.RankedOnly.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "rankedOnly",
			Value: indexer.RankedOnly.Value,
		})
	}

	if !indexer.MinimumSeeders.IsNull() && !indexer.MinimumSeeders.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "minimumSeeders",
			Value: indexer.MinimumSeeders.Value,
		})
	}

	if !indexer.SeasonPackSeedTime.IsNull() && !indexer.SeasonPackSeedTime.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "seasonPackSeedTime",
			Value: indexer.SeasonPackSeedTime.Value,
		})
	}

	if !indexer.SeedTime.IsNull() && !indexer.SeedTime.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "seedTime",
			Value: indexer.SeedTime.Value,
		})
	}

	if !indexer.SeedRatio.IsNull() && !indexer.SeedRatio.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "seedRatio",
			Value: indexer.SeedRatio.Value,
		})
	}

	if !indexer.BaseURL.IsNull() && !indexer.BaseURL.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "baseUrl",
			Value: indexer.BaseURL.Value,
		})
	}

	if !indexer.CaptchaToken.IsNull() && !indexer.CaptchaToken.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "captchaToken",
			Value: indexer.CaptchaToken.Value,
		})
	}

	return output
}
