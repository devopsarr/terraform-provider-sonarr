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

type resourceRootFolderType struct{}

type resourceRootFolder struct {
	provider provider
}

func (t resourceRootFolderType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "RootFolder resource",
		Attributes: map[string]tfsdk.Attribute{
			//TODO: add validator
			"path": {
				MarkdownDescription: "Absolute path of rootFolder",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"accessible": {
				MarkdownDescription: "Access flag",
				Computed:            true,
				Type:                types.BoolType,
			},
			"id": {
				MarkdownDescription: "RootFolder ID",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"unmapped_folders": {
				MarkdownDescription: "List of folders with no associated series",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"path": {
						MarkdownDescription: "Path of unmapped folder",
						Computed:            true,
						Type:                types.StringType,
					},
					"name": {
						MarkdownDescription: "Name of unmapped folder",
						Computed:            true,
						Type:                types.StringType,
					},
				}),
			},
		},
	}, nil
}

func (t resourceRootFolderType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return resourceRootFolder{
		provider: provider,
	}, diags
}

func (r resourceRootFolder) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	// Retrieve values from plan
	var plan string
	diags := req.Plan.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("path"), &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new RootFolder
	request := sonarr.RootFolder{
		Path: plan,
	}
	response, err := r.provider.client.AddRootFolderContext(ctx, &request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create rootFolder, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "created rootFolder: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := writeRootFolder(response)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceRootFolder) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// Get current state
	var state RootFolder
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get rootFolder current value
	response, err := r.provider.client.GetRootFolderContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read rootFolders, got error: %s", err))
		return
	}
	// Map response body to resource schema attribute
	result := writeRootFolder(response)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

// never used
func (r resourceRootFolder) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
}

func (r resourceRootFolder) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state RootFolder

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete rootFolder current value
	err := r.provider.client.DeleteRootFolderContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read rootFolders, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceRootFolder) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
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

func writeRootFolder(rootFolder *sonarr.RootFolder) *RootFolder {
	unmapped := make([]Path, len(rootFolder.UnmappedFolders))
	for i, f := range rootFolder.UnmappedFolders {
		unmapped[i] = Path{
			Name: types.String{Value: f.Name},
			Path: types.String{Value: f.Path},
		}
	}
	return &RootFolder{
		Accessible:      types.Bool{Value: rootFolder.Accessible},
		ID:              types.Int64{Value: rootFolder.ID},
		Path:            types.String{Value: rootFolder.Path},
		UnmappedFolders: unmapped,
	}
}
