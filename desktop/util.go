package main

import (
	"bytes"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func slugifyStr(s string) string {
	b := make([]byte, len(s)*4)
	n, _ := transform.NewReader(strings.NewReader(s), norm.NFD).Read(b)
	var buf bytes.Buffer
	for _, r := range string(b[:n]) {
		if !unicode.Is(unicode.Mn, r) {
			buf.WriteRune(r)
		}
	}
	result := strings.ToLower(buf.String())
	re := regexp.MustCompile(`[^a-z0-9]+`)
	result = re.ReplaceAllString(result, "-")
	result = strings.Trim(result, "-")
	words := strings.Split(result, "-")
	if len(words) > 6 {
		words = words[:6]
	}
	return strings.Join(words, "-")
}
