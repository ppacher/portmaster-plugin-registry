package manager

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/renameio"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/ppacher/portmaster-plugin-registry/installer"
	"github.com/ppacher/portmaster-plugin-registry/structs"
	"github.com/safing/portmaster/plugin/shared"
	"github.com/safing/portmaster/plugin/shared/pluginmanager"
	"github.com/safing/portmaster/plugin/shared/proto"
)

var (
	ErrUnknownPlugin = errors.New("unknown plugin")
)

type (
	// PluginProvider describes the minimum interface required by the manager.
	// It's implemented by registry.Registry.
	PluginProvider interface {
		Fetch() error
		ByName(string) (structs.PluginDesc, bool)
		UpdateAvailable(name, version string) (string, error)
	}

	// Manager manages installed plugins.
	Manager struct {
		stateFile     string
		provider      PluginProvider
		installer     installer.Installer
		pluginManager pluginmanager.Service

		l                 sync.RWMutex
		started           bool
		installedPlugins  []structs.InstalledPlugin
		onFetchDone       []func(err error)
		onUpdateAvailable []func(updates []structs.AvailableUpdate)
	}
)

// NewManager returns a new plugin manager that stores state information in
// stateFile and uses inst for plugin installations and reg for available plugin
// lookups.
func NewManager(stateFile string, inst installer.Installer, reg PluginProvider, service pluginmanager.Service) *Manager {
	return &Manager{
		stateFile:     stateFile,
		installer:     inst,
		provider:      reg,
		pluginManager: service,
	}
}

// OnFetchDone registers a callback function that is invoked when the
// plugin provided fetched new repository data.
//
// The provided error indicates if the fetch was successful or not.
func (mng *Manager) OnFetchDone(fn func(error)) {
	mng.l.Lock()
	defer mng.l.Unlock()

	mng.onFetchDone = append(mng.onFetchDone, fn)
}

// OnUpdateAvailable registers a new callback function that is executed
// when new updates are available.
//
// This is only ever fired after OnFetchDone and only if fetching repositories
// was successful.
func (mng *Manager) OnUpdateAvailable(fn func([]structs.AvailableUpdate)) {
	mng.l.Lock()
	defer mng.l.Unlock()

	mng.onUpdateAvailable = append(mng.onUpdateAvailable, fn)
}

// InstalledPlugins returns a list of all installed plugins.
func (mng *Manager) InstalledPlugins() []structs.InstalledPlugin {
	mng.l.RLock()
	defer mng.l.RUnlock()

	list := make([]structs.InstalledPlugin, len(mng.installedPlugins))
	copy(list, mng.installedPlugins)

	return list
}

// Start starts the plugin manager. The manager will shutdown as soon as
// ctx is cancelled.
func (mng *Manager) Start(ctx context.Context) error {
	mng.l.Lock()
	defer mng.l.Unlock()

	if mng.started {
		return nil
	}
	mng.started = true

	if err := mng.loadStateFile(ctx); err != nil {
		return err
	}

	if err := mng.registerAllPlugins(ctx); err != nil {
		return fmt.Errorf("failed to register plugins: %w", err)
	}

	if err := mng.provider.Fetch(); err != nil {
		return err
	}

	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		defer ticker.Stop()
	L:
		for {

			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				mng.update()
				break L
			}
		}

		// try to read one more time from the ticker channel
		// to make sure we don't leak any goroutines blocking on ticker.C <- ...
		select {
		case <-ticker.C:
		default:
		}
	}()

	return nil
}

// InstallPlugin installs a new plugin, updates the state file, registers it in the
// Portmaster.
func (mng *Manager) InstallPlugin(ctx context.Context, name string) error {
	plg, ok := mng.provider.ByName(name)
	if !ok {
		return ErrUnknownPlugin
	}

	pluginTypes, err := pluginTypesToProto(plg.PluginTypes)
	if err != nil {
		return err
	}

	path, err := mng.installer.InstallPlugin(plg)
	if err != nil {
		return fmt.Errorf("failed to install: %w", err)
	}

	mng.l.Lock()
	defer mng.l.Unlock()

	mng.installedPlugins = append(mng.installedPlugins, structs.InstalledPlugin{
		PluginDesc: plg,
		Path:       path,
	})

	if err := mng.saveStateFile(); err != nil {
		return fmt.Errorf("failed to update state file: %w", err)
	}

	if err := mng.pluginManager.RegisterPlugin(ctx, &proto.PluginConfig{
		Name:             name,
		PluginTypes:      pluginTypes,
		Privileged:       plg.Privileged,
		DisableAutostart: true,
	}); err != nil {
		return fmt.Errorf("failed to register plugin in Portmaster: %w", err)
	}

	return nil
}

