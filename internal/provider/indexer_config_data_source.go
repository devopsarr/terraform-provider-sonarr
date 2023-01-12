package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const indexerConfigDataSourceName = "indexer_config"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IndexerConfigDataSource{}

func NewIndexerConfigDataSource() datasource.DataSource {
	return &IndexerConfigDataSource{}
}

// IndexerConfigDataSource defines the indexer config implementation.
type IndexerConfigDataSource struct {
	client *sonarr.APIClient
}

func (d *IndexerConfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerConfigDataSourceName
}

func (d *IndexerConfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Indexers -->[Indexer Config](../resources/indexer_config).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Delay Profile ID.",
				Computed:            true,
			},
			"maximum_size": schema.Int64Attribute{
				MarkdownDescription: "Maximum size.",
				Computed:            true,
			},
			"minimum_age": schema.Int64Attribute{
				MarkdownDescription: "Minimum age.",
				Computed:            true,
			},
			"retention": schema.Int64Attribute{
				MarkdownDescription: "Retention.",
				Computed:            true,
			},
			"rss_sync_interval": schema.Int64Attribute{
				MarkdownDescription: "RSS sync interval.",
				Computed:            true,
			},
		},
	}
}

func (d *IndexerConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			helpers.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *sonarr.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *IndexerConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get indexer config current value
	response, _, err := d.client.IndexerConfigApi.GetIndexerConfig(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerConfigDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerConfigDataSourceName)

	status := IndexerConfig{}
	status.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, status)...)
}
