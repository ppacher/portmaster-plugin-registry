package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/ppacher/portmaster-plugin-registry/installer"
	"github.com/ppacher/portmaster-plugin-registry/structs"
	"github.com/spf13/cobra"
)

var getArtifactUrl = &cobra.Command{
	Use:   "get-url index-file plugin-name",
	Args:  cobra.ExactArgs(2),
	Short: "Get the download url for a plugin on your system architecture",
	Run: func(cmd *cobra.Command, args []string) {
		indexFile := args[0]
		pluginName := args[1]

		url, _, err := getDownloadURL(indexFile, pluginName)
		if err != nil {
			hclog.L().Error(err.Error())
			os.Exit(1)
		}

		fmt.Println(url)
	},
}

func getPluginDesc(indexFile string, pluginName string) (structs.PluginDesc, error) {
	index, err := loadAndVerifyIndex(indexFile)
	if err != nil {
		return structs.PluginDesc{}, err
	}

	for _, p := range index.Plugins {
		if p.Name == pluginName {
			return p, nil
		}
	}

	return structs.PluginDesc{}, fmt.Errorf("failed to find plugin in index")
}

func getDownloadURL(indexFile string, pluginName string) (string, string, error) {
	plg, err := getPluginDesc(indexFile, pluginName)
	if err != nil {
		return "", "", err
	}

	url, artifactFile, err := installer.FindMatchingArtifact(plg)
	if err != nil {
		return "", "", err
	}

	return url, artifactFile, nil
}
