package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/ppacher/portmaster-plugin-registry/installer"
	"github.com/ppacher/portmaster-plugin-registry/manager"
	"github.com/ppacher/portmaster-plugin-registry/registry"
	"github.com/ppacher/portmaster-plugin-registry/structs"
	"github.com/safing/portmaster/plugin/framework"
	"github.com/safing/portmaster/plugin/framework/cmds"
	"github.com/spf13/cobra"
)

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
			//Privileged: true,
		}),
	)

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}

func bootstrapPlugin(ctx context.Context) error {
	provider := registry.NewRegistry()

	repos, err := loadRepositories()
	if err != nil {
		// TODO(ppacher): create a notification for that?

		return fmt.Errorf("failed to read repositories: %w", err)
	}

	for _, repo := range repos {
		if err := provider.AddRepository(repo); err != nil {
			// TODO(ppacher): create a notification for that?

			continue
		}
	}

	installer := &installer.PluginInstaller{
		TargetDirectory: filepath.Join(
			framework.BaseDirectory(),
			"plugins",
		),
	}

	stateFile := filepath.Join(
		framework.BaseDirectory(),
		"registry.state.hcl",
	)

	manager := manager.NewManager(stateFile, installer, provider, framework.PluginManager())

	// kick of the notification handler that will create error and update notifications.
	NewNotificationHandler(manager, framework.Notify())

	if err := manager.Start(framework.Context()); err != nil {
		return err
	}

	return nil
}

func loadRepositories() ([]structs.Repository, error) {
	repositoryFile := filepath.Join(
		framework.BaseDirectory(),
		"repositories.hcl",
	)

	var repos struct {
		Repositories []structs.Repository `hcl:"repository,block"`
	}

	blob, err := os.ReadFile(repositoryFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if err == nil {
		if err := hclsimple.Decode(repositoryFile, blob, nil, &repos); err != nil {
			return nil, err
		}

		if len(repos.Repositories) > 0 {
			return repos.Repositories, nil
		}
	}

	return []structs.Repository{
		{
			Name: "main",
			URL:  "https://raw.githubusercontent.com/ppacher/portmaster-plugin-registry/main/repository.hcl",
		},
	}, nil
}
