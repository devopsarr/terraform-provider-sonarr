package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const metadataConsumersDataSourceName = "metadata_consumers"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MetadataConsumersDataSource{}

func NewMetadataConsumersDataSource() datasource.DataSource {
	return &MetadataConsumersDataSource{}
}

// MetadataConsumersDataSource defines the metadataConsumers implementation.
type MetadataConsumersDataSource struct {
	client *sonarr.APIClient
	auth   context.Context
}

// MetadataConsumers describes the metadataConsumers data model.
type MetadataConsumers struct {
	MetadataConsumers types.Set    `tfsdk:"metadata_consumers"`
	ID                types.String `tfsdk:"id"`
}

func (d *MetadataConsumersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataConsumersDataSourceName
}

func (d *MetadataConsumersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Metadata -->\nList all available [Metadata Consumers](../resources/metadata).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"metadata_consumers": schema.SetNestedAttribute{
				MarkdownDescription: "MetadataConsumer list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							MarkdownDescription: "Enable flag.",
							Computed:            true,
						},
						"config_contract": schema.StringAttribute{
							MarkdownDescription: "Metadata configuration template.",
							Computed:            true,
						},
						"implementation": schema.StringAttribute{
							MarkdownDescription: "Metadata implementation name.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Metadata name.",
							Computed:            true,
						},
						"tags": schema.SetAttribute{
							MarkdownDescription: "List of associated tags.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
						"id": schema.Int64Attribute{
							MarkdownDescription: "Metadata ID.",
							Computed:            true,
						},
						// Field values
						"episode_metadata": schema.BoolAttribute{
							MarkdownDescription: "Episode metadata flag.",
							Computed:            true,
						},
						"episode_images": schema.BoolAttribute{
							MarkdownDescription: "Episode images flag.",
							Optional:            true,
						},
						"season_images": schema.BoolAttribute{
							MarkdownDescription: "Season images flag.",
							Computed:            true,
						},
						"series_images": schema.BoolAttribute{
							MarkdownDescription: "Series images flag.",
							Computed:            true,
						},
						"series_metadata": schema.BoolAttribute{
							MarkdownDescription: "Series metadata flag.",
							Computed:            true,
						},
						"series_metadata_url": schema.BoolAttribute{
							MarkdownDescription: "Series metadata URL flag.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *MetadataConsumersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *MetadataConsumersDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get metadataConsumers current value
	response, _, err := d.client.MetadataAPI.ListMetadata(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.List, metadataConsumersDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataConsumersDataSourceName)
	// Map response body to resource schema attribute
	consumers := make([]Metadata, len(response))
	for i, p := range response {
		consumers[i].write(ctx, &p, &resp.Diagnostics)
	}

	metadataList, diags := types.SetValueFrom(ctx, Metadata{}.getType(), consumers)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, MetadataConsumers{MetadataConsumers: metadataList, ID: types.StringValue(strconv.Itoa(len(response)))})...)
}
