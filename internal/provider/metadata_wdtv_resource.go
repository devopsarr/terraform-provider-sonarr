package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	metadataWdtvResourceName   = "metadata_wdtv"
	metadataWdtvImplementation = "WdtvMetadata"
	metadataWdtvConfigContract = "WdtvMetadataSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &MetadataWdtvResource{}
	_ resource.ResourceWithImportState = &MetadataWdtvResource{}
)

func NewMetadataWdtvResource() resource.Resource {
	return &MetadataWdtvResource{}
}

// MetadataWdtvResource defines the Wdtv metadata implementation.
type MetadataWdtvResource struct {
	client *sonarr.APIClient
}

// MetadataWdtv describes the Wdtv metadata data model.
type MetadataWdtv struct {
	Tags            types.Set    `tfsdk:"tags"`
	Name            types.String `tfsdk:"name"`
	ID              types.Int64  `tfsdk:"id"`
	Enable          types.Bool   `tfsdk:"enable"`
	SeriesImages    types.Bool   `tfsdk:"series_images"`
	SeasonImages    types.Bool   `tfsdk:"season_images"`
	EpisodeMetadata types.Bool   `tfsdk:"episode_metadata"`
	EpisodeImages   types.Bool   `tfsdk:"episode_images"`
}

func (m MetadataWdtv) toMetadata() *Metadata {
	return &Metadata{
		Tags:            m.Tags,
		Name:            m.Name,
		ID:              m.ID,
		Enable:          m.Enable,
		SeasonImages:    m.SeasonImages,
		SeriesImages:    m.SeriesImages,
		EpisodeImages:   m.EpisodeImages,
		EpisodeMetadata: m.EpisodeMetadata,
		ConfigContract:  types.StringValue(metadataWdtvConfigContract),
		Implementation:  types.StringValue(metadataWdtvImplementation),
	}
}

func (m *MetadataWdtv) fromMetadata(metadata *Metadata) {
	m.ID = metadata.ID
	m.Name = metadata.Name
	m.Tags = metadata.Tags
	m.Enable = metadata.Enable
	m.SeasonImages = metadata.SeasonImages
	m.SeriesImages = metadata.SeriesImages
	m.EpisodeImages = metadata.EpisodeImages
	m.EpisodeMetadata = metadata.EpisodeMetadata
}

func (r *MetadataWdtvResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataWdtvResourceName
}

func (r *MetadataWdtvResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Metadata -->Metadata Wdtv resource.\nFor more information refer to [Metadata](https://wiki.servarr.com/sonarr/settings#metadata) and [WDTV](https://wiki.servarr.com/sonarr/supported#wdtvmetadata).",
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
			"series_images": schema.BoolAttribute{
				MarkdownDescription: "Series images flag.",
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
			"episode_metadata": schema.BoolAttribute{
				MarkdownDescription: "Episode metadata flag.",
				Required:            true,
			},
		},
	}
}

func (r *MetadataWdtvResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *MetadataWdtvResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var metadata *MetadataWdtv

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new MetadataWdtv
	request := metadata.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.MetadataApi.CreateMetadata(ctx).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, metadataWdtvResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+metadataWdtvResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataWdtvResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var metadata *MetadataWdtv

	resp.Diagnostics.Append(req.State.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get MetadataWdtv current value
	response, _, err := r.client.MetadataApi.GetMetadataById(ctx, int32(metadata.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataWdtvResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataWdtvResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataWdtvResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var metadata *MetadataWdtv

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update MetadataWdtv
	request := metadata.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.MetadataApi.UpdateMetadata(ctx, strconv.Itoa(int(request.GetId()))).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to update "+metadataWdtvResourceName+", got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated "+metadataWdtvResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataWdtvResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete MetadataWdtv current value
	_, err := r.client.MetadataApi.DeleteMetadata(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, metadataWdtvResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+metadataWdtvResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *MetadataWdtvResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+metadataWdtvResourceName+": "+req.ID)
}

func (m *MetadataWdtv) write(ctx context.Context, metadata *sonarr.MetadataResource, diags *diag.Diagnostics) {
	genericMetadata := m.toMetadata()
	genericMetadata.write(ctx, metadata, diags)
	m.fromMetadata(genericMetadata)
}

func (m *MetadataWdtv) read(ctx context.Context, diags *diag.Diagnostics) *sonarr.MetadataResource {
	return m.toMetadata().read(ctx, diags)
}
