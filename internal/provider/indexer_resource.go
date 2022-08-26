package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.ResourceType = resourceIndexerType{}
var _ resource.Resource = resourceIndexer{}
var _ resource.ResourceWithImportState = resourceIndexer{}

type resourceIndexerType struct{}

type resourceIndexer struct {
	provider sonarrProvider
}

// Indexer is the indexer resource.
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
	ApiKey                    types.String  `tfsdk:"api_key"`
	ApiPath                   types.String  `tfsdk:"api_path"`
	BaseUrl                   types.String  `tfsdk:"base_url"`
	CaptchaToken              types.String  `tfsdk:"captcha_token"`
	Cookie                    types.String  `tfsdk:"cookie"`
	Passkey                   types.String  `tfsdk:"passkey"`
	Username                  types.String  `tfsdk:"username"`
	AnimeCategories           types.Set     `tfsdk:"anime_categories"`
	Categories                types.Set     `tfsdk:"categories"`
}

func (t resourceIndexerType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Indexer resource",
		Attributes: map[string]tfsdk.Attribute{
			"enable_automatic_search": {
				MarkdownDescription: "Enable automatic search flag",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"enable_interactive_search": {
				MarkdownDescription: "Enable interactive search flag",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"enable_rss": {
				MarkdownDescription: "Enable RSS flag",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"priority": {
				MarkdownDescription: "Priority",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"download_client_id": {
				MarkdownDescription: "Download client ID",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"config_contract": {
				MarkdownDescription: "Indexer configuration template",
				Required:            true,
				Type:                types.StringType,
			},
			"implementation": {
				MarkdownDescription: "Indexer implementation name",
				Required:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "Name",
				Required:            true,
				Type:                types.StringType,
			},
			"protocol": {
				MarkdownDescription: "Protocol. Valid values are 'usenet' and 'torrent'",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					helpers.StringMatch([]string{"usenet", "torrent"}),
				},
			},
			"tags": {
				MarkdownDescription: "List of associated tags",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"id": {
				MarkdownDescription: "Indexer ID",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			// Field values
			"allow_zero_size": {
				MarkdownDescription: "Allow zero size files",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"anime_standard_format_search": {
				MarkdownDescription: "Search anime in standard format",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"ranked_only": {
				MarkdownDescription: "Allow ranked only",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"delay": {
				MarkdownDescription: "Delay before grabbing",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"minimum_seeders": {
				MarkdownDescription: "Minimum seeders",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"season_pack_seed_time": {
				MarkdownDescription: "Season seed time",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"seed_time": {
				MarkdownDescription: "Seed time",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"seed_ratio": {
				MarkdownDescription: "Seed ratio",
				Optional:            true,
				Computed:            true,
				Type:                types.Float64Type,
			},
			"additional_parameters": {
				MarkdownDescription: "Additional parameters",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"api_key": {
				MarkdownDescription: "API key",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"api_path": {
				MarkdownDescription: "API path",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"base_url": {
				MarkdownDescription: "Base URL",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"captcha_token": {
				MarkdownDescription: "Captcha token",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"cookie": {
				MarkdownDescription: "Cookie",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"passkey": {
				MarkdownDescription: "Passkey",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"username": {
				MarkdownDescription: "Username",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"categories": {
				MarkdownDescription: "Series list",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"anime_categories": {
				MarkdownDescription: "Anime list",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
		},
	}, nil
}

func (t resourceIndexerType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return resourceIndexer{
		provider: provider,
	}, diags
}

func (r resourceIndexer) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan Indexer
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Indexer
	request := readIndexer(ctx, &plan)
	response, err := r.provider.client.AddIndexerContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Indexer, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "created Indexer: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := writeIndexer(ctx, response)

	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceIndexer) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state Indexer
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Indexer current value
	response, err := r.provider.client.GetIndexerContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Indexers, got error: %s", err))
		return
	}
	// Map response body to resource schema attribute
	result := writeIndexer(ctx, response)

	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceIndexer) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan Indexer
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update Indexer
	request := readIndexer(ctx, &plan)

	response, err := r.provider.client.UpdateIndexerContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update Indexer, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "update Indexer: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := writeIndexer(ctx, response)

	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceIndexer) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Indexer

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete Indexer current value
	err := r.provider.client.DeleteIndexerContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Indexers, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceIndexer) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	//resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)
		return
	}
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
				output.Delay = types.Int64{Value: f.Value.(int64)}
			case "minimumSeeders":
				output.MinimumSeeders = types.Int64{Value: f.Value.(int64)}
			case "seasonPackSeedTime":
				output.SeasonPackSeedTime = types.Int64{Value: f.Value.(int64)}
			case "seedTime":
				output.SeedTime = types.Int64{Value: f.Value.(int64)}
			case "seedRatio":
				output.SeedRatio = types.Float64{Value: f.Value.(float64)}
			case "additionalParameters":
				output.AdditionalParameters = types.String{Value: f.Value.(string)}
			case "apiKey":
				output.ApiKey = types.String{Value: f.Value.(string)}
			case "apiPath":
				output.ApiPath = types.String{Value: f.Value.(string)}
			case "baseUrl":
				output.BaseUrl = types.String{Value: f.Value.(string)}
			case "captchaToken":
				output.CaptchaToken = types.String{Value: f.Value.(string)}
			case "cookie":
				output.Cookie = types.String{Value: f.Value.(string)}
			case "passkey":
				output.Passkey = types.String{Value: f.Value.(string)}
			case "username":
				output.Username = types.String{Value: f.Value.(string)}
			case "animeCategories":
				output.AnimeCategories = types.Set{ElemType: types.Int64Type}
				tfsdk.ValueFrom(ctx, f.Value, output.AnimeCategories.Type(ctx), &output.AnimeCategories)
			case "categories":
				output.Categories = types.Set{ElemType: types.Int64Type}
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
	if !indexer.ApiKey.IsNull() && !indexer.ApiKey.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "apiKey",
			Value: indexer.ApiKey.Value,
		})
	}
	if !indexer.ApiPath.IsNull() && !indexer.ApiPath.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "apiPath",
			Value: indexer.ApiPath.Value,
		})
	}
	if !indexer.BaseUrl.IsNull() && !indexer.BaseUrl.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "baseUrl",
			Value: indexer.BaseUrl.Value,
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
