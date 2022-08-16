package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.DataSourceType = dataRootFoldersType{}
var _ datasource.DataSource = dataRootFolders{}

type dataRootFoldersType struct{}

type dataRootFolders struct {
	provider sonarrProvider
}

func (t dataRootFoldersType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "List all available rootfolders",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"root_folders": {
				MarkdownDescription: "List of rootfolders",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"path": {
						MarkdownDescription: "Absolute path of rootFolder",
						Computed:            true,
						Type:                types.StringType,
					},
					"accessible": {
						MarkdownDescription: "Access flag",
						Computed:            true,
						Type:                types.BoolType,
					},
					"id": {
						MarkdownDescription: "RootFolder ID",
						Computed:            true,
						Type:                types.Int64Type,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							resource.UseStateForUnknown(),
						},
					},
					"unmapped_folders": {
						MarkdownDescription: "List of folders with no associated series",
						Computed:            true,
						Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
							"path": {
								MarkdownDescription: "Path of unmapped folder",
								Computed:            true,
								Type:                types.StringType,
							},
							"name": {
								MarkdownDescription: "Name of unmapped folder",
								Computed:            true,
								Type:                types.StringType,
							},
						}),
					},
				}),
			},
		},
	}, nil
}

func (t dataRootFoldersType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return dataRootFolders{
		provider: provider,
	}, diags
}

func (d dataRootFolders) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RootFolders
	diags := resp.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get rootfolders current value
	response, err := d.provider.client.GetRootFoldersContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read rootfolders, got error: %s", err))
		return
	}
	// Map response body to resource schema attribute
	data.RootFolders = *writeRootFolders(response)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.String{Value: strconv.Itoa(len(response))}
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func writeRootFolders(folders []*sonarr.RootFolder) *[]RootFolder {
	output := make([]RootFolder, len(folders))
	for i, f := range folders {
		output[i] = *writeRootFolder(f)
	}
	return &output
}
