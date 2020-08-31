package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

func printConflicts(tag string, baseMap map[string]string, targetMap map[string]string) {
	for guid, baseFile := range baseMap {
		if targetFile, ok := targetMap[guid]; ok {
			// 1234abcd module1 Assets/Script1.cs.meta Assets/Script2.cs.meta
			fmt.Printf("%s\t%s\t%s\t%s\n", guid, tag, baseFile, targetFile)
		}
	}
}

func readGuidMap(path string) (map[string]string, error) {
	result := map[string]string{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		cells := strings.Split(line, "\t")

		if len(cells) != 2 {
			continue
		}

		guid := cells[0]
		file := cells[1]
		result[guid] = file
	}

	return result, nil
}

func runCommandConflict(baseFile string, targetFiles []string) {
	baseMap, err := readGuidMap(baseFile)
	if err != nil {
		log.Fatalln("cannot read base file: " + baseFile)
	}

	for i := 0; i < len(targetFiles); i++ {
		targetFile := targetFiles[i]
		targetMap, err := readGuidMap(targetFile)

		if err != nil {
			log.Fatalln("cannot read target file: " + targetFile)
		}

		name := path.Base(targetFile)
		extension := path.Ext(targetFile)
		tag := name[0 : len(name)-len(extension)]
		printConflicts(tag, baseMap, targetMap)
	}
}

func NewCommandConflict() *cobra.Command {
	cmdConflict := &cobra.Command{
		Use:   "conflict [base_guid_filemap_tsv] [target_guidfilemap_tsv]+",
		Short: "list up conflict guids",
		Long:  `list up conflict guids by each tsv list file. the file is created by list command. output format is tsv of [{<guid>,<target_basename>,<conflict_base_filename>,<conflict_target_filename>}+]`,
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			baseFile := args[0]
			targetFiles := args[1:]
			runCommandConflict(baseFile, targetFiles)
		},
	}
	return cmdConflict
}
