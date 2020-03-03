package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	tr "github.com/stts-se/translit"
)

type pair struct {
	s1 string
	s2 string
}

var maptable = []pair{

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
	{s1: "\u064E\u06CC", s2: "ey"}, //?? duplicate key
	{s1: "\u064E\u06CC", s2: "–e"}, // ?? duplicate key
	{s1: "\u06C0", s2: "–ye"},
}

var commonCharsRE = regexp.MustCompile("[A-Za-z0-9()@΄$ï*_]")

var commonChars = map[string]bool{
	" ": true,
	".": true,
	",": true,
	"(": true,
	")": true,
}

var echoInput, failOnError *bool

func upcase(s string) string {
	return strings.ToUpper(s)
}

func upcaseInitial(s string) string {
	runes := []rune(s)
	head := ""
	if len(runes) > 0 {
		head = strings.ToUpper(string(runes[0]))
	}
	tail := ""
	if len(runes) > 1 {
		tail = strings.ToLower(string(runes[1:]))
	}
	//fmt.Println("??? 2", len(runes), s, head, tail)
	return head + tail
}

func upcaseTwoInitials(s string) string {
	runes := []rune(s)
	head := ""
	if len(runes) > 0 {
		head = strings.ToUpper(string(runes[0:2]))
	}
	tail := ""
	if len(runes) > 2 {
		tail = strings.ToLower(string(runes[2:]))
	}
	//fmt.Println("??? 3", len(runes), s, head, tail)
	return head + tail
}

func convert(s string) (string, error) {
	// for _, re := range mapRegexps {
	// 	s = re.from.ReplaceAllString(s, re.to)
	// }

	intAll := []pair{}
	for _, p := range maptable {
		intAll = append(intAll, p)
		intAll = append(intAll, pair{s1: upcaseInitial(p.s1), s2: upcaseInitial(p.s2)})
		intAll = append(intAll, pair{s1: upcase(p.s1), s2: upcase(p.s2)})
		if len([]rune(p.s1)) > 2 {
			intAll = append(intAll, pair{s1: upcaseTwoInitials(p.s1), s2: upcaseInitial(p.s2)})
		}
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

func process(s string) {
	s = tr.NFC(s)
	res, err := convert(s)
	if err != nil {
		if *failOnError {
			log.Fatalf(fmt.Sprintf("%v", err))
		} else {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("ERROR %s\t%v\n", s, err))
			return
		}
	}
	if *echoInput {
		fmt.Printf("%s\t%s\n", s, res)
	} else {
		fmt.Printf("%s\n", res)
	}
}

func main() {

	cmdname := filepath.Base(os.Args[0])
	echoInput = flag.Bool("e", false, "Echo input (default: false)")
	failOnError = flag.Bool("f", false, "Fail on error (default: false)")
	help := flag.Bool("h", false, "Print help and exit")

	var printUsage = func() {
		fmt.Fprintln(os.Stderr, "Transliteration from Farsi to Latin script.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, cmdname+" <input file(s)>")
		fmt.Fprintln(os.Stderr, cmdname+" <input string(s)>")
		fmt.Fprintln(os.Stderr, "cat <input file(s)> | "+cmdname)
		fmt.Fprintln(os.Stderr, "\nOptional flags:")
		flag.PrintDefaults()
	}

	flag.Usage = func() {
		printUsage()
		os.Exit(0)
	}

	flag.Parse()

	if *help { // if flag.NArg() < 1 {
		printUsage()
		os.Exit(0)
	}

	if len(flag.Args()) > 0 {
		for _, arg := range flag.Args() {
			if tr.IsFile(arg) {
				lines, err := tr.ReadFile(arg)
				if err != nil {
					log.Fatalf("Couldn't read file: %v", err)
				}
				for _, line := range lines {
					process(line)
				}
			} else {
				process(arg)
			}
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			s := scanner.Text()
			process(s)
		}
	}
}
