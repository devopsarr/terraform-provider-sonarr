package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ provider.DataSourceType = dataRootFolderType{}
var _ datasource.DataSource = dataRootFolder{}

type dataRootFolderType struct{}

type dataRootFolder struct {
	provider sonarrProvider
}

func (t dataRootFolderType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "Single [Root Folder](../resources/root_folder).",
		Attributes: map[string]tfsdk.Attribute{
			"path": {
				MarkdownDescription: "Root Folder absolute path.",
				Required:            true,
				Type:                types.StringType,
			},
			"accessible": {
				MarkdownDescription: "Access flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"id": {
				MarkdownDescription: "Root Folder ID.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"unmapped_folders": {
				MarkdownDescription: "List of folders with no associated series.",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"path": {
						MarkdownDescription: "Path of unmapped folder.",
						Computed:            true,
						Type:                types.StringType,
					},
					"name": {
						MarkdownDescription: "Name of unmapped folder.",
						Computed:            true,
						Type:                types.StringType,
					},
				}),
			},
		},
	}, nil
}

func (t dataRootFolderType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return dataRootFolder{
		provider: provider,
	}, diags
}

func (d dataRootFolder) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RootFolder
	diags := req.Config.Get(ctx, &data)
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
	rootFolder, err := findRootFolder(data.Path.Value, response)
	if err != nil {
		resp.Diagnostics.AddError("Data Source Error", fmt.Sprintf("Unable to find tags, got error: %s", err))

		return
	}

	result := writeRootFolder(ctx, rootFolder)
	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
}

func findRootFolder(path string, folders []*sonarr.RootFolder) (*sonarr.RootFolder, error) {
	for _, f := range folders {
		if f.Path == path {
			return f, nil
		}
	}

	return nil, fmt.Errorf("no rootfolder with path %s", path)
}
