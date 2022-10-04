package misc

import (
	"os"
	"path/filepath"
	"strings"
)

/*

  File:    imagego
  Author:  Bob Shofner

  Copyright (c) 2022. BSD 3-Clause License
	https://opensource.org/licenses/BSD-3-Clause

  The this permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Read jpeg or png
*/

type ImageConfig struct {
	ImageType FormatType
	Width     int
	Height    int
}
type FormatType int

const (
	Unknown FormatType = iota
	Bmp
	Gif
	Jpeg
	Png
	Tiff
)

//goland:noinspection GoUnusedExportedFunction
func ImageSize(path string) (ImageConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return ImageConfig{ImageType: Unknown}, err
	}
	defer func() {
		_ = f.Close()
	}()
	var ic = ImageConfig{ImageType: Unknown}
	var parser tagParser
	switch strings.ToLower(filepath.Ext(path)) {
	case ".bmp":
		ic.ImageType = Bmp
		parser = new(bmp)
	case ".gif":
		ic.ImageType = Gif
		parser = new(gif)
	case ".jpeg", ".jpg", ".jpe", ".jfif", ".jif":
		ic.ImageType = Jpeg
		parser = new(jpg)
	case ".png":
		ic.ImageType = Png
		parser = new(png)
	case ".tif", ".tiff":
		ic.ImageType = Tiff
		parser = new(tif)
	default:
		return ic, nil
	}
	tags, e := parser.parse(f)
	if e != nil {
		return ic, e
	}
	ic.Width = tags.ImageWidth
	ic.Height = tags.ImageHeight
	return ic, nil
}
