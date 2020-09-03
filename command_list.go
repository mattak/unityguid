package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type OptionListCommand struct {
	ExcludePatterns       []string
	IncludePatterns       []string
	TrimBeforeAssetFolder bool
}

func containsPathPatterns(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}

func trimAssetPath(path string) string {
	return trimBeforeSpecificFolder(path, "/Assets/")
}

func trimProjectSettingsPath(path string) string {
	return trimBeforeSpecificFolder(path, "/ProjectSettings/")
}

// folderName: should be start and endswith "/". e.g. "/Assets/"
func trimBeforeSpecificFolder(path string, folderName string) string {
	index := strings.Index(path, folderName)
	if index == -1 {
		return path
	}

	return path[(index + 1):]
}

func printRecursiveFetchMetaFiles(root string, option OptionListCommand) {
	guidRegex := regexp.MustCompile("guid: (\\w+)\r?\n")

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
		matches := guidRegex.FindStringSubmatch(raw)

		if len(matches) > 1 {
			if containsPathPatterns(path, option.ExcludePatterns) {
				return nil
			}
			if len(option.IncludePatterns) > 0 && !containsPathPatterns(path, option.IncludePatterns) {
				return nil
			}

			shortPath := path
			if option.TrimBeforeAssetFolder {
				shortPath = trimAssetPath(shortPath)
				shortPath = trimProjectSettingsPath(shortPath)
			}

			fmt.Printf("%s\t%s\n", matches[1], shortPath)
		}

		return nil
	})
}

func NewCommandList() *cobra.Command {
	option := OptionListCommand{
		ExcludePatterns:       []string{},
		TrimBeforeAssetFolder: true,
	}
	cmd := &cobra.Command{
		Use:   "list [folder to find meta files]",
		Short: "list up guid meta & file name",
		Long:  `list up guid by find .meta file, output format is tsv of [{<guid>,<filename>}+]`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			rootPath := "."
			if len(args) >= 1 {
				rootPath = args[0]
			}

			printRecursiveFetchMetaFiles(rootPath, option)
		},
	}
	cmd.Flags().StringSliceVarP(&option.ExcludePatterns, "exclude", "e", []string{}, "exclude patterns. It can be split with \",\". eg. \"/Modules/,/Tests/\"")
	cmd.Flags().StringSliceVarP(&option.IncludePatterns, "include", "i", []string{}, "include patterns. It can be split with \",\". eg. \"/Modules/,/Tests/\"")
	cmd.Flags().BoolVarP(&option.TrimBeforeAssetFolder, "trim", "t", true, "trim before directory of asset folder")
	return cmd
}
