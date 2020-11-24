package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const regexStructure string = "^.*/(component|customer)(/[a-zA-Z0-9-]+)?(/(service|component|environment)/[a-zA-Z0-9-]+)?(/(service|component|site)/[a-zA-Z0-9-]+)?(/(service|component)/[a-zA-Z0-9-]+)?(/(service)/[a-zA-Z0-9-]+)?(/[a-zA-Z0-9-]+)(/[a-zA-Z0-9-]+)?(/[a-zA-Z0-9-]+)?\\.(yml|json)$"

// ValidateStructure is
func ValidateStructure(config Config) {
	paths := getPaths(config.Path)

	wrongPaths := checkPaths(paths)

	if len(wrongPaths) > 0 {
		fmt.Printf("Folder structure check failed, wrong paths: %s\n", wrongPaths)
		os.Exit(1)
	} else {
		fmt.Print("Folder structure validated successfully")
		os.Exit(0)
	}
}

func getPaths(rootPath string) []string {
	var results []string

	err := filepath.Walk(rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				shortPath := strings.ReplaceAll(path, rootPath, "config")
				results = append(results, shortPath)
			}

			return nil
		})

	check(err)

	return results
}

func checkPaths(paths []string) []string {
	var wrongPaths []string

	re := regexp.MustCompile(regexStructure)

	for _, path := range paths {
		path := filepath.ToSlash(path)

		if !re.Match([]byte(path)) {
			wrongPaths = append(wrongPaths, path)
		}
	}

	return wrongPaths
}
