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

const (
	metadataKodiResourceName   = "metadata_kodi"
	metadataKodiImplementation = "XbmcMetadata"
	metadataKodiConfigContract = "XbmcMetadataSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &MetadataKodiResource{}
	_ resource.ResourceWithImportState = &MetadataKodiResource{}
)

func NewMetadataKodiResource() resource.Resource {
	return &MetadataKodiResource{}
}

// MetadataKodiResource defines the Kodi metadata implementation.
type MetadataKodiResource struct {
	client *sonarr.APIClient
}

// MetadataKodi describes the Kodi metadata data model.
type MetadataKodi struct {
	Tags              types.Set    `tfsdk:"tags"`
	Name              types.String `tfsdk:"name"`
	ID                types.Int64  `tfsdk:"id"`
	Enable            types.Bool   `tfsdk:"enable"`
	SeriesMetadata    types.Bool   `tfsdk:"series_metadata"`
	SeriesMetadataURL types.Bool   `tfsdk:"series_metadata_url"`
	SeriesImages      types.Bool   `tfsdk:"series_images"`
	SeasonImages      types.Bool   `tfsdk:"season_images"`
	EpisodeMetadata   types.Bool   `tfsdk:"episode_metadata"`
	EpisodeImages     types.Bool   `tfsdk:"episode_images"`
}

func (m MetadataKodi) toMetadata() *Metadata {
	return &Metadata{
		Tags:              m.Tags,
		Name:              m.Name,
		ID:                m.ID,
		EpisodeImages:     m.EpisodeImages,
		Enable:            m.Enable,
		SeriesMetadata:    m.SeriesMetadata,
		SeriesMetadataURL: m.SeriesMetadataURL,
		SeriesImages:      m.SeriesImages,
		SeasonImages:      m.SeasonImages,
		EpisodeMetadata:   m.EpisodeMetadata,
		ConfigContract:    types.StringValue(metadataKodiConfigContract),
		Implementation:    types.StringValue(metadataKodiImplementation),
	}
}

func (m *MetadataKodi) fromMetadata(metadata *Metadata) {
	m.ID = metadata.ID
	m.Name = metadata.Name
	m.Tags = metadata.Tags
	m.EpisodeImages = metadata.EpisodeImages
	m.Enable = metadata.Enable
	m.SeriesMetadata = metadata.SeriesMetadata
	m.SeriesImages = metadata.SeriesImages
	m.SeriesMetadataURL = metadata.SeriesMetadataURL
	m.SeasonImages = metadata.SeasonImages
	m.EpisodeMetadata = metadata.EpisodeMetadata
}

func (r *MetadataKodiResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataKodiResourceName
}

func (r *MetadataKodiResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Metadata -->Metadata Kodi resource.\nFor more information refer to [Metadata](https://wiki.servarr.com/sonarr/settings#metadata) and [KODI](https://wiki.servarr.com/sonarr/supported#xbmcmetadata).",
		Attributes: map[string]schema.Attribute{
			"enable": schema.BoolAttribute{
				MarkdownDescription: "Enable flag.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Metadata name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Metadata ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"episode_metadata": schema.BoolAttribute{
				MarkdownDescription: "Episode metadata flag.",
				Required:            true,
			},
			"episode_images": schema.BoolAttribute{
				MarkdownDescription: "Episode images flag.",
				Required:            true,
			},
			"season_images": schema.BoolAttribute{
				MarkdownDescription: "Season images flag.",
				Required:            true,
			},
			"series_images": schema.BoolAttribute{
				MarkdownDescription: "Series images flag.",
				Required:            true,
			},
			"series_metadata": schema.BoolAttribute{
				MarkdownDescription: "Series metadata flag.",
				Required:            true,
			},
			"series_metadata_url": schema.BoolAttribute{
				MarkdownDescription: "Series metadata URL flag.",
				Required:            true,
			},
		},
	}
}

func (r *MetadataKodiResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *MetadataKodiResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var metadata *MetadataKodi

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new MetadataKodi
	request := metadata.read(ctx)

	response, _, err := r.client.MetadataApi.CreateMetadata(ctx).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, metadataKodiResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+metadataKodiResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	metadata.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataKodiResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var metadata *MetadataKodi

	resp.Diagnostics.Append(req.State.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get MetadataKodi current value
	response, _, err := r.client.MetadataApi.GetMetadataById(ctx, int32(metadata.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataKodiResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataKodiResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	metadata.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataKodiResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var metadata *MetadataKodi

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update MetadataKodi
	request := metadata.read(ctx)

	response, _, err := r.client.MetadataApi.UpdateMetadata(ctx, strconv.Itoa(int(request.GetId()))).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to update "+metadataKodiResourceName+", got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated "+metadataKodiResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	metadata.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataKodiResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var metadata *MetadataKodi

	resp.Diagnostics.Append(req.State.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete MetadataKodi current value
	_, err := r.client.MetadataApi.DeleteMetadata(ctx, int32(metadata.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, metadataKodiResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+metadataKodiResourceName+": "+strconv.Itoa(int(metadata.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *MetadataKodiResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+metadataKodiResourceName+": "+req.ID)
}

func (m *MetadataKodi) write(ctx context.Context, metadata *sonarr.MetadataResource) {
	genericMetadata := m.toMetadata()
	genericMetadata.write(ctx, metadata)
	m.fromMetadata(genericMetadata)
}

func (m *MetadataKodi) read(ctx context.Context) *sonarr.MetadataResource {
	return m.toMetadata().read(ctx)
}
