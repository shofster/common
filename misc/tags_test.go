package misc

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

/*

  File:    tags_test.go
  Author:  Bob Shofner

*/
/*
  Description: test all formats.
*/

func TestInt(t *testing.T) {
	var tests = []struct {
		name   string
		toInt  func(bytes []uint8) int
		buf    []uint8
		result int
	}{
		{"Intel 2", LittleEndianToInt, []byte{2, 1}, 258},
		{"Intel 4", LittleEndianToInt, []byte{4, 3, 2, 1}, 16909060},
		{"Motorola 2", BigEndianToInt, []byte{1, 2}, 258},
		{"Motorola 4", BigEndianToInt, []byte{1, 2, 3, 4}, 16909060},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s,%d", tt.name, len(tt.buf)), func(t *testing.T) {
			result := tt.toInt(tt.buf)
			if result != tt.result {
				t.Errorf("Expected %d: got %d", tt.result, result)
			}
		})
	}
}

type wrapper struct {
	tags Tags
	err  error
}

func TestBmp(t *testing.T) {
	var tests = []struct {
		name     string
		filename string
		w        wrapper
	}{
		{"OK", "./data/gus2.bmp", wrapper{tags: Tags{ImageForm: bmpForm}, err: nil}},
		{"Insufficient", "./data/short.txt", wrapper{tags: Tags{ImageForm: bmpForm}, err: errors.New("")}},
		{"Missing BM", "./data/gus2.jpg", wrapper{tags: Tags{ImageForm: bmpForm}, err: errors.New("")}},
	}
	bm := new(bmp)
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s, %s", tt.name, tt.filename), func(t *testing.T) {
			file, e := os.Open(tt.filename)
			if e != nil {
				t.Errorf("File Open: %v", e)
			}
			defer func() {
				_ = file.Close()
			}()
			tags, err := bm.parse(file)
			t.Log(tags)
			w := wrapper{tags: tags, err: err}
			if w.err == nil && tt.w.err != nil {
				t.Errorf("Expected error '%v': got nil", tt.w.err)
			} else if w.err != nil && tt.w.err == nil {
				t.Errorf("Expected error nil: got '%v'", w.err)
			} else {
				if tags.ImageForm != bmpForm {
					t.Errorf("Expected 'GIF89a': got %s", tags.ImageForm)
				} else {
					if w.err == nil && tags.ImageHeight != 221 {
						t.Errorf("Expected Height = 221: got %d", tags.ImageHeight)
					}
				}
			}

		})
	}
}
func TestGif(t *testing.T) {
	var tests = []struct {
		name     string
		filename string
		w        wrapper
	}{
		{"OK", "./data/gus2.gif", wrapper{tags: Tags{ImageForm: gifForm}, err: nil}},
		{"Insufficient", "./data/short.txt", wrapper{tags: Tags{ImageForm: gifForm}, err: errors.New("")}},
		{"Missing GIF", "./data/gus2.jpg", wrapper{tags: Tags{ImageForm: gifForm}, err: errors.New("")}},
	}
	gf := new(gif)
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s, %s", tt.name, tt.filename), func(t *testing.T) {
			file, e := os.Open(tt.filename)
			if e != nil {
				t.Errorf("File Open: %v", e)
			}
			defer func() {
				_ = file.Close()
			}()
			tags, err := gf.parse(file)
			t.Log(tags)
			w := wrapper{tags: tags, err: err}
			if w.err == nil && tt.w.err != nil {
				t.Errorf("Expected error '%v': got nil", tt.w.err)
			} else if w.err != nil && tt.w.err == nil {
				t.Errorf("Expected error nil: got '%v'", w.err)
			} else {
				if tags.ImageForm != gifForm {
					t.Errorf(fmt.Sprintf("Expected '%s': got %s", gifForm, tags.ImageForm))
				} else {
					if w.err == nil && tags.ImageHeight != 221 {
						t.Errorf("Expected Height = 221: got %d", tags.ImageHeight)
					}
				}
			}
		})
	}
}
func TestTiff(t *testing.T) {
	var tests = []struct {
		name     string
		filename string
		w        wrapper
	}{
		{"OK", "./data/gus2.tif", wrapper{tags: Tags{ImageForm: tifForm}, err: nil}},
		{"Insufficient", "./data/short.txt", wrapper{tags: Tags{ImageForm: tifForm}, err: errors.New("")}},
		{"Missing TIF", "./data/gus2.bmp", wrapper{tags: Tags{ImageForm: tifForm}, err: errors.New("")}},
	}
	tf := new(tif)
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s, %s", tt.name, tt.filename), func(t *testing.T) {
			file, e := os.Open(tt.filename)
			if e != nil {
				t.Errorf("File Open: %v", e)
			}
			defer func() {
				_ = file.Close()
			}()
			tags, err := tf.parse(file)
			t.Log(tags)
			w := wrapper{tags: tags, err: err}
			if w.err == nil && tt.w.err != nil {
				t.Errorf("Expected error '%v': got nil", tt.w.err)
			} else if w.err != nil && tt.w.err == nil {
				t.Errorf("Expected error nil: got '%v'", w.err)
			} else {
				if tags.ImageForm != tifForm {
					t.Errorf(fmt.Sprintf("Expected '%s': got %s", tifForm, tags.ImageForm))
				} else {
					if w.err == nil && tags.ImageHeight != 221 {
						t.Errorf("Expected Height = 221: got %d", tags.ImageHeight)
					}
				}
			}
		})
	}
}
func TestPng(t *testing.T) {
	var tests = []struct {
		name     string
		filename string
		w        wrapper
	}{
		{"OK", "./data/gus2.png", wrapper{tags: Tags{ImageForm: pngForm}, err: nil}},
		{"Insufficient", "./data/short.txt", wrapper{tags: Tags{ImageForm: pngForm}, err: errors.New("")}},
		{"Missing PNG", "./data/gus2.bmp", wrapper{tags: Tags{ImageForm: pngForm}, err: errors.New("")}},
	}
	pf := new(png)
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s, %s", tt.name, tt.filename), func(t *testing.T) {
			file, e := os.Open(tt.filename)
			if e != nil {
				t.Errorf("File Open: %v", e)
			}
			defer func() {
				_ = file.Close()
			}()
			tags, err := pf.parse(file)
			t.Log(tags)
			w := wrapper{tags: tags, err: err}
			if w.err == nil && tt.w.err != nil {
				t.Errorf("Expected error '%v': got nil", tt.w.err)
			} else if w.err != nil && tt.w.err == nil {
				t.Errorf("Expected error nil: got '%v'", w.err)
			} else {
				if tags.ImageForm != pngForm {
					t.Errorf(fmt.Sprintf("Expected '%s': got %s", pngForm, tags.ImageForm))
				} else {
					if w.err == nil && tags.ImageHeight != 221 {
						t.Errorf("Expected Height = 221: got %d", tags.ImageHeight)
					}
				}
			}
		})
	}
}
func TestJpeg(t *testing.T) {
	var tests = []struct {
		name     string
		filename string
		w        wrapper
	}{
		{"OK", "./data/gus2.jpg", wrapper{tags: Tags{ImageForm: jpgForm}, err: nil}},
		{"scanner", "./data/canon-scanner.JPG", wrapper{tags: Tags{ImageForm: jpgForm}, err: nil}},
		{"scanner", "./data/doxieFlip-scanner.JPG", wrapper{tags: Tags{ImageForm: jpgForm}, err: nil}},
		{"Insufficient", "./data/short.txt", wrapper{tags: Tags{ImageForm: jpgForm}, err: errors.New("")}},
		{"Missing JPG", "./data/gus2.bmp", wrapper{tags: Tags{ImageForm: jpgForm}, err: errors.New("")}},
	}
	jf := new(jpg)
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s, %s", tt.name, tt.filename), func(t *testing.T) {
			file, e := os.Open(tt.filename)
			if e != nil {
				t.Errorf("File Open: %v", e)
			}
			defer func() {
				_ = file.Close()
			}()
			tags, err := jf.parse(file)
			t.Log(tags)
			w := wrapper{tags: tags, err: err}
			if w.err == nil && tt.w.err != nil {
				t.Errorf("Expected error '%v': got nil", tt.w.err)
			} else if w.err != nil && tt.w.err == nil {
				t.Errorf("Expected error nil: got '%v'", w.err)
			} else {
				if tags.ImageForm != jpgForm {
					t.Errorf(fmt.Sprintf("Expected Form '%s': got %s",
						jpgForm, tags.ImageForm))
				} else {
					if tt.name != "scanner" { // scanner is different
						if w.err == nil && tags.ImageHeight != 221 {
							t.Errorf("Expected Height = 221: got %d", tags.ImageHeight)
						}
					}
				}
			}
		})
	}
}
