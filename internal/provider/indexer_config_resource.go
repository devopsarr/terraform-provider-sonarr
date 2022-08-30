package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ provider.ResourceType            = resourceIndexerConfigType{}
	_ resource.Resource                = resourceIndexerConfig{}
	_ resource.ResourceWithImportState = resourceIndexerConfig{}
)

type resourceIndexerConfigType struct{}

type resourceIndexerConfig struct {
	provider sonarrProvider
}

// IndexerConfig is the IndexerConfig resource.
type IndexerConfig struct {
	ID              types.Int64 `tfsdk:"id"`
	MaximumSize     types.Int64 `tfsdk:"maximum_size"`
	MinimumAge      types.Int64 `tfsdk:"minimum_age"`
	Retention       types.Int64 `tfsdk:"retention"`
	RssSyncInterval types.Int64 `tfsdk:"rss_sync_interval"`
}

func (t resourceIndexerConfigType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Indexer Config resource.<br/>For more information refer to [Indexer](https://wiki.servarr.com/sonarr/settings#options) documentation.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Indexer Config ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"maximum_size": {
				MarkdownDescription: "Maximum size.",
				Required:            true,
				Type:                types.Int64Type,
			},
			"minimum_age": {
				MarkdownDescription: "Minimum age.",
				Required:            true,
				Type:                types.Int64Type,
			},
			"retention": {
				MarkdownDescription: "Retention.",
				Required:            true,
				Type:                types.Int64Type,
			},
			"rss_sync_interval": {
				MarkdownDescription: "RSS sync interval.",
				Required:            true,
				Type:                types.Int64Type,
			},
		},
	}, nil
}

func (t resourceIndexerConfigType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return resourceIndexerConfig{
		provider: provider,
	}, diags
}

func (r resourceIndexerConfig) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan IndexerConfig
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Init call if we remove this it the very first update on a brand new instance will fail
	init, err := r.provider.client.GetIndexerConfigContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to init indexerConfig, got error: %s", err))

		return
	}

	_, err = r.provider.client.UpdateIndexerConfigContext(ctx, init)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to init indexerConfig, got error: %s", err))

		return
	}

	// Build Create resource
	data := readIndexerConfig(&plan)
	data.ID = 1

	// Create new IndexerConfig
	response, err := r.provider.client.UpdateIndexerConfigContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create indexerConfig, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "created indexerConfig: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := writeIndexerConfig(response)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceIndexerConfig) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state IndexerConfig
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get indexerConfig current value
	response, err := r.provider.client.GetIndexerConfigContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read indexerConfig, got error: %s", err))

		return
	}
	// Map response body to resource schema attribute
	result := writeIndexerConfig(response)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceIndexerConfig) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan IndexerConfig
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := readIndexerConfig(&plan)

	// Update IndexerConfig
	response, err := r.provider.client.UpdateIndexerConfigContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update indexerConfig, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "update indexerConfig: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := writeIndexerConfig(response)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceIndexerConfig) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// IndexerConfig cannot be really deleted just removing configuration
	resp.State.RemoveResource(ctx)
}

func (r resourceIndexerConfig) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), 1)...)
}

func writeIndexerConfig(indexerConfig *sonarr.IndexerConfig) *IndexerConfig {
	return &IndexerConfig{
		ID:              types.Int64{Value: indexerConfig.ID},
		MaximumSize:     types.Int64{Value: indexerConfig.MaximumSize},
		MinimumAge:      types.Int64{Value: indexerConfig.MinimumAge},
		Retention:       types.Int64{Value: indexerConfig.Retention},
		RssSyncInterval: types.Int64{Value: indexerConfig.RssSyncInterval},
	}
}

func readIndexerConfig(indexerConfig *IndexerConfig) *sonarr.IndexerConfig {
	return &sonarr.IndexerConfig{
		ID:              indexerConfig.ID.Value,
		MaximumSize:     indexerConfig.MaximumSize.Value,
		MinimumAge:      indexerConfig.MinimumAge.Value,
		Retention:       indexerConfig.Retention.Value,
		RssSyncInterval: indexerConfig.RssSyncInterval.Value,
	}
}
