---
prev: ./
---

# Getting Started

## Prerequisites

The plugin system of the Portmaster has not yet been released and is still waiting for the Pull Request to be merged.

If you want to try out Portmaster Plugins and/or PECS you must build the Portmaster Core service from source code using the following branch:

- [safing/portmaster#834](https://github.com/safing/portmaster/pull/834)

::: tip
PECS itself adds plugin management features to the Portmaster while actually being implemented as a plugin itself.
:::

### Portmaster Configuration

The Plugin System of the Portmaster is disabled by default and marked as `Experimental` and `Developer-Only`. 

In order to enable the Plugin System you need to change the following settings in the Portmaster settings:

- [Feature Stability](https://docs.safing.io/portmaster/settings#core/releaseLevel):  
   `Experimental`
- [Development Mode](https://docs.safing.io/portmaster/settings#core/devMode):  
   `enabled`

::: warning
Changing the above settings will reveal much more configuration options and detailed information. If you don't feel confortable with changing these it's probably better to wait for the Plugin System to be released as "stable".
:::

Once the settings have been changed you should be able to find a setting called `Enable Plugin System`. Hit that one and restart the Portmaster :tada:

## Installation

There are two possibilities to install the PECS, either [manually (harder)](#manually-from-source) or by using the [prebuilt binaries (easy)](#prebuilt-binaries).

### Prebuilt Binaries

::: danger MISSING

This section still needs to be done

:::

### Manually From Source

To manually build PECS from source you need to have a working [Golang](https://golang.org) environment setup. PECS is built using [Go Modules](https://go.dev/blog/using-go-modules) and requires at least `Go v1.18`.

Once your environment is ready you can clone the repository and build the plugin:

```bash:no-line-numbers
# Clone the repository and enter it
git clone https://github.com/ppacher/portmaster-plugin-registry
cd portmaster-plugin-registry

# Download all required dependencies
go mod download

# Build the registry-util binary
go build ./cmds/registry-util

# Build the actual registry plugin
go build ./cmds/registry-plugin
```

::: tip
The repository contains two different executables for the following purposes:

**registry-plugin**:  
The actual PECS plugin that automates discovery and management of plugins. It provides a neat web-based user interface and tightly integrates with the Portmaster using it's notification and configuration system.

**registry-util**:  
A simple terminal CLI that can be used to install plugins without the need actually run PECS. It supports downloading and installing plugins from PECS plugin repositories and updates the Portmaster configuration file automatically. Though, you need to configure your repositories manually and will not be informed about updates or otherwise track your plugins.
:::

Finally, you can install the PECS plugin using either `registry-util` or `registry-plugin`:

<CodeGroup>
  <CodeGroupItem title="registry-plugin">

```bash:no-line-numbers
# This copies and install the same binary
sudo ./registry-plugin install --data /opt/safing/portmaster
```

  </CodeGroupItem>

  <CodeGroupItem title="registry-util" active>

```bash:no-line-numbers
# This actually bootstraps the registry and downloads the
# pre-built binaries from https://pecs.xyz.
sudo ./registry-util install https://pecs.xyz/bootstrap.hcl pecs
```

  </CodeGroupItem>
</CodeGroup>

**You just installed PECS** :tada:

Just restart your Portmaster and wait for the Notification created by PECS.