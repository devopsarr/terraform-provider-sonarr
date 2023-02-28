package provider

import (
	"context"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mitchellh/hashstructure/v2"
)

const (
	autoTagConditionRootFolderDataSourceName = "auto_tag_condition_root_folder"
	autoTagConditionRootFolderImplementation = "RootFolderSpecification"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AutoTagConditionRootFolderDataSource{}

func NewAutoTagConditionRootFolderDataSource() datasource.DataSource {
	return &AutoTagConditionRootFolderDataSource{}
}

// AutoTagConditionRootFolderDataSource defines the auto_tag_condition_root folder implementation.
type AutoTagConditionRootFolderDataSource struct {
	client *sonarr.APIClient
}

func (d *AutoTagConditionRootFolderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + autoTagConditionRootFolderDataSourceName
}

func (d *AutoTagConditionRootFolderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Tags --> Auto Tag Condition Root Folder data source.\nFor more intagion refer to [Auto Tag Conditions](https://wiki.servarr.com/sonarr/settings#conditions).",
		Attributes: map[string]schema.Attribute{
			"negate": schema.BoolAttribute{
				MarkdownDescription: "Negate flag.",
				Required:            true,
			},
			"required": schema.BoolAttribute{
				MarkdownDescription: "Computed flag.",
				Required:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "Implementation.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Specification name.",
				Required:            true,
			},
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.Int64Attribute{
				MarkdownDescription: "Auto tag condition root folder ID.",
				Computed:            true,
			},
			// Field values
			"value": schema.StringAttribute{
				MarkdownDescription: "Root folder path.",
				Required:            true,
			},
		},
	}
}

func (d *AutoTagConditionRootFolderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *AutoTagConditionRootFolderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *AutoTagCondition

	hash, err := hashstructure.Hash(&data, hashstructure.FormatV2, nil)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, helpers.ParseClientError(helpers.Create, autoTagConditionRootFolderDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+autoTagConditionRootFolderDataSourceName)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("implementation"), autoTagConditionRootFolderImplementation)...)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), int64(hash))...)
}
