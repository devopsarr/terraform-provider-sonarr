package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RootFolderDataSource{}

func NewRootFolderDataSource() datasource.DataSource {
	return &RootFolderDataSource{}
}

// RootFolderDataSource defines the root folders implementation.
type RootFolderDataSource struct {
	client *sonarr.Sonarr
}

func (d *RootFolderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_root_folder"
}

func (d *RootFolderDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (d *RootFolderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RootFolderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RootFolder

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get rootfolders current value
	response, err := d.client.GetRootFoldersContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read root folders, got error: %s", err))

		return
	}

	// Map response body to resource schema attribute
	rootFolder, err := findRootFolder(data.Path.Value, response)
	if err != nil {
		resp.Diagnostics.AddError(DataSourceError, fmt.Sprintf("Unable to find root folders, got error: %s", err))

		return
	}

	result := writeRootFolder(ctx, rootFolder)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func findRootFolder(path string, folders []*sonarr.RootFolder) (*sonarr.RootFolder, error) {
	for _, f := range folders {
		if f.Path == path {
			return f, nil
		}
	}

	return nil, fmt.Errorf("no rootfolder with path %s", path)
}
