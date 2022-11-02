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
		MarkdownDescription: "<!-- subcategory:Profiles -->Quality Profile resource.\nFor more information refer to [Quality Profile](https://wiki.servarr.com/sonarr/settings#quality-profiles) documentation.",
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
	var profile *QualityProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	data := profile.read(ctx)

	// Create new QualityProfile
	response, err := r.client.AddQualityProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+qualityProfileResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *QualityProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var profile *QualityProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get qualityprofile current value
	response, err := r.client.GetQualityProfileContext(ctx, profile.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+qualityProfileResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *QualityProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var profile *QualityProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := profile.read(ctx)

	// Update QualityProfile
	response, err := r.client.UpdateQualityProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+qualityProfileResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *QualityProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var profile *QualityProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete qualityprofile current value
	err := r.client.DeleteQualityProfileContext(ctx, profile.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+qualityProfileResourceName+": "+strconv.Itoa(int(profile.ID.ValueInt64())))
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

func (p *QualityProfile) write(ctx context.Context, profile *sonarr.QualityProfile) {
	p.UpgradeAllowed = types.BoolValue(profile.UpgradeAllowed)
	p.ID = types.Int64Value(profile.ID)
	p.Name = types.StringValue(profile.Name)
	p.Cutoff = types.Int64Value(profile.Cutoff)
	p.QualityGroups = types.SetValueMust(QualityProfileResource{}.getQualityGroupSchema().Type(), nil)

	qualityGroups := make([]QualityGroup, len(profile.Qualities))
	for n, g := range profile.Qualities {
		qualityGroups[n].write(ctx, g)
	}

	tfsdk.ValueFrom(ctx, qualityGroups, p.QualityGroups.Type(ctx), &p.QualityGroups)
}

func (q *QualityGroup) write(ctx context.Context, group *starr.Quality) {
	var (
		name      string
		id        int64
		qualities []Quality
	)

	if len(group.Items) == 0 {
		name = group.Quality.Name
		id = group.Quality.ID
		qualities = []Quality{{
			ID:         types.Int64Value(group.Quality.ID),
			Name:       types.StringValue(group.Quality.Name),
			Source:     types.StringValue(group.Quality.Source),
			Resolution: types.Int64Value(int64(group.Quality.Resolution)),
		}}
	} else {
		name = group.Name
		id = int64(group.ID)
		qualities = make([]Quality, len(group.Items))
		for m, q := range group.Items {
			qualities[m].write(q)
		}
	}

	q.Name = types.StringValue(name)
	q.ID = types.Int64Value(id)
	q.Qualities = types.SetValueMust(QualityProfileResource{}.getQualitySchema().Type(), nil)

	tfsdk.ValueFrom(ctx, qualities, q.Qualities.Type(ctx), &q.Qualities)
}

func (q *Quality) write(quality *starr.Quality) {
	q.ID = types.Int64Value(quality.Quality.ID)
	q.Name = types.StringValue(quality.Quality.Name)
	q.Source = types.StringValue(quality.Quality.Source)
	q.Resolution = types.Int64Value(int64(quality.Quality.Resolution))
}

func (p *QualityProfile) read(ctx context.Context) *sonarr.QualityProfile {
	groups := make([]QualityGroup, len(p.QualityGroups.Elements()))
	tfsdk.ValueAs(ctx, p.QualityGroups, &groups)
	qualities := make([]*starr.Quality, len(groups))

	for n, g := range groups {
		q := make([]Quality, len(g.Qualities.Elements()))
		tfsdk.ValueAs(ctx, g.Qualities, &q)

		if len(q) == 0 {
			qualities[n] = &starr.Quality{
				Allowed: true,
				Quality: &starr.BaseQuality{
					ID:   g.ID.ValueInt64(),
					Name: g.Name.ValueString(),
				},
			}

			continue
		}

		items := make([]*starr.Quality, len(q))
		for m, q := range q {
			items[m] = q.read()
		}

		qualities[n] = &starr.Quality{
			Name:    g.Name.ValueString(),
			ID:      int(g.ID.ValueInt64()),
			Allowed: true,
			Items:   items,
		}
	}

	return &sonarr.QualityProfile{
		UpgradeAllowed: p.UpgradeAllowed.ValueBool(),
		ID:             p.ID.ValueInt64(),
		Cutoff:         p.Cutoff.ValueInt64(),
		Name:           p.Name.ValueString(),
		Qualities:      qualities,
	}
}

func (q *Quality) read() *starr.Quality {
	return &starr.Quality{
		Allowed: true,
		Quality: &starr.BaseQuality{
			Name:       q.Name.ValueString(),
			ID:         q.ID.ValueInt64(),
			Source:     q.Source.ValueString(),
			Resolution: int(q.Resolution.ValueInt64()),
		},
	}
}
