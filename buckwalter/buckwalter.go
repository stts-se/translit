package buckwalter

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/unicode/runenames"
)

type ch struct {
	ar rune
	bw rune
}

func (ch ch) desc() string {
	name := runenames.Name(ch.ar)
	uc := codeFor(ch.ar)
	block := blockFor(ch.ar)
	return fmt.Sprintf("%s %s %s", uc, name, block)
}

type maptable struct {
	from  string
	to    string
	table map[rune]rune
}

func (m maptable) name() string {
	return fmt.Sprintf("%s2%s", m.from, m.to)
}

var defaultChar = '?'
var charset = []ch{
	ch{'ا', 'A'}, // bare alif
	ch{'ب', 'b'},
	ch{'ت', 't'},
	ch{'ث', 'v'},
	ch{'ج', 'j'},
	ch{'ح', 'H'},
	ch{'خ', 'x'},
	ch{'د', 'd'}, // dal \u062F
	ch{'ذ', '*'},
	ch{'ر', 'r'},
	ch{'ز', 'z'},
	ch{'س', 's'},
	ch{'ش', '$'},
	ch{'ص', 'S'},
	ch{'ض', 'D'},
	ch{'ط', 'T'},
	ch{'ظ', 'Z'},
	ch{'ع', 'E'},
	ch{'غ', 'g'},
	ch{'ف', 'f'},
	ch{'ق', 'q'},
	ch{'ك', 'k'},
	ch{'ل', 'l'},
	ch{'م', 'm'},
	ch{'ن', 'n'},
	ch{'ه', 'h'},
	ch{'و', 'w'},
	ch{'ي', 'y'},
	ch{'ة', 'p'}, //teh marbuta

	ch{'\u064E', 'a'}, // fatha
	ch{'\u064f', 'u'}, // damma
	ch{'\u0650', 'i'}, // kasra
	ch{'\u064B', 'F'}, // fathatayn
	ch{'\u064C', 'N'}, // dammatayn
	ch{'\u064D', 'K'}, // kasratayn
	ch{'\u0651', '~'}, // shadda
	ch{'\u0652', 'o'}, // sukun

	ch{'\u0621', '\''}, // lone hamza
	ch{'\u0623', '>'},  // hamza on alif
	ch{'\u0625', '<'},  // hamza below alif
	ch{'\u0624', '&'},  // hamza on wa
	ch{'\u0626', '}'},  // hamza on ya

	ch{'\u0622', '|'}, // madda on alif
	ch{'\u0671', '{'}, // alif al-wasla
	ch{'\u0670', '`'}, // dagger alif
	ch{'\u0649', 'Y'}, // alif maqsura

	// Arabic-indic digits
	ch{'\u0660', '0'},
	ch{'\u0661', '1'},
	ch{'\u0662', '2'},
	ch{'\u0663', '3'},
	ch{'\u0664', '4'},
	ch{'\u0665', '5'},
	ch{'\u0666', '6'},
	ch{'\u0667', '7'},
	ch{'\u0668', '8'},
	ch{'\u0669', '9'},

	// punctuation
	ch{'\u060C', ','},
	ch{'\u061B', ';'},
	ch{'\u061F', '?'},

	// http://www.qamus.org/transliteration.htm
	ch{'\u067e', 'P'}, // peh
	ch{'\u0686', 'J'}, // tcheh
	ch{'\u06a4', 'V'}, // veh
	ch{'\u06af', 'G'}, // gaf
	//ch{'\u0640', '_'}, // tatweel

}

var commonChars = map[rune]bool{
	'\u00A0': true, // non-breaking space
	' ':      true,
	'.':      true,
	',':      true,
	'(':      true,
	')':      true,
}

const alwaysAcceptASCII = false // REVERSE TEST DOES NOT WORK WITH THIS SETTING ON

func isCommonChar(sym rune) bool {
	if _, ok := commonChars[sym]; ok {
		return true
	}
	if alwaysAcceptASCII && int(sym) < 128 {
		return true
	}
	return false
}

func makeAr2bwMap() maptable {
	m := map[rune]rune{}
	for _, ch := range charset {
		m[ch.ar] = ch.bw
	}
	return maptable{"ar", "bw", m}
}

func makeBw2ArMap() maptable {
	m := map[rune]rune{}
	for _, ch := range charset {
		m[ch.bw] = ch.ar
	}
	return maptable{"bw", "ar", m}
}

