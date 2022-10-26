package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const languageProfileDataSourceName = "language_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &LanguageProfileDataSource{}

func NewLanguageProfileDataSource() datasource.DataSource {
	return &LanguageProfileDataSource{}
}

// LanguageProfileDataSource defines the language profile implementation.
type LanguageProfileDataSource struct {
	client *sonarr.Sonarr
}

func (d *LanguageProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + languageProfileDataSourceName
}

func (d *LanguageProfileDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "[subcategory:Profiles]: #\nSingle [Language Profile](../resources/language_profile).",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Language Profile ID.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"name": {
				MarkdownDescription: "Language Profile name.",
				Required:            true,
				Type:                types.StringType,
			},
			"upgrade_allowed": {
				MarkdownDescription: "Upgrade allowed Flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"cutoff_language": {
				MarkdownDescription: "Cutoff Language.",
				Computed:            true,
				Type:                types.StringType,
			},
			"languages": {
				MarkdownDescription: "list of languages in profile.",
				Computed:            true,
				Type:                types.SetType{ElemType: types.StringType},
			},
		},
	}, nil
}

func (d *LanguageProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			helpers.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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
	response, err := d.client.GetLanguageProfilesContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", languageProfileDataSourceName, err))

		return
	}

	profile, err := findLanguageProfile(data.Name.Value, response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", languageProfileDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+languageProfileDataSourceName)
	data.write(ctx, profile)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findLanguageProfile(name string, profiles []*sonarr.LanguageProfile) (*sonarr.LanguageProfile, error) {
	for _, p := range profiles {
		if p.Name == name {
			return p, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(languageProfileDataSourceName, "name", name)
}
