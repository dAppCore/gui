package bytes

import core "dappco.re/go"

type Buffer struct {
	data []byte
}

func (b *Buffer) Write(p []byte) (int, error) {
	b.data = append(b.data, p...)
	return len(p), nil
}

func (b *Buffer) Bytes() []byte {
	return b.data
}

func (b *Buffer) Len() int {
	return len(b.data)
}

func (b *Buffer) String() string {
	return string(b.data)
}

func (b *Buffer) Reset() {
	b.data = nil
}

type Reader struct {
	data []byte
	off  int
}

func NewReader(data []byte) *Reader {
	return &Reader{data: data}
}

func NewBufferString(s string) core.Reader {
	return core.NewReader(s)
}

func (r *Reader) Read(p []byte) (int, error) {
	if r.off >= len(r.data) {
		return 0, core.EOF
	}
	n := copy(p, r.data[r.off:])
	r.off += n
	return n, nil
}

func Equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func Repeat(data []byte, count int) []byte {
	if count <= 0 || len(data) == 0 {
		return nil
	}
	out := make([]byte, 0, len(data)*count)
	for i := 0; i < count; i++ {
		out = append(out, data...)
	}
	return out
}
