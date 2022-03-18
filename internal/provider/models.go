package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Value is generic ID/Name struct applied to a few places.
type Value struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// Tag is the tag resource.
type Tag struct {
	ID    types.Int64  `tfsdk:"id"`
	Label types.String `tfsdk:"label"`
}

//TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
// Tags is a list of Tag.
type Tags struct {
	ID   types.String `tfsdk:"id"`
	Tags []Tag        `tfsdk:"tags"`
}

// LanguageProfile is the language_profile resource.
type LanguageProfile struct {
	UpgradeAllowed types.Bool     `tfsdk:"upgrade_allowed"`
	ID             types.Int64    `tfsdk:"id"`
	Name           types.String   `tfsdk:"name"`
	CutoffLanguage types.String   `tfsdk:"cutoff_language"`
	Languages      []types.String `tfsdk:"languages"`
}

//TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
// LanguageProfiles is a list of LanguageProfile.
type LanguageProfiles struct {
	ID               types.String      `tfsdk:"id"`
	LanguageProfiles []LanguageProfile `tfsdk:"language_profiles"`
}

// DelayProfile is the delay_profile resource.
type DelayProfile struct {
	EnableUsenet           types.Bool    `tfsdk:"enable_usenet"`
	EnableTorrent          types.Bool    `tfsdk:"enable_torrent"`
	BypassIfHighestQuality types.Bool    `tfsdk:"bypass_if_highest_quality"`
	UsenetDelay            types.Int64   `tfsdk:"usenet_delay"`
	TorrentDelay           types.Int64   `tfsdk:"torrent_delay"`
	ID                     types.Int64   `tfsdk:"id"`
	Order                  types.Int64   `tfsdk:"order"`
	Tags                   []types.Int64 `tfsdk:"tags"`
	PreferredProtocol      types.String  `tfsdk:"preferred_protocol"`
}

//TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
// DelayProfiles is a list of DelayProfile.
type DelayProfiles struct {
	ID            types.String   `tfsdk:"id"`
	DelayProfiles []DelayProfile `tfsdk:"delay_profiles"`
}

// QualityProfile is the quality_profile resource.
type QualityProfile struct {
	UpgradeAllowed types.Bool     `tfsdk:"upgrade_allowed"`
	ID             types.Int64    `tfsdk:"id"`
	Cutoff         types.Int64    `tfsdk:"cutoff"`
	Name           types.String   `tfsdk:"name"`
	QualityGroups  []QualityGroup `tfsdk:"quality_groups"`
}

//TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
// QualityProfiles is a list of QualityProfile.
type QualityProfiles struct {
	ID              types.String     `tfsdk:"id"`
	QualityProfiles []QualityProfile `tfsdk:"quality_profiles"`
}

// QualityGroup is part of QualityProfile.
type QualityGroup struct {
	ID        types.Int64  `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Qualities []Quality    `tfsdk:"qualities"`
}

//Quality is part of QualityGroup.
type Quality struct {
	ID         types.Int64  `tfsdk:"id"`
	Resolution types.Int64  `tfsdk:"resolution"`
	Name       types.String `tfsdk:"name"`
	Source     types.String `tfsdk:"source"`
}

// Series is the series resource.
type Series struct {
	Monitored         types.Bool    `tfsdk:"monitored"`
	SeasonFolder      types.Bool    `tfsdk:"season_folder"`
	UseSceneNumbering types.Bool    `tfsdk:"use_scene_numbering"`
	ID                types.Int64   `tfsdk:"id"`
	LanguageProfileID types.Int64   `tfsdk:"language_profile_id"`
	QualityProfileID  types.Int64   `tfsdk:"quality_profile_id"`
	TvdbID            types.Int64   `tfsdk:"tvdb_id"`
	Path              types.String  `tfsdk:"path"`
	Title             types.String  `tfsdk:"title"`
	TitleSlug         types.String  `tfsdk:"title_slug"`
	RootFolderPath    types.String  `tfsdk:"root_folder_path"`
	Tags              []types.Int64 `tfsdk:"tags"`
}

//TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
// QualityProfiles is a list of QualityProfile.
type SeriesList struct {
	ID     types.String `tfsdk:"id"`
	Series []Series     `tfsdk:"series"`
}

// Season is part of Series.
type Season struct {
	Monitored    types.Bool  `tfsdk:"monitored"`
	SeasonNumber types.Int64 `tfsdk:"season_number"`
}

// AddSeriesOptions is used in series creation.
type AddSeriesOptions struct {
	SearchForMissingEpisodes     types.Bool `tfsdk:"search_for_missing_episodes"`
	SearchForCutoffUnmetEpisodes types.Bool `tfsdk:"search_for_cutoff_unmet_episodes"`
	IgnoreEpisodesWithFiles      types.Bool `tfsdk:"ignore_episodes_with_files"`
	IgnoreEpisodesWithoutFiles   types.Bool `tfsdk:"ignore_episodes_without_files"`
}

// Image is part of Series.
type Image struct {
	CoverType types.String `tfsdk:"cover_type"`
	URL       types.String `tfsdk:"url"`
	RemoteURL types.String `tfsdk:"remote_url"`
	Extension types.String `tfsdk:"extension"`
}
