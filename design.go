package depict

import (
	"fmt"
	"reflect"
	"strings"
)

// Design represents a depiction configuration.
type Design struct {
	maxDepth   int
	showString bool
	onlyError  bool
	hideError  bool
	ignoreJson bool

	inclusion      bool
	inclusionPaths map[string]bool
}

// Portray will, when given an interface, return a structure with private fields exported.
func (d *Design) Portray(a interface{}) *Depiction {
	if done, ok := a.(*Depiction); ok {
		// No double depiction!
		return done
	}

	ret, included := d.doInterface(a, "", d.inclusion, 1)
	if !included {
		ret = nil
	}
	if v, ok := a.(error); ok {
		if d.onlyError {
			ret = v.Error()
		} else if m, ok := ret.(map[string]interface{}); ok && !d.hideError {
			m["(error)"] = Extent(v.Error())
		}
	}
	if v, ok := a.(fmt.Stringer); ok && d.showString {
		if m, ok := ret.(map[string]interface{}); ok {
			m["(string)"] = Extent(v.String())
		}
	}

	return &Depiction{
		iFace: ret,
	}
}

func (d *Design) doInterface(a interface{}, path string, inclusion bool, depth int) (ret interface{}, included bool) {
	if depth > d.maxDepth {
		return Extent("..."), inclusion
	}
	depth++

	val := reflect.ValueOf(a)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	vk := val.Kind()
	switch vk {
	case reflect.Struct:
		return d.doStruct(val, path, inclusion, depth)

	case reflect.Map:
		return d.doMap(val, path, inclusion, depth)

	case reflect.Slice:
		if val.IsNil() {
			return nil, inclusion
		}
		fallthrough
	case reflect.Array:
		return d.doArray(val, path, inclusion, depth)
	}
	return d.doStatic(vk, val, path, inclusion, depth)
}

func (d *Design) doStruct(val reflect.Value, path string, inclusion bool, depth int) (ret interface{}, included bool) {
	if depth > d.maxDepth {
		return Extent("{...}"), inclusion
	}
	depth++

	s := map[string]interface{}{}
	ret = s

	if path != "" {
		path += "."
	}
	included = inclusion

	val = reflect.Indirect(val)
	vt := val.Type()
	n := vt.NumField()
	for i := 0; i < n; i++ {
		ft := vt.Field(i)
		f := val.Field(i)
		name := ft.Name

		if tag, ok := ft.Tag.Lookup("depict"); ok {
			name = tag
		} else if tag, ok := ft.Tag.Lookup("json"); ok && !d.ignoreJson {
			name = strings.Split(tag, ",")[0]
		}

		if name == "-" {
			continue
		}

		subPath := path + name
		subInclusion := inclusion
		if include, ok := d.inclusionPaths[subPath]; ok {
			subInclusion = include
		}

		if f.CanInterface() {
			r, inc := d.doInterface(f.Interface(), subPath, subInclusion, depth)
			if inc {
				s[name] = r
				included = true
			}
		} else {
			r, inc := d.doStatic(f.Kind(), f, subPath, subInclusion, depth)
			if inc {
				s[name] = r
				included = true
			}
		}
	}
	return
}

func (d *Design) doMap(val reflect.Value, path string, inclusion bool, depth int) (ret interface{}, included bool) {
	if depth > d.maxDepth {
		return Extent("{...}"), inclusion
	}
	depth++

	s := map[string]interface{}{}
	ret = s

	if path != "" {
		path += "."
	}
	included = inclusion

	val = reflect.Indirect(val)

	keys := val.MapKeys()
	for _, key := range keys {
		f := val.MapIndex(key)

		name := key.String()
		subPath := path + name
		subInclusion := inclusion
		if include, ok := d.inclusionPaths[subPath]; ok {
			subInclusion = include
		}

		if f.CanInterface() {
			r, inc := d.doInterface(f.Interface(), subPath, subInclusion, depth)
			if inc {
				s[name] = r
				included = true
			}
		} else {
			r, inc := d.doStatic(f.Kind(), f, subPath, subInclusion, depth)
			if inc {
				s[name] = r
				included = true
			}
		}
	}
	return
}

func (d *Design) doArray(val reflect.Value, path string, inclusion bool, depth int) (ret interface{}, included bool) {
	if depth > d.maxDepth {
		return Extent("[...]"), inclusion
	}
	depth++

	included = inclusion

	n := val.Len()
	var a []interface{}

	for i := 0; i < n; i++ {
		f := val.Index(i)

		subPath := fmt.Sprintf("%s[%d]", path, i)
		subInclusion := inclusion
		if include, ok := d.inclusionPaths[subPath]; ok {
			subInclusion = include
		}

		if f.CanInterface() {
			r, inc := d.doInterface(f.Interface(), subPath, subInclusion, depth)
			if inc {
				a = append(a, r)
				included = true
			}
		} else {
			r, inc := d.doStatic(f.Kind(), f, subPath, subInclusion, depth)
			if inc {
				a = append(a, r)
				included = true
			}
		}
	}

	ret = a
	return
}

func (d *Design) doStatic(vk reflect.Kind, val reflect.Value, path string, inclusion bool, depth int) (ret interface{}, included bool) {
	if depth > d.maxDepth {
		return Extent("..."), inclusion
	}

	switch vk {
	case reflect.Bool:
		if val.Bool() {
			return true, inclusion
		}
		return false, inclusion

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return val.Int(), inclusion

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return val.Uint(), inclusion

	case reflect.Float32, reflect.Float64:
		return val.Float(), inclusion

	case reflect.Complex64, reflect.Complex128:
		return val.Complex(), inclusion

	case reflect.String:
		return val.String(), inclusion

	case reflect.Invalid:
		return Extent("#invalid: " + val.String()), inclusion

	case reflect.Ptr, reflect.Uintptr, reflect.UnsafePointer, reflect.Chan, reflect.Func:
		return Extent("#pointer: " + val.String()), inclusion
	}

	if val.IsNil() {
		return nil, inclusion
	}
	if val.CanInterface() {
		return val.Interface(), inclusion
	}
	return Extent("#other: " + val.String()), inclusion
}
