package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

type resourceSeriesType struct{}

type resourceSeries struct {
	provider provider
}

func (t resourceSeriesType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	//TODO: waiting to implement seasons and images until empty conversion is managed natively https://www.terraform.io/plugin/framework/accessing-values#conversion-rules
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
					tfsdk.UseStateForUnknown(),
				},
			},
		},
	}, nil
}

func (t resourceSeriesType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return resourceSeries{
		provider: provider,
	}, diags
}

func (r resourceSeries) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	// Retrieve values from plan
	var plan Series
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Series
	request := readSeries(&plan)
	//TODO: can parametrize AddSeriesOptions
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

func (r resourceSeries) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
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

func (r resourceSeries) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
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

func (r resourceSeries) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
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

func (r resourceSeries) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	//tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"), id)...)
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
