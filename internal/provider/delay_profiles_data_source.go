package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golift.io/starr/sonarr"
)

type dataDelayProfilesType struct{}

type dataDelayProfiles struct {
	provider provider
}

func (t dataDelayProfilesType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "List all available delayprofiles",
		Attributes: map[string]tfsdk.Attribute{
			//TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"delay_profiles": {
				MarkdownDescription: "List of delayprofiles",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						MarkdownDescription: "ID of delayprofile",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"enable_usenet": {
						MarkdownDescription: "Usenet allowed Flag",
						Computed:            true,
						Type:                types.BoolType,
					},
					"enable_torrent": {
						MarkdownDescription: "Torrent allowed Flag",
						Computed:            true,
						Type:                types.BoolType,
					},
					"bypass_if_highest_quality": {
						MarkdownDescription: "Bypass for highest quality Flag",
						Computed:            true,
						Type:                types.BoolType,
					},
					"usenet_delay": {
						MarkdownDescription: "Usenet delay",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"torrent_delay": {
						MarkdownDescription: "Torrent Delay",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"order": {
						MarkdownDescription: "Order",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"tags": {
						MarkdownDescription: "List of associated tags",
						Computed:            true,
						Type: types.ListType{
							ElemType: types.Int64Type,
						},
					},
					//TODO: add validation
					"preferred_protocol": {
						MarkdownDescription: "Preferred protocol",
						Computed:            true,
						Type:                types.StringType,
					},
				}, tfsdk.SetNestedAttributesOptions{}),
			},
		},
	}, nil
}

func (t dataDelayProfilesType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return dataDelayProfiles{
		provider: provider,
	}, diags
}

func (d dataDelayProfiles) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data DelayProfiles
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
	// Map response body to resource schema attribute
	data.DelayProfiles = *writeDelayprofiles(response)
	//TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.String{Value: strconv.Itoa(len(response))}
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func writeDelayprofiles(delays []*sonarr.DelayProfile) *[]DelayProfile {
	output := make([]DelayProfile, len(delays))
	for i, p := range delays {
		output[i] = *writeDelayProfile(p)
	}
	return &output
}
