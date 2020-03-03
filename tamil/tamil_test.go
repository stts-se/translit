package tamil

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func imports2() { // keep imports
	fmt.Fprintf(os.Stdout, "")
}

var tlit = NewTranslit()

func testConvertExpectOK(t *testing.T, s string) {
	res1, res2 := tlit.ConvertDebug(s, true), tlit.Convert(s)
	if !res1.OK {
		t.Errorf("TranslitExpected res to be ok for '%s'. Got '%s'", s, res1.Result)
	}
	if !res2.OK {
		t.Errorf("TranslitExpected res to be ok for '%s'. Got '%s'", s, res2.Result)
	}
}

func testConvertExpectFail(t *testing.T, s string) {
	res1, res2 := tlit.ConvertDebug(s, true), tlit.Convert(s)
	if res1.OK {
		t.Errorf("TranslitExpected res to be NOT ok for '%s'. Got '%s'", s, res1.Result)
	}
	if res2.OK {
		t.Errorf("TranslitExpected res to be NOT ok for '%s'. Got '%s'", s, res2.Result)
	}
}

func testConvertExpectOKWithRes(t *testing.T, s string, testConvertExpect string) {
	res1, res2 := tlit.ConvertDebug(s, true), tlit.Convert(s)
	if res1.Result != testConvertExpect {
		t.Errorf("For '%s', testConvertExpected '%s', got '%s'", s, testConvertExpect, res1.Result)
	}
	if !res1.OK {
		t.Errorf("TranslitExpected res to be ok for '%s'. Got '%s'", s, res1.Result)
	}
	if res2.Result != testConvertExpect {
		t.Errorf("For '%s', testConvertExpected '%s', got '%s'", s, testConvertExpect, res2.Result)
	}
	if !res2.OK {
		t.Errorf("TranslitExpected res to be ok for '%s'. Got '%s'", s, res2.Result)
	}
}

func testConvertExpectFailWithRes(t *testing.T, s string, testConvertExpect string) {
	res1, res2 := tlit.ConvertDebug(s, true), tlit.Convert(s)
	if res1.Result != testConvertExpect {
		t.Errorf("For '%s', testConvertExpected '%s', got '%s'", s, testConvertExpect, res1.Result)
	}
	if res1.OK {
		t.Errorf("TranslitExpected res to be NOT ok for '%s'. Got '%s'", s, res1.Result)
	}
	if res2.Result != testConvertExpect {
		t.Errorf("For '%s', testConvertExpected '%s', got '%s'", s, testConvertExpect, res2.Result)
	}
	if res2.OK {
		t.Errorf("TranslitExpected res to be NOT ok for '%s'. Got '%s'", s, res2.Result)
	}
}

func TestTranslitConvert1(t *testing.T) {
	s := "nisse pisse"
	testConvertExpectFail(t, s)
}

func TestTranslitConvert2(t *testing.T) {
	s := "nisse pissar"
	testConvertExpectFailWithRes(t, s, "????? ??????")
}

func TestTranslitConvert3a(t *testing.T) {
	testConvertExpectOKWithRes(t, "மனநலப்", "maṉanalap")
}

func TestTranslitConvert3b(t *testing.T) {
	testConvertExpectOKWithRes(t, "கேள்வி", "kēḷvi")
}

func TestTranslitConvert3c(t *testing.T) {
	s := "ஹெச்.ராஜாவை மனநலப் பரிசோதனைக்கு உட்படுத்தக் கோரிய வழக்கு! - காவல்துறையைக் கேள்வி கேட்ட நீதிமன்றம்"
	testConvertExpectOK(t, s)
}

func TestTranslitConvert4(t *testing.T) {
	s := "\u0BB6\u0BCD\u0BB0\u0BC0"
	testConvertExpectOKWithRes(t, s, "śrī")
}

func TestTranslitConvert5(t *testing.T) {
	s := "\u0BE7: \u0ba3\u0bcb"
	testConvertExpectFailWithRes(t, s, "?: ṇō") // numerals disabled
}

func TestTranslitConvert6(t *testing.T) {
	s := "\u0BAA\u0BBE\u0BBF" // Two modifier vowels in sequence: \u0BBE\u0BBF
	testConvertExpectFail(t, s)
}

func TestTranslitConvert7(t *testing.T) {
	s := "\u0BAA\u0BBE\u0B9A\u0BC1\u0B95\u0BB3\u0BBF\u0BB2\u0BCD" // One modifier vowel in sequence: \u0BBE
	testConvertExpectOKWithRes(t, s, "pācukaḷil")
}

func TestTranslitConvert8(t *testing.T) {
	s := "பாிசுகளில்" // Two modifier vowels in sequence: \u0BBE\u0BBF
	testConvertExpectFailWithRes(t, s, "pā?cukaḷil")
}

