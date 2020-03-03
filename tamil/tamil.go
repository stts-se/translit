package tamil

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/stts-se/translit"
)

// Translit
type Translit struct {
	theTree           *rNode
	revTree           *rNode
	alwaysAcceptASCII bool
	defaultChar       string
}

// Result struct
type Result struct {
	Input string // Input string
	//InputNorm string   // Normalised input string
	Result string   // Converted string
	Msgs   []string // Error messages, if any
	OK     bool     // Conversion success true/false
}

var script2transMap = map[string]string{
	// Extra
	"\u0B83": "ḵ", // Visarga
	// "\u0BD7": "???", // TAMIL AU LENGTH MARK
	// "\u0BB6\u0BCD\u0BB0\u0BC0": "<shrii>", // Shrii
	//"\u200C":                   "", // Zero width non-joiner should not be included

	"\u0B95\u0BCD":             "k",
	"\u0B99\u0BCD":             "ṅ",
	"\u0B9A\u0BCD":             "c",
	"\u0B9E\u0BCD":             "ñ",
	"\u0B9F\u0BCD":             "ṭ",
	"\u0BA3\u0BCD":             "ṇ",
	"\u0BA4\u0BCD":             "t",
	"\u0BA8\u0BCD":             "n",
	"\u0BAA\u0BCD":             "p",
	"\u0BAE\u0BCD":             "m",
	"\u0BAF\u0BCD":             "y",
	"\u0BB0\u0BCD":             "r",
	"\u0BB2\u0BCD":             "l",
	"\u0BB5\u0BCD":             "v",
	"\u0BB4\u0BCD":             "ḻ",
	"\u0BB3\u0BCD":             "ḷ",
	"\u0BB1\u0BCD":             "ṟ",
	"\u0BA9\u0BCD":             "ṉ",
	"\u0B9C\u0BCD":             "j",
	"\u0BB6\u0BCD":             "ś",
	"\u0BB7\u0BCD":             "ṣ",
	"\u0BB8\u0BCD":             "s",
	"\u0BB9\u0BCD":             "h",
	"\u0B95\u0BCD\u0BB7\u0BCD": "kṣ",
	"\u0B85":                   "a",
	"\u0B86":                   "ā",
	"\u0B87":                   "i",
	"\u0B88":                   "ī",
	"\u0B89":                   "u",
	"\u0B8A":                   "ū",
	"\u0B8E":                   "e",
	"\u0B8F":                   "ē",
	"\u0B90":                   "ai",
	"\u0B92":                   "o",
	"\u0B93":                   "ō",
	"\u0B94":                   "au",

	// "\u0BE6": "0",
	// "\u0BE7": "1",
	// "\u0BE8": "2",
	// "\u0BE9": "3",
	// "\u0BEA": "4",
	// "\u0BEB": "5",
	// "\u0BEC": "6",
	// "\u0BED": "7",
	// "\u0BEE": "8",
	// "\u0BEF": "9",
	// "\u0BF0": "10",
	// "\u0BF1": "100",
	// "\u0BF2": "1000",

	// "\u0BF3": "{day}",
	// "\u0BF4": "{month}",
	// "\u0BF5": "{year}",
	// "\u0BF6": "{debit}",
	// "\u0BF7": "{credit//a}",
	// "\u0BF8": "{credit//b}",
	// "\u0BF9": "{rupee}",
	// "\u0BFA": "{numeral}",
	// "\u0BB3": "{time}",
	// "\u0BB5": "{quantity}",

	// TEST WITH COMBINATIONS
	"\u0BB6":                   "śa",
	"\u0B9C":                   "ja",
	"\u0BB7":                   "ṣa",
	"\u0BB8":                   "sa",
	"\u0BB9":                   "ha",
	"\u0B95\u0BCD\u0BB7":       "kṣa",
	"\u0BB6\u0BBE":             "śā",
	"\u0B9C\u0BBE":             "jā",
	"\u0BB7\u0BBE":             "ṣā",
	"\u0BB8\u0BBE":             "sā",
	"\u0BB9\u0BBE":             "hā",
	"\u0B95\u0BCD\u0BB7\u0BBE": "kṣā",
	"\u0BB6\u0BBF":             "śi",
	"\u0B9C\u0BBF":             "ji",
	"\u0BB7\u0BBF":             "ṣi",
	"\u0BB8\u0BBF":             "si",
	"\u0BB9\u0BBF":             "hi",
	"\u0B95\u0BCD\u0BB7\u0BBF": "kṣi",
	"\u0BB6\u0BC0":             "śī",
	"\u0B9C\u0BC0":             "jī",
	"\u0BB7\u0BC0":             "ṣī",
	"\u0BB8\u0BC0":             "sī",
	"\u0BB9\u0BC0":             "hī",
	"\u0B95\u0BCD\u0BB7\u0BC0": "kṣī",
	"\u0BB6\u0BC1":             "śu",
	"\u0B9C\u0BC1":             "ju",
	"\u0BB7\u0BC1":             "ṣu",
	"\u0BB8\u0BC1":             "su",
	"\u0BB9\u0BC1":             "hu",
	"\u0B95\u0BCD\u0BB7\u0BC1": "kṣu",
	"\u0BB6\u0BC2":             "śū",
	"\u0B9C\u0BC2":             "jū",
	"\u0BB7\u0BC2":             "ṣū",
	"\u0BB8\u0BC2":             "sū",
	"\u0BB9\u0BC2":             "hū",
	"\u0B95\u0BCD\u0BB7\u0BC2": "kṣū",
	"\u0BB6\u0BC6":             "śe",
	"\u0B9C\u0BC6":             "je",
	"\u0BB7\u0BC6":             "ṣe",
	"\u0BB8\u0BC6":             "se",
	"\u0BB9\u0BC6":             "he",
	"\u0B95\u0BCD\u0BB7\u0BC6": "kṣe",
	"\u0BB6\u0BC7":             "śē",
	"\u0B9C\u0BC7":             "jē",
	"\u0BB7\u0BC7":             "ṣē",
	"\u0BB8\u0BC7":             "sē",
	"\u0BB9\u0BC7":             "hē",
	"\u0B95\u0BCD\u0BB7\u0BC7": "kṣē",
	"\u0BB6\u0BC8":             "śai",
	"\u0B9C\u0BC8":             "jai",
	"\u0BB7\u0BC8":             "ṣai",
	"\u0BB8\u0BC8":             "sai",
	"\u0BB9\u0BC8":             "hai",
	"\u0B95\u0BCD\u0BB7\u0BC8": "kṣai",
	"\u0BB6\u0BCA":             "śo",
	"\u0B9C\u0BCA":             "jo",
	"\u0BB7\u0BCA":             "ṣo",
	"\u0BB8\u0BCA":             "so",
	"\u0BB9\u0BCA":             "ho",
	"\u0B95\u0BCD\u0BB7\u0BCA": "kṣo",
	"\u0BB6\u0BCB":             "śō",
	"\u0B9C\u0BCB":             "jō",
	"\u0BB7\u0BCB":             "ṣō",
	"\u0BB8\u0BCB":             "sō",
	"\u0BB9\u0BCB":             "hō",
	"\u0B95\u0BCD\u0BB7\u0BCB": "kṣō",
	"\u0BB6\u0BCC":             "śau",
	"\u0B9C\u0BCC":             "jau",
	"\u0BB7\u0BCC":             "ṣau",
	"\u0BB8\u0BCC":             "sau",
	"\u0BB9\u0BCC":             "hau",
	"\u0B95\u0BCD\u0BB7\u0BCC": "kṣau",
	"\u0B95":                   "ka",
	"\u0B99":                   "ṅa",
	"\u0B9A":                   "ca",
	"\u0B9E":                   "ña",
	"\u0B9F":                   "ṭa",
	"\u0BA3":                   "ṇa",
	"\u0BA4":                   "ta",
	"\u0BA8":                   "na",
	"\u0BAA":                   "pa",
	"\u0BAE":                   "ma",
	"\u0BAF":                   "ya",
	"\u0BB0":                   "ra",
	"\u0BB2":                   "la",
	"\u0BB5":                   "va",
	"\u0BB4":                   "ḻa",
	"\u0BB3":                   "ḷa",
	"\u0BB1":                   "ṟa",
	"\u0BA9":                   "ṉa",
	"\u0B95\u0BBE":             "kā",
	"\u0B99\u0BBE":             "ṅā",
	"\u0B9A\u0BBE":             "cā",
	"\u0B9E\u0BBE":             "ñā",
	"\u0B9F\u0BBE":             "ṭā",
	"\u0BA3\u0BBE":             "ṇā",
	"\u0BA4\u0BBE":             "tā",
	"\u0BA8\u0BBE":             "nā",
	"\u0BAA\u0BBE":             "pā",
	"\u0BAE\u0BBE":             "mā",
	"\u0BAF\u0BBE":             "yā",
	"\u0BB0\u0BBE":             "rā",
	"\u0BB2\u0BBE":             "lā",
	"\u0BB5\u0BBE":             "vā",
	"\u0BB4\u0BBE":             "ḻā",
	"\u0BB3\u0BBE":             "ḷā",
	"\u0BB1\u0BBE":             "ṟā",
	"\u0BA9\u0BBE":             "ṉā",
	"\u0B95\u0BBF":             "ki",
	"\u0B99\u0BBF":             "ṅi",
	"\u0B9A\u0BBF":             "ci",
	"\u0B9E\u0BBF":             "ñi",
	"\u0B9F\u0BBF":             "ṭi",
	"\u0BA3\u0BBF":             "ṇi",
	"\u0BA4\u0BBF":             "ti",
	"\u0BA8\u0BBF":             "ni",
	"\u0BAA\u0BBF":             "pi",
	"\u0BAE\u0BBF":             "mi",
	"\u0BAF\u0BBF":             "yi",
	"\u0BB0\u0BBF":             "ri",
	"\u0BB2\u0BBF":             "li",
	"\u0BB5\u0BBF":             "vi",
	"\u0BB4\u0BBF":             "ḻi",
	"\u0BB3\u0BBF":             "ḷi",
	"\u0BB1\u0BBF":             "ṟi",
	"\u0BA9\u0BBF":             "ṉi",
	"\u0B95\u0BC0":             "kī",
	"\u0B99\u0BC0":             "ṅī",
	"\u0B9A\u0BC0":             "cī",
	"\u0B9E\u0BC0":             "ñī",
	"\u0B9F\u0BC0":             "ṭī",
	"\u0BA3\u0BC0":             "ṇī",
	"\u0BA4\u0BC0":             "tī",
	"\u0BA8\u0BC0":             "nī",
	"\u0BAA\u0BC0":             "pī",
	"\u0BAE\u0BC0":             "mī",
	"\u0BAF\u0BC0":             "yī",
	"\u0BB0\u0BC0":             "rī",
	"\u0BB2\u0BC0":             "lī",
	"\u0BB5\u0BC0":             "vī",
	"\u0BB4\u0BC0":             "ḻī",
	"\u0BB3\u0BC0":             "ḷī",
	"\u0BB1\u0BC0":             "ṟī",
	"\u0BA9\u0BC0":             "ṉī",
	"\u0B95\u0BC1":             "ku",
	"\u0B99\u0BC1":             "ṅu",
	"\u0B9A\u0BC1":             "cu",
	"\u0B9E\u0BC1":             "ñu",
	"\u0B9F\u0BC1":             "ṭu",
	"\u0BA3\u0BC1":             "ṇu",
	"\u0BA4\u0BC1":             "tu",
	"\u0BA8\u0BC1":             "nu",
	"\u0BAA\u0BC1":             "pu",
	"\u0BAE\u0BC1":             "mu",
	"\u0BAF\u0BC1":             "yu",
	"\u0BB0\u0BC1":             "ru",
	"\u0BB2\u0BC1":             "lu",
	"\u0BB5\u0BC1":             "vu",
	"\u0BB4\u0BC1":             "ḻu",
	"\u0BB3\u0BC1":             "ḷu",
	"\u0BB1\u0BC1":             "ṟu",
	"\u0BA9\u0BC1":             "ṉu",
	"\u0B95\u0BC2":             "kū",
	"\u0B99\u0BC2":             "ṅū",
	"\u0B9A\u0BC2":             "cū",
	"\u0B9E\u0BC2":             "ñū",
	"\u0B9F\u0BC2":             "ṭū",
	"\u0BA3\u0BC2":             "ṇū",
	"\u0BA4\u0BC2":             "tū",
	"\u0BA8\u0BC2":             "nū",
	"\u0BAA\u0BC2":             "pū",
	"\u0BAE\u0BC2":             "mū",
	"\u0BAF\u0BC2":             "yū",
	"\u0BB0\u0BC2":             "rū",
	"\u0BB2\u0BC2":             "lū",
	"\u0BB5\u0BC2":             "vū",
	"\u0BB4\u0BC2":             "ḻū",
	"\u0BB3\u0BC2":             "ḷū",
	"\u0BB1\u0BC2":             "ṟū",
	"\u0BA9\u0BC2":             "ṉū",
	"\u0B95\u0BC6":             "ke",
	"\u0B99\u0BC6":             "ṅe",
	"\u0B9A\u0BC6":             "ce",
	"\u0B9E\u0BC6":             "ñe",
	"\u0B9F\u0BC6":             "ṭe",
	"\u0BA3\u0BC6":             "ṇe",
	"\u0BA4\u0BC6":             "te",
	"\u0BA8\u0BC6":             "ne",
	"\u0BAA\u0BC6":             "pe",
	"\u0BAE\u0BC6":             "me",
	"\u0BAF\u0BC6":             "ye",
	"\u0BB0\u0BC6":             "re",
	"\u0BB2\u0BC6":             "le",
	"\u0BB5\u0BC6":             "ve",
	"\u0BB4\u0BC6":             "ḻe",
	"\u0BB3\u0BC6":             "ḷe",
	"\u0BB1\u0BC6":             "ṟe",
	"\u0BA9\u0BC6":             "ṉe",
	"\u0B95\u0BC7":             "kē",
	"\u0B99\u0BC7":             "ṅē",
	"\u0B9A\u0BC7":             "cē",
	"\u0B9E\u0BC7":             "ñē",
	"\u0B9F\u0BC7":             "ṭē",
	"\u0BA3\u0BC7":             "ṇē",
	"\u0BA4\u0BC7":             "tē",
	"\u0BA8\u0BC7":             "nē",
	"\u0BAA\u0BC7":             "pē",
	"\u0BAE\u0BC7":             "mē",
	"\u0BAF\u0BC7":             "yē",
	"\u0BB0\u0BC7":             "rē",
	"\u0BB2\u0BC7":             "lē",
	"\u0BB5\u0BC7":             "vē",
	"\u0BB4\u0BC7":             "ḻē",
	"\u0BB3\u0BC7":             "ḷē",
	"\u0BB1\u0BC7":             "ṟē",
	"\u0BA9\u0BC7":             "ṉē",
	"\u0B95\u0BC8":             "kai",
	"\u0B99\u0BC8":             "ṅai",
	"\u0B9A\u0BC8":             "cai",
	"\u0B9E\u0BC8":             "ñai",
	"\u0B9F\u0BC8":             "ṭai",
	"\u0BA3\u0BC8":             "ṇai",
	"\u0BA4\u0BC8":             "tai",
	"\u0BA8\u0BC8":             "nai",
	"\u0BAA\u0BC8":             "pai",
	"\u0BAE\u0BC8":             "mai",
	"\u0BAF\u0BC8":             "yai",
	"\u0BB0\u0BC8":             "rai",
	"\u0BB2\u0BC8":             "lai",
	"\u0BB5\u0BC8":             "vai",
	"\u0BB4\u0BC8":             "ḻai",
	"\u0BB3\u0BC8":             "ḷai",
	"\u0BB1\u0BC8":             "ṟai",
	"\u0BA9\u0BC8":             "ṉai",
	"\u0B95\u0BCA":             "ko",
	"\u0B99\u0BCA":             "ṅo",
	"\u0B9A\u0BCA":             "co",
	"\u0B9E\u0BCA":             "ño",
	"\u0B9F\u0BCA":             "ṭo",
	"\u0BA3\u0BCA":             "ṇo",
	"\u0BA4\u0BCA":             "to",
	"\u0BA8\u0BCA":             "no",
	"\u0BAA\u0BCA":             "po",
	"\u0BAE\u0BCA":             "mo",
	"\u0BAF\u0BCA":             "yo",
	"\u0BB0\u0BCA":             "ro",
	"\u0BB2\u0BCA":             "lo",
	"\u0BB5\u0BCA":             "vo",
	"\u0BB4\u0BCA":             "ḻo",
	"\u0BB3\u0BCA":             "ḷo",
	"\u0BB1\u0BCA":             "ṟo",
	"\u0BA9\u0BCA":             "ṉo",
	"\u0B95\u0BCB":             "kō",
	"\u0B99\u0BCB":             "ṅō",
	"\u0B9A\u0BCB":             "cō",
	"\u0B9E\u0BCB":             "ñō",
	"\u0B9F\u0BCB":             "ṭō",
	"\u0BA3\u0BCB":             "ṇō",
	"\u0BA4\u0BCB":             "tō",
	"\u0BA8\u0BCB":             "nō",
	"\u0BAA\u0BCB":             "pō",
	"\u0BAE\u0BCB":             "mō",
	"\u0BAF\u0BCB":             "yō",
	"\u0BB0\u0BCB":             "rō",
	"\u0BB2\u0BCB":             "lō",
	"\u0BB5\u0BCB":             "vō",
	"\u0BB4\u0BCB":             "ḻō",
	"\u0BB3\u0BCB":             "ḷō",
	"\u0BB1\u0BCB":             "ṟō",
	"\u0BA9\u0BCB":             "ṉō",
	"\u0B95\u0BCC":             "kau",
	"\u0B99\u0BCC":             "ṅau",
	"\u0B9A\u0BCC":             "cau",
	"\u0B9E\u0BCC":             "ñau",
	"\u0B9F\u0BCC":             "ṭau",
	"\u0BA3\u0BCC":             "ṇau",
	"\u0BA4\u0BCC":             "tau",
	"\u0BA8\u0BCC":             "nau",
	"\u0BAA\u0BCC":             "pau",
	"\u0BAE\u0BCC":             "mau",
	"\u0BAF\u0BCC":             "yau",
	"\u0BB0\u0BCC":             "rau",
	"\u0BB2\u0BCC":             "lau",
	"\u0BB5\u0BCC":             "vau",
	"\u0BB4\u0BCC":             "ḻau",
	"\u0BB3\u0BCC":             "ḷau",
	"\u0BB1\u0BCC":             "ṟau",
	"\u0BA9\u0BCC":             "ṉau",
}

