package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func recursiveFetchMetaFiles(root string) {
	re := regexp.MustCompile("guid: (\\w+)\r?\n")
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".meta") {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		raw := string(data)
		matches := re.FindStringSubmatch(raw)
		if len(matches) > 1 {
			fmt.Printf("%s\t%s\n", matches[1], path)
		}
		return nil
	})
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: unityguid <project_root_path>")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	rootPath := "."
	rootPath = os.Args[1]

	recursiveFetchMetaFiles(rootPath)
}
