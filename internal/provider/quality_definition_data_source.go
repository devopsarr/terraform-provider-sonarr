package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const qualityDefinitionDataSourceName = "quality_definition"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &QualityDefinitionDataSource{}

func NewQualityDefinitionDataSource() datasource.DataSource {
	return &QualityDefinitionDataSource{}
}

// QualityDefinitionDataSource defines the quality definitions implementation.
type QualityDefinitionDataSource struct {
	client *sonarr.APIClient
}

func (d *QualityDefinitionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + qualityDefinitionDataSourceName
}

func (d *QualityDefinitionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the quality server.
		MarkdownDescription: "<!-- subcategory:Profiles -->Single [Quality Definition](../resources/quality_definition).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Quality Definition ID.",
				Required:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Quality Definition Title.",
				Computed:            true,
			},
			"min_size": schema.Float64Attribute{
				MarkdownDescription: "Minimum size MB/min.",
				Optional:            true,
				Computed:            true,
			},
			"max_size": schema.Float64Attribute{
				MarkdownDescription: "Maximum size MB/min.",
				Computed:            true,
			},
			"quality_id": schema.Int64Attribute{
				MarkdownDescription: "Quality ID.",
				Computed:            true,
			},
			"resolution": schema.Int64Attribute{
				MarkdownDescription: "Quality Resolution.",
				Computed:            true,
			},
			"quality_name": schema.StringAttribute{
				MarkdownDescription: "Quality Name.",
				Computed:            true,
			},
			"source": schema.StringAttribute{
				MarkdownDescription: "Quality source.",
				Computed:            true,
			},
		},
	}
}

func (d *QualityDefinitionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			helpers.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *sonarr.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *QualityDefinitionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *QualityDefinition

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get qualitydefinitions current value
	response, _, err := d.client.QualityDefinitionApi.ListQualityDefinition(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, qualityDefinitionDataSourceName, err))
		return
	}

	definition, err := findQualityDefinition(data.ID.ValueInt64(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", qualityDefinitionDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+qualityDefinitionDataSourceName)
	data.write(definition)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findQualityDefinition(id int64, definitions []*sonarr.QualityDefinitionResource) (*sonarr.QualityDefinitionResource, error) {
	for _, p := range definitions {
		if int64(p.GetId()) == id {
			return p, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(qualityDefinitionDataSourceName, "ID", fmt.Sprint(id))
}
