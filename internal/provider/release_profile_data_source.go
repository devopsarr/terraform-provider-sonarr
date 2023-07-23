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

const releaseProfileDataSourceName = "release_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ReleaseProfileDataSource{}

func NewReleaseProfileDataSource() datasource.DataSource {
	return &ReleaseProfileDataSource{}
}

// ReleaseProfileDataSource defines the release profile implementation.
type ReleaseProfileDataSource struct {
	client *sonarr.APIClient
}

func (d *ReleaseProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + releaseProfileDataSourceName
}

func (d *ReleaseProfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the release server.
		MarkdownDescription: "<!-- subcategory:Profiles -->Single [Release Profile](../resources/release_profile).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Release Profile ID.",
				Required:            true,
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
	}
}

func (d *ReleaseProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *ReleaseProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *ReleaseProfile

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get releaseprofiles current value
	response, _, err := d.client.ReleaseProfileApi.ListReleaseProfile(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, releaseProfileDataSourceName, err))

		return
	}

	data.find(ctx, data.ID.ValueInt64(), response, &resp.Diagnostics)

	tflog.Trace(ctx, "read "+releaseProfileDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *ReleaseProfile) find(ctx context.Context, id int64, profiles []*sonarr.ReleaseProfileResource, diags *diag.Diagnostics) {
	for _, profile := range profiles {
		if int64(profile.GetId()) == id {
			p.write(ctx, profile, diags)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(releaseProfileDataSourceName, "id", strconv.Itoa(int(id))))
}
