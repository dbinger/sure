package sure

// sure provides a struct that wraps *testing.T and provides only 2 assertion methods.
// Objectives:
// 1. Make go tests a bit easier to read.
// 2. Use minimal set of assert methods.
// 3. Reduce or eliminate the need for putting format strings in test code.
// 4. Make error messages easier to interpret.

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// AnyError is a value that matches any non-nil error.
var AnyError = cmpopts.AnyError

type BeStruct struct {
	*testing.T
	CmpOptions []cmp.Option
	Name       string
	FatalFunc  func(...any)
}

// Be returns a testing struct with two assertion methods: Same and Diff.
// These methods rely on the "cmp" package for comparison,
// including cmpopts.EquateErrors() by default.
func Be(t *testing.T, options ...cmp.Option) *BeStruct {
	options = append(options, cmpopts.EquateErrors())
	return &BeStruct{
		T:          t,
		CmpOptions: options,
		Name:       t.Name(),
		FatalFunc:  t.Fatal,
	}
}

// Same calls Fail if got does not equal want.
// The error message is returned.
func (b *BeStruct) Same(got, want any, notes ...any) string {
	b.Helper()
	return b.failIfDiff(b.eqString(got, want), notes...)
}

// Diff calls Fail if got equals dontwant.
// The error message is returned.
func (b *BeStruct) Diff(got, dontwant any, notes ...any) string {
	b.Helper()
	return b.failIfDiff(b.notEqString(got, dontwant), notes...)
}

// failIfDiff records a testing failure if diff is a non-empty string.
// The diff is combined with the note to produce the log message.
// If now is true, the failure is fatal, so the test exits immediately.
// failIfDiff also leaves the log message on the BeStruct.
// The error message is returned.
func (b *BeStruct) failIfDiff(diff string, notes ...any) string {
	b.Helper()
	if diff == "" {
		return ""
	}
	s := "FAIL in " + b.Name + "\n"
	if len(notes) > 0 {
		s += fmt.Sprintln(notes...)
	}
	s += diff
	b.FatalFunc(s)
	return s
}

// eqString returns a non-empty explanation if got is not the same as want,
// or if the comparison failed.
func (b *BeStruct) eqString(got, want any) string {
	equal, err := b.compare(got, want)
	if err != nil {
		return "error: " + err.Error()
	}
	if equal {
		return ""
	}
	s := b.diff(want, got)
	if got == nil && !strings.Contains(s, "GOT:  ") {
		s = "GOT:  nil\n" + s
	}
	if want == nil && !strings.Contains(s, "WANT: ") {
		s = s + "\nWANT: nil"
	}
	return s
}

// notEqString returns a non-empty explanation if got is the same as dontwant,
// or if the comparison failed.
func (b *BeStruct) notEqString(got, dontwant any) string {
	equal, err := b.compare(got, dontwant)
	if err != nil {
		return "error: " + err.Error()
	}
	if !equal {
		return ""
	}
	if got == nil {
		return "GOT:  nil\nWANT: anything else"
	} else {
		s := b.eqString(got, nil)
		s = strings.Replace(s, "WANT: nil", "WANT: anything else", 1)
		return s
	}
}

// compare is like cmp.Equal, except that it returns an error instead of
// calling panic when the comparison can't be completed.
func (b *BeStruct) compare(got, want any) (bool, error) {
	var equal bool
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("compare: %v", r)
			}
		}()
		equal = cmp.Equal(got, want, b.CmpOptions...)
	}()
	return equal, err
}

// diff is like cmp.Diff, except with some edits on the output string.
func (b *BeStruct) diff(got, want any) string {
	diff := cmp.Diff(want, got, b.CmpOptions...)
	diff = strings.TrimSpace(diff)
	diff = strings.ReplaceAll(diff, "\u00a0", " ")
	diff = strings.ReplaceAll(diff, "\n- \t", "\nGOT:  ")
	diff = strings.ReplaceAll(diff, "\n+ \t", "\nWANT: ")
	i := strings.Index(diff, "(\n")
	if i != -1 {
		diff = diff[i+2 : strings.LastIndex(diff, "\n")]
		diff = strings.Replace(diff, ",\n", "\n", 1)
	}
	diff = strings.ReplaceAll(diff, "\t", "    ")
	diff = strings.Trim(diff, ",\n")
	return diff
}
