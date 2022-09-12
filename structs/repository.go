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

		// ArchiveFile holds the name of the plugin binary if the downloaded artifact is
		// an archive and contains more than one file.
		ArchiveFile string `json:"archiveFile" hcl:"archive_file,optional"`

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

		// ArtifactTemplate contains a template string using github.com/valyala/fasttemplate
		// that is used to craft a download link for the target architecture.
		//
		// If it's not possible to define a download link using templating
		// the artifacts can be defined in the Artifacts member below.
		//
		// If the resulting URL ends in .tar.gz it will be unpacked automatically.
		//
		// For the template, the following substitutions are available:
		//	- {{os}}: The value of runtime.GOOS
		//  - {{arch}}: The value of runtime.GOARCH
		//  - {{version}}: The value of the Version member
		//  - {{stripped_version}}: The value of the Version member but with leading 'v' removed
		//  - {{source}}: The value of the Source member.
		//  - {{plugin_name}}: The value of the Name member
		//
		// For example, a ArtifactTemplate for a plugin released to github via goreleaser might look
		// like:
		//
		//	${source}/releases/download/${version}/${pugin_name}_${stripped_version}_${os}_${arch}.tar.gz
		//
		ArtifactTemplate string `json:"artifact_template" hcl:"artifact_template,optional"`

		// ArchiveFile holds the name of the plugin binary if the downloaded artifact is
		// an archive and contains more than one file.
		//
		// Note that it's also possible to specify the ArchiveFile in dedicated Artifact
		// definitions below as well. If both are specified, the ArchiveFile in the Artifact
		// takes precendence.
		ArchiveFile string `json:"archiveFile" hcl:"archive_file,optional"`

		// Artifacts defines the download URLs for the plugin binary for
		// different architectures and operating systems.
		//
		// If Artifacts and ArtifactTemplate is specified than Artifacts take precedence
		// if there's a matching architecutre definition.
		Artifacts []Artifact `json:"artifacts" hcl:"artifact,block"`

		// PluginTypes defines the list of plugin types implemented
		// by the described plugin.
		PluginTypes []shared.PluginType `json:"pluginTypes" hcl:"pluginTypes"`

		// Author is the name of the plugin author.
		Author string `json:"author" hcl:"author,optional"`

		// License holds the license identifier for the plugin.
		License string `json:"license" hcl:"license,optional"`

		// Privileged specifies if the plugin needs to be enabled as
		// privileged or not.
		Privileged bool `json:"privileged" hcl:"privileged,optional"`

		// Description holds a human readable description of the features and purpose of
		// a plugin.
		Description string `json:"description" hcl:"description,optional"`

		// Tags holds an arbitrary list of tags for the plugin.
		Tags []string `json:"tags" hcl:"tags,optional"`

		// Repository is the name of the repository that contains the
		// plugin.
		Repository string `json:"repository"`
	}

	// RepositoryIndex defines the structure of a repository index file.
	RepositoryIndex struct {
		// Meta holds meta-data about the repository.
		Meta IndexMeta `json:"meta" hcl:"meta,block"`

		// Plugins is the list of plugins available in the repository.
		Plugins []PluginDesc `json:"plugins" hcl:"plugin,block"`
	}
)
