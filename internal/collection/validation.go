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

	n := utf8.RuneCountInString(t)
	if n < collectionTitleLengthMin || n > collectionTitleLengthMax {
		return t, ErrValidationTitleLength
	}

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

	if utf8.RuneCountInString(d) > collectionDescriptionLengthMax {
		return d, ErrValidationDescriptionLength
	}

	for _, r := range d {
		if unicode.IsControl(r) {
			return d, ErrValidationDescriptionInvalidCharacters
		}
	}

	return d, nil
}
