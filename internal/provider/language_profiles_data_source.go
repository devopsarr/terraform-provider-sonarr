package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const languageProfilesDataSourceName = "language_profiles"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &LanguageProfilesDataSource{}

func NewLanguageProfilesDataSource() datasource.DataSource {
	return &LanguageProfilesDataSource{}
}

// LanguageProfilesDataSource defines the tags implementation.
type LanguageProfilesDataSource struct {
	client *sonarr.APIClient
}

// LanguageProfiles is a list of Languag profile.
type LanguageProfiles struct {
	LanguageProfiles types.Set    `tfsdk:"language_profiles"`
	ID               types.String `tfsdk:"id"`
}

func (d *LanguageProfilesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + languageProfilesDataSourceName
}

func (d *LanguageProfilesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		DeprecationMessage:  "This resource is deprecated and it will be removed in provider version 3.0.0",
		MarkdownDescription: "<!-- subcategory:Profiles -->**Deprecated**List all available [Language Profiles](../resources/language_profile).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"language_profiles": schema.SetNestedAttribute{
				MarkdownDescription: "Language Profile list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Language Profile ID.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Language Profile name.",
							Computed:            true,
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
				},
			},
		},
	}
}

func (d *LanguageProfilesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LanguageProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *LanguageProfiles

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get languageprofiles current value
	response, _, err := d.client.LanguageProfileApi.ListLanguageProfile(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, languageProfilesDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+languageProfilesDataSourceName)
	// Map response body to resource schema attribute
	profiles := make([]LanguageProfile, len(response))
	for i, p := range response {
		profiles[i].write(ctx, p)
	}

	tfsdk.ValueFrom(ctx, profiles, data.LanguageProfiles.Type(ctx), &data.LanguageProfiles)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
