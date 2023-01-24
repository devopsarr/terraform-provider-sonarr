package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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
}

// Languages describes the languages data model.
type Languages struct {
	Languages types.Set    `tfsdk:"languages"`
	ID        types.String `tfsdk:"id"`
}

func (d *LanguagesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + languagesDataSourceName
}

func (d *LanguagesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Languages -->List all available [Languages](../data-sources/language).",
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
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *LanguagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Languages

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get languages current value
	response, _, err := d.client.LanguageApi.ListLanguage(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, languagesDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+languagesDataSourceName)
	// Map response body to resource schema attribute
	languages := make([]Language, len(response))
	for i, t := range response {
		languages[i].write(t)
	}

	tfsdk.ValueFrom(ctx, languages, data.Languages.Type(ctx), &data.Languages)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
