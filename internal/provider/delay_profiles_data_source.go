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

const delayProfilesDataSourceName = "delay_profiles"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DelayProfilesDataSource{}

func NewDelayProfilesDataSource() datasource.DataSource {
	return &DelayProfilesDataSource{}
}

// DelayProfilesDataSource defines the delay profiles implementation.
type DelayProfilesDataSource struct {
	client *sonarr.Sonarr
}

// DelayProfiles describes the delay profiles data model.
type DelayProfiles struct {
	DelayProfiles types.Set    `tfsdk:"delay_profiles"`
	ID            types.String `tfsdk:"id"`
}

func (d *DelayProfilesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + delayProfilesDataSourceName
}

func (d *DelayProfilesDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "[subcategory:Profiles]: #\nList all available [Delay Profiles](../resources/delay_profile).",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"delay_profiles": {
				MarkdownDescription: "Delay Profile list.",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						MarkdownDescription: "Delay Profile ID.",
						Computed:            true,
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
				}),
			},
		},
	}, nil
}

func (d *DelayProfilesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DelayProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *DelayProfiles

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get delayprofiles current value
	response, err := d.client.GetDelayProfilesContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", delayProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+delayProfileResourceName)
	// Map response body to resource schema attribute
	profiles := make([]DelayProfile, len(response))
	for i, p := range response {
		profiles[i].write(ctx, p)
	}

	tfsdk.ValueFrom(ctx, profiles, data.DelayProfiles.Type(context.Background()), &data.DelayProfiles)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.String{Value: strconv.Itoa(len(response))}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