func (mng *Manager) update() {
	err := mng.provider.Fetch()

	mng.l.RLock()
	defer mng.l.RUnlock()

	for _, fn := range mng.onFetchDone {
		fn(err)
	}

	// we abort now if there was an error
	if err != nil {
		return
	}

	updates := mng.detectUpdates()
	if len(updates) > 0 {
		for _, fn := range mng.onUpdateAvailable {
			fn(updates)
		}
	}
}

// AvailableUpdates returns a list of available plugin updates.
func (mng *Manager) AvailableUpdates() []structs.AvailableUpdate {
	mng.l.Lock()
	defer mng.l.Unlock()

	return mng.detectUpdates()
}

func (mng *Manager) loadStateFile(ctx context.Context) error {
	content, err := os.ReadFile(mng.stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	var file structs.InstalledPluginsFile
	if err := hclsimple.Decode(mng.stateFile, content, nil, &file); err != nil {
		return err
	}

	if file.Version != "v1.0.0" {
		return fmt.Errorf("unsupported installed plugins file format %q", file.Version)
	}

	mng.installedPlugins = file.Plugins

	return nil
}

func (mng *Manager) registerAllPlugins(ctx context.Context) error {
	multierr := new(multierror.Error)
	for _, plg := range mng.installedPlugins {
		protoTypes, err := pluginTypesToProto(plg.PluginTypes)
		if err != nil {
			multierr.Errors = append(multierr.Errors, fmt.Errorf("plugin %s: %w", plg.Name, err))

			continue
		}

		if err := mng.pluginManager.RegisterPlugin(ctx, &proto.PluginConfig{
			Name:             plg.Name,
			PluginTypes:      protoTypes,
			Privileged:       plg.Privileged,
			DisableAutostart: true,
		}); err != nil {
			multierr.Errors = append(multierr.Errors, fmt.Errorf("plugin %s: failed to register: %w", plg.Name, err))

			continue
		}
	}

	return multierr.ErrorOrNil()
}

func (mng *Manager) saveStateFile() error {
	file := hclwrite.NewEmptyFile()

	fileContent := structs.InstalledPluginsFile{
		Version: "v1.0.0",
		Plugins: mng.installedPlugins,
	}

	gohcl.EncodeIntoBody(fileContent, file.Body())

	buf := new(bytes.Buffer)
	if _, err := file.WriteTo(buf); err != nil {
		return fmt.Errorf("failed to write file body: %w", err)
	}

	return renameio.WriteFile(mng.stateFile, buf.Bytes(), 0444)
}

func (mng *Manager) detectUpdates() []structs.AvailableUpdate {
	var updates []structs.AvailableUpdate

	for _, installedPlugin := range mng.installedPlugins {
		updatedVersion, err := mng.provider.UpdateAvailable(installedPlugin.Name, installedPlugin.Version)
		if err != nil {
			hclog.L().Error("failed to check for available updates", "plugin", installedPlugin.Name, "error", err)

			continue
		}

		if updatedVersion != "" {
			updates = append(updates, structs.AvailableUpdate{
				Name:           installedPlugin.Name,
				CurrentVersion: installedPlugin.Version,
				NewVersion:     updatedVersion,
			})
		}
	}

	return updates
}

func pluginTypesToProto(pTypes []shared.PluginType) ([]proto.PluginType, error) {
	var pluginTypes []proto.PluginType
	for _, pType := range pTypes {
		switch pType {
		case shared.PluginTypeDecider:
			pluginTypes = append(pluginTypes, proto.PluginType_PLUGIN_TYPE_DECIDER)
		case shared.PluginTypeReporter:
			pluginTypes = append(pluginTypes, proto.PluginType_PLUGIN_TYPE_REPORTER)
		case shared.PluginTypeResolver:
			pluginTypes = append(pluginTypes, proto.PluginType_PLUGIN_TYPE_RESOLVER)
		default:
			return nil, fmt.Errorf("unsupported plugin type: %s", pType)
		}
	}

	return pluginTypes, nil
}
