package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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

func (d *RemotePathMappingDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + remotePathMappingDataSourceName
}

func (d *RemotePathMappingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *sonarr.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *RemotePathMappingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var remoteMapping *RemotePathMapping

	resp.Diagnostics.Append(req.Config.Get(ctx, &remoteMapping)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get remote path mapping current value
	response, _, err := d.client.RemotePathMappingApi.ListRemotePathMapping(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", remotePathMappingDataSourceName, err))

		return
	}

	// Map response body to resource schema attribute
	mapping, err := findRemotePathMapping(remoteMapping.ID.ValueInt64(), response)
	if err != nil {
		resp.Diagnostics.AddError(tools.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", remotePathMappingDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+remotePathMappingDataSourceName)

	remoteMapping.write(mapping)
	resp.Diagnostics.Append(resp.State.Set(ctx, &remoteMapping)...)
}

func findRemotePathMapping(id int64, mappings []*sonarr.RemotePathMappingResource) (*sonarr.RemotePathMappingResource, error) {
	for _, m := range mappings {
		if int64(m.GetId()) == id {
			return m, nil
		}
	}

	return nil, tools.ErrDataNotFoundError(remotePathMappingDataSourceName, "id", strconv.Itoa(int(id)))
}
