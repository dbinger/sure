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
	return &BeStruct{T: t, CmpOptions: options, Name: t.Name(),
		FatalFunc: t.Fatal,
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
	if want == nil {
		return fmt.Sprintf("got %#v, wanted nil", got)
	} else if got == nil {
		return fmt.Sprintf("got nil, wanted %#v", want)
	} else {
		return b.diff(want, got)
	}
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
	if dontwant == nil {
		return "got nil, wanted non-nil"
	} else {
		return fmt.Sprintf("got %#v, wanted anything else", got)
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
	// Replace nbsp characters that cmp.Diff uses.
	diff = strings.ReplaceAll(diff, "\u00a0", " ")
	// Trim "any(" ")" wrapper to shorten output.
	if strings.HasPrefix(diff, "any(\n") {
		diff = diff[5:strings.LastIndex(diff, "\n")]
	}
	diff = "mismatch -got +want\n" + diff
	return diff
}
