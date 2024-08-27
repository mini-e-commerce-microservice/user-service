package main

import (
	"github.com/spf13/cobra"
	"user-service/internal/util"
)

func main() {
	var rootCmd = &cobra.Command{}

	rootCmd.AddCommand(restApiCmd)

	err := rootCmd.Execute()
	util.Panic(err)
}