var commonChars = map[rune]bool{
	'\u0027': true, // single quote
	'\u00A0': true, // non-breaking space
	' ':      true,
	'!':      true,
	'"':      true,
	'(':      true,
	')':      true,
	',':      true,
	'-':      true,
	'.':      true,
	':':      true,
	';':      true,
	'‘':      true,
	'’':      true,
	'“':      true,
	'”':      true,
	'?':      true,

	// Numerals
	'0': true,
	'1': true,
	'2': true,
	'3': true,
	'4': true,
	'5': true,
	'6': true,
	'7': true,
	'8': true,
	'9': true,
}

func (tree *rNode) isCommonChar(sym rune, alwaysAcceptASCII bool) bool {
	if _, ok := commonChars[sym]; ok {
		return true
	}
	if alwaysAcceptASCII && int(sym) < 128 {
		return true
	}
	return false
}

// TREE START

type rNode struct {
	r    rune
	daus map[rune]*rNode
	leaf string
}

type arc struct {
	start int
	end   int
	value string
}

func newNode() *rNode {
	return &rNode{daus: map[rune]*rNode{}}
}

func (rn *rNode) add(rs []rune, val string) {
	if len(rs) == 0 {
		return
	}

	r := rs[0]
	if dau, ok := rn.daus[r]; ok {
		if len(rs) == 1 {
			dau.leaf = val
		}
		dau.add(rs[1:], val)
	} else {
		dau := newNode()
		dau.r = r
		if len(rs) == 1 {
			dau.leaf = val
		}
		rn.daus[r] = dau
		dau.add(rs[1:], val)
	}
}

