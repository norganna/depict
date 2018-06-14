package depict

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func testStruct() interface{} {
	return map[string]interface{}{
		"a": map[string]interface{}{
			"b": []interface{}{
				map[string]interface{}{
					"c": 1,
					"d": 2,
				},
				map[string]interface{}{
					"c": 2,
					"d": 3,
				},
			},
			"c": []interface{}{
				map[string]interface{}{
					"c": 3,
					"d": 4,
				},
				map[string]interface{}{
					"c": 4,
					"d": 5,
				},
			},
		},
	}
}

func dump(x interface{}) string {
	data, _ := json.Marshal(x)
	return string(data)
}

func TestDouble_1(t *testing.T) {
	first := Portray(testStruct())
	second := Portray(first)

	if first != second {
		t.Error("Expected second to be copy of first")
	}
}

func TestFormat_1(t *testing.T) {
	depiction := Portray(struct{ a int }{a: 1})

	got := fmt.Sprintf("%#v", depiction)
	expect := fmt.Sprintf("%#v", depiction.Interface())

	if got != expect {
		t.Errorf("formatted text not correct, expected “%s”, got “%s”",
			expect,
			got,
		)
	}

	got = fmt.Sprintf("%+v", depiction)
	expect = fmt.Sprintf("%+v", depiction.Interface())

	if got != expect {
		t.Errorf("formatted text not correct, expected “%s”, got “%s”",
			expect,
			got,
		)
	}

	got = fmt.Sprintf("%-#10v", depiction)
	expect = fmt.Sprintf("%-#10v", depiction.Interface())

	if got != expect {
		t.Errorf("formatted text not correct, expected “%s”, got “%s”",
			expect,
			got,
		)
	}

	got = fmt.Sprintf("%+v", depiction)
	expect = fmt.Sprintf("%+v", depiction.Interface())

	if got != expect {
		t.Errorf("formatted text not correct, expected “%s”, got “%s”",
			expect,
			got,
		)
	}

	got = fmt.Sprintf("%v", depiction)
	expect = fmt.Sprintf("%v", depiction.Interface())

	if got != expect {
		t.Errorf("formatted text not correct, expected “%s”, got “%s”",
			expect,
			got,
		)
	}

	got = fmt.Sprintf("%s", depiction)
	expect = fmt.Sprintf("%s", depiction.Interface())

	if got != expect {
		t.Errorf("formatted text not correct, expected “%s”, got “%s”",
			expect,
			got,
		)
	}
}

func TestFormat_2(t *testing.T) {
	depiction := Portray(197.00234)

	got := fmt.Sprintf("%-+10.03f", depiction)
	expect := fmt.Sprintf("%-+10.03f", depiction.Interface())

	if got != expect {
		t.Errorf("formatted text not correct, expected “%s”, got “%s”",
			expect,
			got,
		)
	}
}

func TestInclusion_1(t *testing.T) {
	got := dump(Portray(
		testStruct(),
		Include("a.b[1].d"),
	))
	expect := `{"a":{"b":[{"d":3}]}}`

	if got != expect {
		t.Errorf("output unexpected, expected %s, got %s", expect, got)
	}
}

func TestInclusion_2(t *testing.T) {
	got := dump(Portray(
		testStruct(),
		Include("a.b"),
		Exclude("a.b[1].d"),
	))
	expect := `{"a":{"b":[{"c":1,"d":2},{"c":2}]}}`

	if got != expect {
		t.Errorf("output unexpected, expected %s, got %s", expect, got)
	}
}

func TestInclusion_3(t *testing.T) {
	got := dump(Portray(
		testStruct(),
		Exclude("a.b[1].d"),
	))
	expect := `{"a":{"b":[{"c":1,"d":2},{"c":2}],"c":[{"c":3,"d":4},{"c":4,"d":5}]}}`

	if got != expect {
		t.Errorf("output unexpected, expected %s, got %s", expect, got)
	}
}

func TestInclusion_4(t *testing.T) {
	got := dump(Portray(
		testStruct(),
		Include("a.b[1].d"),
		Exclude("a.b"),
	))
	expect := `{"a":{"b":[{"d":3}]}}`

	if got != expect {
		t.Errorf("output unexpected, expected %s, got %s", expect, got)
	}
}

func TestInclusion_5(t *testing.T) {
	got := dump(Portray(
		testStruct(),
		Include("a.b"),
		Exclude("a.b[1].d"),
	))
	expect := `{"a":{"b":[{"c":1,"d":2},{"c":2}]}}`

	if got != expect {
		t.Errorf("output unexpected, expected %s, got %s", expect, got)
	}
}

func TestNilInterface_1(t *testing.T) {
	got := dump(Portray((*struct{ a int })(nil)))
	expect := `null`

	if got != expect {
		t.Errorf("output unexpected, expected %s, got %s", expect, got)
	}
}

type B struct {
	c int
}
type A struct {
	b *B
}

func TestNilInterface_2(t *testing.T) {
	a := A{
		b: nil,
	}

	got := dump(Portray(a))
	expect := `{"b":null}`

	if got != expect {
		t.Errorf("output unexpected, expected %s, got %s", expect, got)
	}
}

func TestTime_1(t *testing.T) {
	now := time.Now()
	nowText := now.Format("2006-01-02T15:04:05.000Z07:00")

	a := &struct {
		t time.Time
	}{
		t: now,
	}

	got := dump(Portray(a))
	expect := `{"t":"` + nowText + `"}`

	if got != expect {
		t.Errorf("output unexpected, expected %s, got %s", expect, got)
	}
}
