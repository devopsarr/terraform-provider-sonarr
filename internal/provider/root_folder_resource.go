package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RootFolderResource{}
var _ resource.ResourceWithImportState = &RootFolderResource{}

func NewRootFolderResource() resource.Resource {
	return &RootFolderResource{}
}

// RootFolderResource defines the root folder implementation.
type RootFolderResource struct {
	client *sonarr.Sonarr
}

// RootFolder describes the root folder data model.
type RootFolder struct {
	Accessible      types.Bool   `tfsdk:"accessible"`
	ID              types.Int64  `tfsdk:"id"`
	Path            types.String `tfsdk:"path"`
	UnmappedFolders types.Set    `tfsdk:"unmapped_folders"`
}

// Path part of RootFolder.
type Path struct {
	Name types.String `tfsdk:"name"`
	Path types.String `tfsdk:"path"`
}

func (r *RootFolderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_root_folder"
}

func (r *RootFolderResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Root Folder resource.<br/>For more information refer to [Root Folders](https://wiki.servarr.com/sonarr/settings#root-folders) documentation.",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: add validator
			"path": {
				MarkdownDescription: "Root Folder absolute path.",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
			},
			"accessible": {
				MarkdownDescription: "Access flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"id": {
				MarkdownDescription: "Root Folder ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"unmapped_folders": {
				MarkdownDescription: "List of folders with no associated series.",
				Computed:            true,
				Attributes:          tfsdk.SetNestedAttributes(r.getUnmappedFolderSchema().Attributes),
			},
		},
	}, nil
}

func (r RootFolderResource) getUnmappedFolderSchema() tfsdk.Schema {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"path": {
				MarkdownDescription: "Path of unmapped folder.",
				Computed:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "Name of unmapped folder.",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}
}

func (r *RootFolderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *RootFolderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan string

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("path"), &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new RootFolder
	request := sonarr.RootFolder{
		Path: plan,
	}

	response, err := r.client.AddRootFolderContext(ctx, &request)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to create rootFolder, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "created root_folder: "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeRootFolder(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

func (r *RootFolderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state RootFolder

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get rootFolder current value
	response, err := r.client.GetRootFolderContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read rootFolders, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "read root_folder: "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	result := writeRootFolder(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, result)...)
}

// never used.
func (r *RootFolderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *RootFolderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RootFolder

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete rootFolder current value
	err := r.client.DeleteRootFolderContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read rootFolders, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "deleted root_folder: "+strconv.Itoa(int(state.ID.Value)))
	resp.State.RemoveResource(ctx)
}

func (r *RootFolderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported root_folder: "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func writeRootFolder(ctx context.Context, rootFolder *sonarr.RootFolder) *RootFolder {
	output := RootFolder{
		Accessible:      types.Bool{Value: rootFolder.Accessible},
		ID:              types.Int64{Value: rootFolder.ID},
		Path:            types.String{Value: rootFolder.Path},
		UnmappedFolders: types.Set{ElemType: RootFolderResource{}.getUnmappedFolderSchema().Type()},
	}
	unmapped := writeUnmappedFolders(rootFolder.UnmappedFolders)

	tfsdk.ValueFrom(ctx, unmapped, output.UnmappedFolders.Type(ctx), output.UnmappedFolders)

	return &output
}

func writeUnmappedFolders(folders []*starr.Path) *[]Path {
	output := make([]Path, len(folders))
	for i, f := range folders {
		output[i] = Path{
			Name: types.String{Value: f.Name},
			Path: types.String{Value: f.Path},
		}
	}

	return &output
}
