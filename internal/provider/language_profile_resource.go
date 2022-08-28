package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ provider.ResourceType            = resourceLanguageProfileType{}
	_ resource.Resource                = resourceLanguageProfile{}
	_ resource.ResourceWithImportState = resourceLanguageProfile{}
)

type resourceLanguageProfileType struct{}

type resourceLanguageProfile struct {
	provider sonarrProvider
}

// LanguageProfile is the language_profile resource.
type LanguageProfile struct {
	UpgradeAllowed types.Bool   `tfsdk:"upgrade_allowed"`
	ID             types.Int64  `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	CutoffLanguage types.String `tfsdk:"cutoff_language"`
	Languages      types.Set    `tfsdk:"languages"`
}

func (t resourceLanguageProfileType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "LanguageProfile resource",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "ID of languageprofile",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"name": {
				MarkdownDescription: "Name of languageprofile",
				Required:            true,
				Type:                types.StringType,
			},
			"upgrade_allowed": {
				MarkdownDescription: "Upgrade allowed Flag",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"cutoff_language": {
				MarkdownDescription: "Name of language",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					helpers.StringMatch(helpers.Languages),
				},
			},
			"languages": {
				MarkdownDescription: "list of languages in profile",
				Required:            true,
				Type:                types.SetType{ElemType: types.StringType},
			},
		},
	}, nil
}

func (t resourceLanguageProfileType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return resourceLanguageProfile{
		provider: provider,
	}, diags
}

func (r resourceLanguageProfile) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan LanguageProfile
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	data := readLanguageProfile(ctx, &plan)

	// Create new LanguageProfile
	response, err := r.provider.client.AddLanguageProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create languageprofile, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "created languageprofile: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := writeLanguageProfile(ctx, response)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceLanguageProfile) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state LanguageProfile
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get languageprofile current value
	response, err := r.provider.client.GetLanguageProfileContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read languageprofiles, got error: %s", err))

		return
	}
	// Map response body to resource schema attribute
	result := writeLanguageProfile(ctx, response)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceLanguageProfile) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan LanguageProfile
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := readLanguageProfile(ctx, &plan)

	// Update LanguageProfile
	response, err := r.provider.client.UpdateLanguageProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update languageprofile, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "update languageprofile: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := writeLanguageProfile(ctx, response)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceLanguageProfile) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state LanguageProfile

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete languageprofile current value
	err := r.provider.client.DeleteLanguageProfileContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read languageprofiles, got error: %s", err))

		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceLanguageProfile) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func writeLanguageProfile(ctx context.Context, profile *sonarr.LanguageProfile) *LanguageProfile {
	output := LanguageProfile{
		UpgradeAllowed: types.Bool{Value: profile.UpgradeAllowed},
		ID:             types.Int64{Value: profile.ID},
		Name:           types.String{Value: profile.Name},
		CutoffLanguage: types.String{Value: profile.Cutoff.Name},
		Languages:      types.Set{ElemType: types.StringType},
	}

	languages := make([]string, len(profile.Languages))
	for i, l := range profile.Languages {
		languages[i] = l.Language.Name
	}

	tfsdk.ValueFrom(ctx, languages, output.Languages.Type(ctx), &output.Languages)

	return &output
}

func readLanguageProfile(ctx context.Context, profile *LanguageProfile) *sonarr.LanguageProfile {
	langs := make([]string, len(profile.Languages.Elems))
	tfsdk.ValueAs(ctx, profile.Languages, &langs)

	languages := make([]sonarr.Language, len(langs))
	for i, l := range langs {
		languages[i] = sonarr.Language{
			Allowed: true,
			Language: &starr.Value{
				Name: l,
				ID:   helpers.GetLanguageID(l),
			},
		}
	}

	return &sonarr.LanguageProfile{
		Name:           profile.Name.Value,
		UpgradeAllowed: profile.UpgradeAllowed.Value,
		Cutoff: &starr.Value{
			Name: profile.CutoffLanguage.Value,
			ID:   helpers.GetLanguageID(profile.CutoffLanguage.Value),
		},
		Languages: languages,
		ID:        profile.ID.Value,
	}
}
