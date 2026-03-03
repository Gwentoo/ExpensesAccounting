package encoding

import (
	"golang.org/x/text/encoding/charmap"
	"unicode/utf8"
)

func decodeWindows1251(input string) (string, error) {
	dec := charmap.Windows1251.NewDecoder()

	result, err := dec.String(input)
	if err != nil {
		return "", err
	}

	return result, nil
}

func SafeDecode(raw string) string {
	if utf8.ValidString(raw) {
		return raw
	}

	decoded, err := decodeWindows1251(raw)
	if err != nil {
		return raw
	}
	return decoded
}
