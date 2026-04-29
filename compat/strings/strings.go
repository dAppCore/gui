package strings

import (
	"unicode"
	"unicode/utf8"

	core "dappco.re/go"
)

type Builder struct {
	parts []string
	size  int
}

func (b *Builder) WriteString(s string) (int, error) {
	b.parts = append(b.parts, s)
	b.size += len(s)
	return len(s), nil
}

func (b *Builder) WriteRune(r rune) (int, error) {
	s := string(r)
	return b.WriteString(s)
}

func (b *Builder) WriteByte(c byte) error {
	_, err := b.WriteString(string([]byte{c}))
	return err
}

func (b *Builder) String() string {
	return core.Join("", b.parts...)
}

func (b *Builder) Len() int {
	return b.size
}

func (b *Builder) Reset() {
	b.parts = nil
	b.size = 0
}

func Contains(s, substr string) bool {
	return core.Contains(s, substr)
}

func ContainsAny(s, chars string) bool {
	for _, r := range s {
		if ContainsRune(chars, r) {
			return true
		}
	}
	return false
}

func ContainsRune(s string, r rune) bool {
	for _, candidate := range s {
		if candidate == r {
			return true
		}
	}
	return false
}

func Cut(s, sep string) (before, after string, found bool) {
	if sep == "" {
		return "", s, true
	}
	idx := Index(s, sep)
	if idx < 0 {
		return s, "", false
	}
	return s[:idx], s[idx+len(sep):], true
}

func CutPrefix(s, prefix string) (after string, found bool) {
	if !HasPrefix(s, prefix) {
		return s, false
	}
	return TrimPrefix(s, prefix), true
}

func EqualFold(a, b string) bool {
	return ToLower(a) == ToLower(b)
}

func Fields(s string) []string {
	var fields []string
	start := -1
	for i, r := range s {
		if unicode.IsSpace(r) {
			if start >= 0 {
				fields = append(fields, s[start:i])
				start = -1
			}
			continue
		}
		if start < 0 {
			start = i
		}
	}
	if start >= 0 {
		fields = append(fields, s[start:])
	}
	return fields
}

func HasPrefix(s, prefix string) bool {
	return core.HasPrefix(s, prefix)
}

func HasSuffix(s, suffix string) bool {
	return core.HasSuffix(s, suffix)
}

func Index(s, substr string) int {
	if substr == "" {
		return 0
	}
	if len(substr) > len(s) {
		return -1
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func Join(elems []string, sep string) string {
	return core.Join(sep, elems...)
}

func NewReader(s string) core.Reader {
	return core.NewReader(s)
}

func Repeat(s string, count int) string {
	if count <= 0 || s == "" {
		return ""
	}
	out := make([]string, count)
	for i := range out {
		out[i] = s
	}
	return core.Join("", out...)
}

func ReplaceAll(s, old, new string) string {
	return core.Replace(s, old, new)
}

func Split(s, sep string) []string {
	return core.Split(s, sep)
}

func SplitN(s, sep string, n int) []string {
	return core.SplitN(s, sep, n)
}

func ToLower(s string) string {
	return core.Lower(s)
}

func Trim(s, cutset string) string {
	return trimLeftRight(s, cutset, true, true)
}

func TrimPrefix(s, prefix string) string {
	return core.TrimPrefix(s, prefix)
}

func TrimRight(s, cutset string) string {
	return trimLeftRight(s, cutset, false, true)
}

func TrimSpace(s string) string {
	return core.Trim(s)
}

func TrimSuffix(s, suffix string) string {
	return core.TrimSuffix(s, suffix)
}

func trimLeftRight(s, cutset string, left, right bool) string {
	if cutset == "" || s == "" {
		return s
	}
	start := 0
	end := len(s)
	if left {
		for start < end {
			r, size := utf8.DecodeRuneInString(s[start:end])
			if !ContainsRune(cutset, r) {
				break
			}
			start += size
		}
	}
	if right {
		for start < end {
			r, size := utf8.DecodeLastRuneInString(s[start:end])
			if !ContainsRune(cutset, r) {
				break
			}
			end -= size
		}
	}
	return s[start:end]
}
