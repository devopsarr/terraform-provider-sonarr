package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ provider.DataSourceType = dataIndexerConfigType{}
	_ datasource.DataSource   = dataIndexerConfig{}
)

type dataIndexerConfigType struct{}

type dataIndexerConfig struct {
	provider sonarrProvider
}

func (t dataIndexerConfigType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "[Indexer Config](../resources/indexer_config).",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Delay Profile ID.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"maximum_size": {
				MarkdownDescription: "Maximum size.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"minimum_age": {
				MarkdownDescription: "Minimum age.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"retention": {
				MarkdownDescription: "Retention.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"rss_sync_interval": {
				MarkdownDescription: "RSS sync interval.",
				Computed:            true,
				Type:                types.Int64Type,
			},
		},
	}, nil
}

func (t dataIndexerConfigType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return dataIndexerConfig{
		provider: provider,
	}, diags
}

func (d dataIndexerConfig) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get indexer config current value
	response, err := d.provider.client.GetIndexerConfigContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read indexer cofig, got error: %s", err))

		return
	}

	result := writeIndexerConfig(response)
	diags := resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
}
