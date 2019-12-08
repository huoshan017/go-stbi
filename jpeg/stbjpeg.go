package stbpng

import (
	"io"
	"image"
	"image/color"
	"image/jpeg"

	"github.com/neilpa/go-stbi"
)

const jpegHeader = "\xff\xd8"

// Decode reads a JPEG image from r and returns an image.RGBA.
func Decode(r io.Reader) (image.Image, error) {
	return stbi.LoadReader(r)
}

// DecodeConfig returns the dimensions and an RGBA color model of the JPEG
// backed by reader. Note this simply wraps the stdlib jpeg.DecodeConfig.
func DecodeConfig(r io.Reader) (image.Config, error) {
	c, err := jpeg.DecodeConfig(r)
	c.ColorModel = color.RGBAModel
	return c, err
}

func init() {
	image.RegisterFormat("jpeg", jpegHeader, Decode, DecodeConfig)
}
