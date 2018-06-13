package clone

import (
	"reflect"
	"time"
)

var timeType = reflect.TypeOf(time.Time{})

// IsTime checks if the supplied reflection type is a `time.Time`.
func IsTime(t reflect.Type) bool {
	return t == timeType
}

// Time reads the internal fields of a `time.Time` reflection to create a non-reflection `Time` clone.
func Time(val reflect.Value) time.Time {
	wall := val.FieldByName("wall").Uint()
	ext := val.FieldByName("ext").Int()
	location := val.FieldByName("loc")
	locName := reflect.Indirect(location).FieldByName("name").String()

	var secs int64
	if wall&(1<<63) != 0 {
		secs = int64(wall<<1>>31) - 2682288000
	} else {
		secs = ext
	}
	nanos := int64(int32(wall & (1<<30 - 1)))

	ts := time.Unix(secs, nanos)
	if loc, err := time.LoadLocation(locName); err == nil {
		ts.In(loc)
	}
	return ts
}
