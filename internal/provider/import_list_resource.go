package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/exp/slices"
)

const importListResourceName = "import_list"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListResource{}
	_ resource.ResourceWithImportState = &ImportListResource{}
)

var (
	importListIntFields      = []string{"limit", "traktListType"}
	importListStringFields   = []string{"accessToken", "baseUrl", "apiKey", "refreshToken", "expires", "authUser", "username", "rating", "listname", "genres", "years", "traktAdditionalParameters"}
	importListIntSliceFields = []string{"profileIds", "languageProfileIds", "tagIds"}
)

func NewImportListResource() resource.Resource {
	return &ImportListResource{}
}

// ImportListResource defines the download client implementation.
type ImportListResource struct {
	client *sonarr.APIClient
}

// ImportList describes the download client data model.
type ImportList struct {
	Tags                      types.Set    `tfsdk:"tags"`
	LanguageProfileIds        types.Set    `tfsdk:"language_profile_ids"`
	ProfileIds                types.Set    `tfsdk:"quality_profile_ids"`
	TagIds                    types.Set    `tfsdk:"tag_ids"`
	Implementation            types.String `tfsdk:"implementation"`
	Name                      types.String `tfsdk:"name"`
	ShouldMonitor             types.String `tfsdk:"should_monitor"`
	RootFolderPath            types.String `tfsdk:"root_folder_path"`
	SeriesType                types.String `tfsdk:"series_type"`
	ConfigContract            types.String `tfsdk:"config_contract"`
	AccessToken               types.String `tfsdk:"access_token"`
	RefreshToken              types.String `tfsdk:"refresh_token"`
	Expires                   types.String `tfsdk:"expires"`
	BaseURL                   types.String `tfsdk:"base_url"`
	AuthUser                  types.String `tfsdk:"auth_user"`
	Username                  types.String `tfsdk:"username"`
	Rating                    types.String `tfsdk:"rating"`
	Listname                  types.String `tfsdk:"listname"`
	Genres                    types.String `tfsdk:"genres"`
	Years                     types.String `tfsdk:"years"`
	APIKey                    types.String `tfsdk:"api_key"`
	TraktAdditionalParameters types.String `tfsdk:"trakt_additional_parameters"`
	QualityProfileID          types.Int64  `tfsdk:"quality_profile_id"`
	ID                        types.Int64  `tfsdk:"id"`
	Limit                     types.Int64  `tfsdk:"limit"`
	TraktListType             types.Int64  `tfsdk:"trakt_list_type"`
	EnableAutomaticAdd        types.Bool   `tfsdk:"enable_automatic_add"`
	SeasonFolder              types.Bool   `tfsdk:"season_folder"`
}

func (r *ImportListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListResourceName
}

