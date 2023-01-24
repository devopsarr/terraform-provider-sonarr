package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const releaseProfileResourceName = "release_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ReleaseProfileResource{}
	_ resource.ResourceWithImportState = &ReleaseProfileResource{}
)

func NewReleaseProfileResource() resource.Resource {
	return &ReleaseProfileResource{}
}

// ReleaseProfileResource defines the release profile implementation.
type ReleaseProfileResource struct {
	client *sonarr.APIClient
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
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ReleaseProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var profile *ReleaseProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	request := profile.read(ctx)

	// Create new ReleaseProfile
	response, _, err := r.client.ReleaseProfileApi.CreateReleaseProfile(ctx).ReleaseProfileResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, releaseProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "created"+releaseProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
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
	response, _, err := r.client.ReleaseProfileApi.GetReleaseProfileById(ctx, int32(profile.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, releaseProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+releaseProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
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
	request := profile.read(ctx)

	// Update ReleaseProfile
	response, _, err := r.client.ReleaseProfileApi.UpdateReleaseProfile(ctx, strconv.Itoa(int(request.GetId()))).ReleaseProfileResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, releaseProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+releaseProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
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
	_, err := r.client.ReleaseProfileApi.DeleteReleaseProfile(ctx, int32(profile.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, releaseProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+releaseProfileResourceName+": "+strconv.Itoa(int(profile.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *ReleaseProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+releaseProfileResourceName+": "+req.ID)
}

func (p *ReleaseProfile) write(ctx context.Context, profile *sonarr.ReleaseProfileResource) {
	p.ID = types.Int64Value(int64(profile.GetId()))
	p.Name = types.StringValue(profile.GetName())
	p.Enabled = types.BoolValue(profile.GetEnabled())
	p.IndexerID = types.Int64Value(int64(profile.GetIndexerId()))
	p.Required = types.SetValueMust(types.StringType, nil)
	tfsdk.ValueFrom(ctx, profile.Required, p.Required.Type(ctx), &p.Required)
	p.Ignored = types.SetValueMust(types.StringType, nil)
	tfsdk.ValueFrom(ctx, profile.Ignored, p.Ignored.Type(ctx), &p.Ignored)
	p.Tags = types.SetValueMust(types.Int64Type, nil)
	tfsdk.ValueFrom(ctx, profile.Tags, p.Tags.Type(ctx), &p.Tags)
}

func (p *ReleaseProfile) read(ctx context.Context) *sonarr.ReleaseProfileResource {
	tags := make([]*int32, len(p.Tags.Elements()))
	tfsdk.ValueAs(ctx, p.Tags, &tags)

	required := make([]*string, len(p.Required.Elements()))
	tfsdk.ValueAs(ctx, p.Required, &required)

	ignored := make([]*string, len(p.Required.Elements()))
	tfsdk.ValueAs(ctx, p.Ignored, &ignored)

	profile := sonarr.NewReleaseProfileResource()
	profile.SetEnabled(p.Enabled.ValueBool())
	profile.SetId(int32(p.ID.ValueInt64()))
	profile.SetIgnored(ignored)
	profile.SetIndexerId(int32(p.IndexerID.ValueInt64()))
	profile.SetName(p.Name.ValueString())
	profile.SetRequired(required)
	profile.SetTags(tags)

	return profile
}
