package main

import (
	"context"
	"log"
	"path/filepath"

	"github.com/ppacher/portmaster-plugin-registry/installer"
	"github.com/ppacher/portmaster-plugin-registry/manager"
	"github.com/ppacher/portmaster-plugin-registry/registry"
	"github.com/ppacher/portmaster-plugin-registry/structs"
	"github.com/safing/portmaster/plugin/framework"
	"github.com/safing/portmaster/plugin/framework/cmds"
	"github.com/spf13/cobra"
)

func bootstrapPlugin(ctx context.Context) error {
	provider := registry.NewRegistry()

	// TODO(ppacher): add support to load repository configurations from files
	// and just fallback to the default one if no files are found.
	provider.AddRepository(structs.Repository{
		Name: "main",
		URL:  "https://raw.githubusercontent.com/ppacher/portmaster-plugin-registry/main/repository.hcl",
	})

	installer := &installer.PluginInstaller{
		TargetDirectory: filepath.Join(
			framework.BaseDirectory(),
			"plugins",
		),
	}

	stateFile := filepath.Join(
		framework.BaseDirectory(),
		"plugin-registry.hcl",
	)

	manager := manager.NewManager(stateFile, installer, provider, framework.PluginManager())

	if err := manager.Start(framework.Context()); err != nil {
		return err
	}

	return nil
}

func main() {
	root := &cobra.Command{
		Use: "registry-plugin [install]",
		Run: func(cmd *cobra.Command, args []string) {
			framework.OnInit(bootstrapPlugin)
			framework.Serve()
		},
	}

	root.AddCommand(
		cmds.InstallCommand(&cmds.InstallCommandConfig{
			PluginName: "registry-plugin",
			Privileged: true,
		}),
	)

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
