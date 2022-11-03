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

const qualityProfileDataSourceName = "quality_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &QualityProfileDataSource{}

func NewQualityProfileDataSource() datasource.DataSource {
	return &QualityProfileDataSource{}
}

// QualityProfileDataSource defines the quality profiles implementation.
type QualityProfileDataSource struct {
	client *sonarr.Sonarr
}

func (d *QualityProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + qualityProfileDataSourceName
}

func (d *QualityProfileDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the quality server.
		MarkdownDescription: "<!-- subcategory:Profiles -->Single [Quality Profile](../resources/quality_profile).",
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

func (d *QualityProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *QualityProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *QualityProfile

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get qualityprofiles current value
	response, err := d.client.GetQualityProfilesContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", qualityProfileDataSourceName, err))

		return
	}

	profile, err := findQualityProfile(data.Name.ValueString(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", qualityProfileDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+qualityProfileDataSourceName)
	data.write(ctx, profile)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findQualityProfile(name string, profiles []*sonarr.QualityProfile) (*sonarr.QualityProfile, error) {
	for _, p := range profiles {
		if p.Name == name {
			return p, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(qualityProfileDataSourceName, "name", name)
}
