package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const qualityProfilesDataSourceName = "quality_profiles"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &QualityProfilesDataSource{}

func NewQualityProfilesDataSource() datasource.DataSource {
	return &QualityProfilesDataSource{}
}

// QualityProfilesDataSource defines the qyality profiles implementation.
type QualityProfilesDataSource struct {
	client *sonarr.Sonarr
}

// QualityProfiles describes the qyality profiles data model.
type QualityProfiles struct {
	QualityProfiles types.Set    `tfsdk:"quality_profiles"`
	ID              types.String `tfsdk:"id"`
}

func (d *QualityProfilesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + qualityProfilesDataSourceName
}

func (d *QualityProfilesDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the quality server.
		MarkdownDescription: "<!-- subcategory:Profiles -->List all available [Quality Profiles](../resources/quality_profile).",
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

func (d *QualityProfilesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *QualityProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *QualityProfiles

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get qualityprofiles current value
	response, err := d.client.GetQualityProfilesContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", qualityProfilesDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+qualityProfilesDataSourceName)
	// Map response body to resource schema attribute
	profiles := *writeQualitiyprofiles(ctx, response)
	tfsdk.ValueFrom(ctx, profiles, data.QualityProfiles.Type(context.Background()), &data.QualityProfiles)

	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func writeQualitiyprofiles(ctx context.Context, qualities []*sonarr.QualityProfile) *[]QualityProfile {
	output := make([]QualityProfile, len(qualities))
	for i, p := range qualities {
		output[i].write(ctx, p)
	}

	return &output
}
