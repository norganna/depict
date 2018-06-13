package depict

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strconv"
)

// Depiction type denotes a depiction of an original structure.
type Depiction struct {
	iFace interface{}
}

// Interface returns the depicted interface.
func (d *Depiction) Interface() interface{} {
	return d.iFace
}

// MarshalJSON allows us to be JSON marshaled.
func (d *Depiction) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.iFace)
}

// GobEncode allows us to be gob encoded.
func (d *Depiction) GobEncode() ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(d.iFace)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Format allows us to be fmt.Printed.
func (d *Depiction) Format(state fmt.State, verb rune) {
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

	fmt.Fprintf(state, format, d.iFace)
}
