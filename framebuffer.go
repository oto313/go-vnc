package vnc

import "sync"

type Framebuffer struct {
	Width        int
	Height       int
	BytePerPixel int
	Data         []byte
	mutex        sync.Mutex
}

func (b *Framebuffer) Lock() {
	b.mutex.Lock()
}

func (b *Framebuffer) Unlock() {
	b.mutex.Unlock()
}
