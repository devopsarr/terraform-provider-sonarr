package provider

import (
	"context"
	"slices"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const qualityProfileResourceName = "quality_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &QualityProfileResource{}
	_ resource.ResourceWithImportState = &QualityProfileResource{}
)

func NewQualityProfileResource() resource.Resource {
	return &QualityProfileResource{}
}

// QualityProfileResource defines the quality profile implementation.
type QualityProfileResource struct {
	client *sonarr.APIClient
}

// QualityProfile describes the quality profile data model.
type QualityProfile struct {
	FormatItems       types.Set    `tfsdk:"format_items"`
	QualityGroups     types.List   `tfsdk:"quality_groups"`
	Name              types.String `tfsdk:"name"`
	ID                types.Int64  `tfsdk:"id"`
	Cutoff            types.Int64  `tfsdk:"cutoff"`
	MinFormatScore    types.Int64  `tfsdk:"min_format_score"`
	CutoffFormatScore types.Int64  `tfsdk:"cutoff_format_score"`
	UpgradeAllowed    types.Bool   `tfsdk:"upgrade_allowed"`
}

func (p QualityProfile) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"quality_groups":      types.ListType{}.WithElementType(QualityGroup{}.getType()),
			"format_items":        types.SetType{}.WithElementType(FormatItem{}.getType()),
			"name":                types.StringType,
			"id":                  types.Int64Type,
			"cutoff":              types.Int64Type,
			"min_format_score":    types.Int64Type,
			"cutoff_format_score": types.Int64Type,
			"upgrade_allowed":     types.BoolType,
		})
}

// QualityGroup is part of QualityProfile.
type QualityGroup struct {
	Qualities types.List   `tfsdk:"qualities"`
	Name      types.String `tfsdk:"name"`
	ID        types.Int64  `tfsdk:"id"`
}

func (g QualityGroup) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"qualities": types.ListType{}.WithElementType(Quality{}.getType()),
			"name":      types.StringType,
			"id":        types.Int64Type,
		})
}

// FormatItem is part of QualityProfile.
type FormatItem struct {
	Name   types.String `tfsdk:"name"`
	Format types.Int64  `tfsdk:"format"`
	Score  types.Int64  `tfsdk:"score"`
}

func (f FormatItem) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"name":   types.StringType,
			"format": types.Int64Type,
			"score":  types.Int64Type,
		})
}

func (r *QualityProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + qualityProfileResourceName
}

func (r *QualityProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Profiles -->Quality Profile resource.\nFor more information refer to [Quality Profile](https://wiki.servarr.com/sonarr/settings#quality-profiles) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Quality Profile ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Quality Profile Name.",
				Required:            true,
			},
			"upgrade_allowed": schema.BoolAttribute{
				MarkdownDescription: "Upgrade allowed flag.",
				Optional:            true,
				Computed:            true,
			},
			"cutoff": schema.Int64Attribute{
				MarkdownDescription: "Quality ID to which cutoff.",
				Optional:            true,
				Computed:            true,
			},
			"cutoff_format_score": schema.Int64Attribute{
				MarkdownDescription: "Cutoff format score.",
				Optional:            true,
				Computed:            true,
			},
			"min_format_score": schema.Int64Attribute{
				MarkdownDescription: "Min format score.",
				Optional:            true,
				Computed:            true,
			},
			"quality_groups": schema.ListNestedAttribute{
				MarkdownDescription: "Ordered list of allowed quality groups.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: r.getQualityGroupSchema().Attributes,
				},
			},
			"format_items": schema.SetNestedAttribute{
				MarkdownDescription: "Format items. Only the ones with score > 0 are needed.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: r.getFormatItemsSchema().Attributes,
				},
			},
		},
	}
}

func (r QualityProfileResource) getQualityGroupSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Quality group ID.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Quality group name.",
				Optional:            true,
				Computed:            true,
			},
			"qualities": schema.ListNestedAttribute{
				MarkdownDescription: "Ordered list of qualities in group.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: r.getQualitySchema().Attributes,
				},
			},
		},
	}
}

