package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const languageProfileDataSourceName = "language_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &LanguageProfileDataSource{}

func NewLanguageProfileDataSource() datasource.DataSource {
	return &LanguageProfileDataSource{}
}

// LanguageProfileDataSource defines the language profile implementation.
type LanguageProfileDataSource struct {
	client *sonarr.APIClient
}

func (d *LanguageProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + languageProfileDataSourceName
}

func (d *LanguageProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "<!-- subcategory:Profiles -->Single [Language Profile](../resources/language_profile).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Language Profile ID.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Language Profile name.",
				Required:            true,
			},
			"upgrade_allowed": schema.BoolAttribute{
				MarkdownDescription: "Upgrade allowed Flag.",
				Computed:            true,
			},
			"cutoff_language": schema.StringAttribute{
				MarkdownDescription: "Cutoff Language.",
				Computed:            true,
			},
			"languages": schema.SetAttribute{
				MarkdownDescription: "list of languages in profile.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *LanguageProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *sonarr.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *LanguageProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *LanguageProfile

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get languageprofiles current value
	response, _, err := d.client.LanguageProfileApi.ListLanguageProfile(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", languageProfileDataSourceName, err))

		return
	}

	profile, err := findLanguageProfile(data.Name.ValueString(), response)
	if err != nil {
		resp.Diagnostics.AddError(tools.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", languageProfileDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+languageProfileDataSourceName)
	data.write(ctx, profile)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findLanguageProfile(name string, profiles []*sonarr.LanguageProfileResource) (*sonarr.LanguageProfileResource, error) {
	for _, p := range profiles {
		if p.GetName() == name {
			return p, nil
		}
	}

	return nil, tools.ErrDataNotFoundError(languageProfileDataSourceName, "name", name)
}
