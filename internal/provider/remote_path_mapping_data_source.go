package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const remotePathMappingDataSourceName = "remote_path_mapping"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RemotePathMappingDataSource{}

func NewRemotePathMappingDataSource() datasource.DataSource {
	return &RemotePathMappingDataSource{}
}

// RemotePathMappingDataSource defines the remote path mapping implementation.
type RemotePathMappingDataSource struct {
	client *sonarr.APIClient
}

func (d *RemotePathMappingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + remotePathMappingDataSourceName
}

func (d *RemotePathMappingDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Download Clients -->Single [Remote Path Mapping](../resources/remote_path_mapping).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Remote Path Mapping ID.",
				Required:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "Download Client host.",
				Computed:            true,
			},
			"remote_path": schema.StringAttribute{
				MarkdownDescription: "Download Client remote path.",
				Computed:            true,
			},
			"local_path": schema.StringAttribute{
				MarkdownDescription: "Local path.",
				Computed:            true,
			},
		},
	}
}

func (d *RemotePathMappingDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *RemotePathMappingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *RemotePathMapping

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get remote path mapping current value
	response, _, err := d.client.RemotePathMappingApi.ListRemotePathMapping(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, remotePathMappingDataSourceName, err))

		return
	}

	data.find(data.ID.ValueInt64(), response, &resp.Diagnostics)
	tflog.Trace(ctx, "read "+remotePathMappingDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RemotePathMapping) find(id int64, mappings []*sonarr.RemotePathMappingResource, diags *diag.Diagnostics) {
	for _, m := range mappings {
		if int64(m.GetId()) == id {
			r.write(m)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(remotePathMappingDataSourceName, "id", strconv.Itoa(int(id))))
}
