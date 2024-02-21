package provider

import (
	"context"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const metadataDataSourceName = "metadata"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MetadataDataSource{}

func NewMetadataDataSource() datasource.DataSource {
	return &MetadataDataSource{}
}

// MetadataDataSource defines the metadata implementation.
type MetadataDataSource struct {
	client *sonarr.APIClient
	auth   context.Context
}

func (d *MetadataDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataDataSourceName
}

func (d *MetadataDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Metadata -->Single [Metadata](../resources/metadata).",
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
				Required:            true,
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
	}
}

func (d *MetadataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *MetadataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Metadata

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get metadata current value
	response, _, err := d.client.MetadataAPI.ListMetadata(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataDataSourceName, err))

		return
	}

	data.find(ctx, data.Name.ValueString(), response, &resp.Diagnostics)
	tflog.Trace(ctx, "read "+metadataDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *Metadata) find(ctx context.Context, name string, metadatas []sonarr.MetadataResource, diags *diag.Diagnostics) {
	for _, metadata := range metadatas {
		if metadata.GetName() == name {
			m.write(ctx, &metadata, diags)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(metadataDataSourceName, "name", name))
}
