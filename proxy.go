package depict

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
)

// Proxyer implementors contain another type and can be marshalled and formatted.
type Proxyer interface {
	MarshalJSON() ([]byte, error)
	GobEncode() ([]byte, error)
	Format(state fmt.State, verb rune)
}

// EncodeToGob encodes a rando interface to a gob.
func EncodeToGob(i interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(i)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// FormatFromState builds a Sprintf compatible format string from a format state and verb.
func FormatFromState(state fmt.State, verb rune) string {
	format := "%"
	if state.Flag(' ') {
		format += " "
	}
	if state.Flag('-') {
		format += "-"
	}
	if state.Flag('+') {
		format += "+"
	}
	if state.Flag('#') {
		format += "#"
	}

	if v, ok := state.Width(); ok {
		format += strconv.Itoa(v)
	}
	if v, ok := state.Precision(); ok {
		format += "."
		if state.Flag('0') {
			format += "0"
		}
		format += strconv.Itoa(v)
	}

	format += string(verb)
	return format
}
