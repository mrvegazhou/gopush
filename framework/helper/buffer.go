package helper

import (
"errors"
"io"
)

var (
	ErrNotEnough = errors.New("not enough")
)

type Buffer struct {
	reader io.Reader
	buf    []byte
	start  int
	end    int
}

func NewBuffer(reader io.Reader, len int) Buffer {
	buf := make([]byte, len)
	return Buffer{reader, buf, 0, 0}
}

func (b *Buffer) Len() int {
	return b.end - b.start
}

// grow 将有用的字节前移
func (b *Buffer) Grow() {
	if b.start == 0 {
		return
	}
	copy(b.buf, b.buf[b.start:b.end])
	b.end -= b.start
	b.start = 0
}

// readFromReader 从reader里面读取数据，如果reader阻塞，会发生阻塞
func (b *Buffer) ReadFromReader() (int, error) {
	b.Grow()
	n, err := b.reader.Read(b.buf[b.end:])
	if err != nil {
		return n, err
	}
	b.end += n
	return n, nil
}

// seek 返回n个字节，而不产生移位，如果没有足够字节，返回错误
func (b *Buffer) Seek(start, end int) ([]byte, error) {
	if b.end-b.start >= end-start {
		buf := b.buf[b.start+start : b.start+end]
		return buf, nil
	}
	return nil, ErrNotEnough
}

// read 舍弃offset个字段，读取n个字段,如果没有足够的字节，返回错误
func (b *Buffer) Read(offset, limit int) ([]byte, error) {
	if b.Len() < offset+limit {
		return nil, ErrNotEnough
	}
	b.start += offset
	buf := b.buf[b.start : b.start+limit]
	b.start += limit
	return buf, nil
}

