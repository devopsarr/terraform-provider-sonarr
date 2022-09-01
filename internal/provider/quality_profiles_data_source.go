package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ provider.DataSourceType = dataQualityProfilesType{}
	_ datasource.DataSource   = dataQualityProfiles{}
)

type dataQualityProfilesType struct{}

type dataQualityProfiles struct {
	provider sonarrProvider
}

// QualityProfiles is a list of QualityProfile.
type QualityProfiles struct {
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	ID              types.String `tfsdk:"id"`
	QualityProfiles types.Set    `tfsdk:"quality_profiles"`
}

func (t dataQualityProfilesType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the quality server.
		MarkdownDescription: "List all available [Quality Profiles](../resources/quality_profile).",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"quality_profiles": {
				MarkdownDescription: "Quality Profile list.",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						MarkdownDescription: "Quality Profile ID.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"name": {
						MarkdownDescription: "Quality Profile Name.",
						Computed:            true,
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
	diags := req.Config.Get(ctx, &data)
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
	profiles := *writeQualitiyprofiles(ctx, response)
	tfsdk.ValueFrom(ctx, profiles, data.QualityProfiles.Type(context.Background()), &data.QualityProfiles)

	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.String{Value: strconv.Itoa(len(response))}
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func writeQualitiyprofiles(ctx context.Context, qualities []*sonarr.QualityProfile) *[]QualityProfile {
	output := make([]QualityProfile, len(qualities))
	for i, p := range qualities {
		output[i] = *writeQualityProfile(ctx, p)
	}

	return &output
}
