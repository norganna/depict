package clone

import (
	"reflect"
	"time"
)

const (
	internalToZero int64 = 59453308800 // Jan 1 0001 -> Jan 1 1885 (1,884 years in seconds).
	zeroToInternal       = -internalToZero
	unixToZero           = 62135596800 // Jan 1 0001 -> Jan 1 1970 (1,989 years in seconds).
	zeroToUnix           = -unixToZero
	internalToUnix       = internalToZero + zeroToUnix // Offset to convert an internal time into unix time.
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
	if wall&(1<<63) != 0 { // Monotonic = 1, wall bits > 30 contain seconds since 1885 (internal)
		secs = int64(wall << 1 >> 31) // secs since internal epoch
		secs += internalToUnix        // secs since unix epoch
	} else { // Monotonic = 0, ext contains seconds since 0001
		secs = ext         // secs since zero
		secs += zeroToUnix // secs since unix epoch
	}
	nanos := int64(int32(wall & (1<<30 - 1)))

	ts := time.Unix(secs, nanos)
	if loc, err := time.LoadLocation(locName); err == nil {
		ts = ts.In(loc)
	}
	return ts
}
