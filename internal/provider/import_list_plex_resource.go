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
	importListPlexResourceName   = "import_list_plex"
	importListPlexImplementation = "PlexImport"
	importListPlexConfigContrat  = "PlexListSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListPlexResource{}
	_ resource.ResourceWithImportState = &ImportListPlexResource{}
)

func NewImportListPlexResource() resource.Resource {
	return &ImportListPlexResource{}
}

// ImportListPlexResource defines the import list implementation.
type ImportListPlexResource struct {
	client *sonarr.Sonarr
}

// ImportListPlex describes the import list data model.
type ImportListPlex struct {
	Tags               types.Set    `tfsdk:"tags"`
	Name               types.String `tfsdk:"name"`
	ShouldMonitor      types.String `tfsdk:"should_monitor"`
	RootFolderPath     types.String `tfsdk:"root_folder_path"`
	SeriesType         types.String `tfsdk:"series_type"`
	AccessToken        types.String `tfsdk:"access_token"`
	LanguageProfileID  types.Int64  `tfsdk:"language_profile_id"`
	QualityProfileID   types.Int64  `tfsdk:"quality_profile_id"`
	ID                 types.Int64  `tfsdk:"id"`
	EnableAutomaticAdd types.Bool   `tfsdk:"enable_automatic_add"`
	SeasonFolder       types.Bool   `tfsdk:"season_folder"`
}

func (i ImportListPlex) toImportList() *ImportList {
	return &ImportList{
		Tags:               i.Tags,
		Name:               i.Name,
		ShouldMonitor:      i.ShouldMonitor,
		RootFolderPath:     i.RootFolderPath,
		SeriesType:         i.SeriesType,
		AccessToken:        i.AccessToken,
		LanguageProfileID:  i.LanguageProfileID,
		QualityProfileID:   i.QualityProfileID,
		ID:                 i.ID,
		EnableAutomaticAdd: i.EnableAutomaticAdd,
		SeasonFolder:       i.SeasonFolder,
	}
}

func (i *ImportListPlex) fromImportList(importList *ImportList) {
	i.Tags = importList.Tags
	i.Name = importList.Name
	i.ShouldMonitor = importList.ShouldMonitor
	i.RootFolderPath = importList.RootFolderPath
	i.SeriesType = importList.SeriesType
	i.AccessToken = importList.AccessToken
	i.LanguageProfileID = importList.LanguageProfileID
	i.QualityProfileID = importList.QualityProfileID
	i.ID = importList.ID
	i.EnableAutomaticAdd = importList.EnableAutomaticAdd
	i.SeasonFolder = importList.SeasonFolder
}

func (r *ImportListPlexResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListPlexResourceName
}

func (r *ImportListPlexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->ImportList Plex resource.\nFor more information refer to [Import List](https://wiki.servarr.com/sonarr/settings#import-lists) and [Plex](https://wiki.servarr.com/sonarr/supported#plex).",
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
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access token.",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *ImportListPlexResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ImportListPlexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportListPlex

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportListPlex
	request := importList.read(ctx)

	response, err := r.client.AddImportListContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", importListPlexResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListPlexResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListPlexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportListPlex

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportListPlex current value
	response, err := r.client.GetImportListContext(ctx, importList.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", importListPlexResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListPlexResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListPlexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportListPlex

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportListPlex
	request := importList.read(ctx)

	response, err := r.client.UpdateImportListContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", importListPlexResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListPlexResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListPlexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var importList *ImportListPlex

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportListPlex current value
	err := r.client.DeleteImportListContext(ctx, importList.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", importListPlexResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListPlexResourceName+": "+strconv.Itoa(int(importList.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListPlexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+importListPlexResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (i *ImportListPlex) write(ctx context.Context, importList *sonarr.ImportListOutput) {
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

func (i *ImportListPlex) read(ctx context.Context) *sonarr.ImportListInput {
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
		ConfigContract:     importListPlexConfigContrat,
		Implementation:     importListPlexImplementation,
		ID:                 i.ID.ValueInt64(),
		Name:               i.Name.ValueString(),
		Tags:               tags,
		Fields:             i.toImportList().readFields(ctx),
	}
}
