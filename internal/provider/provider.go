package provider

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// needed for tf debug mode
// var stderr = os.Stderr

// Ensure provider defined types fully satisfy framework interfaces.
var _ provider.Provider = &SonarrProvider{}

// SonarrProvider defines the provider implementation.
type SonarrProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Sonarr describes the provider data model.
type Sonarr struct {
	APIKey types.String `tfsdk:"api_key"`
	URL    types.String `tfsdk:"url"`
}

func (p *SonarrProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sonarr"
	resp.Version = p.version
}

func (p *SonarrProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Sonarr provider is used to interact with any [Sonarr](https://sonarr.tv/) installation.\nYou must configure the provider with the proper [credentials](#api_key) before you can use it.\nUse the left navigation to read about the available resources.\n\nFor more information about Sonarr and its resources, as well as configuration guides and hints, visit the [Servarr wiki](https://wiki.servarr.com/en/sonarr).",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key for Sonarr authentication. Can be specified via the `SONARR_API_KEY` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Full Sonarr URL with protocol and port (e.g. `https://test.sonarr.tv:8989`). You should **NOT** supply any path (`/api`), the SDK will use the appropriate paths. Can be specified via the `SONARR_URL` environment variable.",
				Optional:            true,
			},
		},
	}
}

func (p *SonarrProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data Sonarr

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Extract URL
	APIURL := data.URL.ValueString()
	if APIURL == "" {
		APIURL = os.Getenv("SONARR_URL")
	}

	parsedAPIURL, err := url.Parse(APIURL)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find valid URL",
			"URL cannot parsed",
		)

		return
	}

	// Extract key
	key := data.APIKey.ValueString()
	if key == "" {
		key = os.Getenv("SONARR_API_KEY")
	}

	if key == "" {
		resp.Diagnostics.AddError(
			"Unable to find API key",
			"API key cannot be an empty string",
		)

		return
	}

	// Set context for API calls
	auth := context.WithValue(
		context.Background(),
		sonarr.ContextAPIKeys,
		map[string]sonarr.APIKey{
			"X-Api-Key": {Key: key},
		},
	)
	auth = context.WithValue(auth, sonarr.ContextServerVariables, map[string]string{
		"protocol": parsedAPIURL.Scheme,
		"hostpath": parsedAPIURL.Host,
	})
	resp.DataSourceData = auth
	resp.ResourceData = auth
}

func (p *SonarrProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Download Clients
		NewDownloadClientConfigResource,
		NewDownloadClientResource,
		NewDownloadClientUsenetDownloadStationResource,
		NewDownloadClientNzbgetResource,
		NewDownloadClientNzbvortexResource,
		NewDownloadClientPneumaticResource,
		NewDownloadClientSabnzbdResource,
		NewDownloadClientUsenetBlackholeResource,
		NewDownloadClientAria2Resource,
		NewDownloadClientDelugeResource,
		NewDownloadClientTorrentDownloadStationResource,
		NewDownloadClientFloodResource,
		NewDownloadClientHadoukenResource,
		NewDownloadClientQbittorrentResource,
		NewDownloadClientRtorrentResource,
		NewDownloadClientTorrentBlackholeResource,
		NewDownloadClientTransmissionResource,
		NewDownloadClientUtorrentResource,
		NewDownloadClientVuzeResource,
		NewRemotePathMappingResource,

		// Indexers
		NewIndexerConfigResource,
		NewIndexerResource,
		NewIndexerFanzubResource,
		NewIndexerNewznabResource,
		NewIndexerBroadcastheNetResource,
		NewIndexerFilelistResource,
		NewIndexerHdbitsResource,
		NewIndexerIptorrentsResource,
		NewIndexerNyaaResource,
		NewIndexerTorrentRssResource,
		NewIndexerTorrentleechResource,
		NewIndexerTorznabResource,

		// Import Lists
		NewImportListExclusionResource,
		NewImportListResource,
		NewImportListCustomResource,
		NewImportListSimklUserResource,
		NewImportListSonarrResource,
		NewImportListImdbResource,
		NewImportListPlexResource,
		NewImportListPlexRSSResource,
		NewImportListTraktListResource,
		NewImportListTraktPopularResource,
		NewImportListTraktUserResource,

		// Media Management
		NewMediaManagementResource,
		NewNamingResource,
		NewRootFolderResource,

		// Metadata
		NewMetadataResource,
		NewMetadataKodiResource,
		NewMetadataRoksboxResource,
		NewMetadataWdtvResource,

		// Notifications
		NewNotificationResource,
		NewNotificationAppriseResource,
		NewNotificationCustomScriptResource,
		NewNotificationWebhookResource,
		NewNotificationDiscordResource,
		NewNotificationEmailResource,
		NewNotificationEmbyResource,
		NewNotificationGotifyResource,
		NewNotificationJoinResource,
		NewNotificationKodiResource,
		NewNotificationMailgunResource,
		NewNotificationNtfyResource,
		NewNotificationPlexResource,
		NewNotificationProwlResource,
		NewNotificationPushbulletResource,
		NewNotificationPushoverResource,
		NewNotificationSendgridResource,
		NewNotificationSignalResource,
		NewNotificationSimplepushResource,
		NewNotificationSlackResource,
		NewNotificationSynologyResource,
		NewNotificationTelegramResource,
		NewNotificationTraktResource,
		NewNotificationTwitterResource,

		// Profiles
		NewCustomFormatResource,
		NewDelayProfileResource,
		NewQualityProfileResource,
		NewReleaseProfileResource,
		NewQualityDefinitionResource,

		// Series
		NewSeriesResource,

		// System
		NewHostResource,

		// Tags
		NewTagResource,
		NewAutoTagResource,
	}
}

