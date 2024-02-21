package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const releaseProfilesDataSourceName = "release_profiles"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ReleaseProfilesDataSource{}

func NewReleaseProfilesDataSource() datasource.DataSource {
	return &ReleaseProfilesDataSource{}
}

// ReleaseProfilesDataSource defines the release profiles implementation.
type ReleaseProfilesDataSource struct {
	client *sonarr.APIClient
	auth   context.Context
}

// ReleaseProfiles describes the release profiles data model.
type ReleaseProfiles struct {
	ReleaseProfiles types.Set    `tfsdk:"release_profiles"`
	ID              types.String `tfsdk:"id"`
}

func (d *ReleaseProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + releaseProfilesDataSourceName
}

func (d *ReleaseProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the release server.
		MarkdownDescription: "<!-- subcategory:Profiles -->List all available [Release Profiles](../resources/release_profile).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"release_profiles": schema.SetNestedAttribute{
				MarkdownDescription: "Release Profile list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Release Profile ID.",
							Computed:            true,
						},
						"enabled": schema.BoolAttribute{
							MarkdownDescription: "Enabled",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Release profile name.",
							Computed:            true,
						},
						"indexer_id": schema.Int64Attribute{
							MarkdownDescription: "Indexer ID. Set `0` for all.",
							Computed:            true,
						},
						"required": schema.SetAttribute{
							MarkdownDescription: "Required terms.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"ignored": schema.SetAttribute{
							MarkdownDescription: "Ignored terms.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"tags": schema.SetAttribute{
							MarkdownDescription: "List of associated tags.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
					},
				},
			},
		},
	}
}

func (d *ReleaseProfilesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *ReleaseProfilesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get releaseprofiles current value
	response, _, err := d.client.ReleaseProfileAPI.ListReleaseProfile(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, releaseProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+releaseProfileResourceName)
	// Map response body to resource schema attribute
	profiles := make([]ReleaseProfile, len(response))
	for i, p := range response {
		profiles[i].write(ctx, &p, &resp.Diagnostics)
	}

	profileList, diags := types.SetValueFrom(ctx, ReleaseProfile{}.getType(), profiles)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, ReleaseProfiles{ReleaseProfiles: profileList, ID: types.StringValue(strconv.Itoa(len(response)))})...)
}
