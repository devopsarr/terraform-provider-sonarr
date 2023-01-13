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

const languageDataSourceName = "language"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &LanguageDataSource{}

func NewLanguageDataSource() datasource.DataSource {
	return &LanguageDataSource{}
}

// LanguageDataSource defines the language implementation.
type LanguageDataSource struct {
	client *sonarr.APIClient
}

// Language defines the language data model.
type Language struct {
	Name      types.String `tfsdk:"name"`
	NameLower types.String `tfsdk:"name_lower"`
	ID        types.Int64  `tfsdk:"id"`
}

func (d *LanguageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + languageDataSourceName
}

func (d *LanguageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

func (d *LanguageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var language *Language

	resp.Diagnostics.Append(req.Config.Get(ctx, &language)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get languages current value
	response, _, err := d.client.LanguageApi.ListLanguage(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, languageDataSourceName, err))

		return
	}

	value, err := findLanguage(language.Name.ValueString(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", languageDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+languageDataSourceName)
	language.write(value)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &language)...)
}

func (l *Language) write(language *sonarr.LanguageResource) {
	l.ID = types.Int64Value(int64(language.GetId()))
	l.Name = types.StringValue(language.GetName())
	l.NameLower = types.StringValue(language.GetNameLower())
}

func findLanguage(name string, languages []*sonarr.LanguageResource) (*sonarr.LanguageResource, error) {
	for _, t := range languages {
		if t.GetName() == name {
			return t, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(languageDataSourceName, "name", name)
}
