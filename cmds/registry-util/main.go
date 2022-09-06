package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/ppacher/portmaster-plugin-registry/registry"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use: "registry-util",
}

var verifyCommand = &cobra.Command{
	Use: "verify [path] [path...]",
	Run: func(cmd *cobra.Command, args []string) {
		hasErrors := false

		for _, p := range args {
			f, err := os.Open(p)
			if err != nil {
				hclog.L().Error("failed to open index file", "file", p, "error", err.Error())
				hasErrors = true

				continue
			}

			index, err := registry.DecodeIndex(p, f)
			if err != nil {
				hclog.L().Error("failed to decode index file", "file", p, "error", err.Error())
				hasErrors = true

				continue
			}

			if err := registry.ValidateIndex(index); err != nil {
				hclog.L().Error("failed to decode index file", "file", p, "error", err.Error())
				hasErrors = true

				continue
			}
		}

		if hasErrors {
			os.Exit(1)
		}
	},
}

func main() {
	root.AddCommand(
		verifyCommand,
	)

	if err := root.Execute(); err != nil {
		hclog.L().Error(err.Error())
		os.Exit(1)
	}
}
