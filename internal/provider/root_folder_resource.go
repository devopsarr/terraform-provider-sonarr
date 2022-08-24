package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.ResourceType = resourceRootFolderType{}
var _ resource.Resource = resourceRootFolder{}
var _ resource.ResourceWithImportState = resourceRootFolder{}

type resourceRootFolderType struct{}

type resourceRootFolder struct {
	provider sonarrProvider
}

// RootFolder is the RootFolder resource.
type RootFolder struct {
	Accessible      types.Bool   `tfsdk:"accessible"`
	ID              types.Int64  `tfsdk:"id"`
	Path            types.String `tfsdk:"path"`
	UnmappedFolders []Path       `tfsdk:"unmapped_folders"`
}

// Path part of RootFolder.
type Path struct {
	Name types.String `tfsdk:"name"`
	Path types.String `tfsdk:"path"`
}

func (t resourceRootFolderType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "RootFolder resource",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: add validator
			"path": {
				MarkdownDescription: "Absolute path of rootFolder",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
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
					resource.UseStateForUnknown(),
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

func (t resourceRootFolderType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return resourceRootFolder{
		provider: provider,
	}, diags
}

func (r resourceRootFolder) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan string
	diags := req.Plan.GetAttribute(ctx, path.Root("path"), &plan)
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

func (r resourceRootFolder) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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
func (r resourceRootFolder) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r resourceRootFolder) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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

func (r resourceRootFolder) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