func (r *ImportListResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->Generic Import List resource. When possible use a specific resource instead.\nFor more information refer to [Import List](https://wiki.servarr.com/sonarr/settings#import-lists).",
		Attributes: map[string]schema.Attribute{
			"enable_automatic_add": schema.BoolAttribute{
				MarkdownDescription: "Enable automatic add flag.",
				Optional:            true,
				Computed:            true,
			},
			"season_folder": schema.BoolAttribute{
				MarkdownDescription: "Season folder flag.",
				Optional:            true,
				Computed:            true,
			},
			"quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Quality profile ID.",
				Optional:            true,
				Computed:            true,
			},
			"should_monitor": schema.StringAttribute{
				MarkdownDescription: "Should monitor.",
				Optional:            true,
				Computed:            true,
			},
			"root_folder_path": schema.StringAttribute{
				MarkdownDescription: "Root folder path.",
				Optional:            true,
				Computed:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "ImportList implementation name.",
				Optional:            true,
				Computed:            true,
			},
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "ImportList configuration template.",
				Required:            true,
			},
			"series_type": schema.StringAttribute{
				MarkdownDescription: "Series type.",
				Required:            true,
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
				MarkdownDescription: "Trakt list type.",
				Optional:            true,
				Computed:            true,
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
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"auth_user": schema.StringAttribute{
				MarkdownDescription: "Auth User.",
				Optional:            true,
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
			},
			"rating": schema.StringAttribute{
				MarkdownDescription: "Rating.",
				Optional:            true,
				Computed:            true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
			},
			"expires": schema.StringAttribute{
				MarkdownDescription: "Expires.",
				Optional:            true,
				Computed:            true,
			},
			"listname": schema.StringAttribute{
				MarkdownDescription: "Expires.",
				Optional:            true,
				Computed:            true,
			},
			"genres": schema.StringAttribute{
				MarkdownDescription: "Expires.",
				Optional:            true,
				Computed:            true,
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

func (r *ImportListResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ImportListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportList

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportList
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.CreateImportList(ctx).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state ImportList

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ImportListResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportList

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportList current value
	response, _, err := r.client.ImportListApi.GetImportListById(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	// this is needed because of many empty fields are unknown in both plan and read
	var state ImportList

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ImportListResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportList

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportList
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.UpdateImportList(ctx, strconv.Itoa(int(request.GetId()))).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state ImportList

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ImportListResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var importList *ImportList

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportList current value
	_, err := r.client.ImportListApi.DeleteImportList(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListResourceName+": "+strconv.Itoa(int(importList.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+importListResourceName+": "+req.ID)
}

func (i *ImportList) write(ctx context.Context, importList *sonarr.ImportListResource) {
	i.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, importList.Tags)
	i.EnableAutomaticAdd = types.BoolValue(importList.GetEnableAutomaticAdd())
	i.SeasonFolder = types.BoolValue(importList.GetSeasonFolder())
	i.QualityProfileID = types.Int64Value(int64(importList.GetQualityProfileId()))
	i.ID = types.Int64Value(int64(importList.GetId()))
	i.ConfigContract = types.StringValue(importList.GetConfigContract())
	i.Implementation = types.StringValue(importList.GetImplementation())
	i.ShouldMonitor = types.StringValue(string(importList.GetShouldMonitor()))
	i.RootFolderPath = types.StringValue(importList.GetRootFolderPath())
	i.SeriesType = types.StringValue(string(importList.GetSeriesType()))
	i.Name = types.StringValue(importList.GetName())
	i.LanguageProfileIds = types.SetValueMust(types.Int64Type, nil)
	i.ProfileIds = types.SetValueMust(types.Int64Type, nil)
	i.TagIds = types.SetValueMust(types.Int64Type, nil)
	i.writeFields(ctx, importList.Fields)
}

func (i *ImportList) writeFields(ctx context.Context, fields []*sonarr.Field) {
	for _, f := range fields {
		if f.Value == nil {
			continue
		}

		if slices.Contains(importListStringFields, f.GetName()) {
			helpers.WriteStringField(f, i)

			continue
		}

		if slices.Contains(importListIntFields, f.GetName()) {
			helpers.WriteIntField(f, i)

			continue
		}

		if slices.Contains(importListIntSliceFields, f.GetName()) {
			helpers.WriteIntSliceField(ctx, f, i)

			continue
		}
	}
}

func (i *ImportList) read(ctx context.Context) *sonarr.ImportListResource {
	var tags []*int32

	tfsdk.ValueAs(ctx, i.Tags, &tags)

	list := sonarr.NewImportListResource()
	list.SetEnableAutomaticAdd(i.EnableAutomaticAdd.ValueBool())
	list.SetSeasonFolder(i.SeasonFolder.ValueBool())
	list.SetQualityProfileId(int32(i.QualityProfileID.ValueInt64()))
	list.SetId(int32(i.ID.ValueInt64()))
	list.SetShouldMonitor(sonarr.MonitorTypes(i.ShouldMonitor.ValueString()))
	list.SetRootFolderPath(i.RootFolderPath.ValueString())
	list.SetSeriesType(sonarr.SeriesTypes(i.SeriesType.ValueString()))
	list.SetConfigContract(i.ConfigContract.ValueString())
	list.SetImplementation(i.Implementation.ValueString())
	list.SetName(i.Name.ValueString())
	list.SetTags(tags)
	list.SetFields(i.readFields(ctx))

	return list
}

func (i *ImportList) readFields(ctx context.Context) []*sonarr.Field {
	var output []*sonarr.Field

	for _, j := range importListIntFields {
		if field := helpers.ReadIntField(j, i); field != nil {
			output = append(output, field)
		}
	}

	for _, s := range importListStringFields {
		if field := helpers.ReadStringField(s, i); field != nil {
			output = append(output, field)
		}
	}

	for _, s := range importListIntSliceFields {
		if field := helpers.ReadIntSliceField(ctx, s, i); field != nil {
			output = append(output, field)
		}
	}

	return output
}
