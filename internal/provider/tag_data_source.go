package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

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
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (d *TagDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Single [Tag](../resources/tag).",
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
			UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *TagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data Tag

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get tags current value
	response, err := d.client.GetTagsContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read tags, got error: %s", err))

		return
	}

	tag, err := findTag(data.Label.Value, response)
	if err != nil {
		resp.Diagnostics.AddError(DataSourceError, fmt.Sprintf("Unable to find tags, got error: %s", err))

		return
	}

	result := writeTag(tag)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func findTag(label string, tags []*starr.Tag) (*starr.Tag, error) {
	for _, t := range tags {
		if t.Label == label {
			return t, nil
		}
	}

	return nil, fmt.Errorf("no tag with label %s", label)
}
