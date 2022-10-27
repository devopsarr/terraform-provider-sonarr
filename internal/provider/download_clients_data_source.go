package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const downloadClientsDataSourceName = "download_clients"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DownloadClientsDataSource{}

func NewDownloadClientsDataSource() datasource.DataSource {
	return &DownloadClientsDataSource{}
}

// DownloadClientsDataSource defines the download clients implementation.
type DownloadClientsDataSource struct {
	client *sonarr.Sonarr
}

// DownloadClients describes the download clients data model.
type DownloadClients struct {
	DownloadClients types.Set    `tfsdk:"download_clients"`
	ID              types.String `tfsdk:"id"`
}

func (d *DownloadClientsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientsDataSourceName
}

func (d *DownloadClientsDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "[subcategory:Download Clients]: #\nList all available [DownloadClients](../resources/download_client).",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"download_clients": {
				MarkdownDescription: "Download Client list..",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"enable": {
						MarkdownDescription: "Enable flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"remove_completed_downloads": {
						MarkdownDescription: "Remove completed downloads flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"remove_failed_downloads": {
						MarkdownDescription: "Remove failed downloads flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"priority": {
						MarkdownDescription: "Priority.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"config_contract": {
						MarkdownDescription: "DownloadClient configuration template.",
						Computed:            true,
						Type:                types.StringType,
					},
					"implementation": {
						MarkdownDescription: "DownloadClient implementation name.",
						Computed:            true,
						Type:                types.StringType,
					},
					"name": {
						MarkdownDescription: "Download Client name.",
						Computed:            true,
						Type:                types.StringType,
					},
					"protocol": {
						MarkdownDescription: "Protocol. Valid values are 'usenet' and 'torrent'.",
						Computed:            true,
						Type:                types.StringType,
					},
					"tags": {
						MarkdownDescription: "List of associated tags.",
						Computed:            true,
						Type: types.SetType{
							ElemType: types.Int64Type,
						},
					},
					"id": {
						MarkdownDescription: "Download Client ID.",
						Computed:            true,
						Type:                types.Int64Type,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							resource.UseStateForUnknown(),
						},
					},
					// Field values
					"add_paused": {
						MarkdownDescription: "Add paused flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"use_ssl": {
						MarkdownDescription: "Use SSL flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"start_on_add": {
						MarkdownDescription: "Start on add flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"sequential_order": {
						MarkdownDescription: "Sequential order flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"first_and_last": {
						MarkdownDescription: "First and last flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"add_stopped": {
						MarkdownDescription: "Add stopped flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"save_magnet_files": {
						MarkdownDescription: "Save magnet files flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"read_only": {
						MarkdownDescription: "Read only flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"watch_folder": {
						MarkdownDescription: "Watch folder flag.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"port": {
						MarkdownDescription: "Port.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"recent_tv_priority": {
						MarkdownDescription: "Recent TV priority. `0` Last, `1` First.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"older_tv_priority": {
						MarkdownDescription: "Older TV priority. `0` Last, `1` First.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"initial_state": {
						MarkdownDescription: "Initial state. `0` Start, `1` ForceStart, `2` Pause.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"intial_state": {
						MarkdownDescription: "Initial state, with Stop support. `0` Start, `1` ForceStart, `2` Pause, `3` Stop.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"host": {
						MarkdownDescription: "host.",
						Computed:            true,
						Type:                types.StringType,
					},
					"api_key": {
						MarkdownDescription: "API key.",
						Computed:            true,
						Type:                types.StringType,
					},
					"rpc_path": {
						MarkdownDescription: "RPC path.",
						Computed:            true,
						Type:                types.StringType,
					},
					"url_base": {
						MarkdownDescription: "Base URL.",
						Computed:            true,
						Type:                types.StringType,
					},
					"secret_token": {
						MarkdownDescription: "Secret token.",
						Computed:            true,
						Type:                types.StringType,
					},
					"username": {
						MarkdownDescription: "Username.",
						Computed:            true,
						Type:                types.StringType,
					},
					"password": {
						MarkdownDescription: "Password.",
						Computed:            true,
						Type:                types.StringType,
					},
					"tv_category": {
						MarkdownDescription: "TV category.",
						Computed:            true,
						Type:                types.StringType,
					},
					"tv_imported_category": {
						MarkdownDescription: "TV imported category.",
						Computed:            true,
						Type:                types.StringType,
					},
					"tv_directory": {
						MarkdownDescription: "TV directory.",
						Computed:            true,
						Type:                types.StringType,
					},
					"destination": {
						MarkdownDescription: "Destination.",
						Computed:            true,
						Type:                types.StringType,
					},
					"category": {
						MarkdownDescription: "Category.",
						Computed:            true,
						Type:                types.StringType,
					},
					"nzb_folder": {
						MarkdownDescription: "NZB folder.",
						Computed:            true,
						Type:                types.StringType,
					},
					"strm_folder": {
						MarkdownDescription: "STRM folder.",
						Computed:            true,
						Type:                types.StringType,
					},
					"torrent_folder": {
						MarkdownDescription: "Torrent folder.",
						Computed:            true,
						Type:                types.StringType,
					},
					"magnet_file_extension": {
						MarkdownDescription: "Magnet file extension.",
						Computed:            true,
						Type:                types.StringType,
					},
					"additional_tags": {
						MarkdownDescription: "Additional tags, `0` TitleSlug, `1` Quality, `2` Language, `3` ReleaseGroup, `4` Year, `5` Indexer, `6` Network.",
						Computed:            true,
						Type: types.SetType{
							ElemType: types.Int64Type,
						},
					},
					"field_tags": {
						MarkdownDescription: "Field tags.",
						Computed:            true,
						Type: types.SetType{
							ElemType: types.StringType,
						},
					},
					"post_im_tags": {
						MarkdownDescription: "Post import tags.",
						Computed:            true,
						Type: types.SetType{
							ElemType: types.StringType,
						},
					},
				}),
			},
		},
	}, nil
}

func (d *DownloadClientsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DownloadClientsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *DownloadClients

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get download clients current value
	response, err := d.client.GetDownloadClientsContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", downloadClientsDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientsDataSourceName)
	// Map response body to resource schema attribute
	profiles := make([]DownloadClient, len(response))
	for i, p := range response {
		profiles[i].write(ctx, p)
	}

	tfsdk.ValueFrom(ctx, profiles, data.DownloadClients.Type(context.Background()), &data.DownloadClients)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
