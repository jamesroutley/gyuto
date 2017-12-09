package read

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLex(t *testing.T) {
	testCases := []struct {
		input    string
		expected []item
	}{
		{
			input: "'(hello world)",
			expected: []item{
				{
					typ: itemQuote,
					val: "'",
					col: 1,
				},
				{
					typ: itemLeftParen,
					val: "(",
					col: 2,
				},
				{
					typ: itemAtom,
					val: "hello",
					col: 3,
				},
				{
					typ: itemAtom,
					val: "world",
					col: 9,
				},
				{
					typ: itemRightParen,
					val: ")",
					col: 14,
				},
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			_, items := lex("test", tc.input)
			tokens := itemChanToSlice(items)
			assert.Equal(t, tc.expected, tokens)
		})
	}
}

func TestIsAtom(t *testing.T) {
	testCases := []struct {
		input    rune
		expected bool
	}{
		{'a', true},
		{'\n', false},
		{'\r', false},
		{'\t', false},
		{' ', false},
		{eof, false},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("%s: %v", string(tc.input), tc.expected)
		t.Run(name, func(t *testing.T) {
			actual := isAtom(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
