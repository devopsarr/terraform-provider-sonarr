package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

const tagDataSourceName = "tag"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TagDataSource{}

func NewTagDataSource() datasource.DataSource {
	return &TagDataSource{}
}

// TagDataSource defines the tag implementation.
type TagDataSource struct {
	client *sonarr.Sonarr
}

func (d *TagDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + tagDataSourceName
}

func (d *TagDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "[subcategory:Tags]: #\nSingle [Tag](../resources/tag).",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Tag ID.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"label": {
				MarkdownDescription: "Tag label.",
				Required:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (d *TagDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tag *Tag

	resp.Diagnostics.Append(req.Config.Get(ctx, &tag)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get tags current value
	response, err := d.client.GetTagsContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.UnableToRead(tagDataSourceName, err))

		return
	}

	value, err := findTag(tag.Label.ValueString(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", tagDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+tagDataSourceName)
	tag.write(value)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &tag)...)
}

func findTag(label string, tags []*starr.Tag) (*starr.Tag, error) {
	for _, t := range tags {
		if t.Label == label {
			return t, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(tagDataSourceName, "label", label)
}
