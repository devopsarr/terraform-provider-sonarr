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
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ provider.ResourceType            = resourceQualityProfileType{}
	_ resource.Resource                = resourceQualityProfile{}
	_ resource.ResourceWithImportState = resourceQualityProfile{}
)

type resourceQualityProfileType struct{}

type resourceQualityProfile struct {
	provider sonarrProvider
}

// QualityProfile is the quality_profile resource.
type QualityProfile struct {
	UpgradeAllowed types.Bool   `tfsdk:"upgrade_allowed"`
	ID             types.Int64  `tfsdk:"id"`
	Cutoff         types.Int64  `tfsdk:"cutoff"`
	Name           types.String `tfsdk:"name"`
	QualityGroups  types.Set    `tfsdk:"quality_groups"`
}

// QualityGroup is part of QualityProfile.
type QualityGroup struct {
	ID        types.Int64  `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Qualities types.Set    `tfsdk:"qualities"`
}

// Quality is part of QualityGroup.
type Quality struct {
	ID         types.Int64  `tfsdk:"id"`
	Resolution types.Int64  `tfsdk:"resolution"`
	Name       types.String `tfsdk:"name"`
	Source     types.String `tfsdk:"source"`
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
					resource.UseStateForUnknown(),
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
				Attributes:          tfsdk.SetNestedAttributes(t.getQualityGroupSchema().Attributes),
			},
		},
	}, nil
}

func (t resourceQualityProfileType) getQualityGroupSchema() tfsdk.Schema {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
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
				Attributes:          tfsdk.SetNestedAttributes(t.getQualitySchema().Attributes),
			},
		},
	}
}

func (t resourceQualityProfileType) getQualitySchema() tfsdk.Schema {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
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
		},
	}
}

func (t resourceQualityProfileType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return resourceQualityProfile{
		provider: provider,
	}, diags
}

func (r resourceQualityProfile) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan QualityProfile
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	data := readQualityProfile(ctx, &plan)

	// Create new QualityProfile
	response, err := r.provider.client.AddQualityProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create qualityprofile, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "created qualityprofile: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := writeQualityProfile(ctx, response)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceQualityProfile) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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
	result := writeQualityProfile(ctx, response)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceQualityProfile) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan QualityProfile
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := readQualityProfile(ctx, &plan)

	// Update QualityProfile
	response, err := r.provider.client.UpdateQualityProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update qualityprofile, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "update qualityprofile: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := writeQualityProfile(ctx, response)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceQualityProfile) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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

func (r resourceQualityProfile) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
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

func writeQualityProfile(ctx context.Context, profile *sonarr.QualityProfile) *QualityProfile {
	output := QualityProfile{
		UpgradeAllowed: types.Bool{Value: profile.UpgradeAllowed},
		ID:             types.Int64{Value: profile.ID},
		Name:           types.String{Value: profile.Name},
		Cutoff:         types.Int64{Value: profile.Cutoff},
		QualityGroups:  types.Set{ElemType: resourceQualityProfileType{}.getQualityGroupSchema().AttributeType()},
	}
	qualityGroups := *writeQualityGroups(ctx, profile.Qualities)
	tfsdk.ValueFrom(ctx, qualityGroups, output.QualityGroups.Type(ctx), &output.QualityGroups)

	return &output
}

func writeQualityGroups(ctx context.Context, groups []*starr.Quality) *[]QualityGroup {
	qualityGroups := make([]QualityGroup, len(groups))

	for n, g := range groups {
		var (
			name      string
			id        int64
			qualities []Quality
		)

		if len(g.Items) == 0 {
			name = g.Quality.Name
			id = g.Quality.ID
			qualities = []Quality{{
				ID:         types.Int64{Value: g.Quality.ID},
				Name:       types.String{Value: g.Quality.Name},
				Source:     types.String{Value: g.Quality.Source},
				Resolution: types.Int64{Value: int64(g.Quality.Resolution)},
			}}
		} else {
			name = g.Name
			id = int64(g.ID)
			qualities = *writeQualities(g.Items)
		}

		qualityGroups[n] = QualityGroup{
			Name:      types.String{Value: name},
			ID:        types.Int64{Value: id},
			Qualities: types.Set{ElemType: resourceQualityProfileType{}.getQualitySchema().AttributeType()},
		}
		tfsdk.ValueFrom(ctx, qualities, qualityGroups[n].Qualities.Type(ctx), &qualityGroups[n].Qualities)
	}

	return &qualityGroups
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

func readQualityProfile(ctx context.Context, profile *QualityProfile) *sonarr.QualityProfile {
	groups := make([]QualityGroup, len(profile.QualityGroups.Elems))
	tfsdk.ValueAs(ctx, profile.QualityGroups, &groups)
	qualities := make([]*starr.Quality, len(groups))

	for n, g := range groups {
		q := make([]Quality, len(g.Qualities.Elems))
		tfsdk.ValueAs(ctx, g.Qualities, &q)

		if len(q) == 0 {
			qualities[n] = &starr.Quality{
				Allowed: true,
				Quality: &starr.BaseQuality{
					ID:   g.ID.Value,
					Name: g.Name.Value,
				},
			}

			continue
		}

		items := readQualities(&q)

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
