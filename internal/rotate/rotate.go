package rotate

import (
	"fmt"
	"os"
	"strings"

	"bitmap/config"
	"bitmap/internal/core"
)

var globalHeight, globalWidth int32

var rotationMap = map[string]int{
	"right": 3,
	"90":    3,
	"-270":  3,
	"left":  1,
	"270":   1,
	"-90":   1,
	"180":   2,
	"-180":  2,
}

func RotateBMP(pixels [][]*core.Pixel) [][]*core.Pixel {
	rotated := make([][]*core.Pixel, globalWidth)
	for x := int32(0); x < globalWidth; x++ {
		rotated[x] = make([]*core.Pixel, globalHeight)
		for y := int32(0); y < globalHeight; y++ {
			srcY := globalHeight - 1 - y
			if srcY >= 0 && srcY < globalHeight {
				rotated[x][y] = pixels[srcY][x]
			}
		}
	}
	return rotated
}

func rotateImage(pixels [][]*core.Pixel, rotations int) [][]*core.Pixel {
	for i := 0; i < rotations; i++ {
		pixels = RotateBMP(pixels)
		globalHeight, globalWidth = globalWidth, globalHeight
	}
	return pixels
}

func HandleRotate(b *core.BitMap) {
	if len(config.RotateFlag) == 0 {
		return
	}

	pixels := b.GetPixels()
	globalHeight, globalWidth = b.GetDimensions()

	rotation, exists := rotationMap[strings.ToLower(config.RotateFlag[0])]
	if !exists {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: The rotation flag is specified incorrectly")
		os.Exit(1)
	}

	if rotation > 0 {
		pixels = rotateImage(pixels, rotation)
	}

	config.RotateFlag = config.RotateFlag[1:]
	b.SetPixels(pixels)
	b.SetDimensions(globalHeight, globalWidth)
}
