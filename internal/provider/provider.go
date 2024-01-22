package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure KeycloakUserCacheProvider satisfies various provider interfaces.
var _ provider.Provider = &KeycloakUserCacheProvider{}

// KeycloakUserCacheProvider defines the provider implementation.
type KeycloakUserCacheProvider struct {
	version string
}

// KeycloakUserCacheProviderModel describes the provider schema.
type KeycloakUserCacheProviderModel struct {
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Realm        types.String `tfsdk:"realm"`
	URL          types.String `tfsdk:"url"`
}

func (p *KeycloakUserCacheProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kuc"
	resp.Version = p.version
}

func (p *KeycloakUserCacheProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Optional: true,
			},
			"client_secret": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"realm": schema.StringAttribute{
				Optional: true,
			},
			"url": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *KeycloakUserCacheProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config KeycloakUserCacheProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientId := os.Getenv("KEYCLOAK_CLIENT_ID")
	clientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")
	url := os.Getenv("KEYCLOAK_URL")
	realm := os.Getenv("KEYCLOAK_REALM")

	if config.ClientID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Unknown client_id value",
			"Unknown client_id value",
		)
	}

	if config.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Unknown client_secret value",
			"Unknown client_secret value",
		)
	}

	if config.URL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Unknown url value",
			"Unknown url value",
		)
	}

	if config.Realm.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("realm"),
			"Unknown realm value",
			"Unknown realm value",
		)
	}

	if !config.ClientID.IsNull() {
		clientId = config.ClientID.ValueString()
	}

	if !config.ClientSecret.IsNull() {
		clientSecret = config.ClientSecret.ValueString()
	}

	if !config.URL.IsNull() {
		url = config.URL.ValueString()
	}

	if !config.Realm.IsNull() {
		realm = config.Realm.ValueString()
	}

	if clientId == "" {
		resp.Diagnostics.AddAttributeError(path.Root("client_id"), "client_id is required", "client_id is required")
		return
	}

	if clientSecret == "" {
		resp.Diagnostics.AddAttributeError(path.Root("client_secret"), "client_secret is required", "client_secret is required")
		return
	}

	if url == "" {
		resp.Diagnostics.AddAttributeError(path.Root("url"), "url is required", "url is required")
		return
	}

	if realm == "" {
		resp.Diagnostics.AddAttributeError(path.Root("realm"), "realm is required", "realm is required")
	}

	client := NewClient(config.ClientID.ValueString(), config.ClientSecret.ValueString(), config.URL.ValueString(), config.Realm.ValueString())

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *KeycloakUserCacheProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (p *KeycloakUserCacheProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserResource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &KeycloakUserCacheProvider{
			version: version,
		}
	}
}
