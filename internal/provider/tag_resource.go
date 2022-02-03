package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type resourceTagType struct{}

type resourceTag struct {
	provider provider
}

func (t resourceTagType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Tag resource",

		Attributes: map[string]tfsdk.Attribute{
			"label": {
				MarkdownDescription: "Tag value",
				Required:            true,
				Type:                types.StringType,
			},
			"id": {
				MarkdownDescription: "Tag ID",
				Computed:            true,
				Type:                types.Int64Type,
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
	if !r.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan Tag

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Tag
	id, err := r.provider.client.AddTag(plan.Label.Value)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create tag, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "created tag", "tag_id", id)

	// Generate resource state struct
	var result = Tag{
		ID:    types.Int64{Value: int64(id)},
		Label: types.String{Value: plan.Label.Value},
	}

	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)

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

	// Get order current value
	tags, err := r.provider.client.GetTags()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read tags, got error: %s", err))
		return
	}
	// Map response body to resource schema attribute
	for _, t := range tags {
		if t.ID == int(state.ID.Value) {
			state.Label = types.String{Value: t.Label}
		}
	}

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
	id, err := r.provider.client.UpdateTag(int(state.ID.Value), plan.Label.Value)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update tag, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "update tag", "tag_id", id)

	// Generate resource state struct
	var result = Tag{
		ID:    types.Int64{Value: int64(id)},
		Label: types.String{Value: plan.Label.Value},
	}

	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceTag) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data Tag

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.DeleteExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }

	resp.State.RemoveResource(ctx)
}

func (r resourceTag) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	//TODO implement import
	//	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
	tfsdk.ResourceImportStateNotImplemented(ctx, "", resp)
}
