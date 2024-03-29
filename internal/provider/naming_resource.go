package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const namingResourceName = "naming"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NamingResource{}
	_ resource.ResourceWithImportState = &NamingResource{}
)

func NewNamingResource() resource.Resource {
	return &NamingResource{}
}

// NamingResource defines the naming implementation.
type NamingResource struct {
	client *sonarr.APIClient
	auth   context.Context
}

// Naming describes the naming data model.
type Naming struct {
	DailyEpisodeFormat       types.String `tfsdk:"daily_episode_format"`
	AnimeEpisodeFormat       types.String `tfsdk:"anime_episode_format"`
	SeriesFolderFormat       types.String `tfsdk:"series_folder_format"`
	SeasonFolderFormat       types.String `tfsdk:"season_folder_format"`
	SpecialsFolderFormat     types.String `tfsdk:"specials_folder_format"`
	StandardEpisodeFormat    types.String `tfsdk:"standard_episode_format"`
	ID                       types.Int64  `tfsdk:"id"`
	MultiEpisodeStyle        types.Int64  `tfsdk:"multi_episode_style"`
	ColonReplacementFormat   types.Int64  `tfsdk:"colon_replacement_format"`
	RenameEpisodes           types.Bool   `tfsdk:"rename_episodes"`
	ReplaceIllegalCharacters types.Bool   `tfsdk:"replace_illegal_characters"`
}

func (r *NamingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + namingResourceName
}

func (r *NamingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Media Management -->\nNaming resource.\nFor more information refer to [Naming](https://wiki.servarr.com/sonarr/settings#community-naming-suggestions) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Naming ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"rename_episodes": schema.BoolAttribute{
				MarkdownDescription: "Sonarr will use the existing file name if false.",
				Required:            true,
			},
			"replace_illegal_characters": schema.BoolAttribute{
				MarkdownDescription: "Replace illegal characters. They will be removed if false.",
				Required:            true,
			},
			"multi_episode_style": schema.Int64Attribute{
				MarkdownDescription: "Multi episode style. 0 - 'Extend' 1 - 'Duplicate' 2 - 'Repeat' 3 - 'Scene' 4 - 'Range' 5 - 'Prefixed Range'.",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1, 2, 3, 4, 5),
				},
			},
			"colon_replacement_format": schema.Int64Attribute{
				MarkdownDescription: "Colon replacement format. 0 - 'Delete' 1 - 'Replace with Dash' 2 - 'Replace with Space Dash' 3 - 'Replace with Space Dash Space' 4 - 'Smart Replace'.",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1, 2, 3, 4),
				},
			},
			"daily_episode_format": schema.StringAttribute{
				MarkdownDescription: "Daily episode format.",
				Required:            true,
			},
			"anime_episode_format": schema.StringAttribute{
				MarkdownDescription: "Anime episode format.",
				Required:            true,
			},
			"series_folder_format": schema.StringAttribute{
				MarkdownDescription: "Series folder format.",
				Required:            true,
			},
			"season_folder_format": schema.StringAttribute{
				MarkdownDescription: "Season folder format.",
				Required:            true,
			},
			"specials_folder_format": schema.StringAttribute{
				MarkdownDescription: "Special folder format.",
				Required:            true,
			},
			"standard_episode_format": schema.StringAttribute{
				MarkdownDescription: "Standard episode formatss.",
				Required:            true,
			},
		},
	}
}

func (r *NamingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *NamingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var naming *Naming

	resp.Diagnostics.Append(req.Plan.Get(ctx, &naming)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Init call if we remove this it the very first update on a brand new instance will fail
	if _, _, err := r.client.NamingConfigAPI.GetNamingConfig(r.auth).Execute(); err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError("init", namingResourceName, err))

		return
	}

	// Build Create resource
	request := naming.read()
	request.SetId(1)

	// Create new Naming
	response, _, err := r.client.NamingConfigAPI.UpdateNamingConfig(r.auth, strconv.Itoa(int(request.GetId()))).NamingConfigResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, namingResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+namingResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	naming.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &naming)...)
}

func (r *NamingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var naming *Naming

	resp.Diagnostics.Append(req.State.Get(ctx, &naming)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get naming current value
	response, _, err := r.client.NamingConfigAPI.GetNamingConfig(r.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, namingResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+namingResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	naming.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &naming)...)
}

func (r *NamingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var naming *Naming

	resp.Diagnostics.Append(req.Plan.Get(ctx, &naming)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	request := naming.read()

	// Update Naming
	response, _, err := r.client.NamingConfigAPI.UpdateNamingConfig(r.auth, strconv.Itoa(int(request.GetId()))).NamingConfigResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, namingResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+namingResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	naming.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &naming)...)
}

func (r *NamingResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Naming cannot be really deleted just removing configuration
	tflog.Trace(ctx, "decoupled "+namingResourceName+": 1")
	resp.State.RemoveResource(ctx)
}

func (r *NamingResource) ImportState(ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Trace(ctx, "imported "+namingResourceName+": 1")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), 1)...)
}

func (n *Naming) write(naming *sonarr.NamingConfigResource) {
	n.RenameEpisodes = types.BoolValue(naming.GetRenameEpisodes())
	n.ReplaceIllegalCharacters = types.BoolValue(naming.GetReplaceIllegalCharacters())
	n.ID = types.Int64Value(int64(naming.GetId()))
	n.MultiEpisodeStyle = types.Int64Value(int64(naming.GetMultiEpisodeStyle()))
	n.ColonReplacementFormat = types.Int64Value(int64(naming.GetColonReplacementFormat()))
	n.DailyEpisodeFormat = types.StringValue(naming.GetDailyEpisodeFormat())
	n.AnimeEpisodeFormat = types.StringValue(naming.GetAnimeEpisodeFormat())
	n.SeriesFolderFormat = types.StringValue(naming.GetSeriesFolderFormat())
	n.SeasonFolderFormat = types.StringValue(naming.GetSeasonFolderFormat())
	n.SpecialsFolderFormat = types.StringValue(naming.GetSpecialsFolderFormat())
	n.StandardEpisodeFormat = types.StringValue(naming.GetStandardEpisodeFormat())
}

func (n *Naming) read() *sonarr.NamingConfigResource {
	naming := sonarr.NewNamingConfigResource()
	naming.SetAnimeEpisodeFormat(n.AnimeEpisodeFormat.ValueString())
	naming.SetDailyEpisodeFormat(n.DailyEpisodeFormat.ValueString())
	naming.SetId(int32(n.ID.ValueInt64()))
	naming.SetRenameEpisodes(n.RenameEpisodes.ValueBool())
	naming.SetReplaceIllegalCharacters(n.ReplaceIllegalCharacters.ValueBool())
	naming.SetMultiEpisodeStyle(int32(n.MultiEpisodeStyle.ValueInt64()))
	naming.SetColonReplacementFormat(int32(n.ColonReplacementFormat.ValueInt64()))
	naming.SetSeriesFolderFormat(n.SeriesFolderFormat.ValueString())
	naming.SetSeasonFolderFormat(n.SeasonFolderFormat.ValueString())
	naming.SetSpecialsFolderFormat(n.SpecialsFolderFormat.ValueString())
	naming.SetStandardEpisodeFormat(n.StandardEpisodeFormat.ValueString())

	return naming
}
