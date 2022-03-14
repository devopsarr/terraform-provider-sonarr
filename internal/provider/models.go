package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Value is generic ID/Name struct applied to a few places.
type Value struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// Tag -
type Tag struct {
	ID    types.Int64  `tfsdk:"id"`
	Label types.String `tfsdk:"label"`
}

//TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
// Tags -
type Tags struct {
	ID   types.String `tfsdk:"id"`
	Tags []Tag        `tfsdk:"tags"`
}

// LanguageProfile -
type LanguageProfile struct {
	UpgradeAllowed types.Bool     `tfsdk:"upgrade_allowed"`
	ID             types.Int64    `tfsdk:"id"`
	Name           types.String   `tfsdk:"name"`
	CutoffLanguage types.String   `tfsdk:"cutoff_language"`
	Languages      []types.String `tfsdk:"languages"`
}

//TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
// LanguageProfiles -
type LanguageProfiles struct {
	ID               types.String      `tfsdk:"id"`
	LanguageProfiles []LanguageProfile `tfsdk:"language_profiles"`
}

// DelayProfile -
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
// DelayProfiles -
type DelayProfiles struct {
	ID            types.String   `tfsdk:"id"`
	DelayProfiles []DelayProfile `tfsdk:"delay_profiles"`
}

//TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
// QualityProfiles -
type QualityProfiles struct {
	ID              types.String     `tfsdk:"id"`
	QualityProfiles []QualityProfile `tfsdk:"quality_profiles"`
}

// QualityProfile -
type QualityProfile struct {
	UpgradeAllowed types.Bool     `tfsdk:"upgrade_allowed"`
	ID             types.Int64    `tfsdk:"id"`
	Cutoff         types.Int64    `tfsdk:"cutoff"`
	Name           types.String   `tfsdk:"name"`
	QualityGroups  []QualityGroup `tfsdk:"quality_groups"`
}

// QualityGroup -
type QualityGroup struct {
	ID types.Int64 `tfsdk:"id"`
	//	Resolution types.Int64  `tfsdk:"resolution"`
	Name types.String `tfsdk:"name"`
	//	Source     types.String `tfsdk:"source"`
	Qualities []Quality `tfsdk:"qualities"`
}

//Quality -
type Quality struct {
	ID         types.Int64  `tfsdk:"id"`
	Resolution types.Int64  `tfsdk:"resolution"`
	Name       types.String `tfsdk:"name"`
	Source     types.String `tfsdk:"source"`
}
