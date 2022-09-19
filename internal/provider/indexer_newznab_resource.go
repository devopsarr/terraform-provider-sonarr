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
var _ resource.Resource = &IndexerNewznabResource{}
var _ resource.ResourceWithImportState = &IndexerNewznabResource{}

func NewIndexerNewznabResource() resource.Resource {
	return &IndexerNewznabResource{}
}

// IndexerNewznabResource defines the Newznab indexer implementation.
type IndexerNewznabResource struct {
	client *sonarr.Sonarr
}

const (
	IndexerNewznabImplementation = "Newznab"
	IndexerNewznabConfigContrat  = "NewznabSettings"
	IndexerNewznabProtocol       = "usenet"
)

// IndexerNewznab describes the Newznab indexer data model.
type IndexerNewznab struct {
	AnimeStandardFormatSearch types.Bool   `tfsdk:"anime_standard_format_search"`
	EnableAutomaticSearch     types.Bool   `tfsdk:"enable_automatic_search"`
	EnableInteractiveSearch   types.Bool   `tfsdk:"enable_interactive_search"`
	EnableRss                 types.Bool   `tfsdk:"enable_rss"`
	Priority                  types.Int64  `tfsdk:"priority"`
	DownloadClientID          types.Int64  `tfsdk:"download_client_id"`
	ID                        types.Int64  `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	AdditionalParameters      types.String `tfsdk:"additional_parameters"`
	APIKey                    types.String `tfsdk:"api_key"`
	APIPath                   types.String `tfsdk:"api_path"`
	BaseURL                   types.String `tfsdk:"base_url"`
	AnimeCategories           types.Set    `tfsdk:"anime_categories"`
	Categories                types.Set    `tfsdk:"categories"`
	Tags                      types.Set    `tfsdk:"tags"`
}

func (r *IndexerNewznabResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_indexer_newznab"
}

func (r *IndexerNewznabResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "IndexerNewznab resource.<br/>For more information refer to [IndexerNewznab](https://wiki.servarr.com/sonarr/settings#indexers) documentation.",
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
				MarkdownDescription: "IndexerNewznab name.",
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
				MarkdownDescription: "IndexerNewznab ID.",
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

func (r *IndexerNewznabResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IndexerNewznabResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan IndexerNewznab

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerNewznab
	request := readIndexerNewznab(ctx, &plan)

	response, err := r.client.AddIndexerContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to create IndexerNewznab, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "created indexer: "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeIndexerNewznab(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *IndexerNewznabResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state IndexerNewznab

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IndexerNewznab current value
	response, err := r.client.GetIndexerContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read IndexerNewznabs, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "read indexer: "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	result := writeIndexerNewznab(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *IndexerNewznabResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan IndexerNewznab

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update IndexerNewznab
	request := readIndexerNewznab(ctx, &plan)

	response, err := r.client.UpdateIndexerContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to update IndexerNewznab, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated indexer: "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeIndexerNewznab(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *IndexerNewznabResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state IndexerNewznab

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerNewznab current value
	err := r.client.DeleteIndexerContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read IndexerNewznabs, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "deleted indexer: "+strconv.Itoa(int(state.ID.Value)))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerNewznabResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func writeIndexerNewznab(ctx context.Context, indexer *sonarr.IndexerOutput) *IndexerNewznab {
	output := IndexerNewznab{
		EnableAutomaticSearch:   types.Bool{Value: indexer.EnableAutomaticSearch},
		EnableInteractiveSearch: types.Bool{Value: indexer.EnableInteractiveSearch},
		EnableRss:               types.Bool{Value: indexer.EnableRss},
		Priority:                types.Int64{Value: indexer.Priority},
		DownloadClientID:        types.Int64{Value: indexer.DownloadClientID},
		ID:                      types.Int64{Value: indexer.ID},
		Name:                    types.String{Value: indexer.Name},
		Tags:                    types.Set{ElemType: types.Int64Type},
		AnimeCategories:         types.Set{ElemType: types.Int64Type},
		Categories:              types.Set{ElemType: types.Int64Type},
	}
	tfsdk.ValueFrom(ctx, indexer.Tags, output.Tags.Type(ctx), &output.Tags)

	for _, f := range indexer.Fields {
		if f.Value != nil {
			switch f.Name {
			case "animeStandardFormatSearch":
				output.AnimeStandardFormatSearch = types.Bool{Value: f.Value.(bool)}
			case "additionalParameters":
				output.AdditionalParameters = types.String{Value: f.Value.(string)}
			case "apiKey":
				output.APIKey = types.String{Value: f.Value.(string)}
			case "apiPath":
				output.APIPath = types.String{Value: f.Value.(string)}
			case "baseUrl":
				output.BaseURL = types.String{Value: f.Value.(string)}
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

func readIndexerNewznab(ctx context.Context, indexer *IndexerNewznab) *sonarr.IndexerInput {
	var tags []int

	tfsdk.ValueAs(ctx, indexer.Tags, &tags)

	return &sonarr.IndexerInput{
		EnableAutomaticSearch:   indexer.EnableAutomaticSearch.Value,
		EnableInteractiveSearch: indexer.EnableInteractiveSearch.Value,
		EnableRss:               indexer.EnableRss.Value,
		Priority:                indexer.Priority.Value,
		DownloadClientID:        indexer.DownloadClientID.Value,
		ID:                      indexer.ID.Value,
		ConfigContract:          IndexerNewznabConfigContrat,
		Implementation:          IndexerNewznabImplementation,
		Name:                    indexer.Name.Value,
		Protocol:                IndexerNewznabProtocol,
		Tags:                    tags,
		Fields:                  readIndexerNewznabFields(ctx, indexer),
	}
}

func readIndexerNewznabFields(ctx context.Context, indexer *IndexerNewznab) []*starr.FieldInput {
	var output []*starr.FieldInput

	if !indexer.AnimeStandardFormatSearch.IsNull() && !indexer.AnimeStandardFormatSearch.IsUnknown() {
		output = append(output, &starr.FieldInput{
			Name:  "animeStandardFormatSearch",
			Value: indexer.AnimeStandardFormatSearch.Value,
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
