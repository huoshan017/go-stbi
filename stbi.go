// Package stbi provides go bindings for stb_image.h at v2.23
//
// See subpackages for format specific codecs for use with image.Decode and
// image.DecodeConfig.
package stbi // import "huoshan017/go-stbi"

import (
	"errors"
	"image"
	"io"
	"io/ioutil"
	"os"
	"unsafe"
)

// #cgo LDFLAGS: -lm
// #define STB_IMAGE_IMPLEMENTATION
// #define STBI_FAILURE_USERMSG
// #include "stb_image.h"
import "C"

// Load wraps stbi_load to decode an image into an RGBA pixel struct.
func Load(path string, channelsInFile *int32, desiredChannels int32) (*image.RGBA, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	var x, y C.int
	data := C.stbi_load(cpath, &x, &y, (*C.int)(channelsInFile), C.int(desiredChannels))
	if data == nil {
		msg := C.GoString(C.stbi_failure_reason())
		return nil, errors.New(msg)
	}
	defer C.stbi_image_free(unsafe.Pointer(data))

	return &image.RGBA{
		Pix:    C.GoBytes(unsafe.Pointer(data), y*x*C.int(*channelsInFile)),
		Stride: 4,
		Rect:   image.Rect(0, 0, int(x), int(y)),
	}, nil
}

// LoadFile wraps stbi_load_from_file to decode an image into an RGBA pixel
// struct.
func LoadFile(f *os.File, channelsInFile *int32, desiredChannels int32) (*image.RGBA, error) {
	mode := C.CString("rb")
	defer C.free(unsafe.Pointer(mode))
	fp, err := C.fdopen(C.int(f.Fd()), mode)
	if err != nil {
		return nil, err
	}

	var x, y C.int
	data := C.stbi_load_from_file(fp, &x, &y, (*C.int)(channelsInFile), C.int(desiredChannels))
	if data == nil {
		msg := C.GoString(C.stbi_failure_reason())
		return nil, errors.New(msg)
	}
	defer C.stbi_image_free(unsafe.Pointer(data))

	return &image.RGBA{
		Pix:    C.GoBytes(unsafe.Pointer(data), y*x*C.int(*channelsInFile)),
		Stride: 4,
		Rect:   image.Rect(0, 0, int(x), int(y)),
	}, nil
}

// LoadMemory wraps stbi_load_from_memory to decode an image into an RGBA
// pixel struct.
func LoadMemory(b []byte, channelsInFile *int32, desiredChannels int32) (*image.RGBA, error) {
	var x, y C.int
	mem := (*C.uchar)(unsafe.Pointer(&b[0]))
	data := C.stbi_load_from_memory(mem, C.int(len(b)), &x, &y, (*C.int)(channelsInFile), C.int(desiredChannels))
	if data == nil {
		msg := C.GoString(C.stbi_failure_reason())
		return nil, errors.New(msg)
	}
	defer C.stbi_image_free(unsafe.Pointer(data))

	return &image.RGBA{
		Pix:    C.GoBytes(unsafe.Pointer(data), y*x*C.int(*channelsInFile)),
		Stride: 4,
		Rect:   image.Rect(0, 0, int(x), int(y)),
	}, nil
}

// LoadReader delegates to LoadFile if r is an *os.File, otherwise,
// LoadMemory after reading the contents.
func LoadReader(r io.Reader, channelsInFile *int32, desiredChannels int32) (*image.RGBA, error) {
	if f, ok := r.(*os.File); ok {
		return LoadFile(f, channelsInFile, desiredChannels)
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return LoadMemory(b, channelsInFile, desiredChannels)
}

func SetFlipVerticallyOnLoad(flagTrueIfShouldFlip bool) {
	C.stbi_set_flip_vertically_on_load(func() C.int {
		if flagTrueIfShouldFlip {
			return 1
		} else {
			return 0
		}
	}())
}
