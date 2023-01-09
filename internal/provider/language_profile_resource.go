package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const languageProfileResourceName = "language_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &LanguageProfileResource{}
	_ resource.ResourceWithImportState = &LanguageProfileResource{}
)

func NewLanguageProfileResource() resource.Resource {
	return &LanguageProfileResource{}
}

// LanguageProfileResource defines the language profile implementation.
type LanguageProfileResource struct {
	client *sonarr.APIClient
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

func (r *LanguageProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Profiles -->Language Profile resource.\nFor more information refer to [Language Profile](https://wiki.servarr.com/sonarr/settings#language-profiles) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Language Profile ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Language Profile name.",
				Required:            true,
			},
			"upgrade_allowed": schema.BoolAttribute{
				MarkdownDescription: "Upgrade allowed Flag.",
				Optional:            true,
				Computed:            true,
			},
			"cutoff_language": schema.StringAttribute{
				MarkdownDescription: "Name of language.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(helpers.Languages...),
				},
			},
			"languages": schema.SetAttribute{
				MarkdownDescription: "list of languages in profile.",
				Required:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *LanguageProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *sonarr.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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
	request := profile.read(ctx)

	// Create new LanguageProfile
	response, _, err := r.client.LanguageProfileApi.CreateLanguageProfile(ctx).LanguageProfileResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", languageProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+languageProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
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
	response, _, err := r.client.LanguageProfileApi.GetLanguageProfileById(ctx, int32(profile.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", languageProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+languageProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
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
	request := profile.read(ctx)

	// Update LanguageProfile
	response, _, err := r.client.LanguageProfileApi.UpdateLanguageProfile(ctx, strconv.Itoa(int(request.GetId()))).LanguageProfileResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", languageProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+languageProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
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
	_, err := r.client.LanguageProfileApi.DeleteLanguageProfile(ctx, int32(profile.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", languageProfileResourceName, err))

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
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+languageProfileResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (p *LanguageProfile) write(ctx context.Context, profile *sonarr.LanguageProfileResource) {
	p.UpgradeAllowed = types.BoolValue(profile.GetUpgradeAllowed())
	p.ID = types.Int64Value(int64(profile.GetId()))
	p.Name = types.StringValue(profile.GetName())
	p.CutoffLanguage = types.StringValue(profile.Cutoff.GetName())
	p.Languages = types.SetValueMust(types.StringType, nil)

	languages := make([]string, len(profile.Languages))
	for i, l := range profile.Languages {
		languages[i] = l.Language.GetName()
	}

	tfsdk.ValueFrom(ctx, languages, p.Languages.Type(ctx), &p.Languages)
}

func (p *LanguageProfile) read(ctx context.Context) *sonarr.LanguageProfileResource {
	langs := make([]string, len(p.Languages.Elements()))
	tfsdk.ValueAs(ctx, p.Languages, &langs)

	languages := make([]*sonarr.LanguageProfileItemResource, len(langs))

	for i, l := range langs {
		language := sonarr.NewLanguage()
		language.SetId(int32(helpers.GetLanguageID(l)))
		language.SetName(l)

		item := sonarr.NewLanguageProfileItemResource()
		item.SetAllowed(true)
		item.SetLanguage(*language)

		languages[i] = item
	}

	cutoff := sonarr.NewLanguage()
	cutoff.SetName(p.CutoffLanguage.ValueString())
	cutoff.SetId(int32(helpers.GetLanguageID(cutoff.GetName())))

	profile := sonarr.NewLanguageProfileResource()
	profile.SetName(p.Name.ValueString())
	profile.SetUpgradeAllowed(p.UpgradeAllowed.ValueBool())
	profile.SetLanguages(languages)
	profile.SetId(int32(p.ID.ValueInt64()))
	profile.SetCutoff(*cutoff)

	return profile
}
