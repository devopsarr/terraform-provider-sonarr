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

const languageDataSourceName = "language"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &LanguageDataSource{}

func NewLanguageDataSource() datasource.DataSource {
	return &LanguageDataSource{}
}

// LanguageDataSource defines the language implementation.
type LanguageDataSource struct {
	client *sonarr.APIClient
	auth   context.Context
}

// Language defines the language data model.
type Language struct {
	Name      types.String `tfsdk:"name"`
	NameLower types.String `tfsdk:"name_lower"`
	ID        types.Int64  `tfsdk:"id"`
}

func (l Language) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"id":         types.Int64Type,
			"name":       types.StringType,
			"name_lower": types.StringType,
		})
}

func (d *LanguageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + languageDataSourceName
}

func (d *LanguageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Languages -->Single available Language.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Language ID.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Language.",
				Required:            true,
			},
			"name_lower": schema.StringAttribute{
				MarkdownDescription: "Language in lowercase.",
				Computed:            true,
			},
		},
	}
}

func (d *LanguageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *LanguageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Language

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get languages current value
	response, _, err := d.client.LanguageAPI.ListLanguage(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, languageDataSourceName, err))

		return
	}

	data.find(data.Name.ValueString(), response, &resp.Diagnostics)
	tflog.Trace(ctx, "read "+languageDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (l *Language) write(language *sonarr.LanguageResource) {
	l.ID = types.Int64Value(int64(language.GetId()))
	l.Name = types.StringValue(language.GetName())
	l.NameLower = types.StringValue(language.GetNameLower())
}

func (l *Language) find(name string, languages []sonarr.LanguageResource, diags *diag.Diagnostics) {
	for _, language := range languages {
		if language.GetName() == name {
			l.write(&language)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(languageDataSourceName, "name", name))
}
