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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const importListExclusionResourceName = "import_list_exclusion"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListExclusionResource{}
	_ resource.ResourceWithImportState = &ImportListExclusionResource{}
)

func NewImportListExclusionResource() resource.Resource {
	return &ImportListExclusionResource{}
}

// ImportListExclusionResource defines the importListExclusion implementation.
type ImportListExclusionResource struct {
	client *sonarr.APIClient
}

// ImportListExclusion describes the importListExclusion data model.
type ImportListExclusion struct {
	Title  types.String `tfsdk:"title"`
	TVDBID types.Int64  `tfsdk:"tvdb_id"`
	ID     types.Int64  `tfsdk:"id"`
}

func (r *ImportListExclusionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListExclusionResourceName
}

func (r *ImportListExclusionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->Import List Exclusion resource.\nFor more information refer to [ImportListExclusions](https://wiki.servarr.com/sonarr/settings#list-exclusions) documentation.",
		Attributes: map[string]schema.Attribute{
			"tvdb_id": schema.Int64Attribute{
				MarkdownDescription: "Series TVDB ID.",
				Required:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Series to be excluded.",
				Required:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "ImportListExclusion ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ImportListExclusionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ImportListExclusionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importListExclusion *ImportListExclusion

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importListExclusion)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportListExclusion
	request := importListExclusion.read()

	response, _, err := r.client.ImportListExclusionApi.CreateImportListExclusion(ctx).ImportListExclusionResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListExclusionResourceName, err))

		return
	}

	tflog.Trace(ctx, "created importListExclusion: "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importListExclusion.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importListExclusion)...)
}

func (r *ImportListExclusionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importListExclusion *ImportListExclusion

	resp.Diagnostics.Append(req.State.Get(ctx, &importListExclusion)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get importListExclusion current value
	response, _, err := r.client.ImportListExclusionApi.GetImportListExclusionById(ctx, int32(importListExclusion.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListExclusionResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListExclusionResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	importListExclusion.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importListExclusion)...)
}

func (r *ImportListExclusionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importListExclusion *ImportListExclusion

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importListExclusion)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportListExclusion
	request := importListExclusion.read()

	response, _, err := r.client.ImportListExclusionApi.UpdateImportListExclusion(ctx, strconv.Itoa(int(request.GetId()))).ImportListExclusionResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListExclusionResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListExclusionResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importListExclusion.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importListExclusion)...)
}

func (r *ImportListExclusionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var importListExclusion *ImportListExclusion

	resp.Diagnostics.Append(req.State.Get(ctx, &importListExclusion)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete importListExclusion current value
	_, err := r.client.ImportListExclusionApi.DeleteImportListExclusion(ctx, int32(importListExclusion.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListExclusionResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListExclusionResourceName+": "+strconv.Itoa(int(importListExclusion.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListExclusionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			helpers.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+importListExclusionResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (i *ImportListExclusion) write(importListExclusion *sonarr.ImportListExclusionResource) {
	i.ID = types.Int64Value(int64(importListExclusion.GetId()))
	i.TVDBID = types.Int64Value(int64(importListExclusion.GetTvdbId()))
	i.Title = types.StringValue(importListExclusion.GetTitle())
}

func (i *ImportListExclusion) read() *sonarr.ImportListExclusionResource {
	exclusion := sonarr.NewImportListExclusionResource()
	exclusion.SetId(int32(i.ID.ValueInt64()))
	exclusion.SetTitle(i.Title.ValueString())
	exclusion.SetTvdbId(int32(i.TVDBID.ValueInt64()))

	return exclusion
}
