package gcp_tf_iam_binding_validator

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type Root struct {
	Resources []*struct {
		Kind    string   `hcl:"type,label"`
		ID      string   `hcl:"id,label"`
		HCLBody hcl.Body `hcl:",remain"`
	} `hcl:"resource,block"`

	// Other resources are unused, but must define for decode.
	Data []*struct {
		Kind    string   `hcl:"type,label"`
		ID      string   `hcl:"id,label"`
		HCLBody hcl.Body `hcl:",remain"`
	} `hcl:"data,block"`

	Modules []*struct {
		Kind    string   `hcl:"type,label"`
		HCLBody hcl.Body `hcl:",remain"`
	} `hcl:"module,block"`

	Locals []*struct {
		HCLBody hcl.Body `hcl:",remain"`
	} `hcl:"locals,block"`

	Terraforms []*struct {
		HCLBody hcl.Body `hcl:",remain"`
	} `hcl:"terraform,block"`

	Providers []*struct {
		Kind    string   `hcl:"type,label"`
		HCLBody hcl.Body `hcl:",remain"`
	} `hcl:"provider,block"`

	Outputs []*struct {
		Kind    string   `hcl:"type,label"`
		HCLBody hcl.Body `hcl:",remain"`
	} `hcl:"output,block"`

	Variables []*struct {
		Kind    string   `hcl:"type,label"`
		HCLBody hcl.Body `hcl:",remain"`
	} `hcl:"variable,block"`
}

type GoogleProjectIAMBinding struct {
	ID   string
	Role string
}

func CheckDuplication(files []string) (map[string][]string, error) {
	bodys := make([]hcl.Body, len(files))
	for i := range files {
		b, err := ParseFile(files[i])
		if err != nil {
			return nil, err
		}
		bodys[i] = b
	}

	var googleProjectIAMBindings []*GoogleProjectIAMBinding

	// key: role, value: ids.
	// This is used to make sure every role in google_project_iam_bindings are unique.
	rolesMap := map[string][]string{}

	for _, body := range bodys {
		var root Root
		if diags := gohcl.DecodeBody(body, nil, &root); diags.HasErrors() {
			return nil, fmt.Errorf("decode whole body: %w", diags)
		}

		for _, resource := range root.Resources {
			if resource.Kind != "google_project_iam_binding" {
				continue
			}

			var buff struct {
				Role    string   `hcl:"role"`
				HCLBody hcl.Body `hcl:",remain"` // rest does not matter for validation
			}

			if diags := gohcl.DecodeBody(resource.HCLBody, nil, &buff); diags.HasErrors() {
				return nil, fmt.Errorf("decode google_project_iam_binding: %w", diags)
			}

			googleProjectIAMBindings = append(
				googleProjectIAMBindings,
				&GoogleProjectIAMBinding{
					ID:   resource.ID,
					Role: buff.Role,
				},
			)
		}

	}

	for _, binding := range googleProjectIAMBindings {
		ids, ok := rolesMap[binding.Role]
		if ok {
			rolesMap[binding.Role] = append(ids, binding.ID)
			continue
		}

		rolesMap[binding.Role] = []string{binding.ID}
	}

	return rolesMap, nil
}

func ParseFile(filename string) (hcl.Body, error) {
	parser := hclparse.NewParser()

	var (
		hclF        *hcl.File
		diagnostics hcl.Diagnostics
	)

	if strings.HasSuffix(".json", filename) {
		hclF, diagnostics = parser.ParseJSONFile(filename)
	} else {
		hclF, diagnostics = parser.ParseHCLFile(filename)
	}

	if diagnostics.HasErrors() {
		return nil, fmt.Errorf("parse file: %v", diagnostics.Error())
	}

	return hclF.Body, nil
}
