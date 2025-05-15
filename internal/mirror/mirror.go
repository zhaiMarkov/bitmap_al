package mirror

import (
	"fmt"
	"os"
	"strings"

	"bitmap/config"
	"bitmap/internal/core"
)

var height, width int32

func MirrorHorizontally(pixels [][]*core.Pixel) [][]*core.Pixel {
	newPixels := make([][]*core.Pixel, height)
	for y := int32(0); y < height; y++ {
		newPixels[y] = make([]*core.Pixel, width)
		for x := int32(0); x < width; x++ {
			newX := width - x - 1
			newPixels[y][x] = pixels[y][newX]
		}
	}

	return newPixels
}

func MirrorVertically(pixels [][]*core.Pixel) [][]*core.Pixel {
	newPixels := make([][]*core.Pixel, height)
	for y := int32(0); y < height; y++ {
		newPixels[y] = make([]*core.Pixel, width)
		for x := int32(0); x < width; x++ {
			newY := height - y - 1
			newPixels[y][x] = pixels[newY][x]
		}
	}

	return newPixels
}

func HandleMirror(bm *core.BitMap) {
	if len(config.MirrorFlag) == 0 {
		return
	}

	pixels := bm.GetPixels()
	height, width = bm.GetDimensions()

	cmd := strings.ToLower(config.MirrorFlag[0])
	switch {
	case strings.HasPrefix("horizontally", cmd):
		pixels = MirrorHorizontally(pixels)

	case strings.HasPrefix("vertically", cmd):
		pixels = MirrorVertically(pixels)

	default:
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Invalid mirror command")
		os.Exit(1)
	}

	config.MirrorFlag = config.MirrorFlag[1:]
	bm.SetPixels(pixels)
	bm.SetDimensions(height, width)
}
