package buckwalter

import (
	"testing"
)

var errFmt = "expected '%s', got '%s'"

func Test1(t *testing.T) {

	inp := "Allh"
	exp := "\u0627\u0644\u0644\u0647"
	got, err := Bw2Ar(inp)
	if err != nil {
		t.Errorf("didn't expect error here! got %v", err)
		// no return | keep checking output
	}

	if got != exp {
		t.Errorf(errFmt, exp, got)
		return
	}

	inp2 := exp
	exp2 := inp
	got2, err := Ar2Bw(inp2)
	if err != nil {
		t.Errorf("didn't expect error here! got %v", err)
	}

	if got2 != exp2 {
		t.Errorf(errFmt, exp2, got2)
	}
}

func TestNorm1(t *testing.T) {

	var inp, exp, got, inp2, exp2, got2 string
	var err error

	// shadda ~ \u0651
	// damma  u \u064F

	// AR IN shadda-damma
	// BW UT shadda-damma
	// AR UT damma-shadda
	inp = "\u062D\u064F\u0645\u0651\u064F\u0635"
	exp = "Hum~uS"

	got, err = Ar2Bw(inp)
	if err == nil {
		t.Errorf("expected error here!")
		// no return | keep checking output
	}

	if got != exp {
		t.Errorf(errFmt, exp, got)
		return
	}

	inp2 = exp
	//exp2 = inp
	exp2 = "\u062D\u064F\u0645\u064F\u0651\u0635"

	got2, err = Bw2Ar(inp2)
	if err != nil {
		t.Errorf("didn't expect error here! got %v", err)
	}

	if got2 != exp2 {
		t.Errorf(errFmt, exp2, got2)
	}

}

func TestNorm2(t *testing.T) {

	var inp, exp, got, inp2, exp2, got2 string
	var err error

	// shadda ~ \u0651
	// damma  u \u064F

	// AR IN damma-shadda
	// BW UT shadda-damma
	// AR UT damma-shadda
	inp = "\u062D\u064F\u0645\u064F\u0651\u0635"
	exp = "Hum~uS"

	got, err = Ar2Bw(inp)
	if err != nil {
		t.Errorf("didn't expect error here! got %v", err)
		// no return | keep checking output
	}

	if got != exp {
		t.Errorf(errFmt, exp, got)
		return
	}

	inp2 = exp
	exp2 = inp
	got2, err = Bw2Ar(inp2)
	if err != nil {
		t.Errorf("didn't expect error here! got %v", err)
	}

	if got2 != exp2 {
		t.Errorf(errFmt, exp2, got2)
	}
}

func Test2(t *testing.T) {

	var inp, exp, got, inp2, exp2, got2 string
	var err error

	inp = "\u062D\u064F\u0645\u064F\u0651 \u0635"
	exp = "Hum~u S"

	got, err = Ar2Bw(inp)
	if err != nil {
		t.Errorf("didn't expect error here! got %v", err)
		// no return | keep checking output
	}

	if got != exp {
		t.Errorf(errFmt, exp, got)
		return
	}

	inp2 = exp
	exp2 = inp
	got2, err = Bw2Ar(inp2)
	if err != nil {
		t.Errorf("didn't expect error here! got %v", err)
	}

	if got2 != exp2 {
		t.Errorf(errFmt, exp2, got2)
	}
}

func Test3(t *testing.T) {

	var inp, exp, got, inp2, exp2, got2 string
	var err error

	inp = "\u062D\u0645\u0648\uFEAA"
	exp = "Hmwd"
	exp2 = "\u062D\u0645\u0648\u062F"

	got, err = Ar2Bw(inp)
	if err != nil {
		t.Errorf("didn't expect error here! got %v", err)
		// no return | keep checking output
	}

	if got != exp {
		t.Errorf(errFmt, exp, got)
		return
	}

	inp2 = exp
	//exp2 = inp
	got2, err = Bw2Ar(inp2)
	if err != nil {
		t.Errorf("didn't expect error here! got %v", err)
	}

	if got2 != exp2 {
		t.Errorf(errFmt, exp2, got2)
	}
}

func TestUpdates20180918(t *testing.T) {

	var inp, exp, got, inp2, exp2, got2 string
	var err error

	//ch{'\u067E', 'P'},
	//res = strings.Replace(res, "\u06BE", "\u0647", -1) // HEH DOACHASHMEE => HEH

	inp = "\u067E\u06BE"
	exp = "Ph"
	exp2 = "\u067E\u0647"

	got, err = Ar2Bw(inp)
	if err != nil {
		t.Errorf("didn't expect error here! got %v", err)
		// no return | keep checking output
	}

	if got != exp {
		t.Errorf(errFmt, exp, got)
		return
	}

	inp2 = exp
	//exp2 = inp
	got2, err = Bw2Ar(inp2)
	if err != nil {
		t.Errorf("didn't expect error here! got %v", err)
	}

	if got2 != exp2 {
		t.Errorf(errFmt, exp2, got2)
	}
}
