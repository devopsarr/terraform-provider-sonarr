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
	"golift.io/starr/sonarr"
)

const delayProfileResourceName = "delay_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DelayProfileResource{}
var _ resource.ResourceWithImportState = &DelayProfileResource{}

func NewDelayProfileResource() resource.Resource {
	return &DelayProfileResource{}
}

// DelayProfileResource defines the delay profile implementation.
type DelayProfileResource struct {
	client *sonarr.Sonarr
}

// DelayProfile describes the delay profile data model.
type DelayProfile struct {
	Tags                   types.Set    `tfsdk:"tags"`
	PreferredProtocol      types.String `tfsdk:"preferred_protocol"`
	UsenetDelay            types.Int64  `tfsdk:"usenet_delay"`
	TorrentDelay           types.Int64  `tfsdk:"torrent_delay"`
	ID                     types.Int64  `tfsdk:"id"`
	Order                  types.Int64  `tfsdk:"order"`
	EnableUsenet           types.Bool   `tfsdk:"enable_usenet"`
	EnableTorrent          types.Bool   `tfsdk:"enable_torrent"`
	BypassIfHighestQuality types.Bool   `tfsdk:"bypass_if_highest_quality"`
}

func (r *DelayProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + delayProfileResourceName
}

func (r *DelayProfileResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "[subcategory:Profiles]: #\nDelay Profile resource.\nFor more information refer to [Delay Profiles](https://wiki.servarr.com/sonarr/settings#delay-profiles) documentation.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Delay Profile ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"enable_usenet": {
				MarkdownDescription: "Usenet allowed flag at least one of `enable_usenet` and `enable_torrent` must be defined.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"enable_torrent": {
				MarkdownDescription: "Torrent allowed flag at least one of `enable_usenet` and `enable_torrent` must be defined.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"bypass_if_highest_quality": {
				MarkdownDescription: "Bypass for highest quality flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"usenet_delay": {
				MarkdownDescription: "Usenet delay.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"torrent_delay": {
				MarkdownDescription: "Torrent Delay.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"order": {
				MarkdownDescription: "Order.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"tags": {
				MarkdownDescription: "List of associated tags.",
				Required:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"preferred_protocol": {
				MarkdownDescription: "Preferred protocol.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					helpers.StringMatch([]string{"usenet", "torrent"}),
				},
			},
		},
	}, nil
}

func (r *DelayProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DelayProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan DelayProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	data := readDelayProfile(ctx, &plan)

	// Create new DelayProfile
	response, err := r.client.AddDelayProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", delayProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "created"+delayProfileResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeDelayProfile(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *DelayProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state DelayProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get delayprofile current value
	response, err := r.client.GetDelayProfileContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", delayProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+delayProfileResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	result := writeDelayProfile(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *DelayProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan DelayProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := readDelayProfile(ctx, &plan)

	// Update DelayProfile
	response, err := r.client.UpdateDelayProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", delayProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+delayProfileResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeDelayProfile(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *DelayProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DelayProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete delayprofile current value
	err := r.client.DeleteDelayProfileContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", delayProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+delayProfileResourceName+": "+strconv.Itoa(int(state.ID.Value)))
	resp.State.RemoveResource(ctx)
}

func (r *DelayProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			helpers.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+delayProfileResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func writeDelayProfile(ctx context.Context, profile *sonarr.DelayProfile) *DelayProfile {
	output := DelayProfile{
		ID:                     types.Int64{Value: profile.ID},
		EnableUsenet:           types.Bool{Value: profile.EnableUsenet},
		EnableTorrent:          types.Bool{Value: profile.EnableTorrent},
		BypassIfHighestQuality: types.Bool{Value: profile.BypassIfHighestQuality},
		UsenetDelay:            types.Int64{Value: profile.UsenetDelay},
		TorrentDelay:           types.Int64{Value: profile.TorrentDelay},
		Order:                  types.Int64{Value: profile.Order},
		PreferredProtocol:      types.String{Value: profile.PreferredProtocol},
		Tags:                   types.Set{ElemType: types.Int64Type},
	}

	tfsdk.ValueFrom(ctx, profile.Tags, output.Tags.Type(ctx), &output.Tags)

	return &output
}

func readDelayProfile(ctx context.Context, profile *DelayProfile) *sonarr.DelayProfile {
	tags := make([]int, len(profile.Tags.Elems))
	tfsdk.ValueAs(ctx, profile.Tags, &tags)

	return &sonarr.DelayProfile{
		ID:                     profile.ID.Value,
		EnableUsenet:           profile.EnableUsenet.Value,
		EnableTorrent:          profile.EnableTorrent.Value,
		BypassIfHighestQuality: profile.BypassIfHighestQuality.Value,
		UsenetDelay:            profile.UsenetDelay.Value,
		TorrentDelay:           profile.TorrentDelay.Value,
		Order:                  profile.Order.Value,
		PreferredProtocol:      profile.PreferredProtocol.Value,
		Tags:                   tags,
	}
}
