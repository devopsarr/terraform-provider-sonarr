package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

// needed for tf debug mode
// var stderr = os.Stderr

// Ensure provider defined types fully satisfy framework interfaces.
var _ provider.Provider = &sonarrProvider{}

// provider satisfies the provider.Provider interface and usually is included
// with all Resource and DataSource implementations.
type sonarrProvider struct {
	// client can contain the upstream provider SDK or HTTP client used to
	// communicate with the upstream service. Resource and DataSource
	// implementations can then make calls using this client.
	client sonarr.Sonarr

	// configured is set to true at the end of the Configure method.
	// This can be used in Resource and DataSource implementations to verify
	// that the provider was previously configured.
	configured bool

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	APIKey types.String `tfsdk:"api_key"`
	URL    types.String `tfsdk:"url"`
}

func (p *sonarrProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// User must provide URL to the provider
	if data.URL.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as url",
		)

		return
	}

	var url string
	if data.URL.Null {
		url = os.Getenv("SONARR_URL")
	} else {
		url = data.URL.Value
	}

	if url == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find URL",
			"URL cannot be an empty string",
		)

		return
	}

	// User must provide API key to the provider
	if data.APIKey.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as api_key",
		)

		return
	}

	var key string
	if data.APIKey.Null {
		key = os.Getenv("SONARR_API_KEY")
	} else {
		key = data.APIKey.Value
	}

	if key == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find API key",
			"API key cannot be an empty string",
		)

		return
	}
	// If the upstream provider SDK or HTTP client requires configuration, such
	// as authentication or logging, this is a great opportunity to do so.
	c := *sonarr.New(starr.New(key, url, 0))
	p.client = c
	p.configured = true
}

func (p *sonarrProvider) GetResources(ctx context.Context) (map[string]provider.ResourceType, diag.Diagnostics) {
	return map[string]provider.ResourceType{
		"sonarr_delay_profile":    resourceDelayProfileType{},
		"sonarr_indexer":          resourceIndexerType{},
		"sonarr_indexer_config":   resourceIndexerConfigType{},
		"sonarr_language_profile": resourceLanguageProfileType{},
		"sonarr_media_management": resourceMediaManagementType{},
		"sonarr_naming":           resourceNamingType{},
		"sonarr_quality_profile":  resourceQualityProfileType{},
		"sonarr_root_folder":      resourceRootFolderType{},
		"sonarr_series":           resourceSeriesType{},
		"sonarr_tag":              resourceTagType{},
	}, nil
}

func (p *sonarrProvider) GetDataSources(ctx context.Context) (map[string]provider.DataSourceType, diag.Diagnostics) {
	return map[string]provider.DataSourceType{
		"sonarr_delay_profile":     dataDelayProfileType{},
		"sonarr_delay_profiles":    dataDelayProfilesType{},
		"sonarr_indexer":           dataIndexerType{},
		"sonarr_indexers":          dataIndexersType{},
		"sonarr_language_profile":  dataLanguageProfileType{},
		"sonarr_language_profiles": dataLanguageProfilesType{},
		"sonarr_quality_profile":   dataQualityProfileType{},
		"sonarr_quality_profiles":  dataQualityProfilesType{},
		"sonarr_root_folder":       dataRootFolderType{},
		"sonarr_root_folders":      dataRootFoldersType{},
		"sonarr_series":            dataSeriesType{},
		"sonarr_all_series":        dataAllSeriesType{},
		"sonarr_tag":               dataTagType{},
		"sonarr_tags":              dataTagsType{},
	}, nil
}

func (p *sonarrProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "The Sonarr provider is used to interact with any [Sonarr](https://sonarr.tv/) installation.<br/>You must configure the provider with the proper [credentials](#api_key) before you can use it. <br/>Use the left navigation to read about the available resources.<br/><br/>For more information about Sonarr and its resources, as well as configuration guides and hints, visit the [Servarr wiki](https://wiki.servarr.com/en/sonarr).",
		Attributes: map[string]tfsdk.Attribute{
			"api_key": {
				MarkdownDescription: "API key for Sonarr authentication. Can be specified via the `SONARR_API_KEY` environment variable.",
				Optional:            true,
				Type:                types.StringType,
				Sensitive:           true,
			},
			"url": {
				MarkdownDescription: "Full Sonarr URL with protocol and port (e.g. `https://test.sonarr.tv:8989`). You should **NOT** supply any path (`/api`), the SDK will use the appropriate paths. Can be specified via the `SONARR_URL` environment variable.",
				Optional:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

// New returns the provider with a specific version.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &sonarrProvider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*sonarrProvider)), however using this can prevent
// potential panics.
func convertProviderType(in provider.Provider) (sonarrProvider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*sonarrProvider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)

		return sonarrProvider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)

		return sonarrProvider{}, diags
	}

	return *p, diags
}
