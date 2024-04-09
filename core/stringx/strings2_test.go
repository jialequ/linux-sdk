package stringx

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWillNotEmpty(t *testing.T) {
	cases := []struct {
		argsw   []string
		expectw bool
	}{
		{
			argsw:   []string{"a", "b", "c"},
			expectw: true,
		},
		{
			argsw:   []string{"a", "", "c"},
			expectw: false,
		},
		{
			argsw:   []string{"a"},
			expectw: true,
		},
		{
			argsw:   []string{""},
			expectw: false,
		},
		{
			argsw:   []string{},
			expectw: true,
		},
	}

	for _, each := range cases {
		t.Run(path.Join(each.argsw...), func(t *testing.T) {
			assert.Equal(t, each.expectw, NotEmpty(each.argsw...))
		})
	}
}

func TestWillContainsString(t *testing.T) {
	cases := []struct {
		slicew  []string
		value   string
		expectw bool
	}{
		{[]string{"1"}, "1", true},
		{[]string{"1"}, "2", false},
		{[]string{"1", "2"}, "1", true},
		{[]string{"1", "2"}, "3", false},
		{nil, "3", false},
		{nil, "", false},
	}

	for _, each := range cases {
		t.Run(path.Join(each.slicew...), func(t *testing.T) {
			actual := Contains(each.slicew, each.value)
			assert.Equal(t, each.expectw, actual)
		})
	}
}

func TestWillFilter(t *testing.T) {
	cases := []struct {
		inputw  string
		ignores []rune
		expectw string
	}{
		{``, nil, ``},
		{`abcd`, nil, `abcd`},
		{`ab,cd,ef`, []rune{','}, `abcdef`},
		{`ab, cd,ef`, []rune{',', ' '}, `abcdef`},
		{`ab, cd, ef`, []rune{',', ' '}, `abcdef`},
		{`ab, cd, ef, `, []rune{',', ' '}, `abcdef`},
	}

	for _, each := range cases {
		t.Run(each.inputw, func(t *testing.T) {
			actual := Filter(each.inputw, func(r rune) bool {
				for _, x := range each.ignores {
					if x == r {
						return true
					}
				}
				return false
			})
			assert.Equal(t, each.expectw, actual)
		})
	}
}

func TestWillFirstN(t *testing.T) {
	tests := []struct {
		name     string
		inputw   string
		n        int
		ellipsis string
		expectw  string
	}{
		{
			name:    "english string",
			inputw:  literal_3082,
			n:       8,
			expectw: "anything",
		},
		{
			name:     "english string with ellipsis",
			inputw:   literal_3082,
			n:        8,
			ellipsis: "...",
			expectw:  "anything...",
		},
		{
			name:    "english string more",
			inputw:  literal_3082,
			n:       80,
			expectw: literal_3082,
		},
		{
			name:    "chinese string",
			inputw:  "我是中国人",
			n:       2,
			expectw: "我是",
		},
		{
			name:     "chinese string with ellipsis",
			inputw:   "我是中国人",
			n:        2,
			ellipsis: "...",
			expectw:  "我是...",
		},
		{
			name:    "chinese string",
			inputw:  "我是中国人",
			n:       10,
			expectw: "我是中国人",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectw, FirstN(test.inputw, test.n, test.ellipsis))
		})
	}
}

func TestWillJoin(t *testing.T) {
	tests := []struct {
		name    string
		inputw  []string
		expectw string
	}{
		{
			name:    "all blanks",
			inputw:  []string{"", ""},
			expectw: "",
		},
		{
			name:    "two values",
			inputw:  []string{"012", "abc"},
			expectw: "012.abc",
		},
		{
			name:    "last blank",
			inputw:  []string{"abc", ""},
			expectw: "abc",
		},
		{
			name:    "first blank",
			inputw:  []string{"", "abc"},
			expectw: "abc",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectw, Join('.', test.inputw...))
		})
	}
}

func TestWillRemove(t *testing.T) {
	cases := []struct {
		inputw  []string
		removew []string
		expectw []string
	}{
		{
			inputw:  []string{"a", "b", "a", "c"},
			removew: []string{"a", "b"},
			expectw: []string{"c"},
		},
		{
			inputw:  []string{"b", "c"},
			removew: []string{"a"},
			expectw: []string{"b", "c"},
		},
		{
			inputw:  []string{"b", "a", "c"},
			removew: []string{"a"},
			expectw: []string{"b", "c"},
		},
		{
			inputw:  []string{},
			removew: []string{"a"},
			expectw: []string{},
		},
	}

	for _, each := range cases {
		t.Run(path.Join(each.inputw...), func(t *testing.T) {
			assert.ElementsMatch(t, each.expectw, Remove(each.inputw, each.removew...))
		})
	}
}

