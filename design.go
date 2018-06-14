package depict

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/norganna/depict/clone"
)

// Design represents a depiction configuration.
type Design struct {
	maxDepth    int
	showString  bool
	onlyError   bool
	hideError   bool
	ignoreJson  bool
	ignoreKnown bool

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
			m["(error)"] = &Extent{v.Error()}
		}
	}
	if v, ok := a.(fmt.Stringer); ok && d.showString {
		if m, ok := ret.(map[string]interface{}); ok {
			m["(string)"] = &Extent{v.String()}
		}
	}

	return &Depiction{
		iFace: ret,
	}
}

func (d *Design) doInterface(a interface{}, path string, inclusion bool, depth int) (ret interface{}, included bool) {
	if depth > d.maxDepth {
		return &Extent{"..."}, inclusion
	}
	depth++

	val := reflect.ValueOf(a)
	vk := val.Kind()
	for vk == reflect.Ptr {
		val = reflect.Indirect(val)
		vk = val.Kind()
	}

	return d.doChoose(vk, val, path, inclusion, depth)
}

func (d *Design) doStruct(val reflect.Value, path string, inclusion bool, depth int) (ret interface{}, included bool) {
	if depth > d.maxDepth {
		return &Extent{"{...}"}, inclusion
	}
	depth++

	vk := val.Kind()
	for vk == reflect.Ptr {
		val = reflect.Indirect(val)
		vk = val.Kind()
	}

	if !d.ignoreKnown && vk != reflect.Invalid {
		t := val.Type()
		switch {
		case clone.IsTime(t):
			ts := clone.Time(val)
			return ts.Format("2006-01-02T15:04:05.000-0700"), inclusion
		}

	}

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

		fk := f.Kind()
		for fk == reflect.Ptr {
			f = reflect.Indirect(f)
			fk = f.Kind()
		}

		if fk == reflect.Invalid {
			if subInclusion {
				s[name] = nil
				included = true
			}
		} else {
			r, inc := d.doChoose(fk, f, subPath, subInclusion, depth)
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
		return &Extent{"{...}"}, inclusion
	}
	depth++

	vk := val.Kind()
	for vk == reflect.Ptr {
		val = reflect.Indirect(val)
		vk = val.Kind()
	}

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

		fk := f.Kind()
		for fk == reflect.Ptr {
			f = reflect.Indirect(f)
			fk = f.Kind()
		}

		if fk == reflect.Invalid {
			if subInclusion {
				s[name] = nil
				included = true
			}
		} else {
			r, inc := d.doChoose(fk, f, subPath, subInclusion, depth)
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
		return &Extent{"[...]"}, inclusion
	}
	depth++

	vk := val.Kind()
	for vk == reflect.Ptr {
		val = reflect.Indirect(val)
		vk = val.Kind()
	}

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

		fk := f.Kind()
		for fk == reflect.Ptr {
			f = reflect.Indirect(f)
			fk = f.Kind()
		}

		if fk == reflect.Invalid {
			if subInclusion {
				a = append(a, nil)
				included = true
			}
		} else {
			r, inc := d.doChoose(fk, f, subPath, subInclusion, depth)
			if inc {
				a = append(a, r)
				included = true
			}
		}
	}

	ret = a
	return
}

func (d *Design) doChoose(vk reflect.Kind, val reflect.Value, path string, inclusion bool, depth int) (ret interface{}, included bool) {
	if depth > d.maxDepth {
		return &Extent{"..."}, inclusion
	}

	for vk == reflect.Ptr {
		val = reflect.Indirect(val)
		vk = val.Kind()
	}

	if vk == reflect.Interface {
		return d.doInterface(val.Interface(), path, inclusion, depth)
	}

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
		return nil, inclusion
	}

	vs := val.String()
	if len(vs) > 2 && vs[0] == '<' {
		vs = vs[1 : len(vs)-1]
	}
	return &Extent{"#" + vs}, inclusion
}
