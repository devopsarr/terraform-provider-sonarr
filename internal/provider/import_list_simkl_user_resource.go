package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	importListSimklUserResourceName   = "import_list_simkl_user"
	importListSimklUserImplementation = "SimklUserImport"
	importListSimklUserConfigContract = "SimklUserSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListSimklUserResource{}
	_ resource.ResourceWithImportState = &ImportListSimklUserResource{}
)

func NewImportListSimklUserResource() resource.Resource {
	return &ImportListSimklUserResource{}
}

// ImportListSimklUserResource defines the import list implementation.
type ImportListSimklUserResource struct {
	client *sonarr.APIClient
}

// ImportListSimklUser describes the import list data model.
type ImportListSimklUser struct {
	Tags               types.Set    `tfsdk:"tags"`
	Name               types.String `tfsdk:"name"`
	ShouldMonitor      types.String `tfsdk:"should_monitor"`
	RootFolderPath     types.String `tfsdk:"root_folder_path"`
	SeriesType         types.String `tfsdk:"series_type"`
	AccessToken        types.String `tfsdk:"access_token"`
	RefreshToken       types.String `tfsdk:"refresh_token"`
	Expires            types.String `tfsdk:"expires"`
	AuthUser           types.String `tfsdk:"auth_user"`
	QualityProfileID   types.Int64  `tfsdk:"quality_profile_id"`
	ID                 types.Int64  `tfsdk:"id"`
	ListType           types.Int64  `tfsdk:"list_type"`
	EnableAutomaticAdd types.Bool   `tfsdk:"enable_automatic_add"`
	SeasonFolder       types.Bool   `tfsdk:"season_folder"`
}

func (i ImportListSimklUser) toImportList() *ImportList {
	return &ImportList{
		Tags:               i.Tags,
		Name:               i.Name,
		ShouldMonitor:      i.ShouldMonitor,
		RootFolderPath:     i.RootFolderPath,
		SeriesType:         i.SeriesType,
		AccessToken:        i.AccessToken,
		RefreshToken:       i.RefreshToken,
		Expires:            i.Expires,
		AuthUser:           i.AuthUser,
		ListType:           i.ListType,
		QualityProfileID:   i.QualityProfileID,
		ID:                 i.ID,
		EnableAutomaticAdd: i.EnableAutomaticAdd,
		SeasonFolder:       i.SeasonFolder,
		ConfigContract:     types.StringValue(importListSimklUserConfigContract),
		Implementation:     types.StringValue(importListSimklUserImplementation),
	}
}

func (i *ImportListSimklUser) fromImportList(importList *ImportList) {
	i.Tags = importList.Tags
	i.Name = importList.Name
	i.ShouldMonitor = importList.ShouldMonitor
	i.RootFolderPath = importList.RootFolderPath
	i.SeriesType = importList.SeriesType
	i.AccessToken = importList.AccessToken
	i.RefreshToken = importList.RefreshToken
	i.Expires = importList.Expires
	i.AuthUser = importList.AuthUser
	i.ListType = importList.ListType
	i.QualityProfileID = importList.QualityProfileID
	i.ID = importList.ID
	i.EnableAutomaticAdd = importList.EnableAutomaticAdd
	i.SeasonFolder = importList.SeasonFolder
}

func (r *ImportListSimklUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListSimklUserResourceName
}

func (r *ImportListSimklUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->ImportList SimklUser resource.\nFor more information refer to [Import List](https://wiki.servarr.com/sonarr/settings#import-lists) and [SimklUser](https://wiki.servarr.com/sonarr/supported#simkl_user).",
		Attributes: map[string]schema.Attribute{
			"enable_automatic_add": schema.BoolAttribute{
				MarkdownDescription: "Enable automatic add flag.",
				Required:            true,
			},
			"season_folder": schema.BoolAttribute{
				MarkdownDescription: "Season folder flag.",
				Required:            true,
			},
			"quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Quality profile ID.",
				Required:            true,
			},
			"should_monitor": schema.StringAttribute{
				MarkdownDescription: "Should monitor.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("all", "future", "missing", "existing", "pilot", "firstSeason", "latestSeason", "none"),
				},
			},
			"root_folder_path": schema.StringAttribute{
				MarkdownDescription: "Root folder path.",
				Required:            true,
			},
			"series_type": schema.StringAttribute{
				MarkdownDescription: "Series type.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("standard", "anime", "daily"),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Import List name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Import List ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"list_type": schema.Int64Attribute{
				MarkdownDescription: "Simkl list type. '0' Watching, '1' PlanToWatch, '2' Hold, '1' Completed, '2' Dropped.",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1, 2, 3, 4),
				},
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access token.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"refresh_token": schema.StringAttribute{
				MarkdownDescription: "Refresh token.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"auth_user": schema.StringAttribute{
				MarkdownDescription: "Auth User.",
				Optional:            true,
				Computed:            true,
			},
			"expires": schema.StringAttribute{
				MarkdownDescription: "Expires.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *ImportListSimklUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ImportListSimklUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportListSimklUser

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportListSimklUser
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.CreateImportList(ctx).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListSimklUserResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListSimklUserResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListSimklUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportListSimklUser

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportListSimklUser current value
	response, _, err := r.client.ImportListApi.GetImportListById(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListSimklUserResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListSimklUserResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListSimklUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportListSimklUser

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportListSimklUser
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.UpdateImportList(ctx, strconv.Itoa(int(request.GetId()))).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListSimklUserResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListSimklUserResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListSimklUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportListSimklUser current value
	_, err := r.client.ImportListApi.DeleteImportList(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, importListSimklUserResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListSimklUserResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListSimklUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+importListSimklUserResourceName+": "+req.ID)
}

func (i *ImportListSimklUser) write(ctx context.Context, importList *sonarr.ImportListResource) {
	genericImportList := i.toImportList()
	genericImportList.write(ctx, importList)
	i.fromImportList(genericImportList)
}

func (i *ImportListSimklUser) read(ctx context.Context) *sonarr.ImportListResource {
	return i.toImportList().read(ctx)
}
