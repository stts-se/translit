package rus

// References:
// https://en.wikipedia.org/wiki/Romanization_of_Russian
// https://tt.se/tt-spraket/ord-och-begrepp/internationellt/andra-sprak/ryska/

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

type rPair struct {
	from *regexp.Regexp
	to   string
}

// Translit
type Translit struct {
	SwedishOutput bool
}

func NewTranslit(swedishOutput bool) Translit {
	return Translit{SwedishOutput: swedishOutput}
}

var roadSigns = []pair{ // https://en.wikipedia.org/wiki/Romanization_of_Russian -- Road signs
	{s1: "а", s2: "a"},
	{s1: "б", s2: "b"},
	{s1: "в", s2: "v"},
	{s1: "г", s2: "g"},
	{s1: "д", s2: "d"},
	{s1: "е", s2: "e"},
	{s1: "ё", s2: "e"},
	{s1: "ж", s2: "zh"},
	{s1: "з", s2: "z"},
	{s1: "и", s2: "i"},
	{s1: "й", s2: "y"}, // j
	{s1: "к", s2: "k"},
	{s1: "л", s2: "l"},
	{s1: "м", s2: "m"},
	{s1: "н", s2: "n"},
	{s1: "о", s2: "o"},
	{s1: "п", s2: "p"},
	{s1: "р", s2: "r"},
	{s1: "с", s2: "s"},
	{s1: "т", s2: "t"},
	{s1: "у", s2: "u"},
	{s1: "ф", s2: "f"},
	{s1: "х", s2: "kh"},
	{s1: "ц", s2: "ts"},
	{s1: "ч", s2: "ch"},
	{s1: "ш", s2: "sh"},
	{s1: "щ", s2: "shch"},
	{s1: "ъ", s2: "ie"}, // ’
	{s1: "ы", s2: "y"},
	{s1: "ь", s2: "’"},
	{s1: "э", s2: "e"},
	{s1: "ю", s2: "yu"}, // ju
	{s1: "я", s2: "ya"}, // ja
}

var international = roadSigns

var commonChars = map[string]bool{
	" ":      true,
	",":      true,
	".":      true,
	"?":      true,
	"!":      true,
	"–":      true,
	"-":      true,
	":":      true,
	";":      true,
	"\u0301": true, // Combining acute accent
}

// https://tt.se/tt-spraket/ord-och-begrepp/internationellt/andra-sprak/ryska/
var swePairs = []pair{
	{s1: "zh", s2: "zj"},
	{s1: "kh", s2: "ch"},
	{s1: "ch", s2: "tj"},
	{s1: "sh", s2: "sj"},
	//pair{s1: "shch", s2: "sjtj"}, // not needed
	{s1: "yu", s2: "ju"},
	{s1: "ya", s2: "ja"},
	{s1: "ye", s2: "je"}, // ?? will this ever happen?
}

var sweREs = []rPair{ // https://tt.se/tt-spraket/ord-och-begrepp/internationellt/andra-sprak/ryska/

	// must be unambiguous for output

	// rPair{from: regexp.MustCompile(`(?:i)ev([\p{P}]|$)`), to: "(j)ev$1"},
	// rPair{from: regexp.MustCompile(`(?:i)ov([\p{P}]|$)`), to: "ov$1"},
	// rPair{from: regexp.MustCompile(`(?:i)ev([\p{P}]|$)`), to: "jov$1"}, // Gorbatjov

	{from: regexp.MustCompile(`(?i)ky(\b|$)`), to: "kij$2"},
	{from: regexp.MustCompile(`(?i)gy(\b|$)`), to: "gij$2"},
	{from: regexp.MustCompile(`(?i)ay(\b|$)`), to: "aj$2"},
	{from: regexp.MustCompile(`(?i)ey(\b|$)`), to: "ej$2"},
	{from: regexp.MustCompile(`(?i)y(\b|$)`), to: "yj$2"},
}

func (translit Translit) Convert(s string) (string, error) {
	s = tr.NFC(s)
	intAll := []pair{}
	for _, p := range international {
		intAll = append(intAll, p)
		intAll = append(intAll, pair{s1: tr.UpcaseInitial(p.s1), s2: tr.UpcaseInitial(p.s2)})
		intAll = append(intAll, pair{s1: tr.Upcase(p.s1), s2: tr.Upcase(p.s2)})
	}
	res, err := translit.innerConvert(intAll, s, true)
	if err != nil {
		return "", err
	}
	if translit.SwedishOutput {
		sweAll := []pair{}
		for _, p := range swePairs {
			sweAll = append(sweAll, p)
			sweAll = append(sweAll, pair{s1: tr.UpcaseInitial(p.s1), s2: tr.UpcaseInitial(p.s2)})
			sweAll = append(sweAll, pair{s1: tr.Upcase(p.s1), s2: tr.Upcase(p.s2)})
		}
		res, err := translit.innerConvert(sweAll, res, false)
		if err != nil {
			return "", err
		}
		for _, p := range sweREs {
			res = p.from.ReplaceAllString(res, p.to)
		}
	}
	return res, nil
}

func (translit Translit) innerConvert(chsAll []pair, s string, requireAllMapped bool) (string, error) {
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
