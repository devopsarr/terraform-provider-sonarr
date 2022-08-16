package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.DataSourceType = dataQualityProfilesType{}
var _ datasource.DataSource = dataQualityProfiles{}

type dataQualityProfilesType struct{}

type dataQualityProfiles struct {
	provider sonarrProvider
}

func (t dataQualityProfilesType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the quality server.
		MarkdownDescription: "List all available qualityprofiles",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"quality_profiles": {
				MarkdownDescription: "List of qualityprofiles",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						MarkdownDescription: "ID of qualityprofile",
						Computed:            true,
						Type:                types.Int64Type,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							resource.UseStateForUnknown(),
						},
					},
					"name": {
						MarkdownDescription: "Name",
						Required:            true,
						Type:                types.StringType,
					},
					"upgrade_allowed": {
						MarkdownDescription: "Upgrade allowed flag",
						Optional:            true,
						Computed:            true,
						Type:                types.BoolType,
					},
					"cutoff": {
						MarkdownDescription: "Quality ID to which cutoff",
						Optional:            true,
						Computed:            true,
						Type:                types.Int64Type,
					},
					"quality_groups": {
						MarkdownDescription: "Quality groups",
						Required:            true,
						Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
							"id": {
								MarkdownDescription: "ID of quality group",
								Optional:            true,
								Computed:            true,
								Type:                types.Int64Type,
							},
							"name": {
								MarkdownDescription: "Name of quality group",
								Optional:            true,
								Computed:            true,
								Type:                types.StringType,
							},
							"qualities": {
								MarkdownDescription: "Qualities in group",
								Required:            true,
								Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
									"id": {
										MarkdownDescription: "ID of quality group",
										Optional:            true,
										Computed:            true,
										Type:                types.Int64Type,
									},
									"resolution": {
										MarkdownDescription: "Resolution",
										Optional:            true,
										Computed:            true,
										Type:                types.Int64Type,
									},
									"name": {
										MarkdownDescription: "Name of quality group",
										Optional:            true,
										Computed:            true,
										Type:                types.StringType,
									},
									"source": {
										MarkdownDescription: "Source",
										Optional:            true,
										Computed:            true,
										Type:                types.StringType,
									},
								}),
							},
						}),
					},
				}),
			},
		},
	}, nil
}

func (t dataQualityProfilesType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return dataQualityProfiles{
		provider: provider,
	}, diags
}

func (d dataQualityProfiles) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data QualityProfiles
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

	// Map response body to resource schema attribute
	data.QualityProfiles = *writeQualitiyprofiles(response)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.String{Value: strconv.Itoa(len(response))}
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func writeQualitiyprofiles(qualities []*sonarr.QualityProfile) *[]QualityProfile {
	output := make([]QualityProfile, len(qualities))
	for i, p := range qualities {
		output[i] = *writeQualityProfile(p)
	}
	return &output
}
