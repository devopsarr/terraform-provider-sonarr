package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const (
	downloadClientHadoukenResourceName   = "download_client_hadouken"
	DownloadClientHadoukenImplementation = "Hadouken"
	DownloadClientHadoukenConfigContrat  = "HadoukenSettings"
	DownloadClientHadoukenProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DownloadClientHadoukenResource{}
var _ resource.ResourceWithImportState = &DownloadClientHadoukenResource{}

func NewDownloadClientHadoukenResource() resource.Resource {
	return &DownloadClientHadoukenResource{}
}

// DownloadClientHadoukenResource defines the download client implementation.
type DownloadClientHadoukenResource struct {
	client *sonarr.Sonarr
}

// DownloadClientHadouken describes the download client data model.
type DownloadClientHadouken struct {
	Tags                     types.Set    `tfsdk:"tags"`
	Name                     types.String `tfsdk:"name"`
	Host                     types.String `tfsdk:"host"`
	URLBase                  types.String `tfsdk:"url_base"`
	Username                 types.String `tfsdk:"username"`
	Password                 types.String `tfsdk:"password"`
	Category                 types.String `tfsdk:"category"`
	Priority                 types.Int64  `tfsdk:"priority"`
	Port                     types.Int64  `tfsdk:"port"`
	ID                       types.Int64  `tfsdk:"id"`
	UseSsl                   types.Bool   `tfsdk:"use_ssl"`
	Enable                   types.Bool   `tfsdk:"enable"`
	RemoveFailedDownloads    types.Bool   `tfsdk:"remove_failed_downloads"`
	RemoveCompletedDownloads types.Bool   `tfsdk:"remove_completed_downloads"`
}

func (d DownloadClientHadouken) toDownloadClient() *DownloadClient {
	return &DownloadClient{
		Tags:                     d.Tags,
		Name:                     d.Name,
		Host:                     d.Host,
		URLBase:                  d.URLBase,
		Username:                 d.Username,
		Password:                 d.Password,
		Category:                 d.Category,
		Priority:                 d.Priority,
		Port:                     d.Port,
		ID:                       d.ID,
		UseSsl:                   d.UseSsl,
		Enable:                   d.Enable,
		RemoveFailedDownloads:    d.RemoveFailedDownloads,
		RemoveCompletedDownloads: d.RemoveCompletedDownloads,
	}
}

func (d *DownloadClientHadouken) fromDownloadClient(client *DownloadClient) {
	d.Tags = client.Tags
	d.Name = client.Name
	d.Host = client.Host
	d.URLBase = client.URLBase
	d.Username = client.Username
	d.Password = client.Password
	d.Category = client.Category
	d.Priority = client.Priority
	d.Port = client.Port
	d.ID = client.ID
	d.UseSsl = client.UseSsl
	d.Enable = client.Enable
	d.RemoveFailedDownloads = client.RemoveFailedDownloads
	d.RemoveCompletedDownloads = client.RemoveCompletedDownloads
}

func (r *DownloadClientHadoukenResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientHadoukenResourceName
}

func (r *DownloadClientHadoukenResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Download Client Hadouken resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/sonarr/settings#download-clients) and [Hadouken](https://wiki.servarr.com/sonarr/supported#hadouken).",
		Attributes: map[string]tfsdk.Attribute{
			"enable": {
				MarkdownDescription: "Enable flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"remove_completed_downloads": {
				MarkdownDescription: "Remove completed downloads flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"remove_failed_downloads": {
				MarkdownDescription: "Remove failed downloads flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"priority": {
				MarkdownDescription: "Priority.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"name": {
				MarkdownDescription: "Download Client name.",
				Required:            true,
				Type:                types.StringType,
			},
			"tags": {
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"id": {
				MarkdownDescription: "Download Client ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			// Field values
			"use_ssl": {
				MarkdownDescription: "Use SSL flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"port": {
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"host": {
				MarkdownDescription: "host.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"url_base": {
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"username": {
				MarkdownDescription: "Username.",
				Required:            true,
				Type:                types.StringType,
			},
			"password": {
				MarkdownDescription: "Password.",
				Required:            true,
				Type:                types.StringType,
			},
			"category": {
				MarkdownDescription: "Category.",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (r *DownloadClientHadoukenResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DownloadClientHadoukenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClientHadouken

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClientHadouken
	request := client.read(ctx)

	response, err := r.client.AddDownloadClientContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", downloadClientHadoukenResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientHadoukenResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientHadoukenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClientHadouken

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClientHadouken current value
	response, err := r.client.GetDownloadClientContext(ctx, client.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", downloadClientHadoukenResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientHadoukenResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientHadoukenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClientHadouken

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClientHadouken
	request := client.read(ctx)

	response, err := r.client.UpdateDownloadClientContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", downloadClientHadoukenResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientHadoukenResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientHadoukenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var client *DownloadClientHadouken

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClientHadouken current value
	err := r.client.DeleteDownloadClientContext(ctx, client.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", downloadClientHadoukenResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientHadoukenResourceName+": "+strconv.Itoa(int(client.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientHadoukenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+downloadClientHadoukenResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (d *DownloadClientHadouken) write(ctx context.Context, downloadClient *sonarr.DownloadClientOutput) {
	genericDownloadClient := DownloadClient{
		Enable:                   types.BoolValue(downloadClient.Enable),
		RemoveCompletedDownloads: types.BoolValue(downloadClient.RemoveCompletedDownloads),
		RemoveFailedDownloads:    types.BoolValue(downloadClient.RemoveFailedDownloads),
		Priority:                 types.Int64Value(int64(downloadClient.Priority)),
		ID:                       types.Int64Value(downloadClient.ID),
		Name:                     types.StringValue(downloadClient.Name),
		Tags:                     types.SetValueMust(types.Int64Type, nil),
	}
	tfsdk.ValueFrom(ctx, downloadClient.Tags, genericDownloadClient.Tags.Type(ctx), &genericDownloadClient.Tags)
	genericDownloadClient.writeFields(ctx, downloadClient.Fields)
	d.fromDownloadClient(&genericDownloadClient)
}

func (d *DownloadClientHadouken) read(ctx context.Context) *sonarr.DownloadClientInput {
	var tags []int

	tfsdk.ValueAs(ctx, d.Tags, &tags)

	return &sonarr.DownloadClientInput{
		Enable:                   d.Enable.ValueBool(),
		RemoveCompletedDownloads: d.RemoveCompletedDownloads.ValueBool(),
		RemoveFailedDownloads:    d.RemoveFailedDownloads.ValueBool(),
		Priority:                 int(d.Priority.ValueInt64()),
		ID:                       d.ID.ValueInt64(),
		ConfigContract:           DownloadClientHadoukenConfigContrat,
		Implementation:           DownloadClientHadoukenImplementation,
		Name:                     d.Name.ValueString(),
		Protocol:                 DownloadClientHadoukenProtocol,
		Tags:                     tags,
		Fields:                   d.toDownloadClient().readFields(ctx),
	}
}
