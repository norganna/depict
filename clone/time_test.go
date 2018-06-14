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

func TestTime_1(t *testing.T) {
	now := time.Now().UTC()
	val := reflect.ValueOf(now)

	carbon := Time(val)
	got := carbon.Format(time.RFC3339Nano)
	expect := now.Format(time.RFC3339Nano)

	if got != expect {
		t.Errorf("Expected time format to be the same, expected “%s”, got “%s”",
			expect, got)
	}
}

func TestTime_2(t *testing.T) {
	loc, _ := time.LoadLocation("America/New_York")
	now := time.Now().In(loc)
	val := reflect.ValueOf(now)

	carbon := Time(val)
	got := carbon.Format(time.RFC3339Nano)
	expect := now.Format(time.RFC3339Nano)

	if got != expect {
		t.Errorf("Expected time format to be the same, expected “%s”, got “%s”",
			expect, got)
	}
}
