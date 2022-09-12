package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use: "registry-util",
}

func main() {
	root.AddCommand(
		verifyIndexCommand,
		getArtifactUrl,
		downloadArtifactUrl,
		installCommand,
		listPluginsCommand,
	)

	if err := root.Execute(); err != nil {
		hclog.L().Error(err.Error())
		os.Exit(1)
	}
}
