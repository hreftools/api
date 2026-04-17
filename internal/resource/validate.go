package resource

import (
	"net"
	"net/url"
	"strings"
	"unicode/utf8"
)

const (
	resourceTitleLengthMin = 3
	resourceTitleLengthMax = 255
)

func validateTitle(t string) (string, error) {
	t = strings.TrimSpace(t)

	// Use RuneCountInString instead of len to count human-readable characters,
	// not bytes. Non-ASCII characters (e.g. Polish ąęł, CJK) are multi-byte
	// in UTF-8 and would inflate the byte count, causing valid titles to be
	// rejected or invalid ones to pass.
	if utf8.RuneCountInString(t) < resourceTitleLengthMin || utf8.RuneCountInString(t) > resourceTitleLengthMax {
		return t, ErrValidationTitleLength
	}

	return t, nil
}

const (
	resourceDescriptionLengthMax = 512
)

func validateDescription(d string) (string, error) {
	d = strings.TrimSpace(d)

	// Use RuneCountInString instead of len to count human-readable characters,
	// not bytes. See validateTitle for details.
	if utf8.RuneCountInString(d) > resourceDescriptionLengthMax {
		return d, ErrValidationDescriptionLength
	}

	return d, nil
}

func validateURL(u string) (string, error) {
	u = strings.TrimSpace(u)

	// 2048 characters is the practical URL length limit enforced by most
	// modern browsers (Chrome, Firefox, Safari)
	if len(u) > 2048 {
		return u, ErrValidationURLTooLong
	}

	uParsed, err := url.Parse(u)
	if err != nil || uParsed.Host == "" {
		return u, ErrValidationURLFormat
	}

	// Only allow http and https schemes to prevent XSS via javascript:,
	// data:, or other dangerous URI schemes when rendered as clickable links.
	if uParsed.Scheme != "http" && uParsed.Scheme != "https" {
		return u, ErrValidationURLFormat
	}

	if isPrivateHost(uParsed.Host) {
		return u, ErrValidationURLPrivate
	}

	return uParsed.String(), nil
}

func validateFavourite(f *bool) (bool, error) {
	if f == nil {
		return false, ErrValidationFavouriteRequired
	}

	return *f, nil
}

func validateReadLater(r *bool) (bool, error) {
	if r == nil {
		return false, ErrValidationReadLaterRequired
	}

	return *r, nil
}

// isPrivateHost checks whether a host (with optional port) resolves to a
// loopback, private (RFC 1918), link-local, or IPv6 unique local address.
// These are blocked to prevent SSRF if the backend ever fetches stored URLs
// (e.g. link previews, favicons) and to avoid unsafe links in shared collections.
func isPrivateHost(host string) bool {
	h := host
	if strings.Contains(host, ":") {
		h, _, _ = net.SplitHostPort(host)
	}

	if strings.EqualFold(h, "localhost") {
		return true
	}

	ip := net.ParseIP(h)
	if ip == nil {
		return false
	}

	privateRanges := []string{
		"127.0.0.0/8",    // loopback
		"10.0.0.0/8",     // RFC 1918 private
		"172.16.0.0/12",  // RFC 1918 private
		"192.168.0.0/16", // RFC 1918 private
		"169.254.0.0/16", // link-local (AWS/GCP/Azure metadata endpoint)
		"::1/128",        // IPv6 loopback
		"fc00::/7",       // IPv6 unique local
	}

	for _, cidr := range privateRanges {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return true
		}
	}

	return false
}
