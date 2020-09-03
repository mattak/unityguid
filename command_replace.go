package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type OptionReplaceCommand struct {
	ExcludePatterns       []string
	IncludePatterns       []string
	TrimBeforeAssetFolder bool
}

func newUnityGuid() string {
	newGuid, _ := uuid.NewRandom()
	return strings.ReplaceAll(newGuid.String(), "-", "")
}

func newGuidMap(guids []string) map[string]string {
	guidMap := map[string]string{}
	for _, guid := range guids {
		guidMap[guid] = newUnityGuid()
	}
	return guidMap
}

func replaceGuid(content string, oldGuid string, newGuid string) string {
	return strings.ReplaceAll(content, oldGuid, newGuid)
}

func overwriteFile(file string, content string, mode os.FileMode) {
	err := ioutil.WriteFile(file, []byte(content), mode)

	if err != nil {
		log.Fatalln(err)
	}
}

func runCommandReplace(rootAssetDir string, guids []string, option OptionReplaceCommand) {
	guidMap := newGuidMap(guids)

	err := filepath.Walk(rootAssetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		stat, err := os.Stat(path)
		if err != nil {
			return err
		}

		if stat.IsDir() {
			return nil
		}

		if containsPathPatterns(path, option.ExcludePatterns) {
			return nil
		}

		if len(option.IncludePatterns) > 0 && !containsPathPatterns(path, option.IncludePatterns) {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		content := string(data)
		matched := false

		for _, guid := range guids {
			if strings.Contains(content, guid) {
				shortPath := path
				if option.TrimBeforeAssetFolder {
					shortPath = trimAssetPath(shortPath)
					shortPath = trimProjectSettingsPath(shortPath)
				}

				fmt.Printf("%s => %s\t%s\n", guid, guidMap[guid], shortPath)
				matched = true
				content = replaceGuid(content, guid, guidMap[guid])
			}
		}

		if matched {
			overwriteFile(path, content, stat.Mode())
		}

		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}
}

func NewCommandReplace() *cobra.Command {
	option := OptionReplaceCommand{
		ExcludePatterns: []string{},
	}
	cmd := &cobra.Command{
		Use:   "replace [root_project_dir] [guid]+",
		Short: "replace specified guids",
		Long:  `replace conflict guids from root asset dir. each tsv list file. the file is created by list command. output format is tsv of [{<guid>,<target_basename>,<conflict_base_filename>,<conflict_target_filename>}+]`,
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			rootAssetDir := args[0]
			guids := args[1:]
			runCommandReplace(rootAssetDir, guids, option)
		},
	}
	cmd.Flags().StringSliceVarP(&option.ExcludePatterns, "exclude", "e", []string{}, "exclude patterns. It can be split with \",\". eg. \"/Modules/,/Tests/\"")
	cmd.Flags().StringSliceVarP(&option.IncludePatterns, "include", "i", []string{}, "include patterns. It can be split with \",\". eg. \"/Modules/,/Tests/\"")
	cmd.Flags().BoolVarP(&option.TrimBeforeAssetFolder, "trim", "t", true, "trim before directory of asset folder")

	return cmd
}
