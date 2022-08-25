package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ provider.ResourceType            = resourceDelayProfileType{}
	_ resource.Resource                = resourceDelayProfile{}
	_ resource.ResourceWithImportState = resourceDelayProfile{}
)

type resourceDelayProfileType struct{}

type resourceDelayProfile struct {
	provider sonarrProvider
}

// DelayProfile is the delay_profile resource.
type DelayProfile struct {
	EnableUsenet           types.Bool   `tfsdk:"enable_usenet"`
	EnableTorrent          types.Bool   `tfsdk:"enable_torrent"`
	BypassIfHighestQuality types.Bool   `tfsdk:"bypass_if_highest_quality"`
	UsenetDelay            types.Int64  `tfsdk:"usenet_delay"`
	TorrentDelay           types.Int64  `tfsdk:"torrent_delay"`
	ID                     types.Int64  `tfsdk:"id"`
	Order                  types.Int64  `tfsdk:"order"`
	PreferredProtocol      types.String `tfsdk:"preferred_protocol"`
	Tags                   types.Set    `tfsdk:"tags"`
}

func (t resourceDelayProfileType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "DelayProfile resource",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "ID of delayprofile",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"enable_usenet": {
				MarkdownDescription: "Usenet allowed flag at least one of `enable_usenet` and `enable_torrent` must be defined",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"enable_torrent": {
				MarkdownDescription: "Torrent allowed flag at least one of `enable_usenet` and `enable_torrent` must be defined",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"bypass_if_highest_quality": {
				MarkdownDescription: "Bypass for highest quality flag",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"usenet_delay": {
				MarkdownDescription: "Usenet delay",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"torrent_delay": {
				MarkdownDescription: "Torrent Delay",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"order": {
				MarkdownDescription: "Order",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"tags": {
				MarkdownDescription: "List of associated tags",
				Required:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"preferred_protocol": {
				MarkdownDescription: "Preferred protocol",
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

func (t resourceDelayProfileType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return resourceDelayProfile{
		provider: provider,
	}, diags
}

func (r resourceDelayProfile) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan DelayProfile
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	data := readDelayProfile(ctx, &plan)

	// Create new DelayProfile
	response, err := r.provider.client.AddDelayProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create delayprofile, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "created delayprofile: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := writeDelayProfile(ctx, response)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceDelayProfile) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state DelayProfile
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get delayprofile current value
	response, err := r.provider.client.GetDelayProfileContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read delayprofiles, got error: %s", err))

		return
	}
	// Map response body to resource schema attribute
	result := writeDelayProfile(ctx, response)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceDelayProfile) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan DelayProfile
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := readDelayProfile(ctx, &plan)

	// Update DelayProfile
	response, err := r.provider.client.UpdateDelayProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update delayprofile, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "update delayprofile: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := writeDelayProfile(ctx, response)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceDelayProfile) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DelayProfile

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete delayprofile current value
	err := r.provider.client.DeleteDelayProfileContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read delayprofiles, got error: %s", err))

		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceDelayProfile) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
