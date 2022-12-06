package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const releaseProfileResourceName = "release_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ReleaseProfileResource{}
var _ resource.ResourceWithImportState = &ReleaseProfileResource{}

func NewReleaseProfileResource() resource.Resource {
	return &ReleaseProfileResource{}
}

// ReleaseProfileResource defines the release profile implementation.
type ReleaseProfileResource struct {
	client *sonarr.Sonarr
}

// ReleaseProfile describes the release profile data model.
type ReleaseProfile struct {
	Required  types.Set    `tfsdk:"required"`
	Ignored   types.Set    `tfsdk:"ignored"`
	Tags      types.Set    `tfsdk:"tags"`
	Name      types.String `tfsdk:"name"`
	ID        types.Int64  `tfsdk:"id"`
	IndexerID types.Int64  `tfsdk:"indexer_id"`
	Enabled   types.Bool   `tfsdk:"enabled"`
}

func (r *ReleaseProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + releaseProfileResourceName
}

func (r *ReleaseProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Profiles -->Release Profile resource.\nFor more information refer to [Release Profiles](https://wiki.servarr.com/sonarr/settings#release-profiles) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Release Profile ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enabled",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Release profile name.",
				Optional:            true,
				Computed:            true,
			},
			"indexer_id": schema.Int64Attribute{
				MarkdownDescription: "Indexer ID. Set `0` for all.",
				Optional:            true,
				Computed:            true,
			},
			"required": schema.SetAttribute{
				MarkdownDescription: "Required terms. At least one of `required` and `ignored` must be set.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"ignored": schema.SetAttribute{
				MarkdownDescription: "Ignored terms. At least one of `required` and `ignored` must be set.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (r *ReleaseProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ReleaseProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var profile *ReleaseProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	data := profile.read(ctx)

	// Create new ReleaseProfile
	response, err := r.client.AddReleaseProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", releaseProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "created"+releaseProfileResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *ReleaseProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var profile *ReleaseProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get releaseprofile current value
	response, err := r.client.GetReleaseProfileContext(ctx, profile.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", releaseProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+releaseProfileResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *ReleaseProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var profile *ReleaseProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := profile.read(ctx)

	// Update ReleaseProfile
	response, err := r.client.UpdateReleaseProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", releaseProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+releaseProfileResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *ReleaseProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var profile *ReleaseProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete releaseprofile current value
	err := r.client.DeleteReleaseProfileContext(ctx, profile.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", releaseProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+releaseProfileResourceName+": "+strconv.Itoa(int(profile.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *ReleaseProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+releaseProfileResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (p *ReleaseProfile) write(ctx context.Context, profile *sonarr.ReleaseProfile) {
	p.ID = types.Int64Value(profile.ID)
	p.Name = types.StringValue(profile.Name)
	p.Enabled = types.BoolValue(profile.Enabled)
	p.IndexerID = types.Int64Value(profile.IndexerID)
	p.Required = types.SetValueMust(types.StringType, nil)
	tfsdk.ValueFrom(ctx, profile.Required, p.Required.Type(ctx), &p.Required)
	p.Ignored = types.SetValueMust(types.StringType, nil)
	tfsdk.ValueFrom(ctx, profile.Ignored, p.Ignored.Type(ctx), &p.Ignored)
	p.Tags = types.SetValueMust(types.Int64Type, nil)
	tfsdk.ValueFrom(ctx, profile.Tags, p.Tags.Type(ctx), &p.Tags)
}

func (p *ReleaseProfile) read(ctx context.Context) *sonarr.ReleaseProfile {
	tags := make([]int, len(p.Tags.Elements()))
	tfsdk.ValueAs(ctx, p.Tags, &tags)

	required := make([]string, len(p.Required.Elements()))
	tfsdk.ValueAs(ctx, p.Required, &required)

	ignored := make([]string, len(p.Required.Elements()))
	tfsdk.ValueAs(ctx, p.Ignored, &ignored)

	return &sonarr.ReleaseProfile{
		Name:      p.Name.ValueString(),
		Enabled:   p.Enabled.ValueBool(),
		Required:  required,
		Ignored:   ignored,
		IndexerID: p.IndexerID.ValueInt64(),
		Tags:      tags,
		ID:        p.ID.ValueInt64(),
	}
}
