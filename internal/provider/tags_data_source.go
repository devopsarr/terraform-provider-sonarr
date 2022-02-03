package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dataTagsType struct{}

type dataTags struct {
	provider provider
}

func (t dataTagsType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "List all available tags",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"tags": {
				MarkdownDescription: "List of tags",
				Computed:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
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
				}, tfsdk.ListNestedAttributesOptions{}),
			},
		},
	}, nil
}

func (t dataTagsType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return dataTags{
		provider: provider,
	}, diags
}

func (d dataTags) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {

	var data Tags
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get order current value
	tags, err := d.provider.client.GetTags()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read tags, got error: %s", err))
		return
	}
	// Map response body to resource schema attribute
	for _, t := range tags {
		data.Tags = append(data.Tags, Tag{
			ID:    types.Int64{Value: int64(t.ID)},
			Label: types.String{Value: t.Label},
		})
	}
	//TODO: change id to something meaningful
	data.ID = types.String{Value: "test"}
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
