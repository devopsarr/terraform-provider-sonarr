package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const qualityDefinitionsDataSourceName = "quality_definitions"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &QualityDefinitionsDataSource{}

func NewQualityDefinitionsDataSource() datasource.DataSource {
	return &QualityDefinitionsDataSource{}
}

// QualityDefinitionsDataSource defines the qyality definitions implementation.
type QualityDefinitionsDataSource struct {
	client *sonarr.APIClient
	auth   context.Context
}

// QualityDefinitions describes the qyality definitions data model.
type QualityDefinitions struct {
	QualityDefinitions types.Set    `tfsdk:"quality_definitions"`
	ID                 types.String `tfsdk:"id"`
}

func (d *QualityDefinitionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + qualityDefinitionsDataSourceName
}

func (d *QualityDefinitionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the quality server.
		MarkdownDescription: "<!-- subcategory:Profiles -->List all available [Quality Definitions](../resources/quality_definition).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"quality_definitions": schema.SetNestedAttribute{
				MarkdownDescription: "Quality Definition list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Quality Definition ID.",
							Computed:            true,
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
				},
			},
		},
	}
}

func (d *QualityDefinitionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *QualityDefinitionsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get qualitydefinitions current value
	response, _, err := d.client.QualityDefinitionAPI.ListQualityDefinition(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, qualityDefinitionsDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+qualityDefinitionsDataSourceName)
	// Map response body to resource schema attribute
	definitions := make([]QualityDefinition, len(response))
	for i, p := range response {
		definitions[i].write(&p)
	}

	qualityList, diags := types.SetValueFrom(ctx, QualityDefinition{}.getType(), definitions)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, QualityDefinitions{QualityDefinitions: qualityList, ID: types.StringValue(strconv.Itoa(len(response)))})...)
}
