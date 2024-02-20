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

const remotePathMappingsDataSourceName = "remote_path_mappings"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RemotePathMappingsDataSource{}

func NewRemotePathMappingsDataSource() datasource.DataSource {
	return &RemotePathMappingsDataSource{}
}

// RemotePathMappingsDataSource defines the remote path mappings implementation.
type RemotePathMappingsDataSource struct {
	client *sonarr.APIClient
}

// RemotePathMappings describes the remote path mappings data model.
type RemotePathMappings struct {
	RemotePathMappings types.Set    `tfsdk:"remote_path_mappings"`
	ID                 types.String `tfsdk:"id"`
}

func (d *RemotePathMappingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + remotePathMappingsDataSourceName
}

func (d *RemotePathMappingsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Download Clients -->List all available [Remote Path Mappings](../resources/remote_path_mapping).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"remote_path_mappings": schema.SetNestedAttribute{
				MarkdownDescription: "Remote Path Mapping list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
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
						"id": schema.Int64Attribute{
							MarkdownDescription: "RemotePathMapping ID.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *RemotePathMappingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *RemotePathMappingsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get remotePathMappings current value
	response, _, err := d.client.RemotePathMappingAPI.ListRemotePathMapping(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, remotePathMappingsDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+remotePathMappingsDataSourceName)
	// Map response body to resource schema attribute
	mappings := make([]RemotePathMapping, len(response))
	for i, p := range response {
		mappings[i].write(&p)
	}

	pathList, diags := types.SetValueFrom(ctx, RemotePathMapping{}.getType(), mappings)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, RemotePathMappings{RemotePathMappings: pathList, ID: types.StringValue(strconv.Itoa(len(response)))})...)
}
