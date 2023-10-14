package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const indexerResourceName = "indexer"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerResource{}
	_ resource.ResourceWithImportState = &IndexerResource{}
)

var indexerFields = helpers.Fields{
	IntSlices:        []string{"categories", "animeCategories"},
	Bools:            []string{"allowZeroSize", "animeStandardFormatSearch", "rankedOnly"},
	Ints:             []string{"delay", "minimumSeeders", "seasonPackSeedTime", "seedTime"},
	IntsExceptions:   []string{"seedCriteria.seedTime", "seedCriteria.seasonPackSeedTime"},
	Strings:          []string{"additionalParameters", "apiKey", "apiPath", "baseUrl", "captchaToken", "cookie", "passkey", "username"},
	Floats:           []string{"seedRatio"},
	FloatsExceptions: []string{"seedCriteria.seedRatio"},
}

func NewIndexerResource() resource.Resource {
	return &IndexerResource{}
}

// IndexerResource defines the indexer implementation.
type IndexerResource struct {
	client *sonarr.APIClient
}

// Indexer describes the indexer data model.
type Indexer struct {
	SeedRatio                 types.Float64 `tfsdk:"seed_ratio"`
	Tags                      types.Set     `tfsdk:"tags"`
	Categories                types.Set     `tfsdk:"categories"`
	AnimeCategories           types.Set     `tfsdk:"anime_categories"`
	APIKey                    types.String  `tfsdk:"api_key"`
	Username                  types.String  `tfsdk:"username"`
	ConfigContract            types.String  `tfsdk:"config_contract"`
	Implementation            types.String  `tfsdk:"implementation"`
	Name                      types.String  `tfsdk:"name"`
	Protocol                  types.String  `tfsdk:"protocol"`
	Passkey                   types.String  `tfsdk:"passkey"`
	Cookie                    types.String  `tfsdk:"cookie"`
	CaptchaToken              types.String  `tfsdk:"captcha_token"`
	BaseURL                   types.String  `tfsdk:"base_url"`
	AdditionalParameters      types.String  `tfsdk:"additional_parameters"`
	APIPath                   types.String  `tfsdk:"api_path"`
	SeedTime                  types.Int64   `tfsdk:"seed_time"`
	DownloadClientID          types.Int64   `tfsdk:"download_client_id"`
	Priority                  types.Int64   `tfsdk:"priority"`
	MinimumSeeders            types.Int64   `tfsdk:"minimum_seeders"`
	Delay                     types.Int64   `tfsdk:"delay"`
	ID                        types.Int64   `tfsdk:"id"`
	SeasonPackSeedTime        types.Int64   `tfsdk:"season_pack_seed_time"`
	AnimeStandardFormatSearch types.Bool    `tfsdk:"anime_standard_format_search"`
	AllowZeroSize             types.Bool    `tfsdk:"allow_zero_size"`
	RankedOnly                types.Bool    `tfsdk:"ranked_only"`
	EnableAutomaticSearch     types.Bool    `tfsdk:"enable_automatic_search"`
	EnableRss                 types.Bool    `tfsdk:"enable_rss"`
	EnableInteractiveSearch   types.Bool    `tfsdk:"enable_interactive_search"`
}

func (i Indexer) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"tags":                         types.SetType{}.WithElementType(types.Int64Type),
			"categories":                   types.SetType{}.WithElementType(types.Int64Type),
			"anime_categories":             types.SetType{}.WithElementType(types.Int64Type),
			"api_path":                     types.StringType,
			"additional_parameters":        types.StringType,
			"username":                     types.StringType,
			"config_contract":              types.StringType,
			"implementation":               types.StringType,
			"name":                         types.StringType,
			"protocol":                     types.StringType,
			"passkey":                      types.StringType,
			"cookie":                       types.StringType,
			"captcha_token":                types.StringType,
			"base_url":                     types.StringType,
			"api_key":                      types.StringType,
			"priority":                     types.Int64Type,
			"download_client_id":           types.Int64Type,
			"seed_time":                    types.Int64Type,
			"seed_ratio":                   types.Float64Type,
			"minimum_seeders":              types.Int64Type,
			"delay":                        types.Int64Type,
			"id":                           types.Int64Type,
			"season_pack_seed_time":        types.Int64Type,
			"anime_standard_format_search": types.BoolType,
			"allow_zero_size":              types.BoolType,
			"ranked_only":                  types.BoolType,
			"enable_automatic_search":      types.BoolType,
			"enable_rss":                   types.BoolType,
			"enable_interactive_search":    types.BoolType,
		})
}

func (r *IndexerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerResourceName
}

