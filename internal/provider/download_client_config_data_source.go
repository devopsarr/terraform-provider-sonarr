package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const downloadClientConfigDataSourceName = "download_client_config"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DownloadClientConfigDataSource{}

func NewDownloadClientConfigDataSource() datasource.DataSource {
	return &DownloadClientConfigDataSource{}
}

// DownloadClientConfigDataSource defines the download client config implementation.
type DownloadClientConfigDataSource struct {
	client *sonarr.Sonarr
}

func (d *DownloadClientConfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientConfigDataSourceName
}

func (d *DownloadClientConfigDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "[subcategory:Download Clients]: #\n[Download Client Config](../resources/download_client_config).",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Download Client Config ID.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"enable_completed_download_handling": {
				MarkdownDescription: "Enable Completed Download Handling flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"auto_redownload_failed": {
				MarkdownDescription: "Auto Redownload Failed flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"download_client_working_folders": {
				MarkdownDescription: "Download Client Working Folders.",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (d *DownloadClientConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			helpers.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *DownloadClientConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get indexer config current value
	response, err := d.client.GetDownloadClientConfigContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", downloadClientConfigDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientConfigDataSourceName)

	config := DownloadClientConfig{}
	config.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}
