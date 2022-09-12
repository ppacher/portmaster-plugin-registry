package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/ppacher/portmaster-plugin-registry/installer"
	"github.com/safing/portmaster/plugin/shared"
	"github.com/spf13/cobra"
)

var (
	installTarget string
	pluginsConfig string
)

var installCommand = &cobra.Command{
	Use:   "install index-file plugin-name",
	Short: "Install a plugin",
	Run: func(cmd *cobra.Command, args []string) {
		plg, err := getPluginDesc(args[0], args[1])
		if err != nil {
			hclog.L().Error(err.Error())
			os.Exit(1)
		}

		inst := installer.PluginInstaller{
			TargetDirectory: installTarget,
		}

		path, err := inst.InstallPlugin(context.Background(), plg)
		if err != nil {
			hclog.L().Error("failed to install plugin", "error", err)
			os.Exit(1)
		}

		err = updatePluginsConfig(pluginsConfig, installTarget, shared.PluginConfig{
			Name:             plg.Name,
			Types:            plg.PluginTypes,
			Privileged:       plg.Privileged,
			DisableAutostart: false,
		})
		if err != nil {
			hclog.L().Error("failed to configure plugin", "error", err)
			os.Exit(1)
		}

		hclog.L().Info("plugin successfully installed", "path", path)
	},
}

func init() {
	installCommand.Flags().StringVar(&installTarget, "target", "/opt/safing/portmaster/plugins", "The path where the plugin binary should be installed")
	installCommand.Flags().StringVar(&pluginsConfig, "config", "/opt/safing/portmaster/plugins.json", "The path to the portmaster plugins.json")
}

func updatePluginsConfig(pluginJson, pluginDir string, cfg shared.PluginConfig) error {
	// try to open and read the plugins.json file
	var cfgs []shared.PluginConfig

	blob, err := os.ReadFile(pluginJson)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read plugins.json: %w", err)
	}

	if err == nil {
		if err := json.Unmarshal(blob, &cfgs); err != nil {
			return fmt.Errorf("failed to parse plugins.json: %w", err)
		}
	}

	for idx, existing := range cfgs {
		if existing.Name == cfg.Name {
			cfgs = append(cfgs[:idx], cfgs[idx+1:]...)

			break
		}
	}

	// add the test plugin at the first position
	cfgs = append([]shared.PluginConfig{cfg}, cfgs...)

	// marshal the configuration and write the pluginJson
	blob, err = json.MarshalIndent(cfgs, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON configuration file: %w", err)
	}

	if err := os.WriteFile(pluginJson, blob, 0644); err != nil {
		return fmt.Errorf("failed to write plugins.json: %w", err)
	}

	return nil
}
