package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ provider.DataSourceType = dataDelayProfileType{}
	_ datasource.DataSource   = dataDelayProfile{}
)

type dataDelayProfileType struct{}

type dataDelayProfile struct {
	provider sonarrProvider
}

func (t dataDelayProfileType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "Single [Delay Profile](../resources/delay_profile).",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Delay Profile ID.",
				Required:            true,
				Type:                types.Int64Type,
			},
			"enable_usenet": {
				MarkdownDescription: "Usenet allowed Flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"enable_torrent": {
				MarkdownDescription: "Torrent allowed Flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"bypass_if_highest_quality": {
				MarkdownDescription: "Bypass for highest quality Flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"usenet_delay": {
				MarkdownDescription: "Usenet delay.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"torrent_delay": {
				MarkdownDescription: "Torrent Delay.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"order": {
				MarkdownDescription: "Order.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"tags": {
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"preferred_protocol": {
				MarkdownDescription: "Preferred protocol.",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (t dataDelayProfileType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return dataDelayProfile{
		provider: provider,
	}, diags
}

func (d dataDelayProfile) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DelayProfile
	diags := resp.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get delayprofiles current value
	response, err := d.provider.client.GetDelayProfilesContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read delayprofiles, got error: %s", err))

		return
	}

	profile, err := findDelayProfile(data.ID.Value, response)
	if err != nil {
		resp.Diagnostics.AddError("Data Source Error", fmt.Sprintf("Unable to find delayprofile, got error: %s", err))

		return
	}

	result := writeDelayProfile(ctx, profile)
	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
}

func findDelayProfile(id int64, profiles []*sonarr.DelayProfile) (*sonarr.DelayProfile, error) {
	for _, p := range profiles {
		if p.ID == id {
			return p, nil
		}
	}

	return nil, fmt.Errorf("no delay profile with id %d", id)
}
