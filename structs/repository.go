package structs

import "github.com/safing/portmaster/plugin/shared"

type (
	// Repository holds a repository configuration
	// that is used to fetch available plugins and to check
	// for updates.
	Repository struct {
		// Name is a human-friendly name used in the UI when the repository is
		// referenced.
		Name string `json:"name" hcl:",label"`

		// URL should holds the URL of the repository.[hcl|yaml|json] file that defines
		// all available plugins.
		URL string `json:"url" hcl:"url"`

		// Priority defines the priority for the repository. Usually there's only
		// one repository used by the registry but in case multiple repositories
		// are defined the priority decides which one wins when plugins are listed
		// in multiple repositories.
		Priority int `json:"priority" hcl:"priority"`
	}

	// IndexMeta holds additional information about a index file.
	IndexMeta struct {
		// Version is the version of the index file. This must be set
		// to v1.0.0 right now.
		Version string `json:"version" hcl:"version"`

		// Description may hold a human readable description of the repository.
		Description string `json:"description" hcl:"description,optional"`
	}

	// Artifact defines the download paths for different operating systems and
	// architectures.
	Artifact struct {
		// OS is the name of the operating system this artifact is built for.
		// This should follow the values from runtime.GOOS.
		OS string `json:"os" hcl:",label"`

		// Architecture based download URLs

		AMD64 string `json:"amd64" hcl:"amd64,optional"`
		ARM   string `json:"arm" hcl:"arm,optional"`
		ARM64 string `json:"arm64" hcl:"arm64,optional"`
		I386  string `json:"i386" hcl:"i386,optional"`
	}

	// PluginDesc describes a plugin and additional meta data.
	PluginDesc struct {
		// Name is the name of the plugin and must be unique across all
		// plugins listed in a repository.
		Name string `json:"name" hcl:",label"`

		// SourceURL should point to the source code repository of the plugin.
		SourceURL string `json:"source" hcl:"source"`

		// Version is the current version of the plugin.
		Version string `json:"version" hcl:"version"`

		// Artifacts defines the download URLs for the plugin binary for
		// different architectures and operating systems.
		Artifacts []Artifact `json:"artifacts" hcl:"artifact,block"`

		// PluginTypes defines the list of plugin types implemented
		// by the described plugin.
		PluginTypes []shared.PluginType `json:"pluginTypes" hcl:"pluginTypes"`

		// Author is the name of the plugin author.
		Author string `json:"author" hcl:"author,optional"`

		// License holds the license identifier for the plugin.
		License string `json:"license" hcl:"license,optional"`

		// Description holds a human readable description of the features and purpose of
		// a plugin.
		Description string `json:"description" hcl:"description,optional"`

		// Tags holds an arbitrary list of tags for the plugin.
		Tags []string `json:"tags" hcl:"tags,optional"`

		// Repository is the name of the repository that contains the
		// plugin.
		Repository string `json:"repository" hcl:"-"`
	}

	// RepositoryIndex defines the structure of a repository index file.
	RepositoryIndex struct {
		// Meta holds meta-data about the repository.
		Meta IndexMeta `json:"meta" hcl:"meta,block"`

		// Plugins is the list of plugins available in the repository.
		Plugins []PluginDesc `json:"plugins" hcl:"plugin,block"`
	}
)
