package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	tr "github.com/stts-se/translit"
	"github.com/stts-se/translit/tamil"
)

func readStdinToString() (string, error) {
	stdin := bufio.NewReader(os.Stdin)
	b, err := ioutil.ReadAll(stdin)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func main() {

	var verb = false
	var tlit = tamil.NewTranslit()

	if len(os.Args) == 2 && strings.HasPrefix(os.Args[1], "-h") {
		fmt.Fprintf(os.Stderr, "Usage: translit <strings or files>\n")
		fmt.Fprintf(os.Stderr, "   or: cat <files> | translit\n")
		os.Exit(0)
	}

	var skipInfo = make(map[string]int)

	nPrinted := 0
	nIn := 0

	if len(os.Args) == 1 {
		stdin, err := readStdinToString()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error : %v\n", err)
			os.Exit(1)
		}
		baseName := "<stdin>"

		lines := strings.Split(strings.TrimSuffix(stdin, "\n"), "\n")
		for _, l := range lines {
			nIn++
			translit := tlit.Convert(l)
			if !translit.OK {
				if verb {
					fmt.Fprintf(os.Stderr, "TRANSLIT ERROR\t%s\t%s\t%s\t%v\n", baseName, l, translit.Result, translit.Msgs)
				}
				skipInfo["TRANSLIT ERROR"]++
				continue
			}
			fmt.Printf("%s\t%s\t%s\n", baseName, l, translit.Result)
			nPrinted++
		}
		if nIn%1000 == 0 {
			fmt.Fprintf(os.Stderr, "\rPROCESSED % 7d utterances", nIn)
		}

	} else {
		for _, arg := range os.Args[1:] {
			if _, err := os.Stat(arg); os.IsNotExist(err) {
				nIn++
				translit := tlit.Convert(arg)
				baseName := path.Base("<stdin>")
				if !translit.OK {
					if verb {
						fmt.Fprintf(os.Stderr, "TRANSLIT ERROR\t%s\t%s\t%s\t%v\n", baseName, arg, translit.Result, translit.Msgs)
					}
					skipInfo["TRANSLIT ERROR"]++
					continue
				}
				fmt.Printf("%s\t%s\t%s\n", baseName, arg, translit.Result)
				nPrinted++
			} else {
				baseName := path.Base(arg)

				lines, err := tr.ReadFile(arg)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error : %v\n", err)
					os.Exit(1)
				}
				for _, l := range lines {
					nIn++
					translit := tlit.Convert(l)
					if !translit.OK {
						if verb {
							fmt.Fprintf(os.Stderr, "TRANSLIT ERROR\t%s\t%s\t%s\t%v\n", baseName, l, translit.Result, translit.Msgs)
						}
						skipInfo["TRANSLIT ERROR"]++
						continue
					}
					fmt.Printf("%s\t%s\t%s\n", baseName, l, translit.Result)
					nPrinted++
				}
			}
			if nIn%1000 == 0 {
				fmt.Fprintf(os.Stderr, "\rPROCESSED % 7d utterances", nIn)
			}
		}
	}
	nSkip := 0
	for _, v := range skipInfo {
		nSkip += v
	}

	pluralS := "s"
	if nIn == 1 {
		pluralS = ""
	}
	fmt.Fprintf(os.Stderr, "\rPROCESSED % 7d utterance%s\n", nIn, pluralS)
	fmt.Fprintf(os.Stderr, "  SKIPPED % 7d\n", nSkip)
	for _, label := range tr.SortKeysByFreq(skipInfo) {
		n := skipInfo[label]
		fmt.Fprintf(os.Stderr, "        : %7d %s\n", n, label)
	}
	fmt.Fprintf(os.Stderr, "  PRINTED % 7d\n", nPrinted)
}
