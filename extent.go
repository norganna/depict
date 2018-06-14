package depict

import (
	"encoding/json"
	"fmt"
)

// Extent type denotes a string value that is returned which isn't really a string in the original structure.
type Extent struct {
	Text string
}

var _ Proxyer = (*Extent)(nil)

// MarshalJSON allows us to be JSON marshaled.
func (e *Extent) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Text)
}

// MarshalBinary allows us to be binary marshaled.
func (e *Extent) MarshalBinary() ([]byte, error) {
	return e.MarshalText()
}

// MarshalText allows us to be text marshaled.
func (e *Extent) MarshalText() ([]byte, error) {
	return []byte(e.Text), nil
}

// GobEncode allows us to be gob encoded.
func (e *Extent) GobEncode() ([]byte, error) {
	return EncodeToGob(e.Text)
}

// Format allows us to be fmt.Printed.
func (e *Extent) Format(state fmt.State, verb rune) {
	fmt.Fprintf(state, FormatFromState(state, verb), e.Text)
}

// String allows us to be converted to a string
func (e *Extent) String() string {
	return e.Text
}
