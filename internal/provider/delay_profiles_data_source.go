package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const delayProfilesDataSourceName = "delay_profiles"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DelayProfilesDataSource{}

func NewDelayProfilesDataSource() datasource.DataSource {
	return &DelayProfilesDataSource{}
}

// DelayProfilesDataSource defines the delay profiles implementation.
type DelayProfilesDataSource struct {
	client *sonarr.APIClient
}

// DelayProfiles describes the delay profiles data model.
type DelayProfiles struct {
	DelayProfiles types.Set    `tfsdk:"delay_profiles"`
	ID            types.String `tfsdk:"id"`
}

func (d *DelayProfilesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + delayProfilesDataSourceName
}

func (d *DelayProfilesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Profiles -->List all available [Delay Profiles](../resources/delay_profile).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"delay_profiles": schema.SetNestedAttribute{
				MarkdownDescription: "Delay Profile list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Delay Profile ID.",
							Computed:            true,
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
				},
			},
		},
	}
}

func (d *DelayProfilesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			helpers.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *sonarr.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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
	response, _, err := d.client.DelayProfileApi.ListDelayProfile(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, delayProfileResourceName, err))
		return
	}

	tflog.Trace(ctx, "read "+delayProfileResourceName)
	// Map response body to resource schema attribute
	profiles := make([]DelayProfile, len(response))
	for i, p := range response {
		profiles[i].write(ctx, p)
	}

	tfsdk.ValueFrom(ctx, profiles, data.DelayProfiles.Type(ctx), &data.DelayProfiles)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
