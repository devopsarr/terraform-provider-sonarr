package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

const languageProfileResourceName = "language_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &LanguageProfileResource{}
var _ resource.ResourceWithImportState = &LanguageProfileResource{}

func NewLanguageProfileResource() resource.Resource {
	return &LanguageProfileResource{}
}

// LanguageProfileResource defines the language profile implementation.
type LanguageProfileResource struct {
	client *sonarr.Sonarr
}

// LanguageProfile describes the language profile data model.
type LanguageProfile struct {
	Languages      types.Set    `tfsdk:"languages"`
	Name           types.String `tfsdk:"name"`
	CutoffLanguage types.String `tfsdk:"cutoff_language"`
	ID             types.Int64  `tfsdk:"id"`
	UpgradeAllowed types.Bool   `tfsdk:"upgrade_allowed"`
}

func (r *LanguageProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + languageProfileResourceName
}

func (r *LanguageProfileResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "[subcategory:Profiles]: #\nLanguage Profile resource.\nFor more information refer to [Language Profile](https://wiki.servarr.com/sonarr/settings#language-profiles) documentation.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Language Profile ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"name": {
				MarkdownDescription: "Language Profile name.",
				Required:            true,
				Type:                types.StringType,
			},
			"upgrade_allowed": {
				MarkdownDescription: "Upgrade allowed Flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"cutoff_language": {
				MarkdownDescription: "Name of language.",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					helpers.StringMatch(helpers.Languages),
				},
			},
			"languages": {
				MarkdownDescription: "list of languages in profile.",
				Required:            true,
				Type:                types.SetType{ElemType: types.StringType},
			},
		},
	}, nil
}

func (r *LanguageProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			helpers.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *LanguageProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var profile *LanguageProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	data := profile.read(ctx)

	// Create new LanguageProfile
	response, err := r.client.AddLanguageProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", languageProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+languageProfileResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, profile)...)
}

func (r *LanguageProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var profile *LanguageProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get languageprofile current value
	response, err := r.client.GetLanguageProfileContext(ctx, profile.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", languageProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+languageProfileResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *LanguageProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var profile *LanguageProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := profile.read(ctx)

	// Update LanguageProfile
	response, err := r.client.UpdateLanguageProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", languageProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+languageProfileResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *LanguageProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var profile *LanguageProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete languageprofile current value
	err := r.client.DeleteLanguageProfileContext(ctx, profile.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", languageProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+languageProfileResourceName+": "+strconv.Itoa(int(profile.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *LanguageProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			helpers.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+languageProfileResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (p *LanguageProfile) write(ctx context.Context, profile *sonarr.LanguageProfile) {
	p.UpgradeAllowed = types.BoolValue(profile.UpgradeAllowed)
	p.ID = types.Int64Value(profile.ID)
	p.Name = types.StringValue(profile.Name)
	p.CutoffLanguage = types.StringValue(profile.Cutoff.Name)
	p.Languages = types.SetValueMust(types.StringType, nil)

	languages := make([]string, len(profile.Languages))
	for i, l := range profile.Languages {
		languages[i] = l.Language.Name
	}

	tfsdk.ValueFrom(ctx, languages, p.Languages.Type(ctx), &p.Languages)
}

func (p *LanguageProfile) read(ctx context.Context) *sonarr.LanguageProfile {
	langs := make([]string, len(p.Languages.Elements()))
	tfsdk.ValueAs(ctx, p.Languages, &langs)

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
		Name:           p.Name.ValueString(),
		UpgradeAllowed: p.UpgradeAllowed.ValueBool(),
		Cutoff: &starr.Value{
			Name: p.CutoffLanguage.ValueString(),
			ID:   helpers.GetLanguageID(p.CutoffLanguage.ValueString()),
		},
		Languages: languages,
		ID:        p.ID.ValueInt64(),
	}
}
