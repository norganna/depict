package depict

import (
	"encoding/json"
	"fmt"
)

// Depiction type denotes a depiction of an original structure.
type Depiction struct {
	iFace interface{}
}

var _ Proxyer = (*Depiction)(nil)

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
	return EncodeToGob(d.iFace)
}

// Format allows us to be fmt.Printed.
func (d *Depiction) Format(state fmt.State, verb rune) {
	fmt.Fprintf(state, FormatFromState(state, verb), d.iFace)
}
