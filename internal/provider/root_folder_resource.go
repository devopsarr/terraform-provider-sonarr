package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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

func (r RootFolder) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"unmapped_folders": types.SetType{}.WithElementType(Path{}.getType()),
			"path":             types.StringType,
			"id":               types.Int64Type,
			"accessible":       types.BoolType,
		})
}

// Path part of RootFolder.
type Path struct {
	Name types.String `tfsdk:"name"`
	Path types.String `tfsdk:"path"`
}

func (p Path) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"name": types.StringType,
			"path": types.StringType,
		})
}

func (r *RootFolderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + rootFolderResourceName
}

func (r *RootFolderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
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
	folder.write(ctx, response, &resp.Diagnostics)
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
	folder.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &folder)...)
}

// never used.
func (r *RootFolderResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *RootFolderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete rootFolder current value
	_, err := r.client.RootFolderApi.DeleteRootFolder(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, rootFolderResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+rootFolderResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *RootFolderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+rootFolderResourceName+": "+req.ID)
}

func (r *RootFolder) write(ctx context.Context, rootFolder *sonarr.RootFolderResource, diags *diag.Diagnostics) {
	var tempDiag diag.Diagnostics

	r.Accessible = types.BoolValue(rootFolder.GetAccessible())
	r.ID = types.Int64Value(int64(rootFolder.GetId()))
	r.Path = types.StringValue(rootFolder.GetPath())

	unmapped := make([]Path, len(rootFolder.GetUnmappedFolders()))
	for i, f := range rootFolder.UnmappedFolders {
		unmapped[i].write(f)
	}

	r.UnmappedFolders, tempDiag = types.SetValueFrom(ctx, Path{}.getType(), unmapped)
	diags.Append(tempDiag...)
}

func (p *Path) write(folder *sonarr.UnmappedFolder) {
	p.Name = types.StringValue(folder.GetName())
	p.Path = types.StringValue(folder.GetPath())
}
