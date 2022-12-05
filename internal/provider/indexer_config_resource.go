package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const indexerConfigResourceName = "indexer_config"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IndexerConfigResource{}
var _ resource.ResourceWithImportState = &IndexerConfigResource{}

func NewIndexerConfigResource() resource.Resource {
	return &IndexerConfigResource{}
}

// IndexerConfigResource defines the indexer config implementation.
type IndexerConfigResource struct {
	client *sonarr.Sonarr
}

// IndexerConfig describes the indexer config data model.
type IndexerConfig struct {
	ID              types.Int64 `tfsdk:"id"`
	MaximumSize     types.Int64 `tfsdk:"maximum_size"`
	MinimumAge      types.Int64 `tfsdk:"minimum_age"`
	Retention       types.Int64 `tfsdk:"retention"`
	RssSyncInterval types.Int64 `tfsdk:"rss_sync_interval"`
}

func (r *IndexerConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerConfigResourceName
}

func (r *IndexerConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexers -->Indexer Config resource.\nFor more information refer to [Indexer](https://wiki.servarr.com/sonarr/settings#options) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Indexer Config ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"maximum_size": schema.Int64Attribute{
				MarkdownDescription: "Maximum size.",
				Required:            true,
			},
			"minimum_age": schema.Int64Attribute{
				MarkdownDescription: "Minimum age.",
				Required:            true,
			},
			"retention": schema.Int64Attribute{
				MarkdownDescription: "Retention.",
				Required:            true,
			},
			"rss_sync_interval": schema.Int64Attribute{
				MarkdownDescription: "RSS sync interval.",
				Required:            true,
			},
		},
	}
}

func (r *IndexerConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IndexerConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var config *IndexerConfig

	resp.Diagnostics.Append(req.Plan.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	data := config.read()
	data.ID = 1

	// Create new IndexerConfig
	response, err := r.client.UpdateIndexerConfigContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", indexerConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerConfigResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	config.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (r *IndexerConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var config *IndexerConfig

	resp.Diagnostics.Append(req.State.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get indexerConfig current value
	response, err := r.client.GetIndexerConfigContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", indexerConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerConfigResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	config.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (r *IndexerConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var config *IndexerConfig

	resp.Diagnostics.Append(req.Plan.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := config.read()

	// Update IndexerConfig
	response, err := r.client.UpdateIndexerConfigContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", indexerConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerConfigResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	config.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (r *IndexerConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// IndexerConfig cannot be really deleted just removing configuration
	tflog.Trace(ctx, "decoupled "+indexerConfigResourceName+": 1")
	resp.State.RemoveResource(ctx)
}

func (r *IndexerConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+indexerConfigResourceName+": "+strconv.Itoa(1))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), 1)...)
}

func (c *IndexerConfig) write(indexerConfig *sonarr.IndexerConfig) {
	c.ID = types.Int64Value(indexerConfig.ID)
	c.MaximumSize = types.Int64Value(indexerConfig.MaximumSize)
	c.MinimumAge = types.Int64Value(indexerConfig.MinimumAge)
	c.Retention = types.Int64Value(indexerConfig.Retention)
	c.RssSyncInterval = types.Int64Value(indexerConfig.RssSyncInterval)
}

func (c *IndexerConfig) read() *sonarr.IndexerConfig {
	return &sonarr.IndexerConfig{
		ID:              c.ID.ValueInt64(),
		MaximumSize:     c.MaximumSize.ValueInt64(),
		MinimumAge:      c.MinimumAge.ValueInt64(),
		Retention:       c.Retention.ValueInt64(),
		RssSyncInterval: c.RssSyncInterval.ValueInt64(),
	}
}