func prefix(tree *rNode, rs []rune) arc {
	var res arc

	t := tree

	for i, r := range rs {
		if n, ok := t.daus[r]; ok {
			res.value = n.leaf
			res.end = i + 1
			t = n
		} else {
			break
		}
	}

	return res
}

func (t Translit) reverseTest(tree *rNode, input string, mapped string) error {
	var thisRevTree *rNode
	if reflect.DeepEqual(*tree, *t.theTree) {
		thisRevTree = t.revTree
	} else if reflect.DeepEqual(*tree, *t.revTree) {
		thisRevTree = t.theTree
	} else {
		return fmt.Errorf("translit/reverseTest couldn't compare maptables")
	}
	remapped := t.translit(thisRevTree, []rune(mapped), false)
	if !remapped.OK {
		return fmt.Errorf("%s", strings.Join(remapped.Msgs, "; "))
	}
	if remapped.Result != input {
		return fmt.Errorf("reverse test failed: input '%s', mapped '%s', remapped '%s'", input, mapped, remapped.Result)
	}
	return nil
}

func (t Translit) translit(tree *rNode, rs []rune, doReverseTest bool) Result {
	var trans []string
	var unknown = []string{}
	var result = Result{OK: true, Input: string(rs), Msgs: []string{}}

	for i, n := 0, len(rs); i < n; {
		a := prefix(tree, rs[i:])
		if a.end > 0 {
			i = i + a.end
			trans = append(trans, a.value)
		} else {
			s := string(rs[i])
			if tree.isCommonChar(rs[i], t.alwaysAcceptASCII) {
				trans = append(trans, s)
			} else {
				trans = append(trans, t.defaultChar)
				result.OK = false
				if !translit.StringsContains(unknown, s) {
					unknown = append(unknown, s)
				}
			}
			i++
		}
	}
	result.Result = strings.Join(trans, "")
	if len(unknown) > 0 {
		pluralS := "s"
		if len(unknown) == 1 {
			pluralS = ""
		}
		result.Msgs = []string{fmt.Sprintf("unknown input symbol%s: %v", pluralS, strings.Join(unknown, ","))}
	} else if doReverseTest {
		err := t.reverseTest(tree, result.Input, result.Result)
		if err != nil {
			result.OK = false
			result.Msgs = append(result.Msgs, fmt.Sprintf("%v", err))
			return result
		}
	}

	return result
}

