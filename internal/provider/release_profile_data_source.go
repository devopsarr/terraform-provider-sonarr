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

const releaseProfileDataSourceName = "release_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ReleaseProfileDataSource{}

func NewReleaseProfileDataSource() datasource.DataSource {
	return &ReleaseProfileDataSource{}
}

// ReleaseProfileDataSource defines the release profile implementation.
type ReleaseProfileDataSource struct {
	client *sonarr.Sonarr
}

func (d *ReleaseProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + releaseProfileDataSourceName
}

func (d *ReleaseProfileDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the release server.
		MarkdownDescription: "<!-- subcategory:Profiles -->Single [Release Profile](../resources/release_profile).",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Release Profile ID.",
				Required:            true,
				Type:                types.Int64Type,
			},
			"enabled": {
				MarkdownDescription: "Enabled",
				Computed:            true,
				Type:                types.BoolType,
			},
			"name": {
				MarkdownDescription: "Release profile name.",
				Computed:            true,
				Type:                types.StringType,
			},
			"indexer_id": {
				MarkdownDescription: "Indexer ID. Set `0` for all.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"required": {
				MarkdownDescription: "Required terms.",
				Computed:            true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
			},
			"ignored": {
				MarkdownDescription: "Ignored terms.",
				Computed:            true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
			},
			"tags": {
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
		},
	}, nil
}

func (d *ReleaseProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ReleaseProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var releaseProfile *ReleaseProfile

	resp.Diagnostics.Append(resp.State.Get(ctx, &releaseProfile)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get releaseprofiles current value
	response, err := d.client.GetReleaseProfilesContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", releaseProfileDataSourceName, err))

		return
	}

	profile, err := findReleaseProfile(releaseProfile.ID.ValueInt64(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", releaseProfileDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+releaseProfileDataSourceName)
	releaseProfile.write(ctx, profile)
	resp.Diagnostics.Append(resp.State.Set(ctx, &releaseProfile)...)
}

func findReleaseProfile(id int64, profiles []*sonarr.ReleaseProfile) (*sonarr.ReleaseProfile, error) {
	for _, p := range profiles {
		if p.ID == id {
			return p, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(releaseProfileDataSourceName, "id", strconv.Itoa(int(id)))
}
