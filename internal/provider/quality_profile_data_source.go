package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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

func (d *QualityProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the quality server.
		MarkdownDescription: "<!-- subcategory:Profiles -->Single [Quality Profile](../resources/quality_profile).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Quality Profile ID.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Quality Profile Name.",
				Required:            true,
			},
			"upgrade_allowed": schema.BoolAttribute{
				MarkdownDescription: "Upgrade allowed flag.",
				Computed:            true,
			},
			"cutoff": schema.Int64Attribute{
				MarkdownDescription: "Quality ID to which cutoff.",
				Computed:            true,
			},
			"quality_groups": schema.SetNestedAttribute{
				MarkdownDescription: "Quality groups.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Quality group ID.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Quality group name.",
							Computed:            true,
						},
						"qualities": schema.SetNestedAttribute{
							MarkdownDescription: "Qualities in group.",
							Required:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.Int64Attribute{
										MarkdownDescription: "Quality ID.",
										Computed:            true,
									},
									"resolution": schema.Int64Attribute{
										MarkdownDescription: "Resolution.",
										Computed:            true,
									},
									"name": schema.StringAttribute{
										MarkdownDescription: "Quality name.",
										Computed:            true,
									},
									"source": schema.StringAttribute{
										MarkdownDescription: "Source.",
										Computed:            true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *QualityProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedDataSourceConfigureType,
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
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", qualityProfileDataSourceName, err))

		return
	}

	profile, err := findQualityProfile(data.Name.ValueString(), response)
	if err != nil {
		resp.Diagnostics.AddError(tools.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", qualityProfileDataSourceName, err))

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

	return nil, tools.ErrDataNotFoundError(qualityProfileDataSourceName, "name", name)
}
