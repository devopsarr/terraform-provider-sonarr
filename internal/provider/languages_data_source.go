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

const languagesDataSourceName = "languages"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &LanguagesDataSource{}

func NewLanguagesDataSource() datasource.DataSource {
	return &LanguagesDataSource{}
}

// LanguagesDataSource defines the languages implementation.
type LanguagesDataSource struct {
	client *sonarr.APIClient
	auth   context.Context
}

// Languages describes the languages data model.
type Languages struct {
	Languages types.Set    `tfsdk:"languages"`
	ID        types.String `tfsdk:"id"`
}

func (d *LanguagesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + languagesDataSourceName
}

func (d *LanguagesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Languages -->\nList all available [Languages](../data-sources/language).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"languages": schema.SetNestedAttribute{
				MarkdownDescription: "Language list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Language ID.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Language.",
							Computed:            true,
						},
						"name_lower": schema.StringAttribute{
							MarkdownDescription: "Language in lowercase.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *LanguagesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *LanguagesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get languages current value
	response, _, err := d.client.LanguageAPI.ListLanguage(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, languagesDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+languagesDataSourceName)
	// Map response body to resource schema attribute
	languages := make([]Language, len(response))
	for i, t := range response {
		languages[i].write(&t)
	}

	languageList, diags := types.SetValueFrom(ctx, Language{}.getType(), languages)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, Languages{Languages: languageList, ID: types.StringValue(strconv.Itoa(len(response)))})...)
}
