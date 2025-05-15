package filter

import (
	"fmt"
	"os"
	"strings"

	"bitmap/config"
	"bitmap/internal/core"
)

var filterRegistry = map[string]func(*core.BitMap){
	"blue":      ApplyBlueFilter,
	"red":       ApplyRedFilter,
	"green":     ApplyGreenFilter,
	"grayscale": ApplyGrayscaleFilter,
	"negative":  ApplyNegativeFilter,
	"pixelate":  ApplyPixelateFilter,
	"blur":      ApplyBlurFilter,
}

var pixelateCount = 0

func HandleFilter(b *core.BitMap) {
	if len(config.FilterFlag) == 0 {
		return
	}

	filter := config.FilterFlag[0]
	filterFunc, exists := filterRegistry[strings.ToLower(filter)]
	if !exists {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Filter not found")
		os.Exit(1)
	}

	filterFunc(b)
	config.FilterFlag = config.FilterFlag[1:]
}

// Cycle Helper function to transform pixels
func Cycle(b *core.BitMap, cycleFunc func(pixel *core.Pixel)) {
	pixels := b.GetPixels()
	for i := 0; i < len(pixels); i++ {
		for j := 0; j < len(pixels[i]); j++ {
			cycleFunc(pixels[i][j]) // here sending pixels
		}
	}
}

func ApplyRedFilter(b *core.BitMap) {
	Cycle(b, func(pixel *core.Pixel) {
		pixel.Blue = 0
		pixel.Green = 0
	})
}

func ApplyBlueFilter(b *core.BitMap) {
	Cycle(b, func(pixel *core.Pixel) {
		pixel.Red = 0
		pixel.Green = 0
	})
}

func ApplyGreenFilter(b *core.BitMap) {
	Cycle(b, func(pixel *core.Pixel) {
		pixel.Red = 0
		pixel.Blue = 0
	})
}

func ApplyGrayscaleFilter(b *core.BitMap) {
	Cycle(b, func(pixel *core.Pixel) {
		r, g, b := pixel.Red, pixel.Green, pixel.Blue
		// Convert the pixel to grayscale using wighted coefficients
		// for human eye perception of red, green, blue.
		grayScale := uint8(0.3*float64(r) + 0.59*float64(g) + 0.11*float64(b))
		pixel.Red = grayScale
		pixel.Green = grayScale
		pixel.Blue = grayScale
	})
}

func ApplyNegativeFilter(b *core.BitMap) {
	Cycle(b, func(pixel *core.Pixel) {
		pixel.Red = 255 - pixel.Red
		pixel.Green = 255 - pixel.Green
		pixel.Blue = 255 - pixel.Blue
	})
}

func ApplyPixelateFilter(b *core.BitMap) {
	// Get the 2D array of pixels from the image
	pixels := b.GetPixels()
	// Get the height and width of the image
	height, width := b.GetDimensions()
	// Define the block size for pixelation. 20x20 pixels
	blockSize := 20 + pixelateCount*10
	// Loop through the image with a step equal to the block size
	for x := 0; x < int(height); x += blockSize {
		for y := 0; y < int(width); y += blockSize {
			// Initialize variables to accumulate color values (red, green, and blue)
			var avgRed, avgGreen, avgBlue int
			pixelCount := 0 // Counter for the number of pixels in the block
			// Loop through each pixel inside the current block
			for i := 0; i < blockSize && x+i < int(height); i++ {
				for j := 0; j < blockSize && y+j < int(width); j++ {
					// Add the pixel's color values to the accumulated variables
					avgRed += int(pixels[x+i][y+j].Red)
					avgGreen += int(pixels[x+i][y+j].Green)
					avgBlue += int(pixels[x+i][y+j].Blue)
					// Increase the pixel count for the block
					pixelCount++
				}
			}
			// Apply the average color to all pixels inside the block
			for i := 0; i < blockSize && x+i < int(height); i++ {
				for j := 0; j < blockSize && y+j < int(width); j++ {
					// Replace the pixel's color with the average color and
					// calculate the average color values for the current block
					pixels[x+i][y+j].Red = byte(avgRed / pixelCount)
					pixels[x+i][y+j].Green = byte(avgGreen / pixelCount)
					pixels[x+i][y+j].Blue = byte(avgBlue / pixelCount)
				}
			}
		}
	}
	pixelateCount++
}

func ApplyBlurFilter(b *core.BitMap) {
	pixel := b.GetPixels()
	h, w := b.GetDimensions()
	// Loop to find all indexes
	for y := int32(0); y < h; y++ {
		for x := int32(0); x < w; x++ {
			var count int
			var avgRed, avgGreen, avgBlue int
			// Loop to find all neighbors
			for i := -10; i <= 10; i++ {
				for j := -10; j <= 10; j++ {
					nx := x + int32(i)
					ny := y + int32(j)
					// To check pixels without of range of array
					if nx >= 0 && nx < w && ny >= 0 && ny < h {
						avgRed += int(pixel[ny][nx].Red)
						avgGreen += int(pixel[ny][nx].Green)
						avgBlue += int(pixel[ny][nx].Blue)
						// Count of neighbors
						count++
					}
				}
			}
			// Find average value for pixel
			avgRed /= count
			avgGreen /= count
			avgBlue /= count
			// Set new pixels
			pixel[y][x].Red = byte(avgRed)
			pixel[y][x].Green = byte(avgGreen)
			pixel[y][x].Blue = byte(avgBlue)
		}
	}
}
