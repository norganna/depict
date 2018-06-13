package depict

// Opt is a design configuration function.
type Opt func(*Design)

// MaxDepth sets the maximum recursion depth for the depiction.
func MaxDepth(depth int) Opt {
	return func(d *Design) {
		d.maxDepth = depth
	}
}

// ShowString determines if the value of the `String()` function (if any) should be included.
func ShowString() Opt {
	return func(d *Design) {
		d.showString = true
	}
}

// OnlyError hides all other fields if the value is an error and returns just the `Error()` string.
func OnlyError() Opt {
	return func(d *Design) {
		d.onlyError = true
	}
}

// HideError determines if the value of the `Error()` function (if any) should be excluded.
func HideError() Opt {
	return func(d *Design) {
		d.hideError = true
	}
}

// IgnoreJsonTag causes names in `json` tags on struct fields to be ignored.
func IgnoreJsonTag() Opt {
	return func(d *Design) {
		d.ignoreJson = true
	}
}

// Include restricts fields to only those specified.
// A path looks like "foo[1].bar", which would match only the `bar` field in the second item in th `foo` field of your item.
// E.g. the following would include the bar = "two" item:
//    struct {
//        foo: []interface{}{
//            struct { bar: "one" },
//            struct { bar: "two" },
//        },
//    }
func Include(paths ...string) Opt {
	return func(d *Design) {
		for _, path := range paths {
			if _, ok := d.inclusionPaths[path]; !ok {
				d.inclusion = false
				d.inclusionPaths[path] = true
			}
		}
	}
}

// Exclude prevents fields specified from being included.
// Excludes take priority over includes, but if something isn't included then it won't be able to be excluded.
// See `Include` for details on path matching.
func Exclude(paths ...string) Opt {
	return func(d *Design) {
		for _, path := range paths {
			d.inclusionPaths[path] = false
		}
	}
}
