package main

import (
	"github.com/spf13/cobra"
)

func main() {
	cmdList := NewCommandList()
	cmdConflict := NewCommandConflict()
	cmdReplace := NewCommandReplace()

	rootCmd := &cobra.Command{Use: "unityguid"}
	rootCmd.AddCommand(cmdList, cmdConflict, cmdReplace)
	rootCmd.Version = "1.0.0"
	rootCmd.InitDefaultVersionFlag()
	rootCmd.Execute()
}
