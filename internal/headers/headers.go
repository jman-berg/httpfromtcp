package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}

	headersText := string(data[:idx])
	headersTrimmedSpace := strings.TrimSpace(headersText)
	headersSplitKeyValue := strings.SplitN(headersTrimmedSpace, ":", 2)

	key := headersSplitKeyValue[0]
	if key != strings.TrimSpace(key) {
		return 0, false, fmt.Errorf("Invalid spacing header: %s", headersText)
	}
	if !isValidString(key) {
		return 0, false, fmt.Errorf("Invalid character in header key: %s", key)
	}
	if len(key) < 1 {
		return 0, false, fmt.Errorf("Key should have at least a length of one")
	}

	value := strings.TrimSpace(headersSplitKeyValue[1])

	h.Set(strings.ToLower(key), value)
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	v, ok := h[key]
	if ok {
		value = strings.Join([]string{v, value}, ", ")
	}
	h[key] = value
}

func isValidString(s string) bool {
	allowedChars := map[rune]bool{
		'!': true, '#': true, '$': true, '%': true, '&': true, '\'': true,
		'*': true, '+': true, '-': true, '.': true, '^': true, '_': true,
		'`': true, '|': true, '~': true,
	}
	for _, r := range s {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			continue
		}

		if allowedChars[r] {
			continue
		}
		return false
	}
	return true
}
