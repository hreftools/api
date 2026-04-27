package collection

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	collectionNameLengthMin = 2
	collectionNameLengthMax = 255
)

func ValidateName(n string) (string, error) {
	n = strings.TrimSpace(n)

	// Use RuneCountInString instead of len to count human-readable characters,
	// not bytes. Non-ASCII characters (e.g. Polish ąęł, CJK) are multi-byte
	// in UTF-8 and would inflate the byte count, causing valid names to be
	// rejected or invalid ones to pass.
	count := utf8.RuneCountInString(n)
	if count < collectionNameLengthMin || count > collectionNameLengthMax {
		return n, ErrValidationNameLength
	}

	// Reject control characters (null bytes, tabs, newlines, etc.) which can
	// cause issues in logs, CSV exports, database collation, or rendering.
	for _, r := range n {
		if unicode.IsControl(r) {
			return n, ErrValidationNameInvalidCharacters
		}
	}

	return n, nil
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
