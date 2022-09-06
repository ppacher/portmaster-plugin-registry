package structs

type (
	// InstalledPlugin describes a plugin installed by the manager.
	InstalledPlugin struct {
		PluginDesc `json:",inline" hcl:",inline"`
		Path       string `hcl:"path"`
	}

	AvailableUpdate struct {
		Name           string `json:"name"`
		CurrentVersion string `json:"currentVersion"`
		NewVersion     string `json:"newVersion"`
	}

	InstalledPluginsFile struct {
		Version string            `hcl:"version"`
		Plugins []InstalledPlugin `hcl:"plugins,block"`
	}
)