func TestTranslitConvertFailReverseTest(t *testing.T) {

	var tlitNOASCII = NewTranslit()
	tlitNOASCII.alwaysAcceptASCII = true

	var doReverseTest = true
	var s = "பிசுகளில்a"
	var result, expect Result

	expect = Result{OK: false, Msgs: []string{}, Input: s, Result: "picukaḷila"}
	result = tlitNOASCII.ConvertDebug(s, doReverseTest)
	expErr := "reverse test failed"
	err := strings.Join(result.Msgs, ", ")
	result.Msgs = []string{}
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("Expected %#v, got %#v", expect, result)
	}
	if !strings.Contains(err, expErr) {
		t.Errorf("Expected error \"%s\", got: %#v", expErr, err)
	}

}

func TestTranslitConvertWithReverseTest(t *testing.T) {

	var tlitASCII = NewTranslit()
	tlitASCII.alwaysAcceptASCII = true
	var tlitNOASCII = NewTranslit()
	tlitNOASCII.alwaysAcceptASCII = false

	var doReverseTest bool
	var s = "பிசுகளில்a"
	var result, expect Result
	var err, expErr string

	doReverseTest = false
	expect = Result{OK: true, Msgs: []string{}, Input: s, Result: "picukaḷila"}
	result = tlitASCII.ConvertDebug(s, doReverseTest)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("Expected %#v, got %#v", expect, result)
	}

	doReverseTest = false
	expect = Result{OK: false, Msgs: []string{}, Input: s, Result: "picukaḷil?"}
	result = tlitNOASCII.ConvertDebug(s, doReverseTest)
	expErr = "unknown input symbol"
	err = strings.Join(result.Msgs, ", ")
	result.Msgs = []string{}
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("Expected %#v, got %#v", expect, result)
	}
	if !strings.Contains(err, expErr) {
		t.Errorf("Expected error \"%s\", got: %#v", expErr, err)
	}

	doReverseTest = true
	expect = Result{OK: false, Msgs: []string{}, Input: s, Result: "picukaḷila"}
	result = tlitASCII.ConvertDebug(s, doReverseTest)
	err = strings.Join(result.Msgs, ", ")
	expErr = "reverse test failed"
	result.Msgs = []string{}
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("Expected %#v, got %#v", expect, result)
	}
	if !strings.Contains(err, expErr) {
		t.Errorf("Expected error \"%s\", got: %#v", expErr, err)
	}

	doReverseTest = true
	expect = Result{OK: false, Msgs: []string{}, Input: s, Result: "picukaḷil?"}
	result = tlitNOASCII.ConvertDebug(s, doReverseTest)
	err = strings.Join(result.Msgs, ", ")
	expErr = "unknown input symbol"
	result.Msgs = []string{}
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("Expected %#v, got %#v", expect, result)
	}
	if !strings.Contains(err, expErr) {
		t.Errorf("Expected error \"%s\", got: %#v", expErr, err)
	}

}

func TestTranslitRevertWithReverseTest(t *testing.T) {

	var tlitASCII = NewTranslit()
	tlitASCII.alwaysAcceptASCII = true
	var tlitNOASCII = NewTranslit()
	tlitNOASCII.alwaysAcceptASCII = false

	var doReverseTest bool
	var s = "picukaḷil"
	var result, expect Result

	doReverseTest = false
	expect = Result{OK: true, Msgs: []string{}, Input: s, Result: "பிசுகளில்"}
	result = tlitASCII.RevertDebug(s, doReverseTest)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("Expected %#v, got %#v", expect, result)
	}

	doReverseTest = false
	expect = Result{OK: true, Msgs: []string{}, Input: s, Result: "பிசுகளில்"}
	result = tlitNOASCII.RevertDebug(s, doReverseTest)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("Expected %#v, got %#v", expect, result)
	}

	doReverseTest = true
	expect = Result{OK: true, Msgs: []string{}, Input: s, Result: "பிசுகளில்"}
	result = tlitASCII.RevertDebug(s, doReverseTest)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("Expected %#v, got %#v", expect, result)
	}

	doReverseTest = true
	expect = Result{OK: true, Msgs: []string{}, Input: s, Result: "பிசுகளில்"}
	result = tlitNOASCII.RevertDebug(s, doReverseTest)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("Expected %#v, got %#v", expect, result)
	}

	expect = Result{OK: true, Msgs: []string{}, Input: s, Result: "பிசுகளில்"}
	result = tlitASCII.Revert(s)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("Expected %#v, got %#v", expect, result)
	}

	expect = Result{OK: true, Msgs: []string{}, Input: s, Result: "பிசுகளில்"}
	result = tlitNOASCII.Revert(s)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("Expected %#v, got %#v", expect, result)
	}

}

func TestTranslitConvertUnusualVowelSeq(t *testing.T) {
	var s string

	s = "\u0BAA\u0BC6\u0BBE\u0BB0\u0BC1\u0BB3\u0BBE\u0BA4" // பொருளாத
	testConvertExpectOKWithRes(t, s, "poruḷāta")

	s = "மோடி"
	testConvertExpectOKWithRes(t, s, "mōṭi")
}

func TestTranslitConvertNumerals(t *testing.T) {
	var s string

	s = "1 2 3 45645678012013 54: \u0BAA\u0BC6\u0BBE 15- \u0BB0\u0BC1\u0BB3\u0BBE\u0BA4" // பொருளாத
	testConvertExpectOKWithRes(t, s, "1 2 3 45645678012013 54: po 15- ruḷāta")

}
