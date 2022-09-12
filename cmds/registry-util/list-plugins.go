package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
)

var listPluginsCommand = &cobra.Command{
	Use:   "list-plugins index-file",
	Short: "List all plugins in a repository index",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		index, err := loadAndVerifyIndex(args[0])
		if err != nil {
			hclog.L().Error("failed to get repository index", "error", err)
			os.Exit(1)
		}

		bullet := color.New(color.FgGreen).Sprint("•")
		pluginHeader := color.New(color.Bold, color.FgHiWhite).Sprint
		description := color.New(color.Italic).Sprint

		for _, plg := range index.Plugins {
			fmt.Printf(bullet+" %s %s\n", pluginHeader(plg.Name), description(plg.Version))
			fmt.Println("  " + description(plg.Description))
			fmt.Println("  by " + description(plg.Author))

			fmt.Println()
		}
	},
}
