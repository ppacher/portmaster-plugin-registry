package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/ppacher/portmaster-plugin-registry/structs"
)

// DecodeIndex decodes a repository index file from reader. The path is required to detect
// the correct encoding.
//
// Supported file extensions are .yaml, .json and .hcl.
func DecodeIndex(path string, reader io.Reader) (*structs.RepositoryIndex, error) {
	ext := filepath.Ext(path)

	var repo structs.RepositoryIndex

	blob, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	switch ext {
	case ".yaml":
		blob, err = yaml.YAMLToJSON(blob)
		if err != nil {
			return nil, err
		}

		fallthrough
	case ".json":
		if err := json.Unmarshal(blob, &repo); err != nil {
			return nil, err
		}

	case ".hcl":
		if err := hclsimple.Decode(path, blob, nil, &repo); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported repository index format %q", ext)
	}

	return &repo, nil
}

// ValidateIndex validates all plugin configurations in index and returns a list
// of validation errors. A non-nil error is always of type *multierror.Error.
//
// If no errors are found, nil is returned.
func ValidateIndex(index *structs.RepositoryIndex) error {
	var errs = new(multierror.Error)

	if index.Meta.Version != "v1.0.0" {
		return fmt.Errorf("unsupported index version %q", index.Meta.Version)
	}

	seenPlugins := make(map[string]struct{})

	for _, plg := range index.Plugins {
		plgErrs := new(multierror.Error)

		if plg.Name == "" {
			errs.Errors = append(errs.Errors, fmt.Errorf("plugin name must be specified"))

			continue
		}

		if _, ok := seenPlugins[plg.Name]; ok {
			plgErrs.Errors = append(plgErrs.Errors, fmt.Errorf("duplicated plugin name"))
		}

		hasArtifact := false
		if plg.ArtifactTemplate != "" {
			hasArtifact = true
		}

		for _, a := range plg.Artifacts {
			isValid := true

			if a.OS == "" {
				plgErrs.Errors = append(plgErrs.Errors, fmt.Errorf("artifact OS must be specified"))
				isValid = false
			}

			if a.AMD64 == "" && a.ARM == "" && a.ARM64 == "" && a.I386 == "" {
				plgErrs.Errors = append(plgErrs.Errors, fmt.Errorf("artifact %q: no download URL defined", a.OS))
				isValid = false
			}

			if isValid {
				hasArtifact = true
			}
		}

		if !hasArtifact {
			if len(plg.Artifacts) > 0 {
				plgErrs.Errors = append(plgErrs.Errors, fmt.Errorf("no valid artifacts defined"))
			} else {
				plgErrs.Errors = append(plgErrs.Errors, fmt.Errorf("no artifacts defined"))
			}
		}

		if plg.Version == "" {
			plgErrs.Errors = append(plgErrs.Errors, fmt.Errorf("version not specified"))
		} else {
			_, err := version.NewSemver(plg.Version)
			if err != nil {
				plgErrs.Errors = append(plgErrs.Errors, fmt.Errorf("invalid semver version: %w", err))
			}
		}

		if err := plgErrs.ErrorOrNil(); err != nil {
			errs.Errors = append(errs.Errors, fmt.Errorf("plugin %s: %w", plg.Name, plgErrs))
		}
	}

	return errs.ErrorOrNil()
}