func (r QualityProfileResource) getQualitySchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Quality ID.",
				Optional:            true,
				Computed:            true,
				// plan on uptate is unknown for 1 item array
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"resolution": schema.Int64Attribute{
				MarkdownDescription: "Resolution.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Quality name.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source": schema.StringAttribute{
				MarkdownDescription: "Source.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r QualityProfileResource) getFormatItemsSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"format": schema.Int64Attribute{
				MarkdownDescription: "Format.",
				Optional:            true,
				Computed:            true,
			},
			"score": schema.Int64Attribute{
				MarkdownDescription: "Score.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *QualityProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *QualityProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var profile *QualityProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	request := profile.read(ctx, r.getQualityIDs(ctx, &resp.Diagnostics), r.getFormatsIDs(ctx, &resp.Diagnostics), &resp.Diagnostics)

	// Create new QualityProfile
	response, _, err := r.client.QualityProfileApi.CreateQualityProfile(ctx).QualityProfileResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+qualityProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	profile.write(ctx, response, &resp.Diagnostics)
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
	response, _, err := r.client.QualityProfileApi.GetQualityProfileById(ctx, int32(profile.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+qualityProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	profile.write(ctx, response, &resp.Diagnostics)
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
	request := profile.read(ctx, r.getQualityIDs(ctx, &resp.Diagnostics), r.getFormatsIDs(ctx, &resp.Diagnostics), &resp.Diagnostics)

	// Update QualityProfile
	response, _, err := r.client.QualityProfileApi.UpdateQualityProfile(ctx, strconv.Itoa(int(request.GetId()))).QualityProfileResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+qualityProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	profile.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *QualityProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete qualityprofile current value
	_, err := r.client.QualityProfileApi.DeleteQualityProfile(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+qualityProfileResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *QualityProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+qualityProfileResourceName+": "+req.ID)
}

func (p *QualityProfile) write(ctx context.Context, profile *sonarr.QualityProfileResource, diags *diag.Diagnostics) {
	var tempDiag diag.Diagnostics

	p.UpgradeAllowed = types.BoolValue(profile.GetUpgradeAllowed())
	p.ID = types.Int64Value(int64(profile.GetId()))
	p.Name = types.StringValue(profile.GetName())
	p.Cutoff = types.Int64Value(int64(profile.GetCutoff()))
	p.CutoffFormatScore = types.Int64Value(int64(profile.GetCutoffFormatScore()))
	p.MinFormatScore = types.Int64Value(int64(profile.GetMinFormatScore()))

	qualityGroups := make([]QualityGroup, 0, len(profile.GetItems()))

	for _, g := range profile.GetItems() {
		if g.GetAllowed() {
			group := QualityGroup{}
			group.write(ctx, g, diags)
			qualityGroups = append(qualityGroups, group)
		}
	}

	formatItems := make([]FormatItem, 0, len(profile.GetFormatItems()))

	for _, f := range profile.GetFormatItems() {
		if f.GetScore() != 0 {
			format := FormatItem{}
			format.write(f)
			formatItems = append(formatItems, format)
		}
	}

	// Order groups from higher to lower
	slices.Reverse(qualityGroups)
	p.QualityGroups, tempDiag = types.ListValueFrom(ctx, QualityGroup{}.getType(), qualityGroups)
	diags.Append(tempDiag...)
	p.FormatItems, tempDiag = types.SetValueFrom(ctx, FormatItem{}.getType(), formatItems)
	diags.Append(tempDiag...)
}

func (g *QualityGroup) write(ctx context.Context, group *sonarr.QualityProfileQualityItemResource, diags *diag.Diagnostics) {
	var tempDiag diag.Diagnostics

	name := types.StringValue(group.GetName())
	id := types.Int64Value(int64(group.GetId()))

	qualities := make([]Quality, len(group.GetItems()))
	for m, q := range group.GetItems() {
		qualities[m].write(q)
	}

	if len(group.GetItems()) == 0 {
		name = types.StringNull()
		id = types.Int64Null()
		qualities = []Quality{{
			ID:         types.Int64Value(int64(group.Quality.GetId())),
			Name:       types.StringValue(group.Quality.GetName()),
			Source:     types.StringValue(string(group.Quality.GetSource())),
			Resolution: types.Int64Value(int64(group.Quality.GetResolution())),
		}}
	}

	g.Name = name
	g.ID = id
	g.Qualities, tempDiag = types.ListValueFrom(ctx, Quality{}.getType(), &qualities)
	diags.Append(tempDiag...)
}

func (q *Quality) write(quality *sonarr.QualityProfileQualityItemResource) {
	q.ID = types.Int64Value(int64(quality.Quality.GetId()))
	q.Name = types.StringValue(quality.Quality.GetName())
	q.Source = types.StringValue(string(quality.Quality.GetSource()))
	q.Resolution = types.Int64Value(int64(quality.Quality.GetResolution()))
}

func (f *FormatItem) write(format *sonarr.ProfileFormatItemResource) {
	f.Name = types.StringValue(format.GetName())
	f.Format = types.Int64Value(int64(format.GetFormat()))
	f.Score = types.Int64Value(int64(format.GetScore()))
}

func (p *QualityProfile) read(ctx context.Context, qualitiesIDs []int32, formatIDs []int32, diags *diag.Diagnostics) *sonarr.QualityProfileResource {
	var allowedQualities, allowedFormats []int32

	groups := make([]QualityGroup, len(p.QualityGroups.Elements()))
	diags.Append(p.QualityGroups.ElementsAs(ctx, &groups, false)...)

	// Read allowed qualities
	qualities := make([]*sonarr.QualityProfileQualityItemResource, 0, len(groups))
	for _, g := range groups {
		qualities = append(qualities, g.read(ctx, &allowedQualities, diags))
	}

	// Fill qualities with not allowed ones
	for _, id := range qualitiesIDs {
		if !slices.Contains(allowedQualities, id) {
			quality := sonarr.NewQuality()
			quality.SetId(id)

			item := sonarr.NewQualityProfileQualityItemResource()
			item.SetAllowed(false)
			item.SetItems([]*sonarr.QualityProfileQualityItemResource{})
			item.SetQuality(*quality)

			qualities = append(qualities, item)
		}
	}

	// Order groups from higher to lower
	slices.Reverse(qualities)

	formats := make([]FormatItem, len(p.FormatItems.Elements()))
	diags.Append(p.FormatItems.ElementsAs(ctx, &formats, true)...)

	// Read relevant formats
	formatItems := make([]*sonarr.ProfileFormatItemResource, 0, len(formatIDs))
	for _, f := range formats {
		formatItems = append(formatItems, f.read())
	}

	// Fill with irrelevant formats
	for _, id := range formatIDs {
		if !slices.Contains(allowedFormats, id) {
			format := sonarr.NewProfileFormatItemResource()
			format.SetFormat(id)
			format.SetScore(0)
			formatItems = append(formatItems, format)
		}
	}

	profile := sonarr.NewQualityProfileResource()
	profile.SetUpgradeAllowed(p.UpgradeAllowed.ValueBool())
	profile.SetId(int32(p.ID.ValueInt64()))
	profile.SetCutoff(int32(p.Cutoff.ValueInt64()))
	profile.SetMinFormatScore(int32(p.MinFormatScore.ValueInt64()))
	profile.SetCutoffFormatScore(int32(p.CutoffFormatScore.ValueInt64()))
	profile.SetName(p.Name.ValueString())
	profile.SetItems(qualities)
	profile.SetFormatItems(formatItems)

	return profile
}

func (g *QualityGroup) read(ctx context.Context, allowedQualities *[]int32, diags *diag.Diagnostics) *sonarr.QualityProfileQualityItemResource {
	q := make([]Quality, len(g.Qualities.Elements()))
	diags.Append(g.Qualities.ElementsAs(ctx, &q, false)...)

	if len(q) == 1 {
		quality := sonarr.NewQuality()
		quality.SetId(int32(q[0].ID.ValueInt64()))
		quality.SetName(q[0].Name.ValueString())
		quality.SetSource(sonarr.QualitySource(q[0].Source.ValueString()))
		quality.SetResolution(int32(q[0].Resolution.ValueInt64()))

		item := sonarr.NewQualityProfileQualityItemResource()
		item.SetAllowed(true)
		item.SetQuality(*quality)

		*allowedQualities = append(*allowedQualities, int32(q[0].ID.ValueInt64()))

		return item
	}

	items := make([]*sonarr.QualityProfileQualityItemResource, len(q))
	for m, q := range q {
		items[m] = q.read()
		*allowedQualities = append(*allowedQualities, items[m].Quality.GetId())
	}

	quality := sonarr.NewQualityProfileQualityItemResource()
	quality.SetId(int32(g.ID.ValueInt64()))
	quality.SetName(g.Name.ValueString())
	quality.SetAllowed(true)
	quality.SetItems(items)

	return quality
}

func (q *Quality) read() *sonarr.QualityProfileQualityItemResource {
	quality := sonarr.NewQuality()
	quality.SetName(q.Name.ValueString())
	quality.SetId(int32(q.ID.ValueInt64()))
	quality.SetSource(sonarr.QualitySource(q.Source.ValueString()))
	quality.SetResolution(int32(q.Resolution.ValueInt64()))

	item := sonarr.NewQualityProfileQualityItemResource()
	item.SetAllowed(true)
	item.SetQuality(*quality)

	return item
}

func (f *FormatItem) read() *sonarr.ProfileFormatItemResource {
	formatItem := sonarr.NewProfileFormatItemResource()
	formatItem.SetFormat(int32(f.Format.ValueInt64()))
	formatItem.SetName(f.Name.ValueString())
	formatItem.SetScore(int32(f.Score.ValueInt64()))

	return formatItem
}

func (r QualityProfileResource) getQualityIDs(ctx context.Context, diags *diag.Diagnostics) []int32 {
	// Get qualitydefinitions current value
	qualities, _, err := r.client.QualityDefinitionApi.ListQualityDefinition(ctx).Execute()
	if err != nil {
		diags.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, qualityDefinitionsDataSourceName, err))

		return []int32{}
	}

	// Generate a list of quality IDs
	qualityIDs := make([]int32, len(qualities))
	for i, q := range qualities {
		qualityIDs[i] = q.Quality.GetId()
	}

	// Reverse for better visual
	slices.Reverse(qualityIDs)

	return qualityIDs
}

func (r QualityProfileResource) getFormatsIDs(ctx context.Context, diags *diag.Diagnostics) []int32 {
	// Get qualitydefinitions current value
	formats, _, err := r.client.CustomFormatApi.ListCustomFormat(ctx).Execute()
	if err != nil {
		diags.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, customFormatsDataSourceName, err))

		return []int32{}
	}

	// Generate a list of quality IDs
	formatIDs := make([]int32, len(formats))
	for i, f := range formats {
		formatIDs[i] = f.GetId()
	}

	return formatIDs
}
