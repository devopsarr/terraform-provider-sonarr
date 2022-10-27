package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const languageProfilesDataSourceName = "language_profiles"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &LanguageProfilesDataSource{}

func NewLanguageProfilesDataSource() datasource.DataSource {
	return &LanguageProfilesDataSource{}
}

// LanguageProfilesDataSource defines the tags implementation.
type LanguageProfilesDataSource struct {
	client *sonarr.Sonarr
}

// LanguageProfiles is a list of Languag profile.
type LanguageProfiles struct {
	LanguageProfiles types.Set    `tfsdk:"language_profiles"`
	ID               types.String `tfsdk:"id"`
}

func (d *LanguageProfilesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + languageProfilesDataSourceName
}

func (d *LanguageProfilesDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "[subcategory:Profiles]: #\nList all available [Language Profiles](../resources/language_profile).",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"language_profiles": {
				MarkdownDescription: "Language Profile list.",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						MarkdownDescription: "Language Profile ID.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"name": {
						MarkdownDescription: "Language Profile name.",
						Computed:            true,
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
				}),
			},
		},
	}, nil
}

func (d *LanguageProfilesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LanguageProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *LanguageProfiles

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get languageprofiles current value
	response, err := d.client.GetLanguageProfilesContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", languageProfilesDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+languageProfilesDataSourceName)
	// Map response body to resource schema attribute
	profiles := make([]LanguageProfile, len(response))
	for i, p := range response {
		profiles[i].write(ctx, p)
	}

	tfsdk.ValueFrom(ctx, profiles, data.LanguageProfiles.Type(context.Background()), &data.LanguageProfiles)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
