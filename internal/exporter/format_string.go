package exporter

import "fmt"

// String returns the string representation of a Format.
func (f Format) String() string {
	return string(f)
}

// ParseFormat converts a raw string into a Format, returning an error
// if the value is not recognised.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatJSON, FormatCSV, FormatText:
		return Format(s), nil
	default:
		return "", fmt.Errorf("exporter: unknown format %q (choose json, csv, text)", s)
	}
}

// Formats returns all supported Format values.
func Formats() []Format {
	return []Format{FormatJSON, FormatCSV, FormatText}
}
