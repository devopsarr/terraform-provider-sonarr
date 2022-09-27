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

const namingDataSourceName = "naming"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &NamingDataSource{}

func NewNamingDataSource() datasource.DataSource {
	return &NamingDataSource{}
}

// NamingDataSource defines the naming implementation.
type NamingDataSource struct {
	client *sonarr.Sonarr
}

func (d *NamingDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + namingDataSourceName
}

func (d *NamingDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "[subcategory:Media Management]: #\n[Naming](../resources/naming).",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Delay Profile ID.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"rename_episodes": {
				MarkdownDescription: "Sonarr will use the existing file name if false.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"replace_illegal_characters": {
				MarkdownDescription: "Replace illegal characters. They will be removed if false.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"multi_episode_style": {
				MarkdownDescription: "Multi episode style. 0 - 'Extend' 1 - 'Duplicate' 2 - 'Repeat' 3 - 'Scene' 4 - 'Range' 5 - 'Prefixed Range'.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"daily_episode_format": {
				MarkdownDescription: "Daily episode format.",
				Computed:            true,
				Type:                types.StringType,
			},
			"anime_episode_format": {
				MarkdownDescription: "Anime episode format.",
				Computed:            true,
				Type:                types.StringType,
			},
			"series_folder_format": {
				MarkdownDescription: "Series folder format.",
				Computed:            true,
				Type:                types.StringType,
			},
			"season_folder_format": {
				MarkdownDescription: "Season folder format.",
				Computed:            true,
				Type:                types.StringType,
			},
			"specials_folder_format": {
				MarkdownDescription: "Special folder format.",
				Computed:            true,
				Type:                types.StringType,
			},
			"standard_episode_format": {
				MarkdownDescription: "Standard episode formatss.",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (d *NamingDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NamingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get naming current value
	response, err := d.client.GetNamingContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", namingDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+namingDataSourceName)

	result := writeNaming(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}
