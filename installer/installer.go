package installer

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/ppacher/portmaster-plugin-registry/structs"
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
		InstallPlugin(desc structs.PluginDesc) (string, error)
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
func (installer *PluginInstaller) InstallPlugin(plg structs.PluginDesc) (string, error) {
	downloadURL, err := findMatchingArtifact(plg)
	if err != nil {
		return "", err
	}

	targetFile := filepath.Join(
		installer.TargetDirectory,
		fmt.Sprintf("%s-%s", plg.Name, plg.Version),
	)

	if runtime.GOOS == "windows" {
		targetFile += ".exe"
	}

	res, err := retryablehttp.Get(downloadURL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	f, err := os.Create(targetFile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, res.Body); err != nil {
		defer os.Remove(f.Name())

		return "", err
	}

	if err := os.Chmod(targetFile, 0555); err != nil {
		defer os.Remove(f.Name())

		return "", err
	}

	return targetFile, nil
}

func findMatchingArtifact(plg structs.PluginDesc) (string, error) {
	// find the correct artifact
	var artifact *structs.Artifact

	for _, a := range plg.Artifacts {
		if a.OS == runtime.GOOS {
			artifact = &a

			break
		}
	}

	if artifact == nil {
		return "", ErrNoMatchingArtifact
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
		return "", ErrUnsupportedArch
	}

	if url == "" {
		return "", ErrNoMatchingArtifact
	}

	return url, nil
}

// Interface checks
var _ Installer = new(PluginInstaller)
