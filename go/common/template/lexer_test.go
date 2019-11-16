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
	iEOF                = mkItem(itemEOF, "")
	iError              = mkItem(itemError, "")
	iTemplateLeftDelim  = mkItem(itemTemplateLeftDelim, "")
	iTemplateRightDelim = mkItem(itemTemplateRightDelim, "")
	iCommandLeftDelim  = mkItem(itemCommandLeftDelim, "")
	iCommandRightDelim = mkItem(itemCommandRightDelim, "")
)

// collectItems gathers the emitted items into a slice.
func collectItems(t *lexTest) (items []item) {
	l := lex(t.input)
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}
	return
}

func itemsEqual(i1, i2 []item) bool {
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
	lexTests := []lexTest{
		{
			"empty",
			"",
			[]item{
				iEOF,
			},
		},
		{
			"only text",
			"simple text",
			[]item{
				mkItem(itemText, "simple text"),
				iEOF,
			},
		},
		{
			"simple key between text",
			"text {{ key }} more",
			[]item{
				mkItem(itemText, "text "),
				iTemplateLeftDelim,
				mkItem(itemIdentifier, "key"),
				iTemplateRightDelim,
				mkItem(itemText, " more"),
				iEOF,
			},
		},
		{
			"simple command",
			"{% command %}",
			[]item{
				iCommandLeftDelim,
				mkItem(itemCommand, "command"),
				iCommandRightDelim,
				iEOF,
			},
		},
		{
			"simple command between text",
			"text {% command %} more",
			[]item{
				mkItem(itemText, "text "),
				iCommandLeftDelim,
				mkItem(itemCommand, "command"),
				iCommandRightDelim,
				mkItem(itemText, " more"),
				iEOF,
			},
		},
		{
			"simple command and template between text",
			"text {% command %} more {{ key }}",
			[]item{
				mkItem(itemText, "text "),
				iCommandLeftDelim,
				mkItem(itemCommand, "command"),
				iCommandRightDelim,
				mkItem(itemText, " more "),
				iTemplateLeftDelim,
				mkItem(itemIdentifier, "key"),
				iTemplateRightDelim,
				iEOF,
			},
		},
		{
			"complex key no text",
			"{{ one.two[12].three }}",
			[]item{
				iTemplateLeftDelim,
				mkItem(itemIdentifier, "one.two[12].three"),
				iTemplateRightDelim,
				iEOF,
			},
		},
		{
			"multiple complex keys no text",
			"{{ one.two[12].three }} {{ four.five[6].seven }}",
			[]item{
				iTemplateLeftDelim,
				mkItem(itemIdentifier, "one.two[12].three"),
				iTemplateRightDelim,
				mkItem(itemText, " "),
				iTemplateLeftDelim,
				mkItem(itemIdentifier, "four.five[6].seven"),
				iTemplateRightDelim,
				iEOF,
			},
		},
		{
			"multiple complex keys between text",
			"some {{ one.two[12].three }} random {{ four.five[6].seven }} text",
			[]item{
				mkItem(itemText, "some "),
				iTemplateLeftDelim,
				mkItem(itemIdentifier, "one.two[12].three"),
				iTemplateRightDelim,
				mkItem(itemText, " random "),
				iTemplateLeftDelim,
				mkItem(itemIdentifier, "four.five[6].seven"),
				iTemplateRightDelim,
				mkItem(itemText, " text"),
				iEOF,
			},
		},

		{
			"invalid key",
			"{{ asd!() }}",
			[]item{
				iTemplateLeftDelim,
				mkItem(itemIdentifier, "asd"),
				iError,
			},
		},
		{
			"invalid key with text",
			"text {{ asd!() }}",
			[]item{
				mkItem(itemText, "text "),
				iTemplateLeftDelim,
				mkItem(itemIdentifier, "asd"),
				iError,
			},
		},
		{
			"missing key with text",
			"text {{  }}",
			[]item{
				mkItem(itemText, "text "),
				iTemplateLeftDelim,
				iError,
			},
		},
	}

	for _, test := range lexTests {
		items := collectItems(&test)
		if !itemsEqual(items, test.items) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", test.name, items, test.items)
		}
	}
}
