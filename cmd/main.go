package main

import (
	"github.com/SyaibanAhmadRamadhan/go-collection"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{}

	rootCmd.AddCommand(restApiCmd)

	err := rootCmd.Execute()
	collection.PanicIfErr(err)
}
