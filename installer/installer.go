package installer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/renameio"
	"github.com/hashicorp/go-getter/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/ppacher/portmaster-plugin-registry/structs"
	"github.com/valyala/fasttemplate"
)

var (
	ErrNoMatchingArtifact = errors.New("no artifact matches the current system")
	ErrUnsupportedArch    = errors.New("your system architecture is not yet supported")
)

// PluginInstaller can download and install plugins at a specified
// target directory.
type (
	// Installer defines the interface that is capable of downloading
	// and installing plugins.
	Installer interface {
		// InstallPlugin should install the plugin defined in desc at the
		// local system and return the path to the installed binary.
		InstallPlugin(ctx context.Context, desc structs.PluginDesc) (string, error)
	}

	// PluginInstaller implements the Installer interface and is capable of
	// downloading and installing binary plugins at the local host.
	PluginInstaller struct {
		// TargetDirectory is the directory where plugins should be installed.
		TargetDirectory string
	}
)

// InstallPlugin installs the plugin in the target directory and returns
// the path of the installed plugin binary.
func (installer *PluginInstaller) InstallPlugin(ctx context.Context, plg structs.PluginDesc) (string, error) {
	pluginFile, err := DownloadPlugin(ctx, "", plg)
	if err != nil {
		return "", err
	}
	hclog.L().Info("artifact downloaded successfully", "plugin", plg.Name, "plugin-file", pluginFile)

	targetFile := filepath.Join(
		installer.TargetDirectory,
		fmt.Sprintf("%s-%s", plg.Name, plg.Version),
	)

	if runtime.GOOS == "windows" {
		targetFile += ".exe"
	}

	if err := moveFile(targetFile, pluginFile); err != nil {
		return "", err
	}
	hclog.L().Info("artifact successfully moved", "plugin", plg.Name, "plugin-file", pluginFile, "target", targetFile)

	return targetFile, nil
}

func DownloadPlugin(ctx context.Context, dst string, plg structs.PluginDesc) (string, error) {
	downloadURL, archiveFile, err := FindMatchingArtifact(plg)
	if err != nil {
		return "", err
	}

	if dst == "" {
		var err error
		dst, err = os.MkdirTemp("", plg.Name+"-*")
		if err != nil {
			return "", err
		}
	}

	hclog.L().Info("downloading artifact", "plugin", plg.Name, "url", downloadURL, "dst", dst)

	cli := new(getter.Client)
	res, err := cli.Get(ctx, &getter.Request{
		Src: downloadURL,
		Dst: dst,
	})
	if err != nil {
		return "", err
	}

	pluginFile, err := pluginFileFromArtifact(plg.Name, res.Dst, archiveFile)
	if err != nil {
		return "", err
	}

	return pluginFile, err
}

func moveFile(destination, source string) error {
	f, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open plugin file from %s: %w", source, err)
	}
	defer f.Close()

	target, err := renameio.TempFile("", destination)
	if err != nil {
		return err
	}
	defer target.Cleanup()

	if _, err := io.Copy(target, f); err != nil {
		return fmt.Errorf("failed to copy plugin file: %w", err)
	}

	if err := target.Chmod(0555); err != nil {
		return fmt.Errorf("failed to update file mode: %w", err)
	}

	if err := target.CloseAtomicallyReplace(); err != nil {
		return fmt.Errorf("failed to atuomically rename target: %w", err)
	}

	return nil
}

func pluginFileFromArtifact(plgName string, artifact string, archiveFile string) (string, error) {
	stat, err := os.Stat(artifact)
	if err != nil {
		return "", err
	}

	if !stat.IsDir() {
		return artifact, nil
	}

	if archiveFile != "" {
		return filepath.Join(artifact, archiveFile), nil
	}

	files, err := os.ReadDir(artifact)
	if err != nil {
		return "", err
	}

	if len(files) == 1 {
		// there's only one file in the directory, this must
		// be the plugin
		if files[0].IsDir() {
			return "", fmt.Errorf("failed to find plugin in archive")
		}

		return filepath.Join(artifact, files[0].Name()), nil
	}

	// try to search for the plugin file
	for _, file := range files {
		if file.Name() == plgName {
			return filepath.Join(artifact, file.Name()), nil
		}

		if runtime.GOOS == "windows" {
			if file.Name() == plgName+".exe" {
				return filepath.Join(artifact, file.Name()), nil
			}
		}
	}

	return "", fmt.Errorf("failed to find plugin in archive")
}

func FindMatchingArtifact(plg structs.PluginDesc) (string, string, error) {
	// find the correct artifact
	var artifact *structs.Artifact

	for _, a := range plg.Artifacts {
		if a.OS == runtime.GOOS {
			artifact = &a

			break
		}
	}

	if artifact == nil {

		// if there's an artifact_template try to use that one
		if plg.ArtifactTemplate != "" {
			buf := new(strings.Builder)
			_, err := fasttemplate.Execute(plg.ArtifactTemplate, "{{", "}}", buf, map[string]any{
				"os":               runtime.GOOS,
				"arch":             runtime.GOARCH,
				"version":          plg.Version,
				"stripped_version": strings.TrimPrefix(plg.Version, "v"),
				"plugin_name":      plg.Name,
				"source":           plg.SourceURL,
				"archive_file":     plg.ArchiveFile,
			})
			if err != nil {
				return "", "", err
			}

			return buf.String(), plg.ArchiveFile, nil
		}

		return "", "", ErrNoMatchingArtifact
	}

	archiveFile := artifact.ArchiveFile
	if archiveFile == "" {
		archiveFile = plg.ArchiveFile
	}

	// get the correct architecture download link
	url := ""
	switch runtime.GOARCH {
	case "amd64":
		url = artifact.AMD64
	case "arm":
		url = artifact.ARM
	case "arm64":
		url = artifact.ARM64
	case "i386":
		url = artifact.I386
	default:
		return "", "", ErrUnsupportedArch
	}

	if url == "" {
		return "", "", ErrNoMatchingArtifact
	}

	return url, archiveFile, nil
}

// Interface checks
var _ Installer = new(PluginInstaller)
