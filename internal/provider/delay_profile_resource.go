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

type resourceDelayProfileType struct{}

type resourceDelayProfile struct {
	provider provider
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
					tfsdk.UseStateForUnknown(),
				},
			},
			"enable_usenet": {
				MarkdownDescription: "Usenet allowed flag at least one of enable_usenet and enable_torrent must be defined",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"enable_torrent": {
				MarkdownDescription: "Torrent allowed flag at least one of enable_usenet and enable_torrent must be defined",
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
				Type: types.ListType{
					ElemType: types.Int64Type,
				},
			},
			//TODO: add validation
			"preferred_protocol": {
				MarkdownDescription: "Preferred protocol",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (t resourceDelayProfileType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return resourceDelayProfile{
		provider: provider,
	}, diags
}

func (r resourceDelayProfile) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	// Retrieve values from plan
	var plan DelayProfile
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	data := readDelayProfile(&plan)

	// Create new DelayProfile
	response, err := r.provider.client.AddDelayProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create delayprofile, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "created delayprofile: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := *writeDelayProfile(response)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceDelayProfile) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
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
	result := *writeDelayProfile(response)

	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceDelayProfile) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// Get plan values
	var plan DelayProfile
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get state values
	var state DelayProfile
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	plan.ID.Value = state.ID.Value
	data := readDelayProfile(&plan)

	// Update DelayProfile
	response, err := r.provider.client.UpdateDelayProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update delayprofile, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "update delayprofile: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := *writeDelayProfile(response)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceDelayProfile) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
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

func (r resourceDelayProfile) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
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

func writeDelayProfile(profile *sonarr.DelayProfile) *DelayProfile {
	tags := make([]types.Int64, len(profile.Tags))

	for i, t := range profile.Tags {
		tags[i] = types.Int64{Value: int64(t)}
	}
	return &DelayProfile{
		ID:                     types.Int64{Value: profile.ID},
		EnableUsenet:           types.Bool{Value: profile.EnableUsenet},
		EnableTorrent:          types.Bool{Value: profile.EnableTorrent},
		BypassIfHighestQuality: types.Bool{Value: profile.BypassIfHighestQuality},
		UsenetDelay:            types.Int64{Value: profile.UsenetDelay},
		TorrentDelay:           types.Int64{Value: profile.TorrentDelay},
		Order:                  types.Int64{Value: profile.Order},
		PreferredProtocol:      types.String{Value: profile.PreferredProtocol},
		Tags:                   tags,
	}
}

func readDelayProfile(profile *DelayProfile) *sonarr.DelayProfile {
	tags := make([]int, len(profile.Tags))

	for i, t := range profile.Tags {
		tags[i] = int(t.Value)
	}

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
