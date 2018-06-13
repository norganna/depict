// Package depict is used to get a representation of a private structure into interfaces that can be marshalled.
package depict

const (
	// DefaultMaxDepth is to prevent deep recursion.
	DefaultMaxDepth = 10
)

// Extent type denotes a string value that is returned which isn't really a string in the original structure.
type Extent string

var defaultDesign = New()

// New returns a new `Design` based on the configured options, from which you can call `Portray` to create a
// personalised `Depiction`.
func New(opts ...Opt) *Design {
	d := &Design{
		maxDepth:       DefaultMaxDepth,
		inclusion:      true,
		inclusionPaths: map[string]bool{},
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

// Portray will, when given an interface, return a structure with private fields exported.
// You may pass options to customise the `Design`. This is equivalent to creating a new Design (via `depict.New()`)
// and then calling `Portray` on it, but less efficient.
func Portray(a interface{}, opts ...Opt) *Depiction {
	d := defaultDesign
	if len(opts) > 0 {
		d = New(opts...)
	}
	return d.Portray(a)
}
