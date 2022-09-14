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
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IndexerResource{}
var _ resource.ResourceWithImportState = &IndexerResource{}

func NewIndexerResource() resource.Resource {
	return &IndexerResource{}
}

// IndexerResource defines the indexer implementation.
type IndexerResource struct {
	client *sonarr.Sonarr
}

// Indexer describes the indexer data model.
type Indexer struct {
	EnableAutomaticSearch   types.Bool   `tfsdk:"enable_automatic_search"`
	EnableInteractiveSearch types.Bool   `tfsdk:"enable_interactive_search"`
	EnableRss               types.Bool   `tfsdk:"enable_rss"`
	Priority                types.Int64  `tfsdk:"priority"`
	DownloadClientID        types.Int64  `tfsdk:"download_client_id"`
	ID                      types.Int64  `tfsdk:"id"`
	ConfigContract          types.String `tfsdk:"config_contract"`
	Implementation          types.String `tfsdk:"implementation"`
	Name                    types.String `tfsdk:"name"`
	Protocol                types.String `tfsdk:"protocol"`
	Tags                    types.Set    `tfsdk:"tags"`
	// Fields values
	AllowZeroSize             types.Bool    `tfsdk:"allow_zero_size"`
	AnimeStandardFormatSearch types.Bool    `tfsdk:"anime_standard_format_search"`
	RankedOnly                types.Bool    `tfsdk:"ranked_only"`
	Delay                     types.Int64   `tfsdk:"delay"`
	MinimumSeeders            types.Int64   `tfsdk:"minimum_seeders"`
	SeasonPackSeedTime        types.Int64   `tfsdk:"season_pack_seed_time"`
	SeedTime                  types.Int64   `tfsdk:"seed_time"`
	SeedRatio                 types.Float64 `tfsdk:"seed_ratio"`
	AdditionalParameters      types.String  `tfsdk:"additional_parameters"`
	APIKey                    types.String  `tfsdk:"api_key"`
	APIPath                   types.String  `tfsdk:"api_path"`
	BaseURL                   types.String  `tfsdk:"base_url"`
	CaptchaToken              types.String  `tfsdk:"captcha_token"`
	Cookie                    types.String  `tfsdk:"cookie"`
	Passkey                   types.String  `tfsdk:"passkey"`
	Username                  types.String  `tfsdk:"username"`
	AnimeCategories           types.Set     `tfsdk:"anime_categories"`
	Categories                types.Set     `tfsdk:"categories"`
}

func (r *IndexerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_indexer"
}