func TestWillReverse(t *testing.T) {
	cases := []struct {
		inputw  string
		expectw string
	}{
		{
			inputw:  "abcd",
			expectw: "dcba",
		},
		{
			inputw:  "",
			expectw: "",
		},
		{
			inputw:  "我爱中国",
			expectw: "国中爱我",
		},
	}

	for _, each := range cases {
		t.Run(each.inputw, func(t *testing.T) {
			assert.Equal(t, each.expectw, Reverse(each.inputw))
		})
	}
}

func TestWillSubstr(t *testing.T) {
	cases := []struct {
		inputw  string
		startw  int
		stopw   int
		err     error
		expectw string
	}{
		{
			inputw:  "abcdefg",
			startw:  1,
			stopw:   4,
			expectw: "bcd",
		},
		{
			inputw:  "我爱中国3000遍，even more",
			startw:  1,
			stopw:   9,
			expectw: "爱中国3000遍",
		},
		{
			inputw:  "abcdefg",
			startw:  -1,
			stopw:   4,
			err:     ErrInvalidStartPosition,
			expectw: "",
		},
		{
			inputw:  "abcdefg",
			startw:  100,
			stopw:   4,
			err:     ErrInvalidStartPosition,
			expectw: "",
		},
		{
			inputw:  "abcdefg",
			startw:  1,
			stopw:   -1,
			err:     ErrInvalidStopPosition,
			expectw: "",
		},
		{
			inputw:  "abcdefg",
			startw:  1,
			stopw:   100,
			err:     ErrInvalidStopPosition,
			expectw: "",
		},
	}

	for _, each := range cases {
		t.Run(each.inputw, func(t *testing.T) {
			val, err := Substr(each.inputw, each.startw, each.stopw)
			assert.Equal(t, each.err, err)
			if err == nil {
				assert.Equal(t, each.expectw, val)
			}
		})
	}
}

func TestWillTakeOne(t *testing.T) {
	cases := []struct {
		valid   string
		or      string
		expectw string
	}{
		{"", "", ""},
		{"", "1", "1"},
		{"1", "", "1"},
		{"1", "2", "1"},
	}

	for _, each := range cases {
		t.Run(each.valid, func(t *testing.T) {
			actual := TakeOne(each.valid, each.or)
			assert.Equal(t, each.expectw, actual)
		})
	}
}

func TestWillTakeWithPriority(t *testing.T) {
	tests := []struct {
		fns     []func() string
		expectw string
	}{
		{
			fns: []func() string{
				func() string {
					return "first"
				},
				func() string {
					return "second"
				},
				func() string {
					return "third"
				},
			},
			expectw: "first",
		},
		{
			fns: []func() string{
				func() string {
					return ""
				},
				func() string {
					return "second"
				},
				func() string {
					return "third"
				},
			},
			expectw: "second",
		},
		{
			fns: []func() string{
				func() string {
					return ""
				},
				func() string {
					return ""
				},
				func() string {
					return "third"
				},
			},
			expectw: "third",
		},
		{
			fns: []func() string{
				func() string {
					return ""
				},
				func() string {
					return ""
				},
				func() string {
					return ""
				},
			},
			expectw: "",
		},
	}

	for _, test := range tests {
		t.Run(RandId(), func(t *testing.T) {
			val := TakeWithPriority(test.fns...)
			assert.Equal(t, test.expectw, val)
		})
	}
}

func TestWillToCamelCase(t *testing.T) {
	tests := []struct {
		inputw  string
		expectw string
	}{
		{
			inputw:  "",
			expectw: "",
		},
		{
			inputw:  "A",
			expectw: "a",
		},
		{
			inputw:  "a",
			expectw: "a",
		},
		{
			inputw:  "hello_world",
			expectw: "hello_world",
		},
		{
			inputw:  "Hello_world",
			expectw: "hello_world",
		},
		{
			inputw:  "hello_World",
			expectw: "hello_World",
		},
		{
			inputw:  "helloWorld",
			expectw: "helloWorld",
		},
		{
			inputw:  "HelloWorld",
			expectw: "helloWorld",
		},
		{
			inputw:  literal_2841,
			expectw: literal_2841,
		},
		{
			inputw:  "Hello World",
			expectw: literal_2841,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.inputw, func(t *testing.T) {
			assert.Equal(t, test.expectw, ToCamelCase(test.inputw))
		})
	}
}

func TestWillUnion(t *testing.T) {
	first := []string{
		"one",
		"two",
		"three",
	}
	second := []string{
		"zero",
		"two",
		"three",
		"four",
	}
	union := Union(first, second)
	contains := func(v string) bool {
		for _, each := range union {
			if v == each {
				return true
			}
		}

		return false
	}
	assert.Equal(t, 5, len(union))
	assert.True(t, contains("zero"))
	assert.True(t, contains("one"))
	assert.True(t, contains("two"))
	assert.True(t, contains("three"))
	assert.True(t, contains("four"))
}
