package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hidetatz/gcp_tf_iam_binding_validator"
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
		fmt.Fprintf(os.Stderr, "read dir: %v\n", err)
		os.Exit(1)
	}

	files := []string{}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if strings.HasSuffix(filename, ".tf") {
			files = append(files, filepath.Join(*dir, entry.Name()))
		}
	}

	rolesMap, err := gcp_tf_iam_binding_validator.CheckDuplication(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "check duplication: %v\n", err)
		os.Exit(1)
	}

	duplicated := false
	for role, ids := range rolesMap {
		if len(ids) > 1 {
			duplicated = true
			fmt.Fprintf(os.Stderr, "duplication found: role: %s, resources: %v\n", role, ids)
		}
	}

	if duplicated {
		os.Exit(1)
	}
}
