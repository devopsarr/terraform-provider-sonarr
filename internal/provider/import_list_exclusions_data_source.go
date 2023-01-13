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

const importListExclusionsDataSourceName = "import_list_exclusions"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ImportListExclusionsDataSource{}

func NewImportListExclusionsDataSource() datasource.DataSource {
	return &ImportListExclusionsDataSource{}
}

// ImportListExclusionsDataSource defines the importListExclusions implementation.
type ImportListExclusionsDataSource struct {
	client *sonarr.APIClient
}

// ImportListExclusions describes the importListExclusions data model.
type ImportListExclusions struct {
	ImportListExclusions types.Set    `tfsdk:"import_list_exclusions"`
	ID                   types.String `tfsdk:"id"`
}

func (d *ImportListExclusionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListExclusionsDataSourceName
}

func (d *ImportListExclusionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->List all available [ImportListExclusions](../resources/importListExclusion).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"import_list_exclusions": schema.SetNestedAttribute{
				MarkdownDescription: "ImportListExclusion list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tvdb_id": schema.Int64Attribute{
							MarkdownDescription: "Series TVDB ID.",
							Computed:            true,
						},
						"title": schema.StringAttribute{
							MarkdownDescription: "Series to be excluded.",
							Computed:            true,
						},
						"id": schema.Int64Attribute{
							MarkdownDescription: "ImportListExclusion ID.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *ImportListExclusionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ImportListExclusionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *ImportListExclusions

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get importListExclusions current value
	response, _, err := d.client.ImportListExclusionApi.ListImportListExclusion(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListExclusionsDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListExclusionsDataSourceName)
	// Map response body to resource schema attribute
	importListExclusions := make([]ImportListExclusion, len(response))
	for i, t := range response {
		importListExclusions[i].write(t)
	}

	tfsdk.ValueFrom(ctx, importListExclusions, data.ImportListExclusions.Type(ctx), &data.ImportListExclusions)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
