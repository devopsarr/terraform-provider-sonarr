package provider

import (
	"context"
	"os"

	"github.com/devopsarr/sonarr-go/sonarr"
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

// ScaffoldingProvider defines the provider implementation.
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

func (p *SonarrProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sonarr"
	resp.Version = p.version
}

func (p *SonarrProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
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

	// User must provide URL to the provider
	if data.URL.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as url",
		)

		return
	}

	var url string
	if data.URL.IsNull() {
		url = os.Getenv("SONARR_URL")
	} else {
		url = data.URL.ValueString()
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
	if data.APIKey.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as api_key",
		)

		return
	}

	var key string
	if data.APIKey.IsNull() {
		key = os.Getenv("SONARR_API_KEY")
	} else {
		key = data.APIKey.ValueString()
	}

	if key == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find API key",
			"API key cannot be an empty string",
		)

		return
	}

	// Configuring client. API Key management could be changed once new options avail in sdk.
	config := sonarr.NewConfiguration()
	config.AddDefaultHeader("X-Api-Key", key)
	config.Servers[0].URL = url
	client := sonarr.NewAPIClient(config)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *SonarrProvider) Resources(ctx context.Context) []func() resource.Resource {
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
		NewIndexerRarbgResource,
		NewIndexerTorrentRssResource,
		NewIndexerTorrentleechResource,
		NewIndexerTorznabResource,

		// Import Lists
		NewImportListExclusionResource,
		NewImportListResource,
		NewImportListSonarrResource,
		NewImportListPlexResource,
		NewImportListTraktListResource,
		NewImportListTraktPopularResource,
		NewImportListTraktUserResource,

		// Media Management
		NewMediaManagementResource,
		NewNamingResource,
		NewRootFolderResource,

		// Notifications
		NewNotificationResource,
		NewNotificationBoxcarResource,
		NewNotificationCustomScriptResource,
		NewNotificationWebhookResource,
		NewNotificationDiscordResource,
		NewNotificationEmailResource,
		NewNotificationEmbyResource,
		NewNotificationGotifyResource,
		NewNotificationJoinResource,
		NewNotificationKodiResource,
		NewNotificationMailgunResource,
		NewNotificationPlexResource,
		NewNotificationProwlResource,
		NewNotificationPushbulletResource,
		NewNotificationPushoverResource,
		NewNotificationSendgridResource,
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

		// Tags
		NewTagResource,
	}
}

func (p *SonarrProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
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

		// Series
		NewSeriesDataSource,
		NewAllSeriessDataSource,
		NewSearchSeriesDataSource,

		// System Status
		NewLanguageDataSource,
		NewLanguagesDataSource,
		NewSystemStatusDataSource,

		// Tags
		NewTagDataSource,
		NewTagsDataSource,
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
