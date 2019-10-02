package template

import (
	"testing"
)

type lexTest struct {
	name  string
	input string
	items []item
}

func mkItem(typ itemType, text string) item {
	return item{
		typ: typ,
		val: text,
	}
}

var (
	iEOF        = mkItem(itemEOF, "")
	iError      = mkItem(itemError, "")
	iLeftDelim  = mkItem(itemLeftDelim, "")
	iRightDelim = mkItem(itemRightDelim, "")
)

var lexTests = []lexTest{
	{"empty", "", []item{iEOF}},
	{"only text", "simple text", []item{mkItem(itemText, "simple text"), iEOF}},
	{"simple key between text", "text {{ key }} more", []item{mkItem(itemText, "text "), iLeftDelim, mkItem(itemIdentifier, "key"), iRightDelim, mkItem(itemText, " more"), iEOF}},
	{"complex key no text", "{{ one.two[12].three }}", []item{iLeftDelim, mkItem(itemIdentifier, "one.two[12].three"), iRightDelim, iEOF}},
	{"multiple complex keys no text", "{{ one.two[12].three }} {{ four.five[6].seven }}", []item{iLeftDelim, mkItem(itemIdentifier, "one.two[12].three"), iRightDelim, mkItem(itemText, " "), iLeftDelim, mkItem(itemIdentifier, "four.five[6].seven"), iRightDelim, iEOF}},
	{"multiple complex keys between text", "some {{ one.two[12].three }} random {{ four.five[6].seven }} text", []item{mkItem(itemText, "some "), iLeftDelim, mkItem(itemIdentifier, "one.two[12].three"), iRightDelim, mkItem(itemText, " random "), iLeftDelim, mkItem(itemIdentifier, "four.five[6].seven"), iRightDelim, mkItem(itemText, " text"), iEOF}},

	{"invalid key", "{{ asd!() }}", []item{iLeftDelim, mkItem(itemIdentifier, "asd"), iError}},
	{"invalid key with text", "text {{ asd!() }}", []item{mkItem(itemText, "text "), iLeftDelim, mkItem(itemIdentifier, "asd"), iError}},
	{"missing key with text", "text {{  }}", []item{mkItem(itemText, "text "), iLeftDelim, iError}},
}

// collect gathers the emitted items into a slice.
func collect(t *lexTest) (items []item) {
	l := lex(t.input, leftDelim, rightDelim)
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}
	return
}

func equal(i1, i2 []item, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			return false
		}
		if i1[k].val != i2[k].val && i1[k].typ != itemError {
			return false
		}
	}
	return true
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		items := collect(&test)
		if !equal(items, test.items, false) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", test.name, items, test.items)
		}
	}
}
