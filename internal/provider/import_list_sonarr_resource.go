package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
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
	client *sonarr.Sonarr
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
	LanguageProfileID  types.Int64  `tfsdk:"language_profile_id"`
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
		LanguageProfileID:  i.LanguageProfileID,
		QualityProfileID:   i.QualityProfileID,
		ID:                 i.ID,
		EnableAutomaticAdd: i.EnableAutomaticAdd,
		SeasonFolder:       i.SeasonFolder,
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
	i.LanguageProfileID = importList.LanguageProfileID
	i.QualityProfileID = importList.QualityProfileID
	i.ID = importList.ID
	i.EnableAutomaticAdd = importList.EnableAutomaticAdd
	i.SeasonFolder = importList.SeasonFolder
}

func (r *ImportListSonarrResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListSonarrResourceName
}

func (r *ImportListSonarrResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"language_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Language profile ID.",
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

func (r *ImportListSonarrResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportListSonarr

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportListSonarr
	request := importList.read(ctx)

	response, err := r.client.AddImportListContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", importListSonarrResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListSonarrResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	importList.write(ctx, response)
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
	response, err := r.client.GetImportListContext(ctx, importList.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", importListSonarrResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListSonarrResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	importList.write(ctx, response)
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
	request := importList.read(ctx)

	response, err := r.client.UpdateImportListContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", importListSonarrResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListSonarrResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListSonarrResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var importList *ImportListSonarr

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportListSonarr current value
	err := r.client.DeleteImportListContext(ctx, importList.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", importListSonarrResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListSonarrResourceName+": "+strconv.Itoa(int(importList.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListSonarrResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+importListSonarrResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (i *ImportListSonarr) write(ctx context.Context, importList *sonarr.ImportListOutput) {
	genericImportList := ImportList{
		Name:               types.StringValue(importList.Name),
		ShouldMonitor:      types.StringValue(importList.ShouldMonitor),
		RootFolderPath:     types.StringValue(importList.RootFolderPath),
		SeriesType:         types.StringValue(importList.SeriesType),
		LanguageProfileID:  types.Int64Value(importList.LanguageProfileID),
		QualityProfileID:   types.Int64Value(importList.QualityProfileID),
		ID:                 types.Int64Value(importList.ID),
		EnableAutomaticAdd: types.BoolValue(importList.EnableAutomaticAdd),
		SeasonFolder:       types.BoolValue(importList.SeasonFolder),
	}
	genericImportList.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, importList.Tags)
	genericImportList.writeFields(ctx, importList.Fields)
	i.fromImportList(&genericImportList)
}

func (i *ImportListSonarr) read(ctx context.Context) *sonarr.ImportListInput {
	var tags []int

	tfsdk.ValueAs(ctx, i.Tags, &tags)

	return &sonarr.ImportListInput{
		ShouldMonitor:      i.ShouldMonitor.ValueString(),
		RootFolderPath:     i.RootFolderPath.ValueString(),
		SeriesType:         i.SeriesType.ValueString(),
		LanguageProfileID:  i.LanguageProfileID.ValueInt64(),
		QualityProfileID:   i.QualityProfileID.ValueInt64(),
		EnableAutomaticAdd: i.EnableAutomaticAdd.ValueBool(),
		SeasonFolder:       i.SeasonFolder.ValueBool(),
		ConfigContract:     importListSonarrConfigContract,
		Implementation:     importListSonarrImplementation,
		ID:                 i.ID.ValueInt64(),
		Name:               i.Name.ValueString(),
		Tags:               tags,
		Fields:             i.toImportList().readFields(ctx),
	}
}
