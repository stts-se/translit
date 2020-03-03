package translit

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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
