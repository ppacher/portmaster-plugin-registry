package registry

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/hashicorp/go-getter/v2"
	"github.com/hashicorp/go-version"
	"github.com/ppacher/portmaster-plugin-registry/structs"
	"github.com/safing/portmaster/plugin/shared"
)

// Common errors returned by the registry package.
var (
	ErrRepoDefined   = errors.New("repository is already defined")
	ErrUnknownPlugin = errors.New("unknown plugin name")
)

type (
	// Registry manages one or more plugin repositories, fetches their
	// index files and provides access to available plugins.
	//
	// It also supports detecting available plugin updates.
	Registry struct {
		l sync.RWMutex

		repos   map[string]structs.Repository
		plugins map[string]structs.PluginDesc
	}

	// repoList is a helper to sort repositories by priority.
	repoList []structs.Repository
)

// NewRegistry creates a new plugin registry. Note that the registry
// does not yet contain any plugin repositories, users should call
// AddRepository() and finally update the registry by calling Fetch().
func NewRegistry() *Registry {
	return &Registry{
		repos:   make(map[string]structs.Repository),
		plugins: make(map[string]structs.PluginDesc),
	}
}

// AddRepository adds a new repository to the registry.
func (reg *Registry) AddRepository(repo structs.Repository) error {
	reg.l.Lock()
	defer reg.l.Unlock()

	if _, ok := reg.repos[repo.Name]; ok {
		return ErrRepoDefined
	}

	reg.repos[repo.Name] = repo

	return nil
}

// ListPlugins returns a list of available plugins.
func (reg *Registry) ListPlugins() []structs.PluginDesc {
	reg.l.RLock()
	defer reg.l.RUnlock()

	list := make([]structs.PluginDesc, 0, len(reg.plugins))
	for _, plg := range reg.plugins {
		list = append(list, plg)
	}

	return list
}

// ByName returns the plugin by name.
func (reg *Registry) ByName(name string) (structs.PluginDesc, bool) {
	reg.l.RLock()
	defer reg.l.RUnlock()

	plg, ok := reg.plugins[name]

	return plg, ok
}

// SearchByTag returns a list of plugins that contain searchTag in their tag list.
func (reg *Registry) SearchByTag(searchTag string) []structs.PluginDesc {
	reg.l.RLock()
	defer reg.l.RUnlock()

	var list []structs.PluginDesc
L:
	for _, plg := range reg.plugins {
		for _, tag := range plg.Tags {
			if tag == searchTag {
				list = append(list, plg)

				continue L
			}
		}
	}

	return list
}

// SearchByType returns a list of plugins that implement type.
func (reg *Registry) SearchByType(pType shared.PluginType) []structs.PluginDesc {
	reg.l.RLock()
	defer reg.l.RUnlock()

	var list []structs.PluginDesc
L:
	for _, plg := range reg.plugins {
		for _, plgType := range plg.PluginTypes {
			if plgType == pType {
				list = append(list, plg)

				continue L
			}
		}
	}

	return list
}

// SearchByName returns a list of plugins that match name.
func (reg *Registry) SearchByName(name string) []structs.PluginDesc {
	reg.l.RLock()
	defer reg.l.RUnlock()

	lowerName := strings.ToLower(name)

	var list []structs.PluginDesc
	for _, plg := range reg.plugins {
		if strings.Contains(strings.ToLower(plg.Name), lowerName) {
			list = append(list, plg)
		}
	}

	return list
}

// Fetch fetches the repository index files and update the local
// list of available plugins.
func (reg *Registry) Fetch() error {
	reg.l.Lock()
	defer reg.l.Unlock()

	pluginList := make(map[string]structs.PluginDesc)

	repoList := make(repoList, 0, len(reg.repos))
	for _, repo := range reg.repos {
		repoList = append(repoList, repo)
	}

	// sort the repository list by priority
	sort.Sort(repoList)

	// fetch all index files and parse them
	for _, repo := range repoList {
		index, err := fetchIndex(repo)
		if err != nil {
			return err
		}

		for _, plg := range index.Plugins {
			if _, ok := pluginList[plg.Name]; ok {
				// this plugin has already be defined by a higher-priority
				// repository.
				continue
			}

			plg.Repository = repo.Name

			pluginList[plg.Name] = plg
		}
	}

	reg.plugins = pluginList

	return nil
}

// UpdateAvailable checks if an update to plgName is available. It compares the version
// of the loaded repositories with the current version and returns the available version
// if it is higher than the current one.
//
// If no update is available and empty string and a nil error is returned.
// If there is no such plugin available ErrUnknownPlugin is returned. In case any of
// the version cannot be parsed an error is returned.
func (reg *Registry) UpdateAvailable(plgName string, currentVersion string) (string, error) {
	reg.l.RLock()
	defer reg.l.RUnlock()

	plg, ok := reg.plugins[plgName]
	if !ok {
		return "", ErrUnknownPlugin
	}

	currentSemVer, err := version.NewSemver(currentVersion)
	if err != nil {
		return "", fmt.Errorf("failed to parse current version %q: %w", currentVersion, err)
	}

	availableSemVer, err := version.NewSemver(plg.Version)
	if err != nil {
		return "", fmt.Errorf("failed to parse available version %q: %w", plg.Version, err)
	}

	if availableSemVer.GreaterThan(currentSemVer) {
		return plg.Version, nil
	}

	return "", nil
}

func fetchIndex(repo structs.Repository) (*structs.RepositoryIndex, error) {
	res, err := new(getter.Client).Get(context.Background(), &getter.Request{
		Src: repo.URL,
		Dst: os.TempDir(),
	})
	if err != nil {
		return nil, err
	}

	f, err := os.Open(res.Dst)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	index, err := DecodeIndex(res.Dst, f)
	if err != nil {
		return nil, err
	}

	if err := ValidateIndex(index); err != nil {
		return nil, err
	}

	return index, nil
}

func (list repoList) Len() int           { return len(list) }
func (list repoList) Less(i, j int) bool { return list[i].Priority < list[j].Priority }
func (list repoList) Swap(i, j int)      { list[i], list[j] = list[j], list[i] }
