# koron/giftool

[![PkgGoDev](https://pkg.go.dev/badge/github.com/koron/giftool)](https://pkg.go.dev/github.com/koron/giftool)
[![Actions/Go](https://github.com/koron/giftool/workflows/Go/badge.svg)](https://github.com/koron/giftool/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron/giftool)](https://goreportcard.com/report/github.com/koron/giftool)

## Getting Started

Requirements: Go 1.23 or later

Install and update:

```console
$ go install github.com/koron/giftool
```

Or use pre-compiled binaries from [the latest release](https://github.com/koron/giftool/releases/latest).

Extract a representative (rep) frame from animation GIF.
The representative frame is selected as the frame with the highest image entropy.

```console
$ giftool extract rep path/to/animation.gif
```

Extracted image is saved to path/to/animation\_rep.png
