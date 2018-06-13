package clone

import (
	"reflect"
	"testing"
	"time"
)

func TestIsTime(t *testing.T) {
	now := time.Now()
	if !IsTime(reflect.TypeOf(now)) {
		t.Error("Expected that a time isTime")
	}
}

func TestTime(t *testing.T) {
	now := time.Now()
	val := reflect.ValueOf(now)

	carbon := Time(val)
	got := carbon.Format(time.RFC3339Nano)
	expect := now.Format(time.RFC3339Nano)

	if got != expect {
		t.Errorf("Expected time format to be the same, expected “%s”, got “%s”",
			expect, got)
	}
}
