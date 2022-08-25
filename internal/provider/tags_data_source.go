package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golift.io/starr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ provider.DataSourceType = dataTagsType{}
var _ datasource.DataSource = dataTags{}

type dataTagsType struct{}

type dataTags struct {
	provider sonarrProvider
}

// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
// Tags is a list of Tag.
type Tags struct {
	ID   types.String `tfsdk:"id"`
	Tags types.Set    `tfsdk:"tags"`
}

func (t dataTagsType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "List all available tags",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"tags": {
				MarkdownDescription: "List of tags",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						MarkdownDescription: "ID of tag",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"label": {
						MarkdownDescription: "Actual tag",
						Computed:            true,
						Type:                types.StringType,
					},
				}),
			},
		},
	}, nil
}

func (t dataTagsType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return dataTags{
		provider: provider,
	}, diags
}

func (d dataTags) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data Tags
	diags := resp.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get tags current value
	response, err := d.provider.client.GetTagsContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read tags, got error: %s", err))

		return
	}

	// Map response body to resource schema attribute
	tags := *writeTags(response)
	tfsdk.ValueFrom(ctx, tags, data.Tags.Type(context.Background()), &data.Tags)

	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.String{Value: strconv.Itoa(len(response))}
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func writeTags(tags []*starr.Tag) *[]Tag {
	output := make([]Tag, len(tags))
	for i, t := range tags {
		output[i] = *writeTag(t)
	}

	return &output
}
