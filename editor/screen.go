package editor

import (
	"unicode/utf8"
)

func isWideRune(r rune) bool {
	return r >= 0x4E00 && r <= 0x9FFF ||
		r >= 0xAC00 && r <= 0xD7AF ||
		r >= 0x3000 && r <= 0x303F ||
		r >= 0xFF00 && r <= 0xFFEF
}

func CalculateVisualWidth(s string) int {
	width := 0
	for _, r := range s {
		if r == '\t' {
			width += 8
		} else if isWideRune(r) {
			width += 2
		} else {
			width += 1
		}
	}
	return width
}

func RuneToVisualIndex(s string, visualPos int) int {
	width := 0
	for i := range s {
		r, _ := utf8.DecodeRuneInString(s[i:])
		charWidth := 1
		if r == '\t' {
			charWidth = 8
		} else if isWideRune(r) {
			charWidth = 2
		}
		if width+charWidth > visualPos {
			return i
		}
		width += charWidth
	}
	return len(s)
}
