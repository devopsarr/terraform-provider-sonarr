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
	metadataRoksboxResourceName   = "metadata_roksbox"
	metadataRoksboxImplementation = "RoksboxMetadata"
	metadataRoksboxConfigContract = "RoksboxMetadataSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &MetadataRoksboxResource{}
	_ resource.ResourceWithImportState = &MetadataRoksboxResource{}
)

func NewMetadataRoksboxResource() resource.Resource {
	return &MetadataRoksboxResource{}
}

// MetadataRoksboxResource defines the Roksbox metadata implementation.
type MetadataRoksboxResource struct {
	client *sonarr.APIClient
}

// MetadataRoksbox describes the Roksbox metadata data model.
type MetadataRoksbox struct {
	Tags            types.Set    `tfsdk:"tags"`
	Name            types.String `tfsdk:"name"`
	ID              types.Int64  `tfsdk:"id"`
	Enable          types.Bool   `tfsdk:"enable"`
	SeriesImages    types.Bool   `tfsdk:"series_images"`
	SeasonImages    types.Bool   `tfsdk:"season_images"`
	EpisodeMetadata types.Bool   `tfsdk:"episode_metadata"`
	EpisodeImages   types.Bool   `tfsdk:"episode_images"`
}

func (m MetadataRoksbox) toMetadata() *Metadata {
	return &Metadata{
		Tags:            m.Tags,
		Name:            m.Name,
		ID:              m.ID,
		Enable:          m.Enable,
		SeasonImages:    m.SeasonImages,
		SeriesImages:    m.SeriesImages,
		EpisodeImages:   m.EpisodeImages,
		EpisodeMetadata: m.EpisodeMetadata,
		ConfigContract:  types.StringValue(metadataRoksboxConfigContract),
		Implementation:  types.StringValue(metadataRoksboxImplementation),
	}
}

func (m *MetadataRoksbox) fromMetadata(metadata *Metadata) {
	m.ID = metadata.ID
	m.Name = metadata.Name
	m.Tags = metadata.Tags
	m.Enable = metadata.Enable
	m.SeasonImages = metadata.SeasonImages
	m.SeriesImages = metadata.SeriesImages
	m.EpisodeImages = metadata.EpisodeImages
	m.EpisodeMetadata = metadata.EpisodeMetadata
}

func (r *MetadataRoksboxResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataRoksboxResourceName
}

func (r *MetadataRoksboxResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Metadata -->Metadata Roksbox resource.\nFor more information refer to [Metadata](https://wiki.servarr.com/sonarr/settings#metadata) and [ROKSBOX](https://wiki.servarr.com/sonarr/supported#roksboxmetadata).",
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

func (r *MetadataRoksboxResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *MetadataRoksboxResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var metadata *MetadataRoksbox

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new MetadataRoksbox
	request := metadata.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.MetadataApi.CreateMetadata(ctx).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, metadataRoksboxResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+metadataRoksboxResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataRoksboxResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var metadata *MetadataRoksbox

	resp.Diagnostics.Append(req.State.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get MetadataRoksbox current value
	response, _, err := r.client.MetadataApi.GetMetadataById(ctx, int32(metadata.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataRoksboxResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataRoksboxResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataRoksboxResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var metadata *MetadataRoksbox

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update MetadataRoksbox
	request := metadata.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.MetadataApi.UpdateMetadata(ctx, strconv.Itoa(int(request.GetId()))).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to update "+metadataRoksboxResourceName+", got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated "+metadataRoksboxResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataRoksboxResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete MetadataRoksbox current value
	_, err := r.client.MetadataApi.DeleteMetadata(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, metadataRoksboxResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+metadataRoksboxResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *MetadataRoksboxResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+metadataRoksboxResourceName+": "+req.ID)
}

func (m *MetadataRoksbox) write(ctx context.Context, metadata *sonarr.MetadataResource, diags *diag.Diagnostics) {
	genericMetadata := m.toMetadata()
	genericMetadata.write(ctx, metadata, diags)
	m.fromMetadata(genericMetadata)
}

func (m *MetadataRoksbox) read(ctx context.Context, diags *diag.Diagnostics) *sonarr.MetadataResource {
	return m.toMetadata().read(ctx, diags)
}
