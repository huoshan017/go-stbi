// Package jpeg provides a JPEG decoder via the go bindings for stb_image.h
// and registers for use with image.Decode and image.DecodeConfig.
package jpeg // import "neilpa.me/go-stbi/jpeg"

import (
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"io"

	"huoshan017/go-stbi"
)

// Header is the magic string at the start of a JPEG file.
const Header = "\xff\xd8"

// ErrInvalid is returned from DecodeConfig for non JPEG files.
var ErrInvalid = errors.New("Invalid JPEG")

// Decode reads a JPEG image from r and returns an image.RGBA.
func Decode(r io.Reader) (image.Image, error) {
	return stbi.LoadReader(r)
}

// DecodeConfig returns the dimensions and an RGBA color model of the JPEG
// backed by reader.
func DecodeConfig(r io.Reader) (image.Config, error) {
	cfg := image.Config{ColorModel: color.RGBAModel}

	var magic [2]byte
	err := binary.Read(r, binary.LittleEndian, &magic)
	if err != nil {
		return cfg, err
	}
	if string(magic[:]) != Header {
		return cfg, ErrInvalid
	}

	var h segmentHeader
	var buf []byte
	for {
		err = binary.Read(r, binary.BigEndian, &h)
		if err != nil {
			break
		}
		if h.Sentinel != 0xff {
			err = ErrInvalid
			break
		}
		switch h.Marker {
		// Start of frames
		case 0xc0, 0xc1, 0xc2:
			var dim struct {
				_    byte
				H, W uint16
			}
			err = binary.Read(r, binary.BigEndian, &dim)
			cfg.Width, cfg.Height = int(dim.W), int(dim.H)
			return cfg, err

		default:
			if len(buf) < int(h.Length) {
				buf = make([]byte, int(h.Length))
			}
			// The length above includes the 2 bytes for the length itself
			err = binary.Read(r, binary.BigEndian, buf[:int(h.Length)-2])
			if err != nil {
				return cfg, err
			}
		}
	}

	return cfg, err
}

func init() {
	image.RegisterFormat("jpeg", Header, Decode, DecodeConfig)
}

type segmentHeader struct {
	Sentinel, Marker byte
	Length           uint16
}
