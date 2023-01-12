package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const rootFolderResourceName = "root_folder"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &RootFolderResource{}
	_ resource.ResourceWithImportState = &RootFolderResource{}
)

func NewRootFolderResource() resource.Resource {
	return &RootFolderResource{}
}

// RootFolderResource defines the root folder implementation.
type RootFolderResource struct {
	client *sonarr.APIClient
}

// RootFolder describes the root folder data model.
type RootFolder struct {
	UnmappedFolders types.Set    `tfsdk:"unmapped_folders"`
	Path            types.String `tfsdk:"path"`
	ID              types.Int64  `tfsdk:"id"`
	Accessible      types.Bool   `tfsdk:"accessible"`
}

// Path part of RootFolder.
type Path struct {
	Name types.String `tfsdk:"name"`
	Path types.String `tfsdk:"path"`
}

func (r *RootFolderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + rootFolderResourceName
}

func (r *RootFolderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Media Management -->Root Folder resource.\nFor more information refer to [Root Folders](https://wiki.servarr.com/sonarr/settings#root-folders) documentation.",
		Attributes: map[string]schema.Attribute{
			// TODO: add validator
			"path": schema.StringAttribute{
				MarkdownDescription: "Root Folder absolute path.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"accessible": schema.BoolAttribute{
				MarkdownDescription: "Access flag.",
				Computed:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Root Folder ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"unmapped_folders": schema.SetNestedAttribute{
				MarkdownDescription: "List of folders with no associated series.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: r.getUnmappedFolderSchema().Attributes,
				},
			},
		},
	}
}

func (r RootFolderResource) getUnmappedFolderSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				MarkdownDescription: "Path of unmapped folder.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of unmapped folder.",
				Computed:            true,
			},
		},
	}
}

func (r *RootFolderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			helpers.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *sonarr.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *RootFolderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var folder *RootFolder

	resp.Diagnostics.Append(req.Plan.Get(ctx, &folder)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new RootFolder
	request := *sonarr.NewRootFolderResource()
	request.SetPath(folder.Path.ValueString())

	response, _, err := r.client.RootFolderApi.CreateRootFolder(ctx).RootFolderResource(request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, rootFolderResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+rootFolderResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	folder.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &folder)...)
}

func (r *RootFolderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var folder *RootFolder

	resp.Diagnostics.Append(req.State.Get(ctx, &folder)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get rootFolder current value
	response, _, err := r.client.RootFolderApi.GetRootFolderById(ctx, int32(folder.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, rootFolderResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+rootFolderResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	folder.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &folder)...)
}

// never used.
func (r *RootFolderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *RootFolderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var folder *RootFolder

	resp.Diagnostics.Append(req.State.Get(ctx, &folder)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete rootFolder current value
	_, err := r.client.RootFolderApi.DeleteRootFolder(ctx, int32(folder.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, rootFolderResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+rootFolderResourceName+": "+strconv.Itoa(int(folder.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *RootFolderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			helpers.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+rootFolderResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (r *RootFolder) write(ctx context.Context, rootFolder *sonarr.RootFolderResource) {
	r.Accessible = types.BoolValue(rootFolder.GetAccessible())
	r.ID = types.Int64Value(int64(rootFolder.GetId()))
	r.Path = types.StringValue(rootFolder.GetPath())
	r.UnmappedFolders = types.SetValueMust(RootFolderResource{}.getUnmappedFolderSchema().Type(), nil)

	unmapped := make([]Path, len(rootFolder.GetUnmappedFolders()))
	for i, f := range rootFolder.UnmappedFolders {
		unmapped[i].write(f)
	}

	tfsdk.ValueFrom(ctx, unmapped, r.UnmappedFolders.Type(ctx), r.UnmappedFolders)
}

func (p *Path) write(folder *sonarr.UnmappedFolder) {
	p.Name = types.StringValue(folder.GetName())
	p.Path = types.StringValue(folder.GetPath())
}
