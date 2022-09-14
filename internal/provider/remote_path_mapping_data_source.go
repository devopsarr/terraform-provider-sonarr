package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RemotePathMappingDataSource{}

func NewRemotePathMappingDataSource() datasource.DataSource {
	return &RemotePathMappingDataSource{}
}

// RemotePathMappingDataSource defines the remote path mapping implementation.
type RemotePathMappingDataSource struct {
	client *sonarr.Sonarr
}

func (d *RemotePathMappingDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_path_mapping"
}

func (d *RemotePathMappingDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "Single [Remote Path Mapping](../resources/remote_path_mapping).",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Remote Path Mapping ID.",
				Required:            true,
				Type:                types.Int64Type,
			},
			"host": {
				MarkdownDescription: "Download Client host.",
				Computed:            true,
				Type:                types.StringType,
			},
			"remote_path": {
				MarkdownDescription: "Download Client remote path.",
				Computed:            true,
				Type:                types.StringType,
			},
			"local_path": {
				MarkdownDescription: "Local path.",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (d *RemotePathMappingDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *RemotePathMappingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RemotePathMapping

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get remote path mapping current value
	response, err := d.client.GetRemotePathMappingsContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read remote path mapping, got error: %s", err))

		return
	}

	// Map response body to resource schema attribute
	mapping, err := findRemotePathMapping(data.ID.Value, response)
	if err != nil {
		resp.Diagnostics.AddError(DataSourceError, fmt.Sprintf("Unable to find remote path mappings, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "read remote_path_mapping")

	result := writeRemotePathMapping(mapping)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func findRemotePathMapping(id int64, mappings []*sonarr.RemotePathMapping) (*sonarr.RemotePathMapping, error) {
	for _, m := range mappings {
		if m.ID == id {
			return m, nil
		}
	}

	return nil, fmt.Errorf("no remote path mapping with id %d", id)
}
