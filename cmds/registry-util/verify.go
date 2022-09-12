package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-getter/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/ppacher/portmaster-plugin-registry/registry"
	"github.com/ppacher/portmaster-plugin-registry/structs"
	"github.com/spf13/cobra"
)

var verifyIndexCommand = &cobra.Command{
	Use: "verify-index [path] [path...]",
	Run: func(cmd *cobra.Command, args []string) {
		hasErrors := false

		for _, p := range args {
			_, err := loadAndVerifyIndex(p)
			if err != nil {
				hclog.L().Error(p, "error", err)
				hasErrors = true
			}
		}

		if hasErrors {
			os.Exit(1)
		}
	},
}

func loadAndVerifyIndex(path string) (*structs.RepositoryIndex, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	tempDir, err := os.MkdirTemp("", "*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	file, err := new(getter.Client).Get(context.Background(), &getter.Request{
		Src: path,
		Dst: tempDir,
		Pwd: pwd,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get repository index: %w", err)
	}

	f, err := os.Open(file.Dst)
	if err != nil {
		return nil, fmt.Errorf("failed to open index file: %w", err)
	}

	index, err := registry.DecodeIndex(path, f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode index file: %w", err)
	}

	if err := registry.ValidateIndex(index); err != nil {
		return nil, fmt.Errorf("failed to decode index file: %w", err)
	}

	return index, nil
}
