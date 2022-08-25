package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ provider.DataSourceType = dataLanguageProfilesType{}
var _ datasource.DataSource = dataLanguageProfiles{}

type dataLanguageProfilesType struct{}

type dataLanguageProfiles struct {
	provider sonarrProvider
}

// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
// LanguageProfiles is a list of LanguageProfile.
type LanguageProfiles struct {
	ID               types.String `tfsdk:"id"`
	LanguageProfiles types.Set    `tfsdk:"language_profiles"`
}

func (t dataLanguageProfilesType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "List all available languageprofiles",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
			"language_profiles": {
				MarkdownDescription: "List of languageprofiles",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						MarkdownDescription: "ID of languageprofile",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"name": {
						MarkdownDescription: "Name of languageprofile",
						Computed:            true,
						Type:                types.StringType,
					},
					"upgrade_allowed": {
						MarkdownDescription: "Upgrade allowed Flag",
						Computed:            true,
						Type:                types.BoolType,
					},
					"cutoff_language": {
						MarkdownDescription: "Cutoff Language",
						Computed:            true,
						Type:                types.StringType,
					},
					"languages": {
						MarkdownDescription: "list of languages in profile",
						Computed:            true,
						Type:                types.SetType{ElemType: types.StringType},
					},
				}),
			},
		},
	}, nil
}

func (t dataLanguageProfilesType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return dataLanguageProfiles{
		provider: provider,
	}, diags
}

func (d dataLanguageProfiles) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data LanguageProfiles
	diags := resp.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get languageprofiles current value
	response, err := d.provider.client.GetLanguageProfilesContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read languageprofiles, got error: %s", err))

		return
	}
	// Map response body to resource schema attribute
	profiles := *writeLanguageprofiles(ctx, response)
	tfsdk.ValueFrom(ctx, profiles, data.LanguageProfiles.Type(context.Background()), &data.LanguageProfiles)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.String{Value: strconv.Itoa(len(response))}
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func writeLanguageprofiles(ctx context.Context, languages []*sonarr.LanguageProfile) *[]LanguageProfile {
	output := make([]LanguageProfile, len(languages))
	for i, p := range languages {
		output[i] = *writeLanguageProfile(ctx, p)
	}

	return &output
}
