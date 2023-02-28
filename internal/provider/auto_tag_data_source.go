package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const autoTagDataSourceName = "auto_tag"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AutoTagDataSource{}

func NewAutoTagDataSource() datasource.DataSource {
	return &AutoTagDataSource{}
}

// AutoTagDataSource defines the auto_tag implementation.
type AutoTagDataSource struct {
	client *sonarr.APIClient
}

func (d *AutoTagDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + autoTagDataSourceName
}

func (d *AutoTagDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Profiles -->Single [Auto Tag](../resources/auto_tag).",
		Attributes: map[string]schema.Attribute{
			"remove_tags_automatically": schema.BoolAttribute{
				MarkdownDescription: "Remove tags automatically flag.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Auto Tag name.",
				Required:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Auto Tag ID.",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"specifications": schema.SetNestedAttribute{
				MarkdownDescription: "Specifications.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"negate": schema.BoolAttribute{
							MarkdownDescription: "Negate flag.",
							Computed:            true,
						},
						"required": schema.BoolAttribute{
							MarkdownDescription: "Computed flag.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Specification name.",
							Computed:            true,
						},
						"implementation": schema.StringAttribute{
							MarkdownDescription: "Implementation.",
							Computed:            true,
						},
						// Field values
						"value": schema.StringAttribute{
							MarkdownDescription: "Value.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *AutoTagDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *AutoTagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *AutoTag

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get autoTag current value
	response, _, err := d.client.AutoTaggingApi.ListAutoTagging(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", autoTagDataSourceName, err))

		return
	}

	autoTag, err := findAutoTag(data.Name.ValueString(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", autoTagDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+autoTagDataSourceName)
	data.write(ctx, autoTag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findAutoTag(name string, autoTags []*sonarr.AutoTaggingResource) (*sonarr.AutoTaggingResource, error) {
	for _, i := range autoTags {
		if i.GetName() == name {
			return i, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(autoTagDataSourceName, "name", name)
}
