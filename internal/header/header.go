package header

import (
	"fmt"
	"log"

	"bitmap/internal/core"
)

func PrintHeaderInfo(b *core.BitMap) {
	if b == nil {
		log.Fatal("bitmap is nil")
	}

	if b.GetHeader() == nil {
		log.Fatal("failed to get header")
	}

	if b.GetInfoHeader() == nil {
		log.Fatal("failed to get info header")
	}

	if string(b.GetHeader().FileType[:]) != "BM" {
		log.Fatal("Error: not a bitmap file")
	}

	fmt.Println("BMP Header:")
	fmt.Printf("- FileType: %s\n", string(b.GetHeader().FileType[:]))
	fmt.Printf("- FileSizeInBytes: %d\n", b.GetFileSize())
	fmt.Printf("- HeaderSize: %d\n", 14+b.GetInfoHeader().HeaderSize)
	fmt.Println("DIB Header:")
	fmt.Printf("- DibHeaderSize: %d\n", b.GetInfoHeader().HeaderSize)
	fmt.Printf("- WidthInPixels: %d\n", b.GetInfoHeader().Width)
	fmt.Printf("- HeightInPixels: %d\n", b.GetInfoHeader().Height)
	fmt.Printf("- PixelSizeInBits: %d\n", b.GetInfoHeader().BitsPerPixel)
	fmt.Printf("- ImageSizeInBytes: %d\n", b.GetImageSize())
}
