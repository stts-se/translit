package far

import (
	"fmt"
	"regexp"
	"strings"

	tr "github.com/stts-se/translit"
)

type pair struct {
	s1 string
	s2 string
}

var maptable = []pair{ // https://en.wikipedia.org/wiki/Romanization_of_Persian EI 2012

	// CONSONANTS
	{s1: "\u0627", s2: "’"}, // TODO: not in the beginning of words
	{s1: "\u0628", s2: "b"},
	{s1: "\u067E", s2: "p"},
	{s1: "\u062A", s2: "t"},
	{s1: "\u062B", s2: "ṯ"},
	{s1: "\u062C", s2: "j"},
	{s1: "\u0686", s2: "č"},
	{s1: "\u062D", s2: "ḥ"},
	{s1: "\u062E", s2: "ḵ"},
	{s1: "\u062F", s2: "d"},
	{s1: "\u0630", s2: "ḏ"},
	{s1: "\u0631", s2: "r"},
	{s1: "\u0632", s2: "z"},
	{s1: "\u0698", s2: "ž"},
	{s1: "\u0633", s2: "s"},
	{s1: "\u0634", s2: "š"},
	{s1: "\u0635", s2: "ṣ"},
	{s1: "\u0636", s2: "ż"},
	{s1: "\u0637", s2: "ṭ"},
	{s1: "\u0638", s2: "ẓ"},
	{s1: "\u0639", s2: "‘"},
	{s1: "\u063A", s2: "ḡ"},
	{s1: "\u0641", s2: "f"},
	{s1: "\u0642", s2: "ḳ"},
	{s1: "\u06A9", s2: "k"},
	{s1: "\u06AF", s2: "g"},
	{s1: "\u0644", s2: "l"},
	{s1: "\u0645", s2: "m"},
	{s1: "\u0646", s2: "n"},
	{s1: "\u0648", s2: "v"},
	{s1: "\u0647", s2: "h"},
	{s1: "\u0629", s2: "h"},
	{s1: "\u06CC", s2: "y"},
	{s1: "\u0621", s2: "’"},
	{s1: "\u0624", s2: "’"},
	{s1: "\u0626", s2: "’"},

	// VOWELS
	{s1: "\u064E", s2: "a"},
	{s1: "\u064F", s2: "u"},
	{s1: "\u0648\u064F", s2: "u"},
	{s1: "\u0650", s2: "e"},
	{s1: "\u064E\u0627", s2: "ā"},
	{s1: "\u0622", s2: "ā"},
	{s1: "\u064E\u06CC", s2: "ā"},
	{s1: "\u06CC\u0670", s2: "ā"},
	{s1: "\u064F\u0648", s2: "u"},
	{s1: "\u0650\u06CC", s2: "i"},
	{s1: "\u064E\u0648", s2: "ow"},
	{s1: "\u064E\u06CC", s2: "ey"}, // duplicate key
	{s1: "\u064E\u06CC", s2: "–e"}, // duplicate key
	{s1: "\u06C0", s2: "–ye"},

	// MISC
	//{s1: "\u200c", s2: ""}, // zero width non-joiner
}

var commonCharsRE = regexp.MustCompile("[A-Za-z0-9()@΄$ï*'_]")

var commonChars = map[string]bool{
	" ": true,
	".": true,
	",": true,
	"(": true,
	")": true,
	//"\u200c": true, // zero width non-joiner
}

var echoInput, failOnError *bool

func Convert(s string) (string, error) {
	s = tr.NFC(s)

	// for _, re := range mapRegexps {
	// 	s = re.from.ReplaceAllString(s, re.to)
	// }

	intAll := []pair{}
	for _, p := range maptable {
		intAll = append(intAll, p)
		for _, case1 := range tr.UpcaseInitials(p.s1) {
			for _, case2 := range tr.UpcaseInitials(p.s2) {
				intAll = append(intAll, pair{s1: case1, s2: case2})
			}
		}

		// intAll = append(intAll, p)
		// intAll = append(intAll, pair{s1: tr.UpcaseInitial(p.s1), s2: tr.UpcaseInitial(p.s2)})
		// intAll = append(intAll, pair{s1: tr.Upcase(p.s1), s2: tr.Upcase(p.s2)})
		// if len([]rune(p.s1)) > 2 {
		// 	intAll = append(intAll, pair{s1: tr.UpcaseTwoInitials(p.s1), s2: tr.UpcaseInitial(p.s2)})
		// }
	}
	res, err := innerConvert(intAll, s, true)
	if err != nil {
		return "", err
	}
	return res, nil
}

func innerConvert(chsAll []pair, s string, requireAllMapped bool) (string, error) {
	sOrig := s
	res := []string{}
	for len(s) > 0 {
		sStart := s
		head := string([]rune(s)[0])
		for _, p := range chsAll {
			if strings.HasPrefix(s, p.s1) {
				res = append(res, p.s2)
				s = strings.TrimPrefix(s, p.s1)
				break
			}
			// check for common chars
			if _, ok := commonChars[head]; ok {
				res = append(res, head)
				s = strings.TrimPrefix(s, head)
				break
			}
			if commonCharsRE.MatchString(head) {
				res = append(res, head)
				s = strings.TrimPrefix(s, head)
				break
			}
		}
		if s == sStart { // nothing found to map for this prefix
			if requireAllMapped {
				return "", fmt.Errorf("Couldn't convert '%s'\t%v\tin '%s'", s, tr.UnicodeInfo(s)[0], sOrig)
			} else if s == sStart {
				res = append(res, head)
				s = strings.TrimPrefix(s, head)
			}
		}
	}
	return strings.Join(res, ""), nil
}
