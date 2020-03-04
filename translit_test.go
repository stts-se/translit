package translit

import (
	//"fmt"
	"reflect"
	// "strings"
	"testing"
)

var fsExpGot = "Expected: %v got: %v"

// container to compare variables in testing
type tVar struct {
	name  string
	value string
}

func TestUpcase(t *testing.T) {
	input := "apiö"
	expect := "APIÖ"
	result := Upcase(input)
	if result != expect {
		t.Errorf(fsExpGot, expect, result)
	}
}

func TestUpcaseInitial(t *testing.T) {
	tests := map[string]string{
		"öiÆa": "Öiæa",
		"iöAa": "Iöaa",
		"Ö":    "Ö",
		"å":    "Å",
		"":     "",
	}

	for input, expect := range tests {
		result := UpcaseInitial(input)
		if result != expect {
			t.Errorf(fsExpGot, expect, result)
		}
	}
}

func TestUpcaseInitials(t *testing.T) {
	tests := map[string][]string{
		"öiÆ": []string{
			"Öiæ",
			"ÖIæ",
		},
		"öiÆx": []string{
			"Öiæx",
			"ÖIæx",
			"ÖIÆx",
		},
		"ö": []string{
			"Ö",
		},
		"": []string{
			"",
		},
	}

	for input, expect := range tests {
		result := UpcaseInitials(input)
		if !reflect.DeepEqual(result, expect) {
			t.Errorf(fsExpGot, expect, result)
		}
	}
}
