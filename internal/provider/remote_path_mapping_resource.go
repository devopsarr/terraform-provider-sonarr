package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const remotePathMappingResourceName = "remote_path_mapping"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &RemotePathMappingResource{}
	_ resource.ResourceWithImportState = &RemotePathMappingResource{}
)

func NewRemotePathMappingResource() resource.Resource {
	return &RemotePathMappingResource{}
}

// RemotePathMappingResource defines the remote path mapping implementation.
type RemotePathMappingResource struct {
	client *sonarr.APIClient
}

// RemotePathMapping describes the remote path mapping data model.
type RemotePathMapping struct {
	Host       types.String `tfsdk:"host"`
	RemotePath types.String `tfsdk:"remote_path"`
	LocalPath  types.String `tfsdk:"local_path"`
	ID         types.Int64  `tfsdk:"id"`
}

func (r *RemotePathMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Remote Path Mapping resource.\nFor more information refer to [Remote Path Mapping](https://wiki.servarr.com/sonarr/settings#remote-path-mappings) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Remote Path Mapping ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "Download Client host.",
				Required:            true,
			},
			"remote_path": schema.StringAttribute{
				MarkdownDescription: "Download Client remote path.",
				Required:            true,
			},
			"local_path": schema.StringAttribute{
				MarkdownDescription: "Local path.",
				Required:            true,
			},
		},
	}
}

func (r *RemotePathMappingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + remotePathMappingResourceName
}

func (r *RemotePathMappingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *sonarr.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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

	response, _, err := r.client.RemotePathMappingApi.CreateRemotepathmapping(ctx).RemotePathMappingResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", remotePathMappingResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+remotePathMappingResourceName+": "+strconv.Itoa(int(response.GetId())))
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
	response, _, err := r.client.RemotePathMappingApi.GetRemotepathmappingById(ctx, int32(mapping.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", remotePathMappingResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+remotePathMappingResourceName+": "+strconv.Itoa(int(response.GetId())))
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

	response, _, err := r.client.RemotePathMappingApi.UpdateRemotepathmapping(ctx, strconv.Itoa(int(request.GetId()))).RemotePathMappingResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", remotePathMappingResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+remotePathMappingResourceName+": "+strconv.Itoa(int(response.GetId())))
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
	_, err := r.client.RemotePathMappingApi.DeleteRemotepathmapping(ctx, int32(mapping.ID.ValueInt64())).Execute()
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

func (r *RemotePathMapping) write(remotePathMapping *sonarr.RemotePathMappingResource) {
	r.ID = types.Int64Value(int64(remotePathMapping.GetId()))
	r.Host = types.StringValue(remotePathMapping.GetHost())
	r.RemotePath = types.StringValue(remotePathMapping.GetRemotePath())
	r.LocalPath = types.StringValue(remotePathMapping.GetLocalPath())
}

func (r *RemotePathMapping) read() *sonarr.RemotePathMappingResource {
	mapping := sonarr.NewRemotePathMappingResource()
	mapping.SetHost(r.Host.ValueString())
	mapping.SetLocalPath(r.LocalPath.ValueString())
	mapping.SetRemotePath(r.RemotePath.ValueString())
	mapping.SetId(int32(r.ID.ValueInt64()))

	return mapping
}