func (p *SonarrProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Download Clients
		NewDownloadClientConfigDataSource,
		NewDownloadClientDataSource,
		NewDownloadClientsDataSource,
		NewRemotePathMappingDataSource,
		NewRemotePathMappingsDataSource,

		// Indexers
		NewIndexerConfigDataSource,
		NewIndexerDataSource,
		NewIndexersDataSource,

		// Import Lists
		NewImportListExclusionDataSource,
		NewImportListExclusionsDataSource,
		NewImportListDataSource,
		NewImportListsDataSource,

		// Media Management
		NewMediaManagementDataSource,
		NewNamingDataSource,
		NewRootFolderDataSource,
		NewRootFoldersDataSource,

		// Metadata
		NewMetadataConsumersDataSource,
		NewMetadataDataSource,

		// Notifications
		NewNotificationDataSource,
		NewNotificationsDataSource,

		// Profiles
		NewCustomFormatDataSource,
		NewCustomFormatsDataSource,
		NewDelayProfileDataSource,
		NewDelayProfilesDataSource,
		NewQualityProfileDataSource,
		NewQualityProfilesDataSource,
		NewReleaseProfileDataSource,
		NewReleaseProfilesDataSource,
		NewQualityDefinitionDataSource,
		NewQualityDefinitionsDataSource,
		NewCustomFormatConditionDataSource,
		NewCustomFormatConditionLanguageDataSource,
		NewCustomFormatConditionReleaseGroupDataSource,
		NewCustomFormatConditionReleaseTitleDataSource,
		NewCustomFormatConditionResolutionDataSource,
		NewCustomFormatConditionSizeDataSource,
		NewCustomFormatConditionSourceDataSource,
		NewQualityDataSource,

		// Series
		NewSeriesDataSource,
		NewAllSeriessDataSource,
		NewSearchSeriesDataSource,

		// System
		NewLanguageDataSource,
		NewLanguagesDataSource,
		NewSystemStatusDataSource,
		NewHostDataSource,

		// Tags
		NewTagDataSource,
		NewTagsDataSource,
		NewAutoTagDataSource,
		NewAutoTagsDataSource,
		NewAutoTagConditionDataSource,
		NewAutoTagConditionGenresDataSource,
		NewAutoTagConditionRootFolderDataSource,
		NewAutoTagConditionSeriesTypeDataSource,
	}
}

// New returns the provider with a specific version.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SonarrProvider{
			version: version,
		}
	}
}

// ResourceConfigure is a helper function to set the client for a specific resource.
func resourceConfigure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) (context.Context, *sonarr.APIClient) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return nil, nil
	}

	providerData, ok := req.ProviderData.(context.Context)
	if !ok {
		resp.Diagnostics.AddError(
			helpers.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected context.Context, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return nil, nil
	}

	return providerData, sonarr.NewAPIClient(sonarr.NewConfiguration())
}

func dataSourceConfigure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) (context.Context, *sonarr.APIClient) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return nil, nil
	}

	providerData, ok := req.ProviderData.(context.Context)
	if !ok {
		resp.Diagnostics.AddError(
			helpers.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected context.Context, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return nil, nil
	}

	return providerData, sonarr.NewAPIClient(sonarr.NewConfiguration())
}
