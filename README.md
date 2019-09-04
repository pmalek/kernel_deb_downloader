# kernel_deb_downloader
[![Github Actions Build Status](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fpmalek%2Fkernel_deb_downloader%2Fbadge&style=flat)](https://actions-badge.atrox.dev/pmalek/kernel_deb_downloader/goto) [![Go Report Card](https://goreportcard.com/badge/github.com/pmalek/kernel_deb_downloader)](https://goreportcard.com/report/github.com/pmalek/kernel_deb_downloader) [![codecov.io report badge](https://codecov.io/gh/pmalek/kernel_deb_downloader/branch/master/graph/badge.svg)](https://codecov.io/gh/pmalek/kernel_deb_downloader) [![Maintainability](https://api.codeclimate.com/v1/badges/a96a799303b1171eb5d5/maintainability)](https://codeclimate.com/github/pmalek/kernel_deb_downloader/maintainability)


Simple Go project which downloads newest .deb packages from ubuntu mainline ppa

## Usage

```
kernel_deb_downloader --help

Usage of kernel_deb_downloader:
  -c    Show changes included in particular kernel package
  -n    Print newest version - do not download the .debs
```
