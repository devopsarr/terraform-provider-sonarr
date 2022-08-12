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
	"golift.io/starr"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.ResourceType = resourceTagType{}
var _ resource.Resource = resourceTag{}
var _ resource.ResourceWithImportState = resourceTag{}

type resourceTagType struct{}

type resourceTag struct {
	provider sonarrProvider
}

func (t resourceTagType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Tag resource",
		Attributes: map[string]tfsdk.Attribute{
			"label": {
				MarkdownDescription: "Tag value",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					helpers.StringLowercase(),
				},
			},
			"id": {
				MarkdownDescription: "Tag ID",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
		},
	}, nil
}

func (t resourceTagType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return resourceTag{
		provider: provider,
	}, diags
}

func (r resourceTag) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan Tag
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Tag
	request := starr.Tag{
		Label: plan.Label.Value,
	}
	response, err := r.provider.client.AddTagContext(ctx, &request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create tag, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "created tag: "+strconv.Itoa(response.ID))

	// Generate resource state struct
	var result = Tag{
		ID:    types.Int64{Value: int64(response.ID)},
		Label: types.String{Value: response.Label},
	}

	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceTag) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state Tag
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get tag current value
	response, err := r.provider.client.GetTagContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read tags, got error: %s", err))
		return
	}
	// Map response body to resource schema attribute
	result := Tag{
		ID:    types.Int64{Value: int64(response.ID)},
		Label: types.String{Value: response.Label},
	}

	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceTag) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan Tag
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update Tag
	request := starr.Tag{
		Label: plan.Label.Value,
		ID:    int(plan.ID.Value),
	}
	response, err := r.provider.client.UpdateTagContext(ctx, &request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update tag, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "update tag: "+strconv.Itoa(response.ID))

	// Generate resource state struct
	var result = Tag{
		ID:    types.Int64{Value: int64(response.ID)},
		Label: types.String{Value: response.Label},
	}

	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceTag) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Tag

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete tag current value
	err := r.provider.client.DeleteTagContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read tags, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceTag) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	//resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
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