func (r *IndexerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexers -->Generic Indexer resource. When possible use a specific resource instead.\nFor more information refer to [Indexer](https://wiki.servarr.com/sonarr/settings#indexers) documentation.",
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
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "Indexer configuration template.",
				Required:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "Indexer implementation name.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Indexer name.",
				Required:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "Protocol. Valid values are 'usenet' and 'torrent'.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("usenet", "torrent"),
				},
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Indexer ID.",
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
			"anime_standard_format_search": schema.BoolAttribute{
				MarkdownDescription: "Search anime in standard format.",
				Optional:            true,
				Computed:            true,
			},
			"ranked_only": schema.BoolAttribute{
				MarkdownDescription: "Allow ranked only.",
				Optional:            true,
				Computed:            true,
			},
			"delay": schema.Int64Attribute{
				MarkdownDescription: "Delay before grabbing.",
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
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"api_path": schema.StringAttribute{
				MarkdownDescription: "API path.",
				Optional:            true,
				Computed:            true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
			},
			"captcha_token": schema.StringAttribute{
				MarkdownDescription: "Captcha token.",
				Optional:            true,
				Computed:            true,
			},
			"cookie": schema.StringAttribute{
				MarkdownDescription: "Cookie.",
				Optional:            true,
				Computed:            true,
			},
			"passkey": schema.StringAttribute{
				MarkdownDescription: "Passkey.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
			},
			"categories": schema.SetAttribute{
				MarkdownDescription: "Categories list.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"anime_categories": schema.SetAttribute{
				MarkdownDescription: "Anime categories list.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (r *IndexerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *IndexerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var indexer *Indexer

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Indexer
	request := indexer.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.IndexerApi.CreateIndexer(ctx).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, indexerResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct.
	// this is needed because of many empty fields are unknown in both plan and read
	var state Indexer

	state.writeSensitive(indexer)
	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IndexerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var indexer *Indexer

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get Indexer current value
	response, _, err := r.client.IndexerApi.GetIndexerById(ctx, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct.
	// this is needed because of many empty fields are unknown in both plan and read
	var state Indexer

	state.writeSensitive(indexer)
	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IndexerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var indexer *Indexer

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Indexer
	request := indexer.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.IndexerApi.UpdateIndexer(ctx, strconv.Itoa(int(request.GetId()))).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, indexerResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct.
	// this is needed because of many empty fields are unknown in both plan and read
	var state Indexer

	state.writeSensitive(indexer)
	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IndexerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete Indexer current value
	_, err := r.client.IndexerApi.DeleteIndexer(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, indexerResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+indexerResourceName+": "+req.ID)
}

func (i *Indexer) write(ctx context.Context, indexer *sonarr.IndexerResource, diags *diag.Diagnostics) {
	var localDiag diag.Diagnostics

	i.Tags, localDiag = types.SetValueFrom(ctx, types.Int64Type, indexer.Tags)
	diags.Append(localDiag...)

	i.EnableAutomaticSearch = types.BoolValue(indexer.GetEnableAutomaticSearch())
	i.EnableInteractiveSearch = types.BoolValue(indexer.GetEnableInteractiveSearch())
	i.EnableRss = types.BoolValue(indexer.GetEnableRss())
	i.Priority = types.Int64Value(int64(indexer.GetPriority()))
	i.DownloadClientID = types.Int64Value(int64(indexer.GetDownloadClientId()))
	i.ID = types.Int64Value(int64(indexer.GetId()))
	i.ConfigContract = types.StringValue(indexer.GetConfigContract())
	i.Implementation = types.StringValue(indexer.GetImplementation())
	i.Name = types.StringValue(indexer.GetName())
	i.Protocol = types.StringValue(string(indexer.GetProtocol()))
	i.AnimeCategories = types.SetValueMust(types.Int64Type, nil)
	i.Categories = types.SetValueMust(types.Int64Type, nil)
	helpers.WriteFields(ctx, i, indexer.GetFields(), indexerFields)
}

func (i *Indexer) read(ctx context.Context, diags *diag.Diagnostics) *sonarr.IndexerResource {
	indexer := sonarr.NewIndexerResource()
	indexer.SetEnableAutomaticSearch(i.EnableAutomaticSearch.ValueBool())
	indexer.SetEnableInteractiveSearch(i.EnableInteractiveSearch.ValueBool())
	indexer.SetEnableRss(i.EnableRss.ValueBool())
	indexer.SetPriority(int32(i.Priority.ValueInt64()))
	indexer.SetDownloadClientId(int32(i.DownloadClientID.ValueInt64()))
	indexer.SetId(int32(i.ID.ValueInt64()))
	indexer.SetConfigContract(i.ConfigContract.ValueString())
	indexer.SetImplementation(i.Implementation.ValueString())
	indexer.SetName(i.Name.ValueString())
	indexer.SetProtocol(sonarr.DownloadProtocol(i.Protocol.ValueString()))
	diags.Append(i.Tags.ElementsAs(ctx, &indexer.Tags, true)...)
	indexer.SetFields(helpers.ReadFields(ctx, i, indexerFields))

	return indexer
}

// writeSensitive copy sensitive data from another resource.
func (i *Indexer) writeSensitive(indexer *Indexer) {
	if !indexer.Passkey.IsUnknown() {
		i.Passkey = indexer.Passkey
	}

	if !indexer.APIKey.IsUnknown() {
		i.APIKey = indexer.APIKey
	}
}
