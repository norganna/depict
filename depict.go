// Package depict is used to get a representation of a private structure into interfaces that can be marshalled.
package depict

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	MaxDepth = 10
)

// Interpret will, when given an interface, return a structure with private fields exported.
func Interpret(a interface{}) interface{} {
	fmt.Printf("Interpret: %#v", a)
	return interpretInterface(a, 1)
}

func interpretInterface(a interface{}, depth int) interface{} {
	fmt.Println(strings.Repeat("  ", depth), "iface")
	if depth > MaxDepth {
		return "..."
	}
	depth++

	val := reflect.ValueOf(a)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	vk := val.Kind()
	switch vk {
	case reflect.Struct:
		return interpretStruct(val, depth)

	case reflect.Slice:
		if val.IsNil() {
			return nil
		}
		fallthrough
	case reflect.Array:
		return interpretArray(val, depth)
	}
	return interpretStatic(vk, val, depth)
}

func interpretStruct(val reflect.Value, depth int) (ret interface{}) {
	fmt.Println(strings.Repeat("  ", depth), "struct")
	if depth > MaxDepth {
		return "{...}"
	}
	depth++

	s := map[string]interface{}{}
	ret = s

	val = reflect.Indirect(val)
	vt := val.Type()
	n := vt.NumField()
	for i := 0; i < n; i++ {
		ft := vt.Field(i)
		fmt.Println(strings.Repeat("  ", depth), "->", ft.Name)

		f := val.Field(i)
		if f.CanInterface() {
			s[ft.Name] = interpretInterface(f.Interface(), depth)
		} else {
			s[ft.Name] = interpretStatic(f.Kind(), f, depth)
		}
	}
	return

}

func interpretArray(val reflect.Value, depth int) (ret interface{}) {
	fmt.Println(strings.Repeat("  ", depth), "array")
	if depth > MaxDepth {
		return "[...]"
	}
	depth++

	n := val.Len()
	a := make([]interface{}, n)
	ret = a

	for i := 0; i < n; i++ {
		fmt.Println(strings.Repeat("  ", depth), "[]", i)
		f := val.Index(int(i))
		if f.CanInterface() {
			a[i] = interpretInterface(f.Interface(), depth)
		} else {
			a[i] = interpretStatic(f.Kind(), f, depth)
		}
	}
	return
}

func interpretStatic(vk reflect.Kind, val reflect.Value, depth int) interface{} {
	fmt.Println(strings.Repeat("  ", depth), "value")
	if depth > MaxDepth {
		return "..."
	}
	depth++

	switch vk {
	case reflect.Bool:
		if val.Bool() {
			return true
		}
		return false

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return val.Int()

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return val.Uint()

	case reflect.Float32, reflect.Float64:
		return val.Float()

	case reflect.Complex64, reflect.Complex128:
		return val.Complex()

	case reflect.String:
		return val.String()

	case reflect.Invalid:
		return "#invalid: "+val.String()

	case reflect.Ptr, reflect.Uintptr, reflect.UnsafePointer, reflect.Chan, reflect.Func:
		return "#pointer: "+val.String()
	}

	if val.IsNil() {
		return nil
	}
	if val.CanInterface() {
		return val.Interface()
	}
	return fmt.Sprintf("%v", val.String())
}
