package sure_test

import (
	"errors"
	"testing"

	"github.com/dbinger/sure"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBestruct_Same(t *testing.T) {
	type ex struct {
		A int
	}
	type exbad struct {
		a int
	}
	type testcase struct {
		name string
		got  any
		want any
		msg  string
	}
	var exNil *ex
	tests := []testcase{
		{"", nil, nil, ""},
		{"", 42, nil, "FAIL in tn\nnote\nGOT:  int(42)\nWANT: nil"},
		{"", nil, 42, "FAIL in tn\nnote\nGOT:  nil\nWANT: int(42)"},
		{"", 42.0, 42, "FAIL in tn\nnote\nGOT:  float64(42),\nWANT: int(42)"},
		{"", "a", "a", ""},
		{"", ex{1}, ex{2}, "FAIL in tn\nnote\nsure_test.ex{\nGOT:  A: 1,\nWANT: A: 2,\n  }"},
		{"", ex{1}, ex{1}, ""},
		{"", &ex{1}, &ex{1}, ""},
		{"", ex{1}, &ex{1}, "FAIL in tn\nnote\nGOT:  sure_test.ex{A: 1},\nWANT: &sure_test.ex{A: 1}"},
		{"", exNil, (*ex)(nil), ""},
		{"", (*ex)(nil), exNil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			faket := &testing.T{}
			b := &sure.BeStruct{T: faket, CmpOptions: nil, Name: "tn",
				FatalFunc: func(a ...any) { faket.Fail() },
			}
			if gotmsg := b.Same(tt.got, tt.want, "note"); gotmsg != tt.msg {
				t.Errorf("got %#v, want %#v", gotmsg, tt.msg)
			}
			if (len(tt.msg) == 0) && b.Failed() {
				t.Errorf("reported failed, expected passed")
			}
			if (len(tt.msg) != 0) && !b.Failed() {
				t.Errorf("reported passed, expected failed")
			}
		})
	}
	faket := &testing.T{}
	b := &sure.BeStruct{T: faket, CmpOptions: nil, Name: "tn",
		FatalFunc: func(a ...any) { faket.Fail() },
	}
	if b.Same(exbad{1}, exbad{1}, "n1", "n2"); !b.Failed() {
		t.Errorf("passed test with comparison panic, want fail")
	}
}

func TestBeStruct_Diff(t *testing.T) {
	type ex struct {
		A int
	}
	type exbad struct {
		a int
	}
	type testcase struct {
		name string
		got  any
		want any
		msg  string
	}
	var exNil ex
	tests := []testcase{
		{"1", nil, nil, "FAIL in tn\nnote\nGOT:  nil\nWANT: anything else"},
		{"2", 42, nil, ""},
		{"3", nil, 42, ""},
		{"4", 42.0, 42, ""},
		{"5", 42.0, 42, ""},
		{"6", "a", "a", "FAIL in tn\nnote\nGOT:  string(a)\nWANT: anything else"},
		{"7", ex{1}, ex{2}, ""},
		{"8", ex{1}, ex{1}, "FAIL in tn\nnote\nGOT:  sure_test.ex({1})\nWANT: anything else"},
		{"9", &ex{1}, &ex{1}, "FAIL in tn\nnote\nGOT:  *sure_test.ex(&{1})\nWANT: anything else"},
		{"10", ex{1}, &ex{1}, ""},
		{"11", exNil, nil, ""},
		{"12", nil, exNil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			faket := &testing.T{}
			b := &sure.BeStruct{T: faket, CmpOptions: nil, Name: "tn",
				FatalFunc: func(a ...any) { faket.Fail() }}
			if gotmsg := b.Diff(tt.got, tt.want, "note"); gotmsg != tt.msg {
				t.Errorf("got %#v, want %#v", gotmsg, tt.msg)
			}
			if (len(tt.msg) == 0) && b.Failed() {
				t.Errorf("reported failed, expected passed")
			}
			if (len(tt.msg) != 0) && !b.Failed() {
				t.Errorf("reported passed, expected failed")
			}
		})
	}
	faket := &testing.T{}
	b := &sure.BeStruct{T: faket, CmpOptions: nil, Name: "tn",
		FatalFunc: func(a ...any) { faket.Fail() }}
	if b.Diff(exbad{1}, exbad{1}, "note1", "note2"); !b.Failed() {
		t.Errorf("passed test with comparison panic, want fail")
	}
}

func TestBe(t *testing.T) {
	type Ex struct {
		a int
	}
	option := cmpopts.IgnoreUnexported(Ex{})
	b := sure.Be(t, option)
	b.Same(Ex{}, Ex{})
}

func TestErrors(t *testing.T) {
	b := sure.Be(t)
	e1 := errors.New("1")
	e2 := errors.New("2")
	b.Diff(e1, e2)
	e3 := errors.Join(e1)
	b.Same(e1, e3)
	b.Same(e3, e1)
	b.Diff(e2, e3)
	b.Diff(e3, e2)
	b.Same(e1, sure.AnyError)
	b.Same(e2, sure.AnyError)
	b.Same(e3, sure.AnyError)
	b.Diff(nil, sure.AnyError)
	b.Diff(3, sure.AnyError)
}
