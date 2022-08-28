package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ provider.ResourceType            = resourceNamingType{}
	_ resource.Resource                = resourceNaming{}
	_ resource.ResourceWithImportState = resourceNaming{}
)

type resourceNamingType struct{}

type resourceNaming struct {
	provider sonarrProvider
}

// Naming is the Naming resource.
type Naming struct {
	RenameEpisodes           types.Bool   `tfsdk:"rename_episodes"`
	ReplaceIllegalCharacters types.Bool   `tfsdk:"replace_illegal_characters"`
	ID                       types.Int64  `tfsdk:"id"`
	MultiEpisodeStyle        types.Int64  `tfsdk:"multi_episode_style"`
	DailyEpisodeFormat       types.String `tfsdk:"daily_episode_format"`
	AnimeEpisodeFormat       types.String `tfsdk:"anime_episode_format"`
	SeriesFolderFormat       types.String `tfsdk:"series_folder_format"`
	SeasonFolderFormat       types.String `tfsdk:"season_folder_format"`
	SpecialsFolderFormat     types.String `tfsdk:"specials_folder_format"`
	StandardEpisodeFormat    types.String `tfsdk:"standard_episode_format"`
}

func (t resourceNamingType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Naming resource",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "ID of naming",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"rename_episodes": {
				MarkdownDescription: "Sonarr will use the existing file name if false",
				Required:            true,
				Type:                types.BoolType,
			},
			"replace_illegal_characters": {
				MarkdownDescription: "Replace illegal characters. They will be removed if false",
				Required:            true,
				Type:                types.BoolType,
			},
			"multi_episode_style": {
				MarkdownDescription: "Multi episode style. 0 - 'Extend' 1 - 'Duplicate' 2 - 'Repeat' 3 - 'Scene' 4 - 'Range' 5 - 'Prefixed Range'",
				Required:            true,
				Type:                types.Int64Type,
			},
			"daily_episode_format": {
				MarkdownDescription: "Daily episode format",
				Required:            true,
				Type:                types.StringType,
			},
			"anime_episode_format": {
				MarkdownDescription: "Anime episode format",
				Required:            true,
				Type:                types.StringType,
			},
			"series_folder_format": {
				MarkdownDescription: "Series folder format",
				Required:            true,
				Type:                types.StringType,
			},
			"season_folder_format": {
				MarkdownDescription: "Season folder format",
				Required:            true,
				Type:                types.StringType,
			},
			"specials_folder_format": {
				MarkdownDescription: "Special folder format",
				Required:            true,
				Type:                types.StringType,
			},
			"standard_episode_format": {
				MarkdownDescription: "Standard episode formatss",
				Required:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (t resourceNamingType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return resourceNaming{
		provider: provider,
	}, diags
}

func (r resourceNaming) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan Naming
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Init call if we remove this it the very first update on a brand new instance will fail
	init, err := r.provider.client.GetNamingContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to init naming, got error: %s", err))

		return
	}

	_, err = r.provider.client.UpdateNamingContext(ctx, init)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to init naming, got error: %s", err))

		return
	}

	// Build Create resource
	data := readNaming(&plan)
	data.ID = 1

	// Create new Naming
	response, err := r.provider.client.UpdateNamingContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create naming, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "created naming: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := writeNaming(response)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceNaming) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state Naming
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get naming current value
	response, err := r.provider.client.GetNamingContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read namings, got error: %s", err))

		return
	}
	// Map response body to resource schema attribute
	result := writeNaming(response)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceNaming) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan Naming
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := readNaming(&plan)

	// Update Naming
	response, err := r.provider.client.UpdateNamingContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update naming, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "update naming: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := writeNaming(response)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceNaming) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Naming cannot be really deleted just removing configuration
	resp.State.RemoveResource(ctx)
}

func (r resourceNaming) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
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