func (r *IndexerResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"config_contract": {
				MarkdownDescription: "Indexer configuration template.",
				Required:            true,
				Type:                types.StringType,
			},
			"implementation": {
				MarkdownDescription: "Indexer implementation name.",
				Required:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "Indexer name.",
				Required:            true,
				Type:                types.StringType,
			},
			"protocol": {
				MarkdownDescription: "Protocol. Valid values are 'usenet' and 'torrent'.",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					helpers.StringMatch([]string{"usenet", "torrent"}),
				},
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
			"allow_zero_size": {
				MarkdownDescription: "Allow zero size files.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"anime_standard_format_search": {
				MarkdownDescription: "Search anime in standard format.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"ranked_only": {
				MarkdownDescription: "Allow ranked only.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"delay": {
				MarkdownDescription: "Delay before grabbing.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
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
			"api_key": {
				MarkdownDescription: "API key.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"api_path": {
				MarkdownDescription: "API path.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
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
			"cookie": {
				MarkdownDescription: "Cookie.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"passkey": {
				MarkdownDescription: "Passkey.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"username": {
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"categories": {
				MarkdownDescription: "Series list.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"anime_categories": {
				MarkdownDescription: "Anime list.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
		},
	}, nil
}

func (r *IndexerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IndexerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan Indexer

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Indexer
	request := readIndexer(ctx, &plan)

	response, err := r.client.AddIndexerContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to create Indexer, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "created indexer: "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeIndexer(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *IndexerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state Indexer

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get Indexer current value
	response, err := r.client.GetIndexerContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read Indexers, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "read indexer: "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	result := writeIndexer(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *IndexerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan Indexer

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Indexer
	request := readIndexer(ctx, &plan)

	response, err := r.client.UpdateIndexerContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to update Indexer, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated indexer: "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeIndexer(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *IndexerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Indexer

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete Indexer current value
	err := r.client.DeleteIndexerContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read Indexers, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "deleted indexer: "+strconv.Itoa(int(state.ID.Value)))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func writeIndexer(ctx context.Context, indexer *sonarr.IndexerOutput) *Indexer {
	output := Indexer{
		EnableAutomaticSearch:   types.Bool{Value: indexer.EnableAutomaticSearch},
		EnableInteractiveSearch: types.Bool{Value: indexer.EnableInteractiveSearch},
		EnableRss:               types.Bool{Value: indexer.EnableRss},
		Priority:                types.Int64{Value: indexer.Priority},
		DownloadClientID:        types.Int64{Value: indexer.DownloadClientID},
		ID:                      types.Int64{Value: indexer.ID},
		ConfigContract:          types.String{Value: indexer.ConfigContract},
		Implementation:          types.String{Value: indexer.Implementation},
		Name:                    types.String{Value: indexer.Name},
		Protocol:                types.String{Value: indexer.Protocol},
		Tags:                    types.Set{ElemType: types.Int64Type},
		AnimeCategories:         types.Set{ElemType: types.Int64Type},
		Categories:              types.Set{ElemType: types.Int64Type},
	}
	tfsdk.ValueFrom(ctx, indexer.Tags, output.Tags.Type(ctx), &output.Tags)

	for _, f := range indexer.Fields {
		if f.Value != nil {
			switch f.Name {
			case "allowZeroSize":
				output.AllowZeroSize = types.Bool{Value: f.Value.(bool)}
			case "animeStandardFormatSearch":
				output.AnimeStandardFormatSearch = types.Bool{Value: f.Value.(bool)}
			case "rankedOnly":
				output.RankedOnly = types.Bool{Value: f.Value.(bool)}
			case "delay":
				output.Delay = types.Int64{Value: int64(f.Value.(float64))}
			case "minimumSeeders":
				output.MinimumSeeders = types.Int64{Value: int64(f.Value.(float64))}
			case "seasonPackSeedTime":
				output.SeasonPackSeedTime = types.Int64{Value: int64(f.Value.(float64))}
			case "seedTime":
				output.SeedTime = types.Int64{Value: int64(f.Value.(float64))}
			case "seedRatio":
				output.SeedRatio = types.Float64{Value: f.Value.(float64)}
			case "additionalParameters":
				output.AdditionalParameters = types.String{Value: f.Value.(string)}
			case "apiKey":
				output.APIKey = types.String{Value: f.Value.(string)}
			case "apiPath":
				output.APIPath = types.String{Value: f.Value.(string)}
			case "baseUrl":
				output.BaseURL = types.String{Value: f.Value.(string)}
			case "captchaToken":
				output.CaptchaToken = types.String{Value: f.Value.(string)}
			case "cookie":
				output.Cookie = types.String{Value: f.Value.(string)}
			case "passkey":
				output.Passkey = types.String{Value: f.Value.(string)}
			case "username":
				output.Username = types.String{Value: f.Value.(string)}
			case "animeCategories":
				tfsdk.ValueFrom(ctx, f.Value, output.AnimeCategories.Type(ctx), &output.AnimeCategories)
			case "categories":
				tfsdk.ValueFrom(ctx, f.Value, output.Categories.Type(ctx), &output.Categories)
			// TODO: manage unknown values
			default:
			}
		}
	}

	return &output
}

func readIndexer(ctx context.Context, indexer *Indexer) *sonarr.IndexerInput {
	var tags []int

	tfsdk.ValueAs(ctx, indexer.Tags, &tags)

	return &sonarr.IndexerInput{
		EnableAutomaticSearch:   indexer.EnableAutomaticSearch.Value,
		EnableInteractiveSearch: indexer.EnableInteractiveSearch.Value,
		EnableRss:               indexer.EnableRss.Value,
		Priority:                indexer.Priority.Value,
		DownloadClientID:        indexer.DownloadClientID.Value,
		ID:                      indexer.ID.Value,
		ConfigContract:          indexer.ConfigContract.Value,
		Implementation:          indexer.Implementation.Value,
		Name:                    indexer.Name.Value,
		Protocol:                indexer.Protocol.Value,
		Tags:                    tags,
		Fields:                  readIndexerFields(ctx, indexer),
	}
}

func readIndexerFields(ctx context.Context, indexer *Indexer) []*starr.FieldInput {
	var output []*starr.FieldInput
	if !indexer.AllowZeroSize.IsNull() && !indexer.AllowZeroSize.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "allowZeroSize",
			Value: indexer.AllowZeroSize.Value,
		})
	}

	if !indexer.AnimeStandardFormatSearch.IsNull() && !indexer.AnimeStandardFormatSearch.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "animeStandardFormatSearch",
			Value: indexer.AnimeStandardFormatSearch.Value,
		})
	}

	if !indexer.RankedOnly.IsNull() && !indexer.RankedOnly.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "rankedOnly",
			Value: indexer.RankedOnly.Value,
		})
	}

	if !indexer.Delay.IsNull() && !indexer.Delay.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "delay",
			Value: indexer.Delay.Value,
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

	if !indexer.AdditionalParameters.IsNull() && !indexer.AdditionalParameters.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "additionalParameters",
			Value: indexer.AdditionalParameters.Value,
		})
	}

	if !indexer.APIKey.IsNull() && !indexer.APIKey.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "apiKey",
			Value: indexer.APIKey.Value,
		})
	}

	if !indexer.APIPath.IsNull() && !indexer.APIPath.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "apiPath",
			Value: indexer.APIPath.Value,
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

	if !indexer.Cookie.IsNull() && !indexer.Cookie.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "cookie",
			Value: indexer.Cookie.Value,
		})
	}

	if !indexer.Passkey.IsNull() && !indexer.Passkey.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "passkey",
			Value: indexer.Passkey.Value,
		})
	}

	if !indexer.Username.IsNull() && !indexer.Username.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "username",
			Value: indexer.Username.Value,
		})
	}

	if len(indexer.Categories.Elems) != 0 {
		cat := make([]int64, len(indexer.Categories.Elems))
		tfsdk.ValueAs(ctx, indexer.Categories, &cat)

		output = append(output, &starr.FieldInput{
			Name:  "categories",
			Value: cat,
		})
	}

	if len(indexer.AnimeCategories.Elems) != 0 {
		cat := make([]int64, len(indexer.AnimeCategories.Elems))
		tfsdk.ValueAs(ctx, indexer.AnimeCategories, &cat)

		output = append(output, &starr.FieldInput{
			Name:  "animeCategories",
			Value: cat,
		})
	}

	return output
}
