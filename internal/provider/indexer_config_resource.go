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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const indexerConfigResourceName = "indexer_config"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerConfigResource{}
	_ resource.ResourceWithImportState = &IndexerConfigResource{}
)

func NewIndexerConfigResource() resource.Resource {
	return &IndexerConfigResource{}
}

// IndexerConfigResource defines the indexer config implementation.
type IndexerConfigResource struct {
	client *sonarr.APIClient
}

// IndexerConfig describes the indexer config data model.
type IndexerConfig struct {
	ID              types.Int64 `tfsdk:"id"`
	MaximumSize     types.Int64 `tfsdk:"maximum_size"`
	MinimumAge      types.Int64 `tfsdk:"minimum_age"`
	Retention       types.Int64 `tfsdk:"retention"`
	RssSyncInterval types.Int64 `tfsdk:"rss_sync_interval"`
}

func (r *IndexerConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerConfigResourceName
}

func (r *IndexerConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *IndexerConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var config *IndexerConfig

	resp.Diagnostics.Append(req.Plan.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	request := config.read()
	request.SetId(1)

	// Create new IndexerConfig
	response, _, err := r.client.IndexerConfigAPI.UpdateIndexerConfig(ctx, strconv.Itoa(int(request.GetId()))).IndexerConfigResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, indexerConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerConfigResourceName+": "+strconv.Itoa(int(response.GetId())))
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
	response, _, err := r.client.IndexerConfigAPI.GetIndexerConfig(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerConfigResourceName+": "+strconv.Itoa(int(response.GetId())))
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
	request := config.read()

	// Update IndexerConfig
	response, _, err := r.client.IndexerConfigAPI.UpdateIndexerConfig(ctx, strconv.Itoa(int(request.GetId()))).IndexerConfigResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, indexerConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerConfigResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	config.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (r *IndexerConfigResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	// IndexerConfig cannot be really deleted just removing configuration
	tflog.Trace(ctx, "decoupled "+indexerConfigResourceName+": 1")
	resp.State.RemoveResource(ctx)
}

func (r *IndexerConfigResource) ImportState(ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Trace(ctx, "imported "+indexerConfigResourceName+": 1")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), 1)...)
}

func (c *IndexerConfig) write(indexerConfig *sonarr.IndexerConfigResource) {
	c.ID = types.Int64Value(int64(indexerConfig.GetId()))
	c.MaximumSize = types.Int64Value(int64(indexerConfig.GetMaximumSize()))
	c.MinimumAge = types.Int64Value(int64(indexerConfig.GetMinimumAge()))
	c.Retention = types.Int64Value(int64(indexerConfig.GetRetention()))
	c.RssSyncInterval = types.Int64Value(int64(indexerConfig.GetRssSyncInterval()))
}

func (c *IndexerConfig) read() *sonarr.IndexerConfigResource {
	config := sonarr.NewIndexerConfigResource()
	config.SetId(int32(c.ID.ValueInt64()))
	config.SetMaximumSize(int32(c.MaximumSize.ValueInt64()))
	config.SetMinimumAge(int32(c.MinimumAge.ValueInt64()))
	config.SetRetention(int32(c.Retention.ValueInt64()))
	config.SetRssSyncInterval(int32(c.RssSyncInterval.ValueInt64()))

	return config
}
