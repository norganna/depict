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
