package main

// References:
// https://en.wikipedia.org/wiki/Romanization_of_Russian
// https://tt.se/tt-spraket/ord-och-begrepp/internationellt/andra-sprak/ryska/

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tr "github.com/stts-se/translit"
	"github.com/stts-se/translit/rus"
)

var echoInput, failOnError *bool

func process(translit rus.Translit, s string) {
	res, err := translit.Convert(s)
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
	swedishOutput := flag.Bool("s", false, "Swedish (TT style) output (default: international output)")
	echoInput = flag.Bool("e", false, "Echo input (default: false)")
	failOnError = flag.Bool("f", false, "Fail on error (default: false)")
	help := flag.Bool("h", false, "Print help and exit")

	var printUsage = func() {
		fmt.Fprintln(os.Stderr, "Transliteration from Russian to Latin script.")
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

	translit := rus.NewTranslit(*swedishOutput)

	if len(flag.Args()) > 0 {
		for _, arg := range flag.Args() {
			if tr.IsFile(arg) {
				lines, err := tr.ReadFile(arg)
				if err != nil {
					log.Fatalf("Couldn't read file: %v", err)
				}
				for _, line := range lines {
					process(translit, line)
				}
			} else {
				process(translit, arg)
			}
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			s := scanner.Text()
			process(translit, s)
		}
	}
}
