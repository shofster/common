package misc

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

/*

  File:    tags.go
  Author:  Bob Shofner

  Copyright (c) 2022. BSD 3-Clause License
	https://opensource.org/licenses/BSD-3-Clause

  The this permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

type Tags struct {
	ImageForm     string
	ImageWidth    int
	ImageHeight   int
	ImageModel    string
	ImageDateTime time.Time
}

// tagParser is the function to parse th image tags.
type tagParser interface {
	parse(file *os.File) (Tags, error)
}

func (t Tags) String() string {
	return fmt.Sprintf("Form:%s, Width:%d, Height:%d",
		t.ImageForm, t.ImageWidth, t.ImageHeight)
}

const bmpForm = "BMP"
const gifForm = "GIF"
const jpgForm = "JPG"
const pngForm = "PNG"
const tifForm = "TIF"

type bmp struct {
}
type gif struct {
}
type jpg struct {
}
type png struct {
}
type tif struct {
}

var _ tagParser = (*bmp)(nil)
var _ tagParser = (*gif)(nil)
var _ tagParser = (*jpg)(nil)
var _ tagParser = (*png)(nil)
var _ tagParser = (*tif)(nil)

var bmpCode = []byte{'B', 'M'}

func (b *bmp) parse(file *os.File) (Tags, error) {
	buf := make([]byte, 31)
	l, err := file.Read(buf)
	if err != nil {
		fmt.Println(err)
	}
	if l > 30 {
		if buf[0] == bmpCode[0] && buf[1] == bmpCode[1] {
			tags := Tags{ImageForm: bmpForm,
				ImageWidth:  LittleEndianToInt(buf[18:22]),
				ImageHeight: LittleEndianToInt(buf[22:26]),
			}
			return tags, nil
		} else {
			return Tags{ImageForm: bmpForm}, errors.New(fmt.Sprintf("missing %s '%s' code",
				bmpForm, string(bmpCode)))
		}
	}
	return Tags{ImageForm: bmpForm}, errors.New("insufficient bytes")
}

const gifCode = "GIF89a"

func (g *gif) parse(file *os.File) (Tags, error) {
	buf := make([]byte, 21)
	l, err := file.Read(buf)
	if err != nil {
		fmt.Println(err)
	}
	if l > 20 {
		code := string(buf[0:6])
		if code == gifCode {
			tags := Tags{ImageForm: gifForm,
				ImageWidth:  LittleEndianToInt(buf[6:8]),
				ImageHeight: LittleEndianToInt(buf[8:10]),
			}
			return tags, nil
		} else {
			return Tags{ImageForm: gifForm}, errors.New(fmt.Sprintf("missing %s code '%s'",
				gifForm, gifCode))
		}
	}
	return Tags{ImageForm: gifForm}, errors.New("insufficient bytes")
}

var jpegSoi = []uint8{0xff, 0xd8}  // 216.D
var jpegEoi = []uint8{0xff, 0xd9}  // 217.D
var jpegApp0 = []uint8{0xff, 0xE0} // 224.D
var jpegSof0 = []uint8{0xff, 0xC0} // 192.D
var jpegSof2 = []uint8{0xff, 0xC2} // 194.D
var jpegCode = "JFIF"

func (t *jpg) parse(file *os.File) (Tags, error) {
	segment := make([]byte, 2)
	_, _ = file.Read(segment)
	if string(segment) != string(jpegSoi) {
		return Tags{ImageForm: jpgForm},
			errors.New(fmt.Sprintf("missing %s SOI '%v'", jpgForm, jpegSoi))
	}
	needApp0 := true
	needSof := true
	tags := Tags{ImageForm: jpgForm}
	var err = errors.New(fmt.Sprintf("missing %s EXIF", jpgForm))
	for {
		if !needSof && !needApp0 {
			return tags, err
		}
		_, eof := file.Read(segment)
		if eof != nil {
			return tags, err
		}
		if string(segment) == string(jpegEoi) {
			return tags, err
		}
		if segment[0] == 0xff {
			size := make([]byte, 2)
			n, e := file.Read(size)
			if n != 2 || e != nil {
				return tags, err
			}
			segSize := BigEndianToInt(size[0:2])
			// fmt.Printf("segment: %d %d\n", segment[1], segSize)
			segSize -= 2
			data := make([]byte, segSize)
			n, e = file.Read(data)
			if n != segSize || e != nil {
				return tags, err
			}
			if needApp0 {
				if segment[1] == jpegApp0[1] {
					code := string(data[0:4])
					if code == jpegCode {
						needApp0 = false
					}
				}
				continue
			}
			if needSof && (segment[1] == jpegSof0[1] || segment[1] == jpegSof2[1]) {
				needSof = false
				tags.ImageHeight = BigEndianToInt(data[1:3])
				tags.ImageWidth = BigEndianToInt(data[3:5])
				err = nil
			}
			continue
		}
		skip := make([]byte, 1)
		_, _ = file.Read(skip)
	}
}

var pngSignature = []uint8{137, 80, 78, 71, 13, 10, 26, 10}

const pngHead = "IHDR"

func (t *png) parse(file *os.File) (Tags, error) {
	buf := make([]byte, 8)
	l, err := file.Read(buf)
	if err != nil {
		fmt.Println(err)
	}
	if l == 8 {
		if string(pngSignature) != string(buf[0:8]) {
			return Tags{ImageForm: pngForm}, errors.New(fmt.Sprintf("missing %s '%v' signature",
				pngForm, pngSignature))
		}
		if err == nil {
			tags := Tags{ImageForm: pngForm}
			for {
				chunk := make([]byte, 8)
				_, eof := file.Read(chunk)
				if eof != nil {
					break
				}
				chunkLength := BigEndianToInt(chunk[0:4])
				chunkType := strings.ToUpper(string(chunk[4:8]))
				switch chunkType {
				case pngHead:
					head := make([]byte, 13)
					_, eof := file.Read(head)
					if eof == nil {
						tags.ImageWidth = BigEndianToInt(head[0:4])
						tags.ImageHeight = BigEndianToInt(head[4:8])
					}
				default:
					_, _ = file.Seek(int64(chunkLength), 1)
				}
			}
			return tags, nil
		}
	}
	return Tags{ImageForm: pngForm}, errors.New("insufficient bytes")

}

const tiffImageWidth = 256  // (100.H)
const tiffImageHeight = 257 // Length(101.h)
const tiffDatetime = 306    // (132.H) “YYYY:MM:DD HH:MM:SS”
const tiffMake = 271        // (10f.H)
const tiffModel = 272       // (110.H)
const tiffOrientation = 274 // (112.H)

func (t *tif) parse(file *os.File) (Tags, error) {
	buf := make([]byte, 8)
	l, err := file.Read(buf)
	if err != nil {
		fmt.Println(err)
	}
	toInt := func(bytes []uint8) int { return 0 }
	if l == 8 {
		switch string(buf[0:2]) {
		case "II":
			toInt = func(bytes []uint8) int {
				return LittleEndianToInt(bytes)
			}
		case "MM":
			toInt = func(bytes []uint8) int {
				return BigEndianToInt(bytes)
			}
		default:
			return Tags{ImageForm: tifForm}, errors.New(fmt.Sprintf("missing %s 'XX' endian code",
				tifForm))
		}
		tags := Tags{ImageForm: tifForm}
		offset := toInt(buf[4:8])
		_, err = file.Seek(int64(offset), 0)
		if err == nil {
			for {
				ifd := make([]byte, 12)
				_, eof := file.Read(ifd)
				if eof != nil {
					break
				}
				//field := toInt(ifd[0:2])
				fieldType := toInt(ifd[2:4])
				//fieldCount := toInt(ifd[4:8])
				// fmt.Printf("field: %d type %d count %d\n", field, fieldType, fieldCount)
				switch fieldType {
				case tiffImageWidth:
					tags.ImageWidth = toInt(ifd[10:12])
				case tiffImageHeight:
					tags.ImageHeight = toInt(ifd[10:12])
				case tiffDatetime:
				case tiffMake:
				case tiffModel:
				case tiffOrientation:
				}
			}
			return tags, nil
		}
	}
	return Tags{ImageForm: tifForm}, errors.New("insufficient bytes")
}
