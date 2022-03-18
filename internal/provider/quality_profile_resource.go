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
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

type resourceQualityProfileType struct{}

type resourceQualityProfile struct {
	provider provider
}

func (t resourceQualityProfileType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "QualityProfile resource",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "ID of qualityprofile",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"name": {
				MarkdownDescription: "Name",
				Required:            true,
				Type:                types.StringType,
			},
			"upgrade_allowed": {
				MarkdownDescription: "Upgrade allowed flag",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"cutoff": {
				MarkdownDescription: "Quality ID to which cutoff",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"quality_groups": {
				MarkdownDescription: "Quality groups",
				Required:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						MarkdownDescription: "ID of quality group",
						Optional:            true,
						Computed:            true,
						Type:                types.Int64Type,
					},
					"name": {
						MarkdownDescription: "Name of quality group",
						Optional:            true,
						Computed:            true,
						Type:                types.StringType,
					},
					"qualities": {
						MarkdownDescription: "Qualities in group",
						Required:            true,
						Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
							"id": {
								MarkdownDescription: "ID of quality group",
								Optional:            true,
								Computed:            true,
								Type:                types.Int64Type,
							},
							"resolution": {
								MarkdownDescription: "Resolution",
								Optional:            true,
								Computed:            true,
								Type:                types.Int64Type,
							},
							"name": {
								MarkdownDescription: "Name of quality group",
								Optional:            true,
								Computed:            true,
								Type:                types.StringType,
							},
							"source": {
								MarkdownDescription: "Source",
								Optional:            true,
								Computed:            true,
								Type:                types.StringType,
							},
						}, tfsdk.SetNestedAttributesOptions{}),
					},
				}, tfsdk.SetNestedAttributesOptions{}),
			},
		},
	}, nil
}

func (t resourceQualityProfileType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return resourceQualityProfile{
		provider: provider,
	}, diags
}

func (r resourceQualityProfile) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	// Retrieve values from plan
	var plan QualityProfile
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	data := readQualityProfile(&plan)

	// Create new QualityProfile
	response, err := r.provider.client.AddQualityProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create qualityprofile, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "created qualityprofile: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := *writeQualityProfile(response)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceQualityProfile) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// Get current state
	var state QualityProfile
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get qualityprofile current value
	response, err := r.provider.client.GetQualityProfileContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read qualityprofiles, got error: %s", err))
		return
	}
	// Map response body to resource schema attribute
	result := *writeQualityProfile(response)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceQualityProfile) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// Get plan values
	var plan QualityProfile
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := readQualityProfile(&plan)

	// Update QualityProfile
	response, err := r.provider.client.UpdateQualityProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update qualityprofile, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "update qualityprofile: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := *writeQualityProfile(response)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceQualityProfile) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state QualityProfile

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete qualityprofile current value
	err := r.provider.client.DeleteQualityProfileContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read qualityprofiles, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceQualityProfile) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
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

func writeQualityProfile(profile *sonarr.QualityProfile) *QualityProfile {
	qualityGroups := make([]QualityGroup, len(profile.Qualities))
	for n, g := range profile.Qualities {
		if len(g.Items) == 0 {
			qualityGroups[n] = QualityGroup{
				Name: types.String{Value: g.Quality.Name},
				ID:   types.Int64{Value: g.Quality.ID},
				Qualities: []Quality{{
					ID:         types.Int64{Value: g.Quality.ID},
					Name:       types.String{Value: g.Quality.Name},
					Source:     types.String{Value: g.Quality.Source},
					Resolution: types.Int64{Value: int64(g.Quality.Resolution)},
				}},
			}
			continue
		}

		qualities := writeQualities(g.Items)
		qualityGroups[n] = QualityGroup{
			Name:      types.String{Value: g.Name},
			ID:        types.Int64{Value: int64(g.ID)},
			Qualities: *qualities,
		}
	}

	return &QualityProfile{
		UpgradeAllowed: types.Bool{Value: profile.UpgradeAllowed},
		ID:             types.Int64{Value: profile.ID},
		Name:           types.String{Value: profile.Name},
		Cutoff:         types.Int64{Value: profile.Cutoff},
		QualityGroups:  qualityGroups,
	}
}

func writeQualities(qualities []*starr.Quality) *[]Quality {
	output := make([]Quality, len(qualities))
	for m, q := range qualities {
		output[m] = Quality{
			ID:         types.Int64{Value: q.Quality.ID},
			Name:       types.String{Value: q.Quality.Name},
			Source:     types.String{Value: q.Quality.Source},
			Resolution: types.Int64{Value: int64(q.Quality.Resolution)},
		}
	}
	return &output
}

func readQualityProfile(profile *QualityProfile) *sonarr.QualityProfile {
	qualities := make([]*starr.Quality, len(profile.QualityGroups))
	for n, g := range profile.QualityGroups {
		if len(g.Qualities) == 0 {
			qualities[n] = &starr.Quality{
				Allowed: true,
				Quality: &starr.BaseQuality{
					ID:   g.ID.Value,
					Name: g.Name.Value,
				},
			}
			continue
		}

		items := readQualities(&g.Qualities)

		qualities[n] = &starr.Quality{
			Name:    g.Name.Value,
			ID:      int(g.ID.Value),
			Allowed: true,
			Items:   items,
		}
	}

	return &sonarr.QualityProfile{
		UpgradeAllowed: profile.UpgradeAllowed.Value,
		ID:             profile.ID.Value,
		Cutoff:         profile.Cutoff.Value,
		Name:           profile.Name.Value,
		Qualities:      qualities,
	}
}

func readQualities(qualities *[]Quality) []*starr.Quality {
	output := make([]*starr.Quality, len(*qualities))
	for m, q := range *qualities {
		output[m] = &starr.Quality{
			Allowed: true,
			Quality: &starr.BaseQuality{
				Name:       q.Name.Value,
				ID:         q.ID.Value,
				Source:     q.Source.Value,
				Resolution: int(q.Resolution.Value),
			},
		}
	}
	return output
}
