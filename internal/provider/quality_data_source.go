package provider

import (
	"context"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const qualityDataSourceName = "quality"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &QualityDataSource{}

func NewQualityDataSource() datasource.DataSource {
	return &QualityDataSource{}
}

// QualityDataSource defines the quality implementation.
type QualityDataSource struct {
	client *sonarr.APIClient
}

// Quality is part of QualityGroup.
type Quality struct {
	Name       types.String `tfsdk:"name"`
	Source     types.String `tfsdk:"source"`
	ID         types.Int64  `tfsdk:"id"`
	Resolution types.Int64  `tfsdk:"resolution"`
}

func (q Quality) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"name":       types.StringType,
			"source":     types.StringType,
			"id":         types.Int64Type,
			"resolution": types.Int64Type,
		})
}

func (d *QualityDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + qualityDataSourceName
}

func (d *QualityDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the quality server.
		MarkdownDescription: "<!-- subcategory:Profiles -->Single Quality.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Quality  ID.",
				Computed:            true,
			},
			"resolution": schema.Int64Attribute{
				MarkdownDescription: "Quality Resolution.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Quality Name.",
				Required:            true,
			},
			"source": schema.StringAttribute{
				MarkdownDescription: "Quality source.",
				Computed:            true,
			},
		},
	}
}

func (d *QualityDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *QualityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Quality

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get qualitys current value
	response, _, err := d.client.QualityDefinitionApi.ListQualityDefinition(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, qualityDataSourceName, err))

		return
	}

	data.find(data.Name.ValueString(), response, &resp.Diagnostics)
	tflog.Trace(ctx, "read "+qualityDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (q *Quality) find(name string, definitions []*sonarr.QualityDefinitionResource, diags *diag.Diagnostics) {
	for _, def := range definitions {
		if def.Quality.GetName() == name {
			q.writeFromDefinition(def)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(qualityDataSourceName, "name", name))
}

func (q *Quality) writeFromDefinition(quality *sonarr.QualityDefinitionResource) {
	q.ID = types.Int64Value(int64(quality.Quality.GetId()))
	q.Name = types.StringValue(quality.Quality.GetName())
	q.Source = types.StringValue(string(quality.Quality.GetSource()))
	q.Resolution = types.Int64Value(int64(quality.Quality.GetResolution()))
}
