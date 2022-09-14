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
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RemotePathMappingResource{}
var _ resource.ResourceWithImportState = &RemotePathMappingResource{}

func NewRemotePathMappingResource() resource.Resource {
	return &RemotePathMappingResource{}
}

// RemotePathMappingResource defines the remote path mapping implementation.
type RemotePathMappingResource struct {
	client *sonarr.Sonarr
}

// RemotePathMapping describes the remote path mapping data model.
type RemotePathMapping struct {
	ID         types.Int64  `tfsdk:"id"`
	Host       types.String `tfsdk:"host"`
	RemotePath types.String `tfsdk:"remote_path"`
	LocalPath  types.String `tfsdk:"local_path"`
}

func (r *RemotePathMappingResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Remote Path Mapping resource.<br/>For more information refer to [Remote Path Mapping](https://wiki.servarr.com/sonarr/settings#remote-path-mappings) documentation.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Remote Path Mapping ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"host": {
				MarkdownDescription: "Download Client host.",
				Required:            true,
				Type:                types.StringType,
			},
			"remote_path": {
				MarkdownDescription: "Download Client remote path.",
				Required:            true,
				Type:                types.StringType,
			},
			"local_path": {
				MarkdownDescription: "Local path.",
				Required:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (r *RemotePathMappingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_path_mapping"
}

func (r *RemotePathMappingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RemotePathMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan RemotePathMapping

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new RemotePathMapping
	request := readRemotePathMapping(&plan)

	response, err := r.client.AddRemotePathMappingContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to create remotePathMapping, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "created remote_path_mapping: "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeRemotePathMapping(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *RemotePathMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state RemotePathMapping

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get remotePathMapping current value
	response, err := r.client.GetRemotePathMappingContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read remotePathMappings, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "read remote_path_mapping: "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	result := writeRemotePathMapping(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *RemotePathMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan RemotePathMapping

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update RemotePathMapping
	request := readRemotePathMapping(&plan)

	response, err := r.client.UpdateRemotePathMappingContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to update remotePathMapping, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated remote_path_mapping: "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	result := writeRemotePathMapping(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *RemotePathMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RemotePathMapping

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete remotePathMapping current value
	err := r.client.DeleteRemotePathMappingContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read remotePathMappings, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "deleted remote_path_mapping: "+strconv.Itoa(int(state.ID.Value)))
	resp.State.RemoveResource(ctx)
}

func (r *RemotePathMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported remote_path_mapping: "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func writeRemotePathMapping(remotePathMapping *sonarr.RemotePathMapping) *RemotePathMapping {
	return &RemotePathMapping{
		ID:         types.Int64{Value: remotePathMapping.ID},
		Host:       types.String{Value: remotePathMapping.Host},
		RemotePath: types.String{Value: remotePathMapping.RemotePath},
		LocalPath:  types.String{Value: remotePathMapping.LocalPath},
	}
}

func readRemotePathMapping(mapping *RemotePathMapping) *sonarr.RemotePathMapping {
	return &sonarr.RemotePathMapping{
		ID:         mapping.ID.Value,
		Host:       mapping.Host.Value,
		RemotePath: mapping.RemotePath.Value,
		LocalPath:  mapping.LocalPath.Value,
	}
}
