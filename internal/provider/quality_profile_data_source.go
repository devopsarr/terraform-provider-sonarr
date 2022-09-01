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
	_ provider.DataSourceType = dataQualityProfileType{}
	_ datasource.DataSource   = dataQualityProfile{}
)

type dataQualityProfileType struct{}

type dataQualityProfile struct {
	provider sonarrProvider
}

func (t dataQualityProfileType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the quality server.
		MarkdownDescription: "Single [Quality Profile](../resources/quality_profile).",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Quality Profile ID.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"name": {
				MarkdownDescription: "Quality Profile Name.",
				Required:            true,
				Type:                types.StringType,
			},
			"upgrade_allowed": {
				MarkdownDescription: "Upgrade allowed flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"cutoff": {
				MarkdownDescription: "Quality ID to which cutoff.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"quality_groups": {
				MarkdownDescription: "Quality groups.",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						MarkdownDescription: "Quality group ID.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"name": {
						MarkdownDescription: "Quality group name.",
						Computed:            true,
						Type:                types.StringType,
					},
					"qualities": {
						MarkdownDescription: "Qualities in group.",
						Required:            true,
						Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
							"id": {
								MarkdownDescription: "Quality ID.",
								Computed:            true,
								Type:                types.Int64Type,
							},
							"resolution": {
								MarkdownDescription: "Resolution.",
								Computed:            true,
								Type:                types.Int64Type,
							},
							"name": {
								MarkdownDescription: "Quality name.",
								Computed:            true,
								Type:                types.StringType,
							},
							"source": {
								MarkdownDescription: "Source.",
								Computed:            true,
								Type:                types.StringType,
							},
						}),
					},
				}),
			},
		},
	}, nil
}

func (t dataQualityProfileType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return dataQualityProfile{
		provider: provider,
	}, diags
}

func (d dataQualityProfile) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data QualityProfile
	diags := resp.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get qualityprofiles current value
	response, err := d.provider.client.GetQualityProfilesContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read qualityprofiles, got error: %s", err))

		return
	}

	profile, err := findQualityProfile(data.Name.Value, response)
	if err != nil {
		resp.Diagnostics.AddError("Data Source Error", fmt.Sprintf("Unable to find qualityprofile, got error: %s", err))

		return
	}

	result := writeQualityProfile(ctx, profile)
	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
}

func findQualityProfile(name string, profiles []*sonarr.QualityProfile) (*sonarr.QualityProfile, error) {
	for _, p := range profiles {
		if p.Name == name {
			return p, nil
		}
	}

	return nil, fmt.Errorf("no quality profile with name %s", name)
}
