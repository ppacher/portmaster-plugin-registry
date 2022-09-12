package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/ppacher/portmaster-plugin-registry/installer"
	"github.com/spf13/cobra"
)

var downloadArtifactUrl = &cobra.Command{
	Use:   "download-plugin index-file plugin-name",
	Args:  cobra.ExactArgs(2),
	Short: "Download the plugin binary for your system architecture",
	Run: func(cmd *cobra.Command, args []string) {
		indexFile := args[0]
		pluginName := args[1]

		plg, err := getPluginDesc(indexFile, pluginName)
		if err != nil {
			hclog.L().Error("failed to find plugin", "error", err)
			os.Exit(1)
		}

		dst, err := installer.DownloadPlugin(context.Background(), "", plg)
		if err != nil {
			hclog.L().Error("failed to download plugin", "error", err)
			os.Exit(1)
		}

		fmt.Println(dst)
	},
}
