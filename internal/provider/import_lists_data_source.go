package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const importListsDataSourceName = "import_lists"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ImportListsDataSource{}

func NewImportListsDataSource() datasource.DataSource {
	return &ImportListsDataSource{}
}

// ImportListsDataSource defines the import lists implementation.
type ImportListsDataSource struct {
	client *sonarr.Sonarr
}

// ImportLists describes the import lists data model.
type ImportLists struct {
	ImportLists types.Set    `tfsdk:"import_lists"`
	ID          types.String `tfsdk:"id"`
}

func (d *ImportListsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListsDataSourceName
}

func (d *ImportListsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Import Lists -->List all available [Import Lists](../resources/import_list).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"import_lists": schema.SetNestedAttribute{
				MarkdownDescription: "Import List list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"enable_automatic_add": schema.BoolAttribute{
							MarkdownDescription: "Enable automatic add flag.",
							Computed:            true,
						},
						"season_folder": schema.BoolAttribute{
							MarkdownDescription: "Season folder flag.",
							Computed:            true,
						},
						"language_profile_id": schema.Int64Attribute{
							MarkdownDescription: "Language profile ID.",
							Computed:            true,
						},
						"quality_profile_id": schema.Int64Attribute{
							MarkdownDescription: "Quality profile ID.",
							Computed:            true,
						},
						"should_monitor": schema.StringAttribute{
							MarkdownDescription: "Should monitor.",
							Computed:            true,
						},
						"root_folder_path": schema.StringAttribute{
							MarkdownDescription: "Root folder path.",
							Computed:            true,
						},
						"implementation": schema.StringAttribute{
							MarkdownDescription: "ImportList implementation name.",
							Computed:            true,
						},
						"config_contract": schema.StringAttribute{
							MarkdownDescription: "ImportList configuration template.",
							Computed:            true,
						},
						"series_type": schema.StringAttribute{
							MarkdownDescription: "Series type.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Import List name.",
							Computed:            true,
						},
						"tags": schema.SetAttribute{
							MarkdownDescription: "List of associated tags.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
						"id": schema.Int64Attribute{
							MarkdownDescription: "Import List ID.",
							Computed:            true,
						},
						// Field values
						"limit": schema.Int64Attribute{
							MarkdownDescription: "Limit.",
							Computed:            true,
						},
						"trakt_list_type": schema.Int64Attribute{
							MarkdownDescription: "Trakt list type.",
							Computed:            true,
						},
						"access_token": schema.StringAttribute{
							MarkdownDescription: "Access token.",
							Computed:            true,
							Sensitive:           true,
						},
						"refresh_token": schema.StringAttribute{
							MarkdownDescription: "Refresh token.",
							Computed:            true,
							Sensitive:           true,
						},
						"api_key": schema.StringAttribute{
							MarkdownDescription: "API key.",
							Computed:            true,
							Sensitive:           true,
						},
						"auth_user": schema.StringAttribute{
							MarkdownDescription: "Auth User.",
							Computed:            true,
						},
						"username": schema.StringAttribute{
							MarkdownDescription: "Username.",
							Computed:            true,
						},
						"rating": schema.StringAttribute{
							MarkdownDescription: "Rating.",
							Computed:            true,
						},
						"base_url": schema.StringAttribute{
							MarkdownDescription: "Base URL.",
							Computed:            true,
						},
						"expires": schema.StringAttribute{
							MarkdownDescription: "Expires.",
							Computed:            true,
						},
						"listname": schema.StringAttribute{
							MarkdownDescription: "Expires.",
							Computed:            true,
						},
						"genres": schema.StringAttribute{
							MarkdownDescription: "Expires.",
							Computed:            true,
						},
						"years": schema.StringAttribute{
							MarkdownDescription: "Expires.",
							Computed:            true,
						},
						"trakt_additional_parameters": schema.StringAttribute{
							MarkdownDescription: "Trakt additional parameters.",
							Computed:            true,
						},
						"language_profile_ids": schema.SetAttribute{
							MarkdownDescription: "Language profile IDs.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
						"quality_profile_ids": schema.SetAttribute{
							MarkdownDescription: "Quality profile IDs.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
						"tag_ids": schema.SetAttribute{
							MarkdownDescription: "Tag IDs.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
					},
				},
			},
		},
	}
}

func (d *ImportListsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ImportListsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *ImportLists

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get import lists current value
	response, err := d.client.GetImportListsContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", importListsDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListsDataSourceName)
	// Map response body to resource schema attribute
	importLists := make([]ImportList, len(response))
	for i, d := range response {
		importLists[i].Tags = types.SetNull(types.Int64Type)
		importLists[i].write(ctx, d)
	}

	tfsdk.ValueFrom(ctx, importLists, data.ImportLists.Type(context.Background()), &data.ImportLists)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
