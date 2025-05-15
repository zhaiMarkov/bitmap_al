package core

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

var (
	optionalHeader []byte
	lastData       []byte
)

type BitMap struct {
	header     *BMPHeader
	infoHeader *DIBHeader
	pixels     [][]*Pixel
}

type BMPHeader struct {
	FileType     [2]byte
	FileSize     uint32
	Reserved1    uint16
	Reserved2    uint16
	BitmapOffset uint32
}

type DIBHeader struct {
	HeaderSize      uint32
	Width           int32
	Height          int32
	Planes          uint16
	BitsPerPixel    uint16
	Compression     uint32
	ImageSize       uint32
	XPixelsPerMeter int32
	YPixelsPerMeter int32
	ColorsUsed      uint32
	ColorsImportant uint32
}

type Pixel struct {
	Blue  byte
	Green byte
	Red   byte
}

func NewBitMap() *BitMap {
	return &BitMap{
		header:     &BMPHeader{},
		infoHeader: &DIBHeader{},
		pixels:     nil,
	}
}

func (b *BitMap) Read(r io.Reader) {
	var err error

	err = b.header.Read(r)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %v \n", err)
		os.Exit(1)
	}

	err = b.infoHeader.Read(r)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %v \n", err)
		os.Exit(1)
	}

	if int(b.header.BitmapOffset) > 54 {
		paddingSize := int(b.header.BitmapOffset) - 54
		temp := make([]byte, paddingSize) // Allocate a slice of the required size
		_, err := r.Read(temp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(1)
		}

		optionalHeader = append(optionalHeader, temp...) // Append all bytes read to data
	}

	h, w := b.GetDimensions()
	arr := make([][]*Pixel, h)
	for i := int32(0); i < h; i++ {
		arr[i] = make([]*Pixel, w)
		for j := int32(0); j < w; j++ {
			p := &Pixel{}
			err = p.Read(r)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "ERROR: %v \n", err)
				os.Exit(1)
			}
			arr[i][j] = p
		}
		if int(w)%4 != 0 {
			for i := 0; i < (4-(int((w*3))%4))%4; i++ {
				var p byte
				err := binary.Read(r, binary.LittleEndian, &p)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "ERROR: %v \n", err)
					os.Exit(1)
				}
			}
		}

	}

	b.pixels = arr
	temp := make([]byte, 1024)
	for {
		n, err := r.Read(temp)
		if err != nil && err != io.EOF {
			_, _ = fmt.Fprintf(os.Stderr, "ERROR: %v \n", err)
			os.Exit(1)
		}
		if n == 0 {
			break
		}
		lastData = append(lastData, temp[:n]...)
	}
}

func (b *BitMap) GetInfoHeader() *DIBHeader {
	return b.infoHeader
}

func (b *BitMap) GetHeader() *BMPHeader {
	return b.header
}

func (b *BitMap) GetPixels() [][]*Pixel {
	return b.pixels
}

func (b *BitMap) SetPixels(pixels [][]*Pixel) {
	b.pixels = pixels
}

func (b *BitMap) GetDimensions() (int32, int32) {
	return b.infoHeader.Height, b.infoHeader.Width
}

func (b *BitMap) SetDimensions(height, width int32) {
	b.infoHeader.Height, b.infoHeader.Width = height, width
}

func (b *BitMap) GetImageSize() uint32 {
	return b.infoHeader.ImageSize
}

func (b *BitMap) SetImageSize(imageSize uint32) {
	b.infoHeader.ImageSize = imageSize
}

func (b *BitMap) GetFileSize() uint32 {
	return b.header.FileSize
}

func (b *BitMap) SetFileSize(fileSize uint32) {
	b.header.FileSize = fileSize
}

func (b *BitMap) Save(w io.Writer) {
	err := binary.Write(w, binary.LittleEndian, b.header)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %v \n", err)
		os.Exit(1)
	}

	err = binary.Write(w, binary.LittleEndian, b.infoHeader)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %v \n", err)
		os.Exit(1)
	}

	for _, v := range optionalHeader {
		err = binary.Write(w, binary.LittleEndian, v)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}
	}

	for _, row := range b.pixels {
		for _, pixel := range row {
			err = binary.Write(w, binary.LittleEndian, pixel)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "ERROR: %v \n", err)
				os.Exit(1)
			}
		}
		if len(row)%4 != 0 {
			for i := 0; i < (4-(len(row)*3)%4)%4; i++ {
				err = binary.Write(w, binary.LittleEndian, byte(0))
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "ERROR: %v \n", err)
					os.Exit(1)
				}
			}
		}
	}
	for _, v := range lastData {
		err = binary.Write(w, binary.LittleEndian, v)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "ERROR: %v \n", err)
			os.Exit(1)
		}
	}
}

func (b *BMPHeader) Read(r io.Reader) (err error) {
	err = binary.Read(r, binary.LittleEndian, b)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	if b.FileType != [2]byte{'B', 'M'} {
		return fmt.Errorf("invalid file type")
	}
	if b.Reserved1 != 0 || b.Reserved2 != 0 {
		return fmt.Errorf("invalid reserved field")
	}
	if b.BitmapOffset < 54 {
		return fmt.Errorf("invalid bitmap offset")
	}
	return nil
}

func (d *DIBHeader) Read(r io.Reader) (err error) {
	err = binary.Read(r, binary.LittleEndian, d)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if d.BitsPerPixel != 24 {
		return fmt.Errorf("unsupported bits per pixel: %d", d.BitsPerPixel)
	}
	return nil
}

func (p *Pixel) Read(r io.Reader) (err error) {
	err = binary.Read(r, binary.LittleEndian, p)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}
