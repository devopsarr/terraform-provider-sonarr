package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const delayProfileResourceName = "delay_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DelayProfileResource{}
	_ resource.ResourceWithImportState = &DelayProfileResource{}
)

func NewDelayProfileResource() resource.Resource {
	return &DelayProfileResource{}
}

// DelayProfileResource defines the delay profile implementation.
type DelayProfileResource struct {
	client *sonarr.APIClient
}

// DelayProfile describes the delay profile data model.
type DelayProfile struct {
	Tags                           types.Set    `tfsdk:"tags"`
	PreferredProtocol              types.String `tfsdk:"preferred_protocol"`
	UsenetDelay                    types.Int64  `tfsdk:"usenet_delay"`
	TorrentDelay                   types.Int64  `tfsdk:"torrent_delay"`
	ID                             types.Int64  `tfsdk:"id"`
	Order                          types.Int64  `tfsdk:"order"`
	MinimumCustomFormatScore       types.Int64  `tfsdk:"minimum_custom_format_score"`
	EnableUsenet                   types.Bool   `tfsdk:"enable_usenet"`
	EnableTorrent                  types.Bool   `tfsdk:"enable_torrent"`
	BypassIfHighestQuality         types.Bool   `tfsdk:"bypass_if_highest_quality"`
	BypassIfAboveCustomFormatScore types.Bool   `tfsdk:"bypass_if_above_custom_format_score"`
}

func (r *DelayProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + delayProfileResourceName
}

func (r *DelayProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Profiles -->Delay Profile resource.\nFor more information refer to [Delay Profiles](https://wiki.servarr.com/sonarr/settings#delay-profiles) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Delay Profile ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"enable_usenet": schema.BoolAttribute{
				MarkdownDescription: "Usenet allowed flag at least one of `enable_usenet` and `enable_torrent` must be defined.",
				Optional:            true,
				Computed:            true,
			},
			"enable_torrent": schema.BoolAttribute{
				MarkdownDescription: "Torrent allowed flag at least one of `enable_usenet` and `enable_torrent` must be defined.",
				Optional:            true,
				Computed:            true,
			},
			"bypass_if_highest_quality": schema.BoolAttribute{
				MarkdownDescription: "Bypass for highest quality flag.",
				Optional:            true,
				Computed:            true,
			},
			"bypass_if_above_custom_format_score": schema.BoolAttribute{
				MarkdownDescription: "Bypass for higher custom format score flag.",
				Optional:            true,
				Computed:            true,
			},
			"usenet_delay": schema.Int64Attribute{
				MarkdownDescription: "Usenet delay.",
				Optional:            true,
				Computed:            true,
			},
			"torrent_delay": schema.Int64Attribute{
				MarkdownDescription: "Torrent Delay.",
				Optional:            true,
				Computed:            true,
			},
			"order": schema.Int64Attribute{
				MarkdownDescription: "Order.",
				Optional:            true,
				Computed:            true,
			},
			"minimum_custom_format_score": schema.Int64Attribute{
				MarkdownDescription: "Minimum custom format score.",
				Optional:            true,
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Required:            true,
				ElementType:         types.Int64Type,
			},
			"preferred_protocol": schema.StringAttribute{
				MarkdownDescription: "Preferred protocol.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("usenet", "torrent"),
				},
			},
		},
	}
}

func (r *DelayProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *DelayProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var profile *DelayProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	request := profile.read(ctx)

	// Create new DelayProfile
	response, _, err := r.client.DelayProfileApi.CreateDelayProfile(ctx).DelayProfileResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, delayProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "created"+delayProfileResourceName+": "+strconv.Itoa(int(response.GetId())))

	// Set order on create
	if !profile.Order.IsUnknown() {
		response.Order = request.Order

		response, _, err = r.client.DelayProfileApi.UpdateDelayProfile(ctx, strconv.Itoa(int(response.GetId()))).DelayProfileResource(*response).Execute()
		if err != nil {
			resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, delayProfileResourceName, err))

			return
		}
	}

	// Generate resource state struct
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *DelayProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var profile *DelayProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get delayprofile current value
	response, _, err := r.client.DelayProfileApi.GetDelayProfileById(ctx, int32(profile.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, delayProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+delayProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *DelayProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var profile *DelayProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	request := profile.read(ctx)

	// Update DelayProfile
	response, _, err := r.client.DelayProfileApi.UpdateDelayProfile(ctx, strconv.Itoa(int(request.GetId()))).DelayProfileResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, delayProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+delayProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *DelayProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var profile *DelayProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete delayprofile current value
	_, err := r.client.DelayProfileApi.DeleteDelayProfile(ctx, int32(profile.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, delayProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+delayProfileResourceName+": "+strconv.Itoa(int(profile.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *DelayProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+delayProfileResourceName+": "+req.ID)
}

func (p *DelayProfile) write(ctx context.Context, profile *sonarr.DelayProfileResource) {
	p.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, profile.GetTags())
	p.ID = types.Int64Value(int64(profile.GetId()))
	p.EnableUsenet = types.BoolValue(profile.GetEnableUsenet())
	p.EnableTorrent = types.BoolValue(profile.GetEnableTorrent())
	p.BypassIfHighestQuality = types.BoolValue(profile.GetBypassIfHighestQuality())
	p.BypassIfAboveCustomFormatScore = types.BoolValue(profile.GetBypassIfAboveCustomFormatScore())
	p.UsenetDelay = types.Int64Value(int64(profile.GetUsenetDelay()))
	p.TorrentDelay = types.Int64Value(int64(profile.GetTorrentDelay()))
	p.Order = types.Int64Value(int64(profile.GetOrder()))
	p.MinimumCustomFormatScore = types.Int64Value(int64(profile.GetMinimumCustomFormatScore()))
	p.PreferredProtocol = types.StringValue(string(profile.GetPreferredProtocol()))
}

func (p *DelayProfile) read(ctx context.Context) *sonarr.DelayProfileResource {
	tags := make([]*int32, len(p.Tags.Elements()))
	tfsdk.ValueAs(ctx, p.Tags, &tags)

	profile := sonarr.NewDelayProfileResource()
	profile.SetId(int32(p.ID.ValueInt64()))
	profile.SetBypassIfHighestQuality(p.BypassIfHighestQuality.ValueBool())
	profile.SetBypassIfAboveCustomFormatScore(p.BypassIfAboveCustomFormatScore.ValueBool())
	profile.SetEnableTorrent(p.EnableTorrent.ValueBool())
	profile.SetEnableUsenet(p.EnableUsenet.ValueBool())
	profile.SetOrder(int32(p.Order.ValueInt64()))
	profile.SetMinimumCustomFormatScore(int32(p.MinimumCustomFormatScore.ValueInt64()))
	profile.SetPreferredProtocol(sonarr.DownloadProtocol(p.PreferredProtocol.ValueString()))
	profile.SetTags(tags)
	profile.SetTorrentDelay(int32(p.TorrentDelay.ValueInt64()))
	profile.SetUsenetDelay(int32(p.UsenetDelay.ValueInt64()))

	return profile
}
