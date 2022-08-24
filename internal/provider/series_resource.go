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

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.ResourceType = resourceSeriesType{}
var _ resource.Resource = resourceSeries{}
var _ resource.ResourceWithImportState = resourceSeries{}

type resourceSeriesType struct{}

type resourceSeries struct {
	provider sonarrProvider
}

// Series is the series resource.
type Series struct {
	Monitored         types.Bool    `tfsdk:"monitored"`
	SeasonFolder      types.Bool    `tfsdk:"season_folder"`
	UseSceneNumbering types.Bool    `tfsdk:"use_scene_numbering"`
	ID                types.Int64   `tfsdk:"id"`
	LanguageProfileID types.Int64   `tfsdk:"language_profile_id"`
	QualityProfileID  types.Int64   `tfsdk:"quality_profile_id"`
	TvdbID            types.Int64   `tfsdk:"tvdb_id"`
	Path              types.String  `tfsdk:"path"`
	Title             types.String  `tfsdk:"title"`
	TitleSlug         types.String  `tfsdk:"title_slug"`
	RootFolderPath    types.String  `tfsdk:"root_folder_path"`
	Tags              []types.Int64 `tfsdk:"tags"`
}

// Season is part of Series.
type Season struct {
	Monitored    types.Bool  `tfsdk:"monitored"`
	SeasonNumber types.Int64 `tfsdk:"season_number"`
}

// AddSeriesOptions is used in series creation.
type AddSeriesOptions struct {
	SearchForMissingEpisodes     types.Bool `tfsdk:"search_for_missing_episodes"`
	SearchForCutoffUnmetEpisodes types.Bool `tfsdk:"search_for_cutoff_unmet_episodes"`
	IgnoreEpisodesWithFiles      types.Bool `tfsdk:"ignore_episodes_with_files"`
	IgnoreEpisodesWithoutFiles   types.Bool `tfsdk:"ignore_episodes_without_files"`
}

// Image is part of Series.
type Image struct {
	CoverType types.String `tfsdk:"cover_type"`
	URL       types.String `tfsdk:"url"`
	RemoteURL types.String `tfsdk:"remote_url"`
	Extension types.String `tfsdk:"extension"`
}

func (t resourceSeriesType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	// TODO: waiting to implement seasons and images until empty conversion is managed natively https://www.terraform.io/plugin/framework/accessing-values#conversion-rules
	return tfsdk.Schema{
		MarkdownDescription: "Series resource",
		Attributes: map[string]tfsdk.Attribute{
			"title": {
				MarkdownDescription: "Series Title",
				Required:            true,
				Type:                types.StringType,
			},
			"title_slug": {
				MarkdownDescription: "Series Title in kebab format",
				Required:            true,
				Type:                types.StringType,
			},
			"monitored": {
				MarkdownDescription: "Monitored flag",
				Required:            true,
				Type:                types.BoolType,
			},
			"season_folder": {
				MarkdownDescription: "Season Folder flag",
				Required:            true,
				Type:                types.BoolType,
			},
			"use_scene_numbering": {
				MarkdownDescription: "Scene numbering flag",
				Required:            true,
				Type:                types.BoolType,
			},
			"language_profile_id": {
				MarkdownDescription: "Language Profile ID ",
				Required:            true,
				Type:                types.Int64Type,
			},
			"quality_profile_id": {
				MarkdownDescription: "Quality Profile ID",
				Required:            true,
				Type:                types.Int64Type,
			},
			"tvdb_id": {
				MarkdownDescription: "TVDB ID",
				Required:            true,
				Type:                types.Int64Type,
			},
			"path": {
				MarkdownDescription: "Series Path",
				Required:            true,
				Type:                types.StringType,
			},
			"root_folder_path": {
				MarkdownDescription: "Series Root Folder",
				Required:            true,
				Type:                types.StringType,
			},
			"tags": {
				MarkdownDescription: "Tags",
				Optional:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"id": {
				MarkdownDescription: "Series ID",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
		},
	}, nil
}

func (t resourceSeriesType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return resourceSeries{
		provider: provider,
	}, diags
}

func (r resourceSeries) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan Series
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Series
	request := readSeries(&plan)
	// TODO: can parametrize AddSeriesOptions
	request.AddOptions = &sonarr.AddSeriesOptions{
		SearchForMissingEpisodes:     true,
		SearchForCutoffUnmetEpisodes: true,
		IgnoreEpisodesWithFiles:      false,
		IgnoreEpisodesWithoutFiles:   false,
	}

	response, err := r.provider.client.AddSeriesContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create series, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "created series: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	var result = *writeSeries(response)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceSeries) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state Series
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get series current value
	response, err := r.provider.client.GetSeriesByIDContext(ctx, state.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read series, got error: %s", err))
		return
	}
	// Map response body to resource schema attribute
	var result = *writeSeries(response)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceSeries) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan Series
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update Series
	request := *readSeries(&plan)
	response, err := r.provider.client.UpdateSeriesContext(ctx, &request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update series, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "update series: "+strconv.Itoa(int(response.ID)))

	// Map response body to resource schema attribute
	result := writeSeries(response)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceSeries) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Series

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete tag current value
	err := r.provider.client.DeleteSeriesContext(ctx, int(state.ID.Value), true, false)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read tags, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceSeries) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func readSeries(series *Series) *sonarr.AddSeriesInput {
	tags := make([]int, len(series.Tags))

	for i, t := range series.Tags {
		tags[i] = int(t.Value)
	}

	return &sonarr.AddSeriesInput{
		ID:                series.ID.Value,
		TvdbID:            series.TvdbID.Value,
		Title:             series.Title.Value,
		TitleSlug:         series.TitleSlug.Value,
		QualityProfileID:  series.QualityProfileID.Value,
		LanguageProfileID: series.LanguageProfileID.Value,
		Monitored:         series.Monitored.Value,
		SeasonFolder:      series.SeasonFolder.Value,
		Path:              series.Path.Value,
		RootFolderPath:    series.Path.Value,
		UseSceneNumbering: series.UseSceneNumbering.Value,
		Tags:              tags,
	}
}

func writeSeries(series *sonarr.Series) *Series {
	tags := make([]types.Int64, len(series.Tags))
	for i, t := range series.Tags {
		tags[i] = types.Int64{Value: int64(t)}
	}

	return &Series{
		Monitored:         types.Bool{Value: series.Monitored},
		SeasonFolder:      types.Bool{Value: series.SeasonFolder},
		UseSceneNumbering: types.Bool{Value: series.UseSceneNumbering},
		ID:                types.Int64{Value: series.ID},
		LanguageProfileID: types.Int64{Value: series.LanguageProfileID},
		QualityProfileID:  types.Int64{Value: series.QualityProfileID},
		TvdbID:            types.Int64{Value: series.TvdbID},
		Path:              types.String{Value: series.Path},
		Title:             types.String{Value: series.Title},
		TitleSlug:         types.String{Value: series.TitleSlug},
		RootFolderPath:    types.String{Value: series.RootFolderPath},
		Tags:              tags,
	}
}