// TREE END

func NewTranslit() Translit {
	var theTree = newNode()
	var revTree = newNode()
	for k, v := range script2transMap {
		theTree.add([]rune(k), v)
		revTree.add([]rune(v), k)
	}
	return Translit{
		theTree:           theTree,
		revTree:           revTree,
		alwaysAcceptASCII: false,
		defaultChar:       "?",
	}
}

// Convert - transliterate from Tamil script to transliteration alphabet
func (t Translit) Convert(input string) Result {
	input = translit.NFC(input)
	return t.translit(t.theTree, []rune(input), false)
}

// ConvertDebug - transliterate from Tamil script to transliteration alphabet
func (t Translit) ConvertDebug(input string, debug bool) Result {
	input = translit.NFC(input)
	return t.translit(t.theTree, []rune(input), debug)
}

// Revert - transliterate from transliteration alphabet to Tamil script
func (t Translit) Revert(input string) Result {
	input = translit.NFC(input)
	return t.translit(t.revTree, []rune(input), false)
}

// RevertDebug - transliterate from transliteration alphabet to Tamil script
func (t Translit) RevertDebug(input string, debug bool) Result {
	input = translit.NFC(input)
	return t.translit(t.revTree, []rune(input), debug)
}

// var TranslitCharsRE = buildTranslitCharsRE()

// func buildTranslitCharsRE() *regexp.Regexp {
// 	chars := []string{}
// 	for _, v := range script2transMap {
// 		chars = append(chars, v)
// 	}
// 	return regexp.MustCompile(fmt.Sprintf("^(%v)+$", strings.Join(chars, "|")))
// }
