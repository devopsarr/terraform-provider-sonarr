package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const seriesResourceName = "series"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &SeriesResource{}
	_ resource.ResourceWithImportState = &SeriesResource{}
)

func NewSeriesResource() resource.Resource {
	return &SeriesResource{}
}

// SeriesResource defines the series implementation.
type SeriesResource struct {
	client *sonarr.APIClient
}

// Series describes the series data model.
type Series struct {
	Tags              types.Set    `tfsdk:"tags"`
	Path              types.String `tfsdk:"path"`
	Title             types.String `tfsdk:"title"`
	TitleSlug         types.String `tfsdk:"title_slug"`
	RootFolderPath    types.String `tfsdk:"root_folder_path"`
	ID                types.Int64  `tfsdk:"id"`
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

func (r *SeriesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// TODO: waiting to implement seasons and images until empty conversion is managed natively https://www.terraform.io/plugin/framework/accessing-values#conversion-rules
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Series -->Series resource.\nFor more information refer to [Series](https://wiki.servarr.com/sonarr/library#series) documentation.",
		Attributes: map[string]schema.Attribute{
			"title": schema.StringAttribute{
				MarkdownDescription: "Series Title.",
				Required:            true,
			},
			"title_slug": schema.StringAttribute{
				MarkdownDescription: "Series Title in kebab format.",
				Required:            true,
			},
			"monitored": schema.BoolAttribute{
				MarkdownDescription: "Monitored flag.",
				Required:            true,
			},
			"season_folder": schema.BoolAttribute{
				MarkdownDescription: "Season Folder flag.",
				Required:            true,
			},
			"use_scene_numbering": schema.BoolAttribute{
				MarkdownDescription: "Scene numbering flag.",
				Required:            true,
			},
			"quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Quality Profile ID.",
				Required:            true,
			},
			"tvdb_id": schema.Int64Attribute{
				MarkdownDescription: "TVDB ID.",
				Required:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Series Path.",
				Required:            true,
			},
			"root_folder_path": schema.StringAttribute{
				MarkdownDescription: "Series Root Folder.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Series ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *SeriesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *SeriesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var series *Series

	resp.Diagnostics.Append(req.Plan.Get(ctx, &series)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Series
	request := series.read(ctx, &resp.Diagnostics)
	// TODO: can parametrize AddSeriesOptions
	options := sonarr.NewAddSeriesOptions()
	options.SetSearchForMissingEpisodes(true)
	options.SetSearchForCutoffUnmetEpisodes(true)
	options.SetIgnoreEpisodesWithFiles(false)
	options.SetIgnoreEpisodesWithoutFiles(false)

	request.SetAddOptions(*options)

	response, _, err := r.client.SeriesApi.CreateSeries(ctx).SeriesResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, seriesResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+seriesResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	series.write(ctx, response, &resp.Diagnostics)
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
	response, _, err := r.client.SeriesApi.GetSeriesById(ctx, int32(series.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, seriesResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+seriesResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	series.write(ctx, response, &resp.Diagnostics)
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
	request := series.read(ctx, &resp.Diagnostics)

	// TODO: manage movefiles on sdk
	response, _, err := r.client.SeriesApi.UpdateSeries(ctx, strconv.Itoa(int(request.GetId()))).SeriesResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, seriesResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+seriesResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	series.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &series)...)
}

func (r *SeriesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete series current value
	// TODO: manage delete parameters on SDK
	_, err := r.client.SeriesApi.DeleteSeries(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, seriesResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+seriesResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *SeriesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+seriesResourceName+": "+req.ID)
}

func (s *Series) write(ctx context.Context, series *sonarr.SeriesResource, diags *diag.Diagnostics) {
	var tempDiag diag.Diagnostics

	s.Tags, tempDiag = types.SetValueFrom(ctx, types.Int64Type, series.GetTags())
	diags.Append(tempDiag...)
	s.Monitored = types.BoolValue(series.GetMonitored())
	s.SeasonFolder = types.BoolValue(series.GetSeasonFolder())
	s.UseSceneNumbering = types.BoolValue(series.GetUseSceneNumbering())
	s.ID = types.Int64Value(int64(series.GetId()))
	s.QualityProfileID = types.Int64Value(int64(series.GetQualityProfileId()))
	s.TvdbID = types.Int64Value(int64(series.GetTvdbId()))
	s.Path = types.StringValue(series.GetPath())
	s.Title = types.StringValue(series.GetTitle())
	s.TitleSlug = types.StringValue(series.GetTitleSlug())
	s.RootFolderPath = types.StringValue(series.GetRootFolderPath())
}

func (s *Series) read(ctx context.Context, diags *diag.Diagnostics) *sonarr.SeriesResource {
	series := sonarr.NewSeriesResource()
	series.SetId(int32(s.ID.ValueInt64()))
	series.SetTvdbId(int32(s.TvdbID.ValueInt64()))
	series.SetTitle(s.Title.ValueString())
	series.SetTitleSlug(s.TitleSlug.ValueString())
	series.SetQualityProfileId(int32(s.QualityProfileID.ValueInt64()))
	series.SetMonitored(s.Monitored.ValueBool())
	series.SetSeasonFolder(s.SeasonFolder.ValueBool())
	series.SetPath(s.Path.ValueString())
	series.SetRootFolderPath(s.Path.ValueString())
	series.SetUseSceneNumbering(s.UseSceneNumbering.ValueBool())
	diags.Append(s.Tags.ElementsAs(ctx, &series.Tags, true)...)

	return series
}
