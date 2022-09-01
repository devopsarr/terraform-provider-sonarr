package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ provider.DataSourceType = dataNamingType{}
	_ datasource.DataSource   = dataNaming{}
)

type dataNamingType struct{}

type dataNaming struct {
	provider sonarrProvider
}

func (t dataNamingType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "[Naming](../resources/naming).",
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

func (t dataNamingType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return dataNaming{
		provider: provider,
	}, diags
}

func (d dataNaming) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get naming current value
	response, err := d.provider.client.GetNamingContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read naming, got error: %s", err))

		return
	}

	result := writeNaming(response)
	diags := resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
}
