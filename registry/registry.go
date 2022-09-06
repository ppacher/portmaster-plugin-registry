package registry

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"sync"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/go-version"
	"github.com/ppacher/portmaster-plugin-registry/structs"
)

// Common errors returned by the registry package.
var (
	ErrRepoDefined   = errors.New("repository is already defined")
	ErrUnknownPlugin = errors.New("unknown plugin name")
)

type (
	Registry struct {
		l sync.RWMutex

		repos   map[string]structs.Repository
		plugins map[string]pluginDesc
	}

	pluginDesc struct {
		structs.PluginDesc

		Repository string
	}

	// repoList is a helper to sort repositories by priority.
	repoList []structs.Repository
)

func NewRegistry() *Registry {
	return &Registry{
		repos:   make(map[string]structs.Repository),
		plugins: make(map[string]pluginDesc),
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

// Fetch fetches the repository index files and update the local
// list of available plugins.
func (reg *Registry) Fetch() error {
	reg.l.Lock()
	defer reg.l.Unlock()

	pluginList := make(map[string]pluginDesc)

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

			pluginList[plg.Name] = pluginDesc{
				PluginDesc: plg,
				Repository: repo.Name,
			}
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
	u, err := url.Parse(repo.URL)
	if err != nil {
		return nil, err
	}

	res, err := retryablehttp.Get(repo.URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	index, err := DecodeIndex(u.Path, res.Body)
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
