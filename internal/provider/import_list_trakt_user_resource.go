package provider

import (
	"context"
	"fmt"
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
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	importListTraktUserResourceName   = "import_list_trakt_user"
	importListTraktUserImplementation = "TraktUserImport"
	importListTraktUserConfigContract = "TraktUserSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListTraktUserResource{}
	_ resource.ResourceWithImportState = &ImportListTraktUserResource{}
)

func NewImportListTraktUserResource() resource.Resource {
	return &ImportListTraktUserResource{}
}

// ImportListTraktUserResource defines the import list implementation.
type ImportListTraktUserResource struct {
	client *sonarr.APIClient
}

// ImportListTraktUser describes the import list data model.
type ImportListTraktUser struct {
	Tags                      types.Set    `tfsdk:"tags"`
	Name                      types.String `tfsdk:"name"`
	ShouldMonitor             types.String `tfsdk:"should_monitor"`
	RootFolderPath            types.String `tfsdk:"root_folder_path"`
	SeriesType                types.String `tfsdk:"series_type"`
	Username                  types.String `tfsdk:"username"`
	AccessToken               types.String `tfsdk:"access_token"`
	RefreshToken              types.String `tfsdk:"refresh_token"`
	Expires                   types.String `tfsdk:"expires"`
	AuthUser                  types.String `tfsdk:"auth_user"`
	Rating                    types.String `tfsdk:"rating"`
	Genres                    types.String `tfsdk:"genres"`
	Years                     types.String `tfsdk:"years"`
	TraktAdditionalParameters types.String `tfsdk:"trakt_additional_parameters"`
	LanguageProfileID         types.Int64  `tfsdk:"language_profile_id"`
	QualityProfileID          types.Int64  `tfsdk:"quality_profile_id"`
	ID                        types.Int64  `tfsdk:"id"`
	Limit                     types.Int64  `tfsdk:"limit"`
	TraktListType             types.Int64  `tfsdk:"trakt_list_type"`
	EnableAutomaticAdd        types.Bool   `tfsdk:"enable_automatic_add"`
	SeasonFolder              types.Bool   `tfsdk:"season_folder"`
}

func (i ImportListTraktUser) toImportList() *ImportList {
	return &ImportList{
		Tags:                      i.Tags,
		Name:                      i.Name,
		ShouldMonitor:             i.ShouldMonitor,
		RootFolderPath:            i.RootFolderPath,
		SeriesType:                i.SeriesType,
		Username:                  i.Username,
		AccessToken:               i.AccessToken,
		RefreshToken:              i.RefreshToken,
		Expires:                   i.Expires,
		AuthUser:                  i.AuthUser,
		Rating:                    i.Rating,
		TraktListType:             i.TraktListType,
		Genres:                    i.Genres,
		Years:                     i.Years,
		TraktAdditionalParameters: i.TraktAdditionalParameters,
		Limit:                     i.Limit,
		LanguageProfileID:         i.LanguageProfileID,
		QualityProfileID:          i.QualityProfileID,
		ID:                        i.ID,
		EnableAutomaticAdd:        i.EnableAutomaticAdd,
		SeasonFolder:              i.SeasonFolder,
	}
}

func (i *ImportListTraktUser) fromImportList(importList *ImportList) {
	i.Tags = importList.Tags
	i.Name = importList.Name
	i.ShouldMonitor = importList.ShouldMonitor
	i.RootFolderPath = importList.RootFolderPath
	i.SeriesType = importList.SeriesType
	i.Username = importList.Username
	i.AccessToken = importList.AccessToken
	i.RefreshToken = importList.RefreshToken
	i.Expires = importList.Expires
	i.AuthUser = importList.AuthUser
	i.Rating = importList.Rating
	i.TraktListType = importList.TraktListType
	i.Genres = importList.Genres
	i.Years = importList.Years
	i.TraktAdditionalParameters = importList.TraktAdditionalParameters
	i.Limit = importList.Limit
	i.LanguageProfileID = importList.LanguageProfileID
	i.QualityProfileID = importList.QualityProfileID
	i.ID = importList.ID
	i.EnableAutomaticAdd = importList.EnableAutomaticAdd
	i.SeasonFolder = importList.SeasonFolder
}

func (r *ImportListTraktUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListTraktUserResourceName
}

