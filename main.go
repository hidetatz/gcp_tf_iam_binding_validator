package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

func main() {
	var dir = flag.String("dir", "", "target directory (All files under the dir will be target, non-recursive)")
	flag.Parse()

	if *dir == "" {
		fmt.Fprintln(os.Stderr, "dir must not be empty")
		os.Exit(1)
	}

	dirEntries, err := os.ReadDir(*dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read dir: %v", err)
		os.Exit(1)
	}

	files := []string{}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if strings.HasSuffix(filename, ".tf") || strings.HasSuffix(filename, ".json") {
			files = append(files, filepath.Join(*dir, entry.Name()))
		}
	}

	if err := validateFiles(files); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Root struct {
	Resources []*struct {
		Kind    string   `hcl:"type,label"`
		ID      string   `hcl:"id,label"`
		HCLBody hcl.Body `hcl:",remain"`
	} `hcl:"resource,block"`

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
}

type GoogleProjectIAMBinding struct {
	ID   string
	Role string
}

func validateFiles(files []string) error {
	bodys := make([]hcl.Body, len(files))
	for i := range files {
		b, err := ParseFile(files[i])
		if err != nil {
			return err
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
			return fmt.Errorf("decode whole body: %w", diags)
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
				return fmt.Errorf("decode google_project_iam_binding: %w", diags)
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

	duplicated := false
	for role, ids := range rolesMap {
		if len(ids) > 1 {
			duplicated = true
			fmt.Fprintf(os.Stderr, "duplication found: role: %s, resources: %v\n", role, ids)
		}
	}

	if duplicated {
		return fmt.Errorf("validation failed")
	}

	return nil
}

func validateDirs(files []string) error {
	return nil
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
