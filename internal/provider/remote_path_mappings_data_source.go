package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const remotePathMappingsDataSourceName = "remote_path_mappings"

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
	RemotePathMappings types.Set    `tfsdk:"remote_path_mappings"`
	ID                 types.String `tfsdk:"id"`
}

func (d *RemotePathMappingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + remotePathMappingsDataSourceName
}

func (d *RemotePathMappingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *RemotePathMappingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *RemotePathMappings

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get remotePathMappings current value
	response, err := d.client.GetRemotePathMappingsContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", remotePathMappingsDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+remotePathMappingsDataSourceName)
	// Map response body to resource schema attribute
	mappings := make([]RemotePathMapping, len(response))
	for i, p := range response {
		mappings[i].write(p)
	}

	tfsdk.ValueFrom(ctx, mappings, data.RemotePathMappings.Type(context.Background()), &data.RemotePathMappings)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
