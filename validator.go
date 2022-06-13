package gcp_tf_iam_binding_validator

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/tmccombs/hcl2json/convert"
)

func CheckDuplication(files []string) (map[string][]string, error) {
	ret := map[string][]string{}
	for i := range files {
		f, err := ParseFile(files[i])
		if err != nil {
			return nil, err
		}

		j, err := convert.File(f, convert.Options{})
		if err != nil {
			return nil, err
		}

		var mapped interface{}
		if err := json.Unmarshal(j, &mapped); err != nil {
			return nil, err
		}

		r := mapped.(map[string]interface{})["resource"]
		if r == nil {
			// not a resource
			continue
		}

		g := r.(map[string]interface{})["google_project_iam_binding"]
		if g == nil {
			// not a google_project_iam_binding
			continue
		}

		b, ok := g.(map[string]interface{})
		if !ok {
			continue
		}

		// the map is used to find duplication in the set of google_project_iam_binding.
		// The key is a formatted string: "$role_$project_$conditionTitle_$conditionDescription_$conditionExpression".
		// The value is a name of the resource.
		bindings := map[string]struct{}{}

		for name, body := range b {
			content := body.([]interface{})[0].(map[string]interface{})
			role := content["role"].(string)
			project := content["project"].(string)
			condition, ok := content["condition"].([]interface{})
			var condTitle, condDesc, condExpr string
			if ok {
				condTitle, ok = condition[0].(map[string]interface{})["title"].(string)
				if !ok {
					condTitle = ""
				}

				condExpr, ok = condition[0].(map[string]interface{})["expression"].(string)
				if !ok {
					condExpr = ""
				}

				condDesc, ok = condition[0].(map[string]interface{})["description"].(string)
				if !ok {
					condDesc = ""
				}
			}

			key := fmt.Sprintf("%s_%s_%s_%s_%s", role, project, condTitle, condDesc, condExpr)

			fmt.Println(key)

			if _, found := bindings[key]; found {
				ret[role] = append(ret[role], name)
			} else {
				bindings[key] = struct{}{}
			}
		}
	}

	return ret, nil

}

func ParseFile(filename string) (*hcl.File, error) {
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

	return hclF, nil
}
