# Portmaster Extension Community Store (_PECS_)

[![Go Report Card](https://goreportcard.com/badge/github.com/ppacher/portmaster-plugin-registry)](https://goreportcard.com/report/github.com/ppacher/portmaster-plugin-registry)
![Lint & Build](https://github.com/ppacher/portmaster-plugin-registry/actions/workflows/lint.yml/badge.svg)
![Lint & Build](https://github.com/ppacher/portmaster-plugin-registry/actions/workflows/codeql.yml/badge.svg)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/nextdhcp/nextdhcp)](https://github.com/nextdhcp/nextdhcp/releases)

_Manage (Download, install, track and upgrade) and discover third-party plugins for the [Portmaster Application Firewall](https://safing.io/portmaster?source=ppacher-pecs)_



## :confused: What?

With Pull-Request [safing/portmaster#834](https://github.com/safing/portmaster/pull/834) the Safing Portmaster Application Firewall gained support for third-party plugins that can extend the feature set the privacy suite.

The Portmaster Extension Community Store (short _PECS_) integrates with the Portmaster and gives the user an easy way to discover and install third-party Portmaster plugins.

**Feature Highlights:**  

- Discover third-party plugins using a central plugin repository
- Add additional plugin repositories
- Automated download and installation of plugins
  - including automatic plugin upgrades
- Manage your plugins from CLI or use a neat web-based UI

## :bug: Issues

Found a issue/bug or want to suggest an improvement/feature? The best way to report them is to open an issue in **this** repo.

[Issue link](https://github.com/ppacher/portmaster-plugin-registry/issues)

## :warning: Disclaimer

While the author of _Portmaster Extension Community Store_ (@ppacher)  is a member of the Safing Portmaster Core Team, **this is not an official Safing product and you will not get support from Safing**.

**:lock: Important Security Note:**  
The _PECS_ provides easy acces to **third-party** code that will be executed with highest privileges on your system! There is no code review process for plugins added to the plugin repository hosted here. :warning: Be careful! :warning:
