package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	importListSonarrResourceName   = "import_list_sonarr"
	importListSonarrImplementation = "SonarrImport"
	importListSonarrConfigContract = "SonarrSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListSonarrResource{}
	_ resource.ResourceWithImportState = &ImportListSonarrResource{}
)

func NewImportListSonarrResource() resource.Resource {
	return &ImportListSonarrResource{}
}

// ImportListSonarrResource defines the import list implementation.
type ImportListSonarrResource struct {
	client *sonarr.APIClient
}

// ImportListSonarr describes the import list data model.
type ImportListSonarr struct {
	Tags               types.Set    `tfsdk:"tags"`
	LanguageProfileIds types.Set    `tfsdk:"language_profile_ids"`
	ProfileIds         types.Set    `tfsdk:"quality_profile_ids"`
	TagIds             types.Set    `tfsdk:"tag_ids"`
	Name               types.String `tfsdk:"name"`
	ShouldMonitor      types.String `tfsdk:"should_monitor"`
	RootFolderPath     types.String `tfsdk:"root_folder_path"`
	SeriesType         types.String `tfsdk:"series_type"`
	BaseURL            types.String `tfsdk:"base_url"`
	APIKey             types.String `tfsdk:"api_key"`
	QualityProfileID   types.Int64  `tfsdk:"quality_profile_id"`
	ID                 types.Int64  `tfsdk:"id"`
	EnableAutomaticAdd types.Bool   `tfsdk:"enable_automatic_add"`
	SeasonFolder       types.Bool   `tfsdk:"season_folder"`
}

func (i ImportListSonarr) toImportList() *ImportList {
	return &ImportList{
		Tags:               i.Tags,
		LanguageProfileIds: i.LanguageProfileIds,
		ProfileIds:         i.ProfileIds,
		TagIds:             i.TagIds,
		Name:               i.Name,
		ShouldMonitor:      i.ShouldMonitor,
		RootFolderPath:     i.RootFolderPath,
		SeriesType:         i.SeriesType,
		BaseURL:            i.BaseURL,
		APIKey:             i.APIKey,
		QualityProfileID:   i.QualityProfileID,
		ID:                 i.ID,
		EnableAutomaticAdd: i.EnableAutomaticAdd,
		SeasonFolder:       i.SeasonFolder,
		ConfigContract:     types.StringValue(importListSonarrConfigContract),
		Implementation:     types.StringValue(importListSonarrImplementation),
	}
}

func (i *ImportListSonarr) fromImportList(importList *ImportList) {
	i.Tags = importList.Tags
	i.LanguageProfileIds = importList.LanguageProfileIds
	i.ProfileIds = importList.ProfileIds
	i.TagIds = importList.TagIds
	i.Name = importList.Name
	i.ShouldMonitor = importList.ShouldMonitor
	i.RootFolderPath = importList.RootFolderPath
	i.SeriesType = importList.SeriesType
	i.BaseURL = importList.BaseURL
	i.APIKey = importList.APIKey
	i.QualityProfileID = importList.QualityProfileID
	i.ID = importList.ID
	i.EnableAutomaticAdd = importList.EnableAutomaticAdd
	i.SeasonFolder = importList.SeasonFolder
}

func (r *ImportListSonarrResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListSonarrResourceName
}

func (r *ImportListSonarrResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->ImportList Sonarr resource.\nFor more information refer to [Import List](https://wiki.servarr.com/sonarr/settings#import-lists) and [Sonarr](https://wiki.servarr.com/sonarr/supported#sonarr).",
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
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Required:            true,
				Sensitive:           true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Required:            true,
			},
			"language_profile_ids": schema.SetAttribute{
				MarkdownDescription: "Language profile IDs.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"quality_profile_ids": schema.SetAttribute{
				MarkdownDescription: "Quality profile IDs.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"tag_ids": schema.SetAttribute{
				MarkdownDescription: "Tag IDs.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (r *ImportListSonarrResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ImportListSonarrResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportListSonarr

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportListSonarr
	request := importList.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ImportListAPI.CreateImportList(ctx).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListSonarrResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListSonarrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListSonarrResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportListSonarr

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportListSonarr current value
	response, _, err := r.client.ImportListAPI.GetImportListById(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListSonarrResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListSonarrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListSonarrResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportListSonarr

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportListSonarr
	request := importList.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ImportListAPI.UpdateImportList(ctx, strconv.Itoa(int(request.GetId()))).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListSonarrResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListSonarrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListSonarrResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportListSonarr current value
	_, err := r.client.ImportListAPI.DeleteImportList(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, importListSonarrResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListSonarrResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListSonarrResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+importListSonarrResourceName+": "+req.ID)
}

func (i *ImportListSonarr) write(ctx context.Context, importList *sonarr.ImportListResource, diags *diag.Diagnostics) {
	genericImportList := i.toImportList()
	genericImportList.write(ctx, importList, diags)
	i.fromImportList(genericImportList)
}

func (i *ImportListSonarr) read(ctx context.Context, diags *diag.Diagnostics) *sonarr.ImportListResource {
	return i.toImportList().read(ctx, diags)
}
