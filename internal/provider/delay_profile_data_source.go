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

const delayProfileDataSourceName = "delay_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DelayProfileDataSource{}

func NewDelayProfileDataSource() datasource.DataSource {
	return &DelayProfileDataSource{}
}

// DelayProfileDataSource defines the delay profile implementation.
type DelayProfileDataSource struct {
	client *sonarr.Sonarr
}

func (d *DelayProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + delayProfileDataSourceName
}

func (d *DelayProfileDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Profiles -->Single [Delay Profile](../resources/delay_profile).",
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

func (d *DelayProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DelayProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var delayProfile *DelayProfile

	resp.Diagnostics.Append(resp.State.Get(ctx, &delayProfile)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get delayprofiles current value
	response, err := d.client.GetDelayProfilesContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", delayProfileDataSourceName, err))

		return
	}

	profile, err := findDelayProfile(delayProfile.ID.ValueInt64(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", delayProfileDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+delayProfileDataSourceName)
	delayProfile.write(ctx, profile)
	resp.Diagnostics.Append(resp.State.Set(ctx, &delayProfile)...)
}

func findDelayProfile(id int64, profiles []*sonarr.DelayProfile) (*sonarr.DelayProfile, error) {
	for _, p := range profiles {
		if p.ID == id {
			return p, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(delayProfileDataSourceName, "id", strconv.Itoa(int(id)))
}
