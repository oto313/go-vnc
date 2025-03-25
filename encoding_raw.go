package vnc

import (
	"MediaProCon/internal/common"
	"fmt"
	"io"
	"time"
)

// An Encoding implements a method for encoding pixel data that is
// sent by the server to the client.
type Encoding interface {
	// The number that uniquely identifies this encoding type.
	Type() int32

	// Read reads the contents of the encoded pixel data from the reader.
	// This should return a new Encoding implementation that contains
	// the proper data.
	Read(*ClientConn, *Rectangle, io.Reader) (Encoding, error)
}

// RawEncoding is raw pixel data sent by the server.
//
// See RFC 6143 Section 7.7.1
type RawEncoding struct {
	Framebuffer common.Framebuffer
	rectBuffer  []byte
}

func (e *RawEncoding) SetBufferSize(width int, height int) {
	e.Framebuffer.Lock()
	defer e.Framebuffer.Unlock()
	e.Framebuffer.Width = width
	e.Framebuffer.Height = height
	e.Framebuffer.BytePerPixel = 3
	e.Framebuffer.Data = make([]byte, width*width*e.Framebuffer.BytePerPixel)

	e.rectBuffer = make([]uint8, len(e.Framebuffer.Data))
}

func (*RawEncoding) Type() int32 {
	return 0
}

func (e *RawEncoding) Read(c *ClientConn, rect *Rectangle, r io.Reader) (Encoding, error) {
	bytesPerPixel := int(c.PixelFormat.BPP / 8)
	startTime := time.Now().Local()
	rectLen := int(rect.Height) * int(rect.Width) * int(bytesPerPixel)
	_, err := io.ReadFull(r, e.rectBuffer[:rectLen])
	if err != nil {
		return nil, err
	}
	e.Framebuffer.Lock()
	defer e.Framebuffer.Unlock()
	buffer := e.Framebuffer.Data
	for y := 0; y < int(rect.Height); y++ {
		for x := 0; x < int(rect.Width); x++ {
			srcColorIndex := (y*int(rect.Width) + x) * bytesPerPixel
			pixelBytes := e.rectBuffer[srcColorIndex : srcColorIndex+3]

			dstColorIndex := ((int(rect.Y)+y)*int(c.FrameBufferWidth) + int(rect.X) + x) * 3

			buffer[dstColorIndex+2] = pixelBytes[0]
			buffer[dstColorIndex+1] = pixelBytes[1]
			buffer[dstColorIndex+0] = pixelBytes[2]
		}
	}
	fmt.Printf("Time elapsed: %v\n", time.Since(startTime))
	return e, nil
}
