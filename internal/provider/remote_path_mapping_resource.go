package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const remotePathMappingResourceName = "remote_path_mapping"

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
	Host       types.String `tfsdk:"host"`
	RemotePath types.String `tfsdk:"remote_path"`
	LocalPath  types.String `tfsdk:"local_path"`
	ID         types.Int64  `tfsdk:"id"`
}

func (r *RemotePathMappingResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Remote Path Mapping resource.\nFor more information refer to [Remote Path Mapping](https://wiki.servarr.com/sonarr/settings#remote-path-mappings) documentation.",
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
	resp.TypeName = req.ProviderTypeName + "_" + remotePathMappingResourceName
}

func (r *RemotePathMappingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *RemotePathMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var mapping *RemotePathMapping

	resp.Diagnostics.Append(req.Plan.Get(ctx, &mapping)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new RemotePathMapping
	request := mapping.read()

	response, err := r.client.AddRemotePathMappingContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", remotePathMappingResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+remotePathMappingResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	mapping.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &mapping)...)
}

func (r *RemotePathMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var mapping *RemotePathMapping

	resp.Diagnostics.Append(req.State.Get(ctx, &mapping)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get remotePathMapping current value
	response, err := r.client.GetRemotePathMappingContext(ctx, mapping.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", remotePathMappingResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+remotePathMappingResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	mapping.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &mapping)...)
}

func (r *RemotePathMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var mapping *RemotePathMapping

	resp.Diagnostics.Append(req.Plan.Get(ctx, &mapping)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update RemotePathMapping
	request := mapping.read()

	response, err := r.client.UpdateRemotePathMappingContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", remotePathMappingResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+remotePathMappingResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	mapping.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &mapping)...)
}

func (r *RemotePathMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var mapping *RemotePathMapping

	diags := req.State.Get(ctx, &mapping)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete remotePathMapping current value
	err := r.client.DeleteRemotePathMappingContext(ctx, mapping.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", remotePathMappingResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+remotePathMappingResourceName+": "+strconv.Itoa(int(mapping.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *RemotePathMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+remotePathMappingResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (r *RemotePathMapping) write(remotePathMapping *sonarr.RemotePathMapping) {
	r.ID = types.Int64Value(remotePathMapping.ID)
	r.Host = types.StringValue(remotePathMapping.Host)
	r.RemotePath = types.StringValue(remotePathMapping.RemotePath)
	r.LocalPath = types.StringValue(remotePathMapping.LocalPath)
}

func (r *RemotePathMapping) read() *sonarr.RemotePathMapping {
	return &sonarr.RemotePathMapping{
		ID:         r.ID.ValueInt64(),
		Host:       r.Host.ValueString(),
		RemotePath: r.RemotePath.ValueString(),
		LocalPath:  r.LocalPath.ValueString(),
	}
}
