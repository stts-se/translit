package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tr "github.com/stts-se/translit"
	"github.com/stts-se/translit/buckwalter"
)

var reverse, echoInput, failOnError *bool

func process(s string) {
	s = tr.NFC(s)
	var res string
	var err error
	if *reverse {
		res, err = buckwalter.Bw2Ar(s)
	} else {
		res, err = buckwalter.Ar2Bw(s)
	}
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
	reverse = flag.Bool("r", false, "Reverse conversion (Buckwalter to Arabic)")
	help := flag.Bool("h", false, "Print help and exit")

	var printUsage = func() {
		fmt.Fprintln(os.Stderr, "Transliteration from Arabic to Latin script.")
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
