package collection

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	collectionTitleLengthMin = 3
	collectionTitleLengthMax = 255
)

func ValidateTitle(t string) (string, error) {
	t = strings.TrimSpace(t)

	// Use RuneCountInString instead of len to count human-readable characters,
	// not bytes. Non-ASCII characters (e.g. Polish ąęł, CJK) are multi-byte
	// in UTF-8 and would inflate the byte count, causing valid titles to be
	// rejected or invalid ones to pass.
	n := utf8.RuneCountInString(t)
	if n < collectionTitleLengthMin || n > collectionTitleLengthMax {
		return t, ErrValidationTitleLength
	}

	// Reject control characters (null bytes, tabs, newlines, etc.) which can
	// cause issues in logs, CSV exports, database collation, or rendering.
	for _, r := range t {
		if unicode.IsControl(r) {
			return t, ErrValidationTitleInvalidCharacters
		}
	}

	return t, nil
}

const (
	collectionDescriptionLengthMax = 512
)

func ValidateDescription(d string) (string, error) {
	d = strings.TrimSpace(d)

	// Use RuneCountInString instead of len to count human-readable characters,
	// not bytes. Non-ASCII characters (e.g. Polish ąęł, CJK) are multi-byte
	// in UTF-8 and would inflate the byte count, causing valid descriptions to be
	// rejected or invalid ones to pass.
	if utf8.RuneCountInString(d) > collectionDescriptionLengthMax {
		return d, ErrValidationDescriptionLength
	}

	// Reject control characters (null bytes, tabs, newlines, etc.) which can
	// cause issues in logs, CSV exports, database collation, or rendering.
	for _, r := range d {
		if unicode.IsControl(r) {
			return d, ErrValidationDescriptionInvalidCharacters
		}
	}

	return d, nil
}
