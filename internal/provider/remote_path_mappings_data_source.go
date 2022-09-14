package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RemotePathMappingsDataSource{}

func NewRemotePathMappingsDataSource() datasource.DataSource {
	return &RemotePathMappingsDataSource{}
}

// RemotePathMappingsDataSource defines the remote path mappings implementation.
type RemotePathMappingsDataSource struct {
	client *sonarr.Sonarr
}

// RemotePathMappings describes the remote path mappings data model.
type RemotePathMappings struct {
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	ID                 types.String `tfsdk:"id"`
	RemotePathMappings types.Set    `tfsdk:"remote_path_mappings"`
}

func (d *RemotePathMappingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_path_mappings"
}

func (d *RemotePathMappingsDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "List all available [Remote Path Mappings](../resources/remote_path_mapping).",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"remote_path_mappings": {
				MarkdownDescription: "Remote Path Mapping list.",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
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
					"id": {
						MarkdownDescription: "RemotePathMapping ID.",
						Computed:            true,
						Type:                types.Int64Type,
					},
				}),
			},
		},
	}, nil
}

func (d *RemotePathMappingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RemotePathMappingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RemotePathMappings

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get remotePathMappings current value
	response, err := d.client.GetRemotePathMappingsContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read remotePathMappings, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "read remote_path_mappings")
	// Map response body to resource schema attribute
	profiles := *writeRemotePathMappings(response)
	tfsdk.ValueFrom(ctx, profiles, data.RemotePathMappings.Type(context.Background()), &data.RemotePathMappings)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.String{Value: strconv.Itoa(len(response))}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func writeRemotePathMappings(mappings []*sonarr.RemotePathMapping) *[]RemotePathMapping {
	output := make([]RemotePathMapping, len(mappings))
	for i, p := range mappings {
		output[i] = *writeRemotePathMapping(p)
	}

	return &output
}
