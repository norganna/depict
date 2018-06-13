# Depict

Builds a depiction of a variable (any `interface{}`), using reflection to customisably export both public and private
fields within.

Although generally this package was designed to allow JSON encoding of private fields, it may be gainfully employed for
other tasks as well.

## Example:

```go
package main

import (
	"encoding/json"

	"github.com/norganna/depict"
)

// Foo is our test structure.
type Foo struct {
	ID          int    `json:"id,omitempty"`
	secret      string `depict:"value"`
	superSecret string `depict:"-"`
}

func main() {
	foo := &Foo{
		ID:     1,
		secret: "test",
	}

	data, _ := json.Marshal(foo)
	println(string(data))
	// Will output: {"id":1}

	data, _ = json.Marshal(depict.Portray(foo))
	println(string(data))
	// Will output: {"id":1,"value":"test"}
}
```
