package translit

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/unicode/runenames"
)

func ReadFile(fn string) ([]string, error) {
	fn = filepath.Clean(fn)
	var res []string
	var scanner *bufio.Scanner
	fh, err := os.Open(fn)
	if err != nil {
		return res, fmt.Errorf("failed to read '%s' : %v", fn, err)
	}

	if strings.HasSuffix(fn, ".gz") {
		gz, err := gzip.NewReader(fh)
		if err != nil {
			return res, fmt.Errorf("failed to read '%s' : %v", fn, err)
		}
		scanner = bufio.NewScanner(gz)
	} else {
		scanner = bufio.NewScanner(fh)
	}
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return res, fmt.Errorf("failed to read '%s' : %v", fn, err)
	}
	return res, nil
}

// func ReadFile(fName string) ([]string, error) {
// 	b, err := ioutil.ReadFile(filepath.Clean(fName))
// 	if err != nil {
// 		return []string{}, err
// 	}
// 	s := strings.TrimSuffix(string(b), "\n")
// 	s = strings.Replace(s, "\r", "", -1)
// 	return strings.Split(s, "\n"), nil
// }

func SortKeysByFreq(m map[string]int) []string {
	res := []string{}
	for k := range m {
		res = append(res, k)
	}

	sort.Slice(res, func(i, j int) bool { return m[res[i]] > m[res[j]] })
	return res
}

func StringsContains(slice []string, elem string) bool {
	for _, e0 := range slice {
		if e0 == elem {
			return true
		}
	}
	return false
}

//NFC convert string
func NFC(s string) string {
	normed, _, _ := transform.String(norm.NFC, s)
	return normed
}

func IsFile(fName string) bool {
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		return false
	}
	return true
}

// Upcase string
func Upcase(s string) string {
	return strings.ToUpper(s)
}

// UpcaseInitial upcase first character, downcase the rest
func UpcaseInitial(s string) string {
	runes := []rune(s)
	head := ""
	if len(runes) > 0 {
		head = strings.ToUpper(string(runes[0]))
	}
	tail := ""
	if len(runes) > 0 {
		tail = strings.ToLower(string(runes[1:]))
	}
	return head + tail
}

// func UpcaseTwoInitials(s string) string {
// 	runes := []rune(s)
// 	head := ""
// 	if len(runes) > 0 {
// 		head = strings.ToUpper(string(runes[0:2]))
// 	}
// 	tail := ""
// 	if len(runes) > 2 {
// 		tail = strings.ToLower(string(runes[2:]))
// 	}
// 	//fmt.Println("??? 3", len(runes), s, head, tail)
// 	return head + tail
// }

// UpcaseInitials generate all variations upcasing a first group of characters, downcasing the rest
func UpcaseInitials(s string) []string {
	var res []string

	runes := []rune(s)
	len := len(runes)
	if len == 0 {
		return []string{s}
	}
	if len == 1 {
		return []string{UpcaseInitial(s)}
	}

	for i := 1; i < len; i++ {
		head := strings.ToUpper(string(runes[0:i]))
		tail := strings.ToLower(string(runes[i:]))
		res = append(res, head+tail)
	}
	return res
}

// Unicode info stuff

func unicodeBlockFor(r rune) string {
	for s, t := range unicode.Scripts {
		if unicode.In(r, t) {
			return s
		}
	}
	return "<UNDEF>"
}

func codeFor(r rune) string {
	uc := fmt.Sprintf("%U", r)
	return fmt.Sprintf("\\u%s", uc[2:])
}

var ucNumberRe = regexp.MustCompile(`^(?:\\u|[uU][+])([a-fA-F0-9]{4})$`)

const newline rune = '\n'

var hardwiredUnicodeNames = map[rune]string{
	newline: "NEWLINE",
	'	': "TAB",
}

func unicodeNameFor(r rune) string {
	if name, ok := hardwiredUnicodeNames[r]; ok {
		return name
	}
	return runenames.Name(r)
}

func inhibitSpecialChar(r rune) bool {
	_, ok := hardwiredUnicodeNames[r]
	return ok
}

type UnicodeChar struct {
	Char, Name, Code, Block string
}

func UnicodeInfo(s string) []UnicodeChar {
	var res []UnicodeChar
	for _, r := range []rune(s) {
		thisS := string(r)
		if inhibitSpecialChar(r) {
			thisS = ""
		}
		res = append(res, UnicodeChar{
			Char:  thisS,
			Name:  unicodeNameFor(r),
			Code:  codeFor(r),
			Block: unicodeBlockFor(r),
		})
	}
	return res
}