func (r *ImportListTraktUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->ImportList TraktUser resource.\nFor more information refer to [Import List](https://wiki.servarr.com/sonarr/settings#import-lists) and [TraktUser](https://wiki.servarr.com/sonarr/supported#trakt_user).",
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
			"limit": schema.Int64Attribute{
				MarkdownDescription: "Limit.",
				Optional:            true,
				Computed:            true,
			},
			"trakt_list_type": schema.Int64Attribute{
				MarkdownDescription: "Trakt list type. '0' UserWatchList, '1' UserWatchedList, '2' UserCollectionList.",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1, 2),
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
			"rating": schema.StringAttribute{
				MarkdownDescription: "Rating.",
				Optional:            true,
				Computed:            true,
			},
			"expires": schema.StringAttribute{
				MarkdownDescription: "Expires.",
				Optional:            true,
				Computed:            true,
			},
			"genres": schema.StringAttribute{
				MarkdownDescription: "Expires.",
				Optional:            true,
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Required:            true,
			},
			"years": schema.StringAttribute{
				MarkdownDescription: "Expires.",
				Optional:            true,
				Computed:            true,
			},
			"trakt_additional_parameters": schema.StringAttribute{
				MarkdownDescription: "Trakt additional parameters.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *ImportListTraktUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ImportListTraktUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportListTraktUser

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportListTraktUser
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.CreateImportList(ctx).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListTraktUserResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListTraktUserResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListTraktUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportListTraktUser

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportListTraktUser current value
	response, _, err := r.client.ImportListApi.GetImportListById(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListTraktUserResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListTraktUserResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListTraktUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportListTraktUser

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportListTraktUser
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.UpdateImportList(ctx, strconv.Itoa(int(request.GetId()))).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListTraktUserResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListTraktUserResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListTraktUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var importList *ImportListTraktUser

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportListTraktUser current value
	_, err := r.client.ImportListApi.DeleteImportList(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListTraktUserResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListTraktUserResourceName+": "+strconv.Itoa(int(importList.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListTraktUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			helpers.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+importListTraktUserResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (i *ImportListTraktUser) write(ctx context.Context, importList *sonarr.ImportListResource) {
	genericImportList := ImportList{
		Name:               types.StringValue(importList.GetName()),
		ShouldMonitor:      types.StringValue(string(importList.GetShouldMonitor())),
		RootFolderPath:     types.StringValue(importList.GetRootFolderPath()),
		SeriesType:         types.StringValue(string(importList.GetSeriesType())),
		LanguageProfileID:  types.Int64Value(int64(importList.GetLanguageProfileId())),
		QualityProfileID:   types.Int64Value(int64(importList.GetQualityProfileId())),
		ID:                 types.Int64Value(int64(importList.GetId())),
		EnableAutomaticAdd: types.BoolValue(importList.GetEnableAutomaticAdd()),
		SeasonFolder:       types.BoolValue(importList.GetSeasonFolder()),
	}
	genericImportList.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, importList.Tags)
	genericImportList.writeFields(ctx, importList.Fields)
	i.fromImportList(&genericImportList)
}

func (i *ImportListTraktUser) read(ctx context.Context) *sonarr.ImportListResource {
	var tags []*int32

	tfsdk.ValueAs(ctx, i.Tags, &tags)

	list := sonarr.NewImportListResource()
	list.SetShouldMonitor(sonarr.MonitorTypes(i.ShouldMonitor.ValueString()))
	list.SetRootFolderPath(i.RootFolderPath.ValueString())
	list.SetSeriesType(sonarr.SeriesTypes(i.SeriesType.ValueString()))
	list.SetLanguageProfileId(int32(i.LanguageProfileID.ValueInt64()))
	list.SetQualityProfileId(int32(i.QualityProfileID.ValueInt64()))
	list.SetEnableAutomaticAdd(i.EnableAutomaticAdd.ValueBool())
	list.SetSeasonFolder(i.SeasonFolder.ValueBool())
	list.SetConfigContract(importListTraktUserConfigContract)
	list.SetImplementation(importListTraktUserImplementation)
	list.SetId(int32(i.ID.ValueInt64()))
	list.SetName(i.Name.ValueString())
	list.SetTags(tags)
	list.SetFields(i.toImportList().readFields(ctx))

	return list
}
