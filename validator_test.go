package gcp_tf_iam_binding_validator

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestCheckDuplication(t *testing.T) {
	tests := []struct {
		dir  string
		want map[string][]string
	}{
		{
			dir: "./test/1",
			want: map[string][]string{
				"roles/storage.admin":        []string{"binding_1", "binding_2"},
				"roles/storage.objectViewer": []string{"binding_3"},
			},
		},
		{
			dir: "./test/2",
			want: map[string][]string{
				"roles/storage.admin":        []string{"binding_1"},
				"roles/storage.editor":       []string{"binding_2"},
				"roles/storage.objectViewer": []string{"binding_3"},
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

			rolesMap, err := CheckDuplication(files)
			if err != nil {
				t.Fatalf("check duplication: %v", err)
			}

			if !reflect.DeepEqual(rolesMap, tt.want) {
				t.Fatalf("fail: want: %v, got: %v", tt.want, rolesMap)
			}
		})
	}
}
