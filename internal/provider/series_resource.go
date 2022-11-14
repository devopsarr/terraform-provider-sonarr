package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const seriesResourceName = "series"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SeriesResource{}
var _ resource.ResourceWithImportState = &SeriesResource{}

func NewSeriesResource() resource.Resource {
	return &SeriesResource{}
}

// SeriesResource defines the series implementation.
type SeriesResource struct {
	client *sonarr.Sonarr
}

// Series describes the series data model.
type Series struct {
	Tags              types.Set    `tfsdk:"tags"`
	Path              types.String `tfsdk:"path"`
	Title             types.String `tfsdk:"title"`
	TitleSlug         types.String `tfsdk:"title_slug"`
	RootFolderPath    types.String `tfsdk:"root_folder_path"`
	ID                types.Int64  `tfsdk:"id"`
	LanguageProfileID types.Int64  `tfsdk:"language_profile_id"`
	QualityProfileID  types.Int64  `tfsdk:"quality_profile_id"`
	TvdbID            types.Int64  `tfsdk:"tvdb_id"`
	Monitored         types.Bool   `tfsdk:"monitored"`
	SeasonFolder      types.Bool   `tfsdk:"season_folder"`
	UseSceneNumbering types.Bool   `tfsdk:"use_scene_numbering"`
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

func (r *SeriesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + seriesResourceName
}

func (r *SeriesResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	// TODO: waiting to implement seasons and images until empty conversion is managed natively https://www.terraform.io/plugin/framework/accessing-values#conversion-rules
	return tfsdk.Schema{
		MarkdownDescription: "<!-- subcategory:Series -->Series resource.\nFor more information refer to [Series](https://wiki.servarr.com/sonarr/library#series) documentation.",
		Attributes: map[string]tfsdk.Attribute{
			"title": {
				MarkdownDescription: "Series Title.",
				Required:            true,
				Type:                types.StringType,
			},
			"title_slug": {
				MarkdownDescription: "Series Title in kebab format.",
				Required:            true,
				Type:                types.StringType,
			},
			"monitored": {
				MarkdownDescription: "Monitored flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"season_folder": {
				MarkdownDescription: "Season Folder flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"use_scene_numbering": {
				MarkdownDescription: "Scene numbering flag.",
				Required:            true,
				Type:                types.BoolType,
			},
			"language_profile_id": {
				MarkdownDescription: "Language Profile ID .",
				Required:            true,
				Type:                types.Int64Type,
			},
			"quality_profile_id": {
				MarkdownDescription: "Quality Profile ID.",
				Required:            true,
				Type:                types.Int64Type,
			},
			"tvdb_id": {
				MarkdownDescription: "TVDB ID.",
				Required:            true,
				Type:                types.Int64Type,
			},
			"path": {
				MarkdownDescription: "Series Path.",
				Required:            true,
				Type:                types.StringType,
			},
			"root_folder_path": {
				MarkdownDescription: "Series Root Folder.",
				Required:            true,
				Type:                types.StringType,
			},
			"tags": {
				MarkdownDescription: "Tags.",
				Optional:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"id": {
				MarkdownDescription: "Series ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
		},
	}, nil
}

func (r *SeriesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *SeriesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var series *Series

	resp.Diagnostics.Append(req.Plan.Get(ctx, &series)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Series
	request := series.read(ctx)
	// TODO: can parametrize AddSeriesOptions
	request.AddOptions = &sonarr.AddSeriesOptions{
		SearchForMissingEpisodes:     true,
		SearchForCutoffUnmetEpisodes: true,
		IgnoreEpisodesWithFiles:      false,
		IgnoreEpisodesWithoutFiles:   false,
	}

	response, err := r.client.AddSeriesContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", seriesResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+seriesResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	series.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &series)...)
}

func (r *SeriesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var series *Series

	resp.Diagnostics.Append(req.State.Get(ctx, &series)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get series current value
	response, err := r.client.GetSeriesByIDContext(ctx, series.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", seriesResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+seriesResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	series.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &series)...)
}

func (r *SeriesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var series *Series

	resp.Diagnostics.Append(req.Plan.Get(ctx, &series)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Series
	request := series.read(ctx)

	response, err := r.client.UpdateSeriesContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", seriesResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+seriesResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	series.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &series)...)
}

func (r *SeriesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var series *Series

	resp.Diagnostics.Append(req.State.Get(ctx, &series)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete series current value
	err := r.client.DeleteSeriesContext(ctx, int(series.ID.ValueInt64()), true, false)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to delete %s, got error: %s", seriesResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+seriesResourceName+": "+strconv.Itoa(int(series.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *SeriesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+seriesResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (s *Series) write(ctx context.Context, series *sonarr.Series) {
	s.Monitored = types.BoolValue(series.Monitored)
	s.SeasonFolder = types.BoolValue(series.SeasonFolder)
	s.UseSceneNumbering = types.BoolValue(series.UseSceneNumbering)
	s.ID = types.Int64Value(series.ID)
	s.LanguageProfileID = types.Int64Value(series.LanguageProfileID)
	s.QualityProfileID = types.Int64Value(series.QualityProfileID)
	s.TvdbID = types.Int64Value(series.TvdbID)
	s.Path = types.StringValue(series.Path)
	s.Title = types.StringValue(series.Title)
	s.TitleSlug = types.StringValue(series.TitleSlug)
	s.RootFolderPath = types.StringValue(series.RootFolderPath)
	s.Tags = types.SetValueMust(types.Int64Type, nil)
	tfsdk.ValueFrom(ctx, series.Tags, s.Tags.Type(ctx), &s.Tags)
}

func (s *Series) read(ctx context.Context) *sonarr.AddSeriesInput {
	tags := make([]int, len(s.Tags.Elements()))
	tfsdk.ValueAs(ctx, s.Tags, &tags)

	return &sonarr.AddSeriesInput{
		ID:                s.ID.ValueInt64(),
		TvdbID:            s.TvdbID.ValueInt64(),
		Title:             s.Title.ValueString(),
		TitleSlug:         s.TitleSlug.ValueString(),
		QualityProfileID:  s.QualityProfileID.ValueInt64(),
		LanguageProfileID: s.LanguageProfileID.ValueInt64(),
		Monitored:         s.Monitored.ValueBool(),
		SeasonFolder:      s.SeasonFolder.ValueBool(),
		Path:              s.Path.ValueString(),
		RootFolderPath:    s.Path.ValueString(),
		UseSceneNumbering: s.UseSceneNumbering.ValueBool(),
		Tags:              tags,
	}
}