var ar2bwMap = makeAr2bwMap()
var bw2arMap = makeBw2ArMap()

var bwDenormRe = regexp.MustCompile("([aiuoFKN])(~)")
var bwDenormReTo = "$2$1"

func arPostNorm(s string) string {
	res, _, _ := transform.String(norm.NFC, s)
	return res
}

func bwPostNorm(s string) string {
	return bwDenormRe.ReplaceAllString(s, bwDenormReTo)
}
func arPreNorm(s string) string {
	var res = s
	res = strings.Replace(res, "\uFEAA", "\u062F", -1) // DAL FINAL FORM => DAL
	res = strings.Replace(res, "\u06BE", "\u0647", -1) // HEH DOACHASHMEE => HEH
	res = strings.Replace(res, "\u200F", "", -1)       // RTL MARK
	return res
}

func bwPreNorm(s string) string {
	return s
}

func postNormalise(outputName, s string) string {
	if outputName == "bw" {
		return bwPostNorm(s)
	} else {
		return arPostNorm(s)
	}
}

func preNormalise(outputName, s string) string {
	if outputName == "bw" {
		return bwPreNorm(s)
	} else {
		return arPreNorm(s)
	}
}

func reverseTest(mapTo string, input string, mapped string) error {
	remaptable := maptable{}
	if mapTo == ar2bwMap.to {
		remaptable = ar2bwMap
	} else if mapTo == bw2arMap.to {
		remaptable = bw2arMap
	}
	remapped, err := convert(remaptable, mapped, false)
	if err != nil {
		return err
	}
	if remapped != input {
		return fmt.Errorf("reverse test failed: input '%s', mapped '%s', remapped '%s'", input, mapped, remapped)
	}
	return nil
}

func convert(maptable maptable, input string, doReverseTest bool) (string, error) {
	//fmt.Fprintf(os.Stderr, "convert from %s | input: %s\n", mapName, input)
	input = preNormalise(maptable.from, input)
	res := []rune{}
	errs := []string{}
	for _, sym := range input {
		mapped, exists := maptable.table[rune(sym)]
		if !exists {
			if ok := isCommonChar(sym); ok {
				mapped = sym
			} else {
				mapped = defaultChar
				errs = append(errs, fmt.Sprintf("no mapping for %s symbol '%s'", maptable.name(), string(sym)))
			}
		}
		res = append(res, mapped)
		//fmt.Fprintf(os.Stderr, "convert | '%s' -> '%s'\n", string(sym), string(mapped))
	}
	mapped := string(res)
	mapped = postNormalise(maptable.to, mapped)

	if len(errs) > 0 {
		err := strings.Join(errs, "; ")
		return mapped, fmt.Errorf("%s", err)
	}
	if doReverseTest {
		err := reverseTest(maptable.from, input, mapped)
		if err != nil {
			return mapped, err
		}
	}
	return mapped, nil
}

// Bw2Ar converts an input Buckwalter string into Arabic alphabet. An error is returned if there are unknown symbols in the input string. The resulting string always contains the full converted string (but with '?' for unknown chars). The Arabic output is NFC normalised (cons + vowel + cons length).
func Bw2Ar(s string) (string, error) {
	return convert(bw2arMap, s, true)
}

// Ar2Bw converts an input Arabic string into Buckwalter. An error is returned if there are unknown symbols in the input string. The resulting string always contains the full converted string (but with '?' for unknown chars). The output is in Buckwalter order (cons + cons length + vowel) -- i.e., not matching Arabic script NFC normalisation.
func Ar2Bw(s string) (string, error) {
	return convert(ar2bwMap, s, true)
}

func blockFor(r rune) string {
	for s, t := range unicode.Scripts {
		if unicode.In(r, t) {
			return s
		}
	}
	return "<UNDEF>"
}

func codeFor(r rune) string {
	uc := fmt.Sprintf("%U", r)
	return fmt.Sprintf("U+%s", uc[2:])
}

type CharEntry struct {
	Ar   string
	Bw   string
	Desc string
}

// BuildCharTable creates a list of the character mappings, for use in human readable docs
func BuildCharTable() []CharEntry {
	res := []CharEntry{}
	for _, ch := range charset {
		entry := CharEntry{string(ch.bw), string(ch.ar), ch.desc()}
		res = append(res, entry)
	}
	return res
}
