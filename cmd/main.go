package main

import (
	"github.com/mini-e-commerce-microservice/user-service/internal/util"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{}

	rootCmd.AddCommand(restApiCmd)

	err := rootCmd.Execute()
	util.Panic(err)
}
