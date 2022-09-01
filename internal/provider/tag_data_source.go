package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golift.io/starr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ provider.DataSourceType = dataTagType{}
	_ datasource.DataSource   = dataTag{}
)

type dataTagType struct{}

type dataTag struct {
	provider sonarrProvider
}

func (t dataTagType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (t dataTagType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return dataTag{
		provider: provider,
	}, diags
}

func (d dataTag) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data Tag
	diags := req.Config.Get(ctx, &data)
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

	tag, err := findTag(data.Label.Value, response)
	if err != nil {
		resp.Diagnostics.AddError("Data Source Error", fmt.Sprintf("Unable to find tags, got error: %s", err))

		return
	}

	result := writeTag(tag)
	// Map response body to resource schema attribute
	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
}

func findTag(label string, tags []*starr.Tag) (*starr.Tag, error) {
	for _, t := range tags {
		if t.Label == label {
			return t, nil
		}
	}

	return nil, fmt.Errorf("no tag with label %s", label)
}
