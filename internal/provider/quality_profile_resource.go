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
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

const qualityProfileResourceName = "quality_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &QualityProfileResource{}
var _ resource.ResourceWithImportState = &QualityProfileResource{}

func NewQualityProfileResource() resource.Resource {
	return &QualityProfileResource{}
}

// QualityProfileResource defines the quality profile implementation.
type QualityProfileResource struct {
	client *sonarr.Sonarr
}

// QualityProfile describes the quality profile data model.
type QualityProfile struct {
	QualityGroups  types.Set    `tfsdk:"quality_groups"`
	Name           types.String `tfsdk:"name"`
	ID             types.Int64  `tfsdk:"id"`
	Cutoff         types.Int64  `tfsdk:"cutoff"`
	UpgradeAllowed types.Bool   `tfsdk:"upgrade_allowed"`
}

// QualityGroup is part of QualityProfile.
type QualityGroup struct {
	Qualities types.Set    `tfsdk:"qualities"`
	Name      types.String `tfsdk:"name"`
	ID        types.Int64  `tfsdk:"id"`
}

// Quality is part of QualityGroup.
type Quality struct {
	Name       types.String `tfsdk:"name"`
	Source     types.String `tfsdk:"source"`
	ID         types.Int64  `tfsdk:"id"`
	Resolution types.Int64  `tfsdk:"resolution"`
}

func (r *QualityProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + qualityProfileResourceName
}

func (r *QualityProfileResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "[subcategory:Profiles]: #\nQuality Profile resource.\nFor more information refer to [Quality Profile](https://wiki.servarr.com/sonarr/settings#quality-profiles) documentation.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Quality Profile ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"name": {
				MarkdownDescription: "Quality Profile Name.",
				Required:            true,
				Type:                types.StringType,
			},
			"upgrade_allowed": {
				MarkdownDescription: "Upgrade allowed flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"cutoff": {
				MarkdownDescription: "Quality ID to which cutoff.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"quality_groups": {
				MarkdownDescription: "Quality groups.",
				Required:            true,
				Attributes:          tfsdk.SetNestedAttributes(r.getQualityGroupSchema().Attributes),
			},
		},
	}, nil
}

func (r QualityProfileResource) getQualityGroupSchema() tfsdk.Schema {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Quality group ID.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"name": {
				MarkdownDescription: "Quality group name.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"qualities": {
				MarkdownDescription: "Qualities in group.",
				Required:            true,
				Attributes:          tfsdk.SetNestedAttributes(r.getQualitySchema().Attributes),
			},
		},
	}
}

func (r QualityProfileResource) getQualitySchema() tfsdk.Schema {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Quality ID.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"resolution": {
				MarkdownDescription: "Resolution.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"name": {
				MarkdownDescription: "Quality name.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"source": {
				MarkdownDescription: "Source.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}
}

func (r *QualityProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *QualityProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan QualityProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	data := readQualityProfile(ctx, &plan)

	// Create new QualityProfile
	response, err := r.client.AddQualityProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+qualityProfileResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeQualityProfile(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *QualityProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state QualityProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get qualityprofile current value
	response, err := r.client.GetQualityProfileContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+qualityProfileResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	result := writeQualityProfile(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *QualityProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan QualityProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := readQualityProfile(ctx, &plan)

	// Update QualityProfile
	response, err := r.client.UpdateQualityProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+qualityProfileResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeQualityProfile(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *QualityProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state QualityProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete qualityprofile current value
	err := r.client.DeleteQualityProfileContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+qualityProfileResourceName+": "+strconv.Itoa(int(state.ID.Value)))
	resp.State.RemoveResource(ctx)
}

func (r *QualityProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			helpers.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+qualityProfileResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func writeQualityProfile(ctx context.Context, profile *sonarr.QualityProfile) *QualityProfile {
	output := QualityProfile{
		UpgradeAllowed: types.Bool{Value: profile.UpgradeAllowed},
		ID:             types.Int64{Value: profile.ID},
		Name:           types.String{Value: profile.Name},
		Cutoff:         types.Int64{Value: profile.Cutoff},
		QualityGroups:  types.Set{ElemType: QualityProfileResource{}.getQualityGroupSchema().Type()},
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
			Qualities: types.Set{ElemType: QualityProfileResource{}.getQualitySchema().Type()},
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
