package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr"
)

type resourceTagType struct{}

type resourceTag struct {
	provider provider
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
					stringLowercase(),
				},
			},
			"id": {
				MarkdownDescription: "Tag ID",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
		},
	}, nil
}

func (t resourceTagType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return resourceTag{
		provider: provider,
	}, diags
}

func (r resourceTag) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
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

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceTag) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
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
	state.Label = types.String{Value: response.Label}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r resourceTag) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// Get plan values
	var plan Tag
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get state values
	var state Tag
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update Tag
	request := starr.Tag{
		Label: plan.Label.Value,
		ID:    int(state.ID.Value),
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

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceTag) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
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

func (r resourceTag) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
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

//TODO: move into validators file
type stringLowercaseValidator struct {
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v stringLowercaseValidator) Description(ctx context.Context) string {
	return "string must be lowercase"
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v stringLowercaseValidator) MarkdownDescription(ctx context.Context) string {
	return "string must be lowercase"
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v stringLowercaseValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var str types.String
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &str)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	if str.Unknown || str.Null {
		return
	}
	upper, _ := regexp.Match(`^.*[A-Z]+.*$`, []byte(str.Value))
	if upper {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid String Content",
			"String cannot contains uppercase values",
		)
		return
	}
}

func stringLowercase() stringLowercaseValidator {
	return stringLowercaseValidator{}
}
