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

## 説明: 代表フレームの取り出し方 (Description in Japanese)

アニメーションGIFから代表フレームを取り出し保存するコマンドは以下の通りです。

```console
$ giftool extract rep path/to/animation.gif
```

この時、取り出した画像は `path/to/animation_rep.png` として保存されます。

より一般的な使い方は以下のようになります。

```console
$ giftool extract rep [OPTIONS] {INPUT ANIMATION GIF}
```

サポートしているオプションは以下の通りです。

* `-grayedentropy` エントロピーの計算にグレースケールに変換した画像を使う。デフォルトは入力画像そのままのカラーモードを使う。
* `-output {出力ファイル名}` 出力ファイル名を指定する。
    省略時は `{入力ファイル名のベース部分}_rep.png` となる。
