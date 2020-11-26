package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const regexBase string = "/component/[a-zA-Z0-9-.]+(/service/[a-zA-Z0-9-.]+(/[a-zA-Z0-9-.]+)*)?"
const regexCustomer string = "/customer/[a-zA-Z0-9-.]+(/environment/[a-zA-Z0-9-.]+)?(/site/[a-zA-Z0-9-.]+)?"
const regexDefaultConfig string = "^.*/config" + regexBase + "$"
const regexCustomerCofig string = "^.*/config" + regexCustomer + "(" + regexBase + ")?$"
const regexAllowedExtensions string = "(\\.yml|json)?"

func checkPaths(paths []string) []string {
	var wrongPaths []string

	defaultConfig := regexp.MustCompile(regexDefaultConfig)
	customerConfig := regexp.MustCompile(regexCustomerCofig)
	allowedExtensions := regexp.MustCompile(regexAllowedExtensions)

	for _, path := range paths {
		extension := filepath.Ext(path)
		path := filepath.ToSlash(path)

		switch {
		case !allowedExtensions.Match([]byte(extension)):
			fmt.Printf("Wrong extension: %s\n", path)
			wrongPaths = append(wrongPaths, path)
		case defaultConfig.Match([]byte(path)):
			continue
		case customerConfig.Match([]byte(path)):
			continue
		default:
			fmt.Printf("Wrong default: %s\n", path)
			wrongPaths = append(wrongPaths, path)
		}
	}

	return wrongPaths
}

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
