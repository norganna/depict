package clone

import (
	"reflect"
	"time"
)

const (
	internalToUnix int64 = 2682288000
	wallToUnix = 62135596800
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

	locName := "UTC"
	location := val.FieldByName("loc")
	location = reflect.Indirect(location)
	if location.Kind() != reflect.Invalid {
		locName = location.FieldByName("name").String()
	}

	var secs int64
	if wall&(1<<63) != 0 {
		secs = int64(wall<<1>>31) - internalToUnix
	} else {
		secs = ext - wallToUnix
	}
	nanos := int64(int32(wall & (1<<30 - 1)))

	ts := time.Unix(secs, nanos)
	if loc, err := time.LoadLocation(locName); err == nil {
		ts = ts.In(loc)
	}
	return ts
}
