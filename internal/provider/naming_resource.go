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

const namingResourceName = "naming"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NamingResource{}
var _ resource.ResourceWithImportState = &NamingResource{}

func NewNamingResource() resource.Resource {
	return &NamingResource{}
}

// NamingResource defines the naming implementation.
type NamingResource struct {
	client *sonarr.Sonarr
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
	RenameEpisodes           types.Bool   `tfsdk:"rename_episodes"`
	ReplaceIllegalCharacters types.Bool   `tfsdk:"replace_illegal_characters"`
}

func (r *NamingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + namingResourceName
}

func (r *NamingResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "[subcategory:Media Management]: #\nNaming resource.\nFor more information refer to [Naming](https://wiki.servarr.com/sonarr/settings#community-naming-suggestions) documentation.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Naming ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"rename_episodes": {
				MarkdownDescription: "Sonarr will use the existing file name if false.",
				Required:            true,
				Type:                types.BoolType,
			},
			"replace_illegal_characters": {
				MarkdownDescription: "Replace illegal characters. They will be removed if false.",
				Required:            true,
				Type:                types.BoolType,
			},
			"multi_episode_style": {
				MarkdownDescription: "Multi episode style. 0 - 'Extend' 1 - 'Duplicate' 2 - 'Repeat' 3 - 'Scene' 4 - 'Range' 5 - 'Prefixed Range'.",
				Required:            true,
				Type:                types.Int64Type,
			},
			"daily_episode_format": {
				MarkdownDescription: "Daily episode format.",
				Required:            true,
				Type:                types.StringType,
			},
			"anime_episode_format": {
				MarkdownDescription: "Anime episode format.",
				Required:            true,
				Type:                types.StringType,
			},
			"series_folder_format": {
				MarkdownDescription: "Series folder format.",
				Required:            true,
				Type:                types.StringType,
			},
			"season_folder_format": {
				MarkdownDescription: "Season folder format.",
				Required:            true,
				Type:                types.StringType,
			},
			"specials_folder_format": {
				MarkdownDescription: "Special folder format.",
				Required:            true,
				Type:                types.StringType,
			},
			"standard_episode_format": {
				MarkdownDescription: "Standard episode formatss.",
				Required:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (r *NamingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NamingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan Naming

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Init call if we remove this it the very first update on a brand new instance will fail
	if _, err := r.client.GetNamingContext(ctx); err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to init %s, got error: %s", namingResourceName, err))

		return
	}

	// Build Create resource
	data := readNaming(&plan)
	data.ID = 1

	// Create new Naming
	response, err := r.client.UpdateNamingContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", namingResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+namingResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeNaming(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *NamingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state Naming

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get naming current value
	response, err := r.client.GetNamingContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", namingResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+namingResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	result := writeNaming(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *NamingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan Naming

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := readNaming(&plan)

	// Update Naming
	response, err := r.client.UpdateNamingContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", namingResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+namingResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeNaming(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *NamingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Naming cannot be really deleted just removing configuration
	tflog.Trace(ctx, "decoupled "+namingResourceName+": 1")
	resp.State.RemoveResource(ctx)
}

func (r *NamingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+namingResourceName+": 1")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), 1)...)
}

func writeNaming(naming *sonarr.Naming) *Naming {
	return &Naming{
		RenameEpisodes:           types.Bool{Value: naming.RenameEpisodes},
		ReplaceIllegalCharacters: types.Bool{Value: naming.ReplaceIllegalCharacters},
		ID:                       types.Int64{Value: naming.ID},
		MultiEpisodeStyle:        types.Int64{Value: naming.MultiEpisodeStyle},
		DailyEpisodeFormat:       types.String{Value: naming.DailyEpisodeFormat},
		AnimeEpisodeFormat:       types.String{Value: naming.AnimeEpisodeFormat},
		SeriesFolderFormat:       types.String{Value: naming.SeriesFolderFormat},
		SeasonFolderFormat:       types.String{Value: naming.SeasonFolderFormat},
		SpecialsFolderFormat:     types.String{Value: naming.SpecialsFolderFormat},
		StandardEpisodeFormat:    types.String{Value: naming.StandardEpisodeFormat},
	}
}

func readNaming(naming *Naming) *sonarr.Naming {
	return &sonarr.Naming{
		RenameEpisodes:           naming.RenameEpisodes.Value,
		ReplaceIllegalCharacters: naming.ReplaceIllegalCharacters.Value,
		ID:                       naming.ID.Value,
		MultiEpisodeStyle:        naming.MultiEpisodeStyle.Value,
		DailyEpisodeFormat:       naming.DailyEpisodeFormat.Value,
		AnimeEpisodeFormat:       naming.AnimeEpisodeFormat.Value,
		SeriesFolderFormat:       naming.SeriesFolderFormat.Value,
		SeasonFolderFormat:       naming.SeasonFolderFormat.Value,
		SpecialsFolderFormat:     naming.SpecialsFolderFormat.Value,
		StandardEpisodeFormat:    naming.StandardEpisodeFormat.Value,
	}
}
