package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const delayProfileDataSourceName = "delay_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DelayProfileDataSource{}

func NewDelayProfileDataSource() datasource.DataSource {
	return &DelayProfileDataSource{}
}

// DelayProfileDataSource defines the delay profile implementation.
type DelayProfileDataSource struct {
	client *sonarr.APIClient
	auth   context.Context
}

func (d *DelayProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + delayProfileDataSourceName
}

func (d *DelayProfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Profiles -->\nSingle [Delay Profile](../resources/delay_profile).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Delay Profile ID.",
				Required:            true,
			},
			"enable_usenet": schema.BoolAttribute{
				MarkdownDescription: "Usenet allowed Flag.",
				Computed:            true,
			},
			"enable_torrent": schema.BoolAttribute{
				MarkdownDescription: "Torrent allowed Flag.",
				Computed:            true,
			},
			"bypass_if_highest_quality": schema.BoolAttribute{
				MarkdownDescription: "Bypass for highest quality Flag.",
				Computed:            true,
			},
			"bypass_if_above_custom_format_score": schema.BoolAttribute{
				MarkdownDescription: "Bypass for higher custom format score flag.",
				Computed:            true,
			},
			"usenet_delay": schema.Int64Attribute{
				MarkdownDescription: "Usenet delay.",
				Computed:            true,
			},
			"torrent_delay": schema.Int64Attribute{
				MarkdownDescription: "Torrent Delay.",
				Computed:            true,
			},
			"order": schema.Int64Attribute{
				MarkdownDescription: "Order.",
				Computed:            true,
			},
			"minimum_custom_format_score": schema.Int64Attribute{
				MarkdownDescription: "Minimum custom format score.",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"preferred_protocol": schema.StringAttribute{
				MarkdownDescription: "Preferred protocol.",
				Computed:            true,
			},
		},
	}
}

func (d *DelayProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *DelayProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *DelayProfile

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get delayprofiles current value
	response, _, err := d.client.DelayProfileAPI.ListDelayProfile(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, delayProfileDataSourceName, err))

		return
	}

	data.find(ctx, data.ID.ValueInt64(), response, &resp.Diagnostics)

	tflog.Trace(ctx, "read "+delayProfileDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *DelayProfile) find(ctx context.Context, id int64, profiles []sonarr.DelayProfileResource, diags *diag.Diagnostics) {
	for _, profile := range profiles {
		if int64(profile.GetId()) == id {
			p.write(ctx, &profile, diags)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(delayProfileDataSourceName, "id", strconv.Itoa(int(id))))
}
