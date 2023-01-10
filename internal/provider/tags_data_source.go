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

const tagsDataSourceName = "tags"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TagsDataSource{}

func NewTagsDataSource() datasource.DataSource {
	return &TagsDataSource{}
}

// TagsDataSource defines the tags implementation.
type TagsDataSource struct {
	client *sonarr.APIClient
}

// Tags describes the tags data model.
type Tags struct {
	Tags types.Set    `tfsdk:"tags"`
	ID   types.String `tfsdk:"id"`
}

func (d *TagsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + tagsDataSourceName
}

func (d *TagsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "<!-- subcategory:Tags -->List all available [Tags](../resources/tag).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "Tag list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Tag ID.",
							Computed:            true,
						},
						"label": schema.StringAttribute{
							MarkdownDescription: "Tag label.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *TagsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TagsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Tags

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get tags current value
	response, _, err := d.client.TagApi.ListTag(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, tagsDataSourceName, err))
		return
	}

	tflog.Trace(ctx, "read "+tagsDataSourceName)
	// Map response body to resource schema attribute
	tags := make([]Tag, len(response))
	for i, t := range response {
		tags[i].write(t)
	}

	tfsdk.ValueFrom(ctx, tags, data.Tags.Type(ctx), &data.Tags)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
