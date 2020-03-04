package grc

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

type repair struct {
	from *regexp.Regexp
	to   string
}

// https://en.wikipedia.org/wiki/Romanization_of_Greek#Modern_Greek
// Simplified version of ALA-LC [3]

var mapRegexps = []repair{
	{from: regexp.MustCompile(`(^|[\s/()'".!?-])Γ[Κκ](.+)`), to: "${1}G${2}"},
	{from: regexp.MustCompile(`(^|[\s/()'".!?-])Μ[Ππ](.+)`), to: "${1}B${2}"},
	{from: regexp.MustCompile(`(^|[\s/()'".!?-])Ν[Ττ](.+)`), to: "${1}D${2}"},
	{from: regexp.MustCompile(`(^|[\s/()'".!?-])(?i)γκ(.+)`), to: "${1}g${2}"},
	{from: regexp.MustCompile(`(^|[\s/()'".!?-])(?i)μπ(.+)`), to: "${1}b${2}"},
	{from: regexp.MustCompile(`(^|[\s/()'".!?-])(?i)ντ(.+)`), to: "${1}d${2}"},
}

var maptable = []pair{
	{s1: "αι", s2: "ai"},
	{s1: "ει", s2: "ei"},
	{s1: "οι", s2: "oi"},
	{s1: "υι", s2: "yi"},

	{s1: "αυ", s2: "au"},
	{s1: "ευ", s2: "eu"},
	{s1: "ου", s2: "ou"},

	{s1: "αύ", s2: "au"},
	{s1: "εύ", s2: "eú"},
	{s1: "ού", s2: "oú"},
	{s1: "άυ", s2: "áu"},
	{s1: "έυ", s2: "éu"},
	{s1: "όυ", s2: "óu"},

	{s1: "ήυ", s2: "íy"},
	{s1: "υί", s2: "yí"},
	{s1: "ηυ", s2: "iy"},

	{s1: "ωυ", s2: "oy"},
	{s1: "ώυ", s2: "óy"},

	{s1: "μμπ", s2: "mb"},
	{s1: "νντ", s2: "nd"},

	{s1: "ά", s2: "á"},
	{s1: "έ", s2: "é"},
	{s1: "ή", s2: "í"},
	{s1: "ί", s2: "í"},
	{s1: "ύ", s2: "í"},
	{s1: "ό", s2: "ó"},
	{s1: "ώ", s2: "ó"},
	{s1: "ϊ", s2: "ï"},
	{s1: "ΐ", s2: "ḯ"},
	{s1: "ϋ", s2: "ü"},
	{s1: "ΰ", s2: "ǘ"},

	{s1: "α", s2: "a"},
	{s1: "β", s2: "v"},
	{s1: "γγ", s2: "ng"},
	{s1: "γκ", s2: "nk"},
	{s1: "γξ", s2: "nx"},
	{s1: "γχ", s2: "nch"},
	{s1: "γ", s2: "g"},
	{s1: "δ", s2: "d"},
	{s1: "ε", s2: "e"},
	{s1: "ζ", s2: "z"},

	{s1: "η", s2: "i"},
	{s1: "θ", s2: "th"},
	{s1: "ι", s2: "i"},
	{s1: "κ", s2: "k"},
	{s1: "λ", s2: "l"},
	{s1: "μ", s2: "m"},
	{s1: "ν", s2: "n"},
	{s1: "ξ", s2: "x"},
	{s1: "ο", s2: "o"},
	{s1: "π", s2: "p"},
	{s1: "ρ", s2: "r"},
	{s1: "ς", s2: "s"},
	{s1: "σ", s2: "s"},
	{s1: "τ", s2: "t"},
	{s1: "υ", s2: "y"},
	{s1: "φ", s2: "f"},
	{s1: "χ", s2: "ch"},
	{s1: "ψ", s2: "ps"},
	{s1: "ω", s2: "o"},
}

var commonCharsRE = regexp.MustCompile("[A-Za-z0-9()@΄$ï*_]")

var commonChars = map[string]bool{
	" ":  true,
	"\t": true,
	",":  true,
	".":  true,
	"?":  true,
	"!":  true,
	"–":  true,
	"-":  true,
	":":  true,
	";":  true,
	"&":  true,
	"/":  true,
	"'":  true,
}

func Convert(s string) (string, error) {
	s = tr.NFC(s)
	for _, re := range mapRegexps {
		s = re.from.ReplaceAllString(s, re.to)
	}

	intAll := []pair{}
	for _, p := range maptable {
		intAll = append(intAll, p)
		for _, case1 := range tr.UpcaseInitials(p.s1) {
			for _, case2 := range tr.UpcaseInitials(p.s2) {
				intAll = append(intAll, pair{s1: case1, s2: case2})
			}
		}
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
				return "", fmt.Errorf("Couldn't convert '%s' in '%s'", s, sOrig)
			} else if s == sStart {
				res = append(res, head)
				s = strings.TrimPrefix(s, head)
			}
		}
	}
	return strings.Join(res, ""), nil
}
