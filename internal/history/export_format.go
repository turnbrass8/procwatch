package history

import (
	"fmt"
	"strings"
)

// Format represents an output format for history export.
type Format string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

// ParseFormat parses a string into a Format, returning an error for unknown values.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "json":
		return FormatJSON, nil
	case "text", "":
		return FormatText, nil
	default:
		return "", fmt.Errorf("history: unknown export format %q (supported: json, text)", s)
	}
}

// String implements the Stringer interface.
func (f Format) String() string {
	return string(f)
}

// IsValid reports whether the format is a recognised value.
func (f Format) IsValid() bool {
	switch f {
	case FormatJSON, FormatText:
		return true
	}
	return false
}
