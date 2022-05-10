package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	cmdList := NewCommandList()
	cmdConflict := NewCommandConflict()
	cmdReplace := NewCommandReplace()

	rootCmd := &cobra.Command{Use: "unityguid"}
	rootCmd.AddCommand(cmdList, cmdConflict, cmdReplace)
	rootCmd.Version = fmt.Sprintf("%s (%s)", VERSION, REVISION)
	rootCmd.InitDefaultVersionFlag()

	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
