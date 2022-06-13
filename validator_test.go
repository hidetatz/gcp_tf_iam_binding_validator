package gcp_tf_iam_binding_validator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCheckDuplication(t *testing.T) {
	tests := []struct {
		dir  string
		want []*GoogleProjectIAMBinding
	}{
		{
			dir: "./test/1",
			want: []*GoogleProjectIAMBinding{
				{
					Names:          []string{"binding_1", "binding_2"},
					Role:           "roles/storage.admin",
					Project:        "${var.project}",
					ConditionTitle: "",
					ConditionDesc:  "",
					ConditionExpr:  "",
				},
				{
					Names:          []string{"binding_3"},
					Role:           "roles/storage.objectViewer",
					Project:        "${var.project}",
					ConditionTitle: "",
					ConditionDesc:  "",
					ConditionExpr:  "",
				},
				{
					Names:          []string{"binding_4", "binding_5"},
					Role:           "roles/storage.admin",
					Project:        "${var.project}",
					ConditionTitle: "expires_after_2019_12_31",
					ConditionDesc:  "",
					ConditionExpr:  `request.time < timestamp("2020-01-01T00:00:00Z")`,
				},
			},
		},
		{
			dir: "./test/2",
			want: []*GoogleProjectIAMBinding{
				{
					Names:          []string{"binding_1"},
					Role:           "roles/storage.admin",
					Project:        "${var.project}",
					ConditionTitle: "",
					ConditionDesc:  "",
					ConditionExpr:  "",
				},
				{
					Names:          []string{"binding_2"},
					Role:           "roles/storage.editor",
					Project:        "${var.project}",
					ConditionTitle: "",
					ConditionDesc:  "",
					ConditionExpr:  "",
				},
				{
					Names:          []string{"binding_3"},
					Role:           "roles/storage.objectViewer",
					Project:        "${var.project}",
					ConditionTitle: "temporary",
					ConditionDesc:  "temporary permission",
					ConditionExpr:  `request.time < timestamp("2020-01-01T00:00:00Z")`,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.dir, func(t *testing.T) {
			dirEntries, err := os.ReadDir(tt.dir)
			if err != nil {
				t.Fatalf("read dir: %v", err)
			}

			files := []string{}
			for _, entry := range dirEntries {
				if entry.IsDir() {
					continue
				}

				filename := entry.Name()
				if strings.HasSuffix(filename, ".tf") || strings.HasSuffix(filename, ".json") {
					files = append(files, filepath.Join(tt.dir, entry.Name()))
				}
			}

			bindings, err := FindGoogleProjectIAMBindings(files)
			if err != nil {
				t.Fatalf("check duplication: %v", err)
			}

			if diff := cmp.Diff(bindings, tt.want); diff != "" {
				t.Fatalf("%s\n", diff)
			}
		})
	}
}
