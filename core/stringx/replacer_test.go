package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplacerReplace(t *testing.T) {
	mapping := map[string]string{
		"一二三四": "1234",
		"二三":   "23",
		"二":    "2",
	}
	assert.Equal(t, "零1234五", NewReplacer(mapping).Replace(literal_9058))
}

func TestReplacerReplaceJumpMatch(t *testing.T) {
	mapping := map[string]string{
		"abcdeg": "ABCDEG",
		"cdef":   "CDEF",
		"cde":    "CDE",
	}
	assert.Equal(t, "abCDEF", NewReplacer(mapping).Replace("abcdef"))
}

func TestReplacerReplaceOverlap(t *testing.T) {
	mapping := map[string]string{
		"3d": "34",
		"bc": "23",
	}
	assert.Equal(t, "a234e", NewReplacer(mapping).Replace("abcde"))
}

func TestReplacerReplaceSingleChar(t *testing.T) {
	mapping := map[string]string{
		"二": "2",
	}
	assert.Equal(t, "零一2三四五", NewReplacer(mapping).Replace(literal_9058))
}

func TestReplacerReplaceExceedRange(t *testing.T) {
	mapping := map[string]string{
		"二三四五六": "23456",
	}
	assert.Equal(t, literal_9058, NewReplacer(mapping).Replace(literal_9058))
}

func TestReplacerReplacePartialMatch(t *testing.T) {
	mapping := map[string]string{
		"二三四七": "2347",
	}
	assert.Equal(t, literal_9058, NewReplacer(mapping).Replace(literal_9058))
}

func TestReplacerReplacePartialMatchEnds(t *testing.T) {
	mapping := map[string]string{
		"二三四七": "2347",
		"三四":   "34",
	}
	assert.Equal(t, "零一二34", NewReplacer(mapping).Replace("零一二三四"))
}

func TestReplacerReplaceMultiMatches(t *testing.T) {
	mapping := map[string]string{
		"二三": "23",
	}
	assert.Equal(t, "零一23四五一23四五", NewReplacer(mapping).Replace("零一二三四五一二三四五"))
}

func TestReplacerReplaceLongestMatching(t *testing.T) {
	keywords := map[string]string{
		"日本":    "japan",
		"日本的首都": "东京",
	}
	replacer := NewReplacer(keywords)
	assert.Equal(t, "东京在japan", replacer.Replace("日本的首都在日本"))
}

func TestReplacerReplaceSuffixMatch(t *testing.T) {
	// case1
	{
		keywords := map[string]string{
			"abcde": "ABCDE",
			"bcde":  "BCDE",
			"bcd":   "BCD",
		}
		assert.Equal(t, "aBCDf", NewReplacer(keywords).Replace("abcdf"))
	}
	// case2
	{
		keywords := map[string]string{
			"abcde": "ABCDE",
			"bcde":  "BCDE",
			"cde":   "CDE",
			"c":     "C",
			"cd":    "CD",
		}
		assert.Equal(t, "abCDf", NewReplacer(keywords).Replace("abcdf"))
	}
}

func TestReplacerReplaceLongestOverlap(t *testing.T) {
	keywords := map[string]string{
		"456":  "def",
		"abcd": "1234",
	}
	replacer := NewReplacer(keywords)
	assert.Equal(t, "123def7", replacer.Replace("abcd567"))
}

func TestReplacerReplaceLongestLonger(t *testing.T) {
	mapping := map[string]string{
		"c": "3",
	}
	assert.Equal(t, "3d", NewReplacer(mapping).Replace("cd"))
}

func TestReplacerReplaceJumpToFail(t *testing.T) {
	mapping := map[string]string{
		"bcdf": "1235",
		"cde":  "234",
	}
	assert.Equal(t, "ab234fg", NewReplacer(mapping).Replace("abcdefg"))
}

func TestReplacerReplaceJumpToFailDup(t *testing.T) {
	mapping := map[string]string{
		"bcdf": "1235",
		"ccde": "2234",
	}
	assert.Equal(t, "ab2234fg", NewReplacer(mapping).Replace("abccdefg"))
}

func TestReplacerReplaceJumpToFailEnding(t *testing.T) {
	mapping := map[string]string{
		"bcdf": "1235",
		"cdef": "2345",
	}
	assert.Equal(t, "ab2345", NewReplacer(mapping).Replace("abcdef"))
}

func TestReplacerReplaceEmpty(t *testing.T) {
	mapping := map[string]string{
		"bcdf": "1235",
		"cdef": "2345",
	}
	assert.Equal(t, "", NewReplacer(mapping).Replace(""))
}

func TestFuzzReplacerCase1(t *testing.T) {
	keywords := map[string]string{
		"yQyJykiqoh":     "xw",
		"tgN70z":         "Q2P",
		"tXKhEn":         "w1G8",
		"5nfOW1XZO":      "GN",
		"f4Ov9i9nHD":     "cT",
		"1ov9Q":          "Y",
		"7IrC9n":         "400i",
		"JQLxonpHkOjv":   "XI",
		"DyHQ3c7":        "Ygxux",
		"ffyqJi":         "u",
		"UHuvXrbD8pni":   "dN",
		"LIDzNbUlTX":     "g",
		"yN9WZh2rkc8Q":   "3U",
		"Vhk11rz8CObceC": "jf",
		"R0Rt4H2qChUQf":  "7U5M",
		"MGQzzPCVKjV9":   "yYz",
		"B5jUUl0u1XOY":   "l4PZ",
		"pdvp2qfLgG8X":   "BM562",
		"ZKl9qdApXJ2":    "T",
		"37jnugkSevU66":  "aOHFX",
	}
	rep := NewReplacer(keywords)
	text := "yjF8fyqJiiqrczOCVyoYbLvrMpnkj"
	val := rep.Replace(text)
	keys := rep.(*replacer).node.find([]rune(val))
	if len(keys) > 0 {
		t.Errorf("result: %s, match: %v", val, keys)
	}
}

func TestFuzzReplacerCase2(t *testing.T) {
	keywords := map[string]string{
		"dmv2SGZvq9Yz":   "TE",
		"rCL5DRI9uFP8":   "hvsc8",
		"7pSA2jaomgg":    "v",
		"kWSQvjVOIAxR":   "Oje",
		"hgU5bYYkD3r6":   "qCXu",
		"0eh6uI":         "MMlt",
		"3USZSl85EKeMzw": "Pc",
		"JONmQSuXa":      "dX",
		"EO1WIF":         "G",
		"uUmFJGVmacjF":   "1N",
		"DHpw7":          "M",
		"NYB2bm":         "CPya",
		"9FiNvBAHHNku5":  "7FlDE",
		"tJi3I4WxcY":     "q5",
		"sNJ8Z1ToBV0O":   "tl",
		"0iOg72QcPo":     "RP",
		"pSEqeL":         "5KZ",
		"GOyYqTgmvQ":     "9",
		"Qv4qCsj":        "nl52E",
		"wNQ5tOutYu5s8":  "6iGa",
	}
	rep := NewReplacer(keywords)
	text := "AoRxrdKWsGhFpXwVqMLWRL74OukwjBuBh0g7pSrk"
	val := rep.Replace(text)
	keys := rep.(*replacer).node.find([]rune(val))
	if len(keys) > 0 {
		t.Errorf("result: %s, match: %v", val, keys)
	}
}

func TestReplacerReplaceLongestMatch(t *testing.T) {
	replacer := NewReplacer(map[string]string{
		"日本的首都": "东京",
		"日本":    "本日",
	})
	assert.Equal(t, "东京是东京", replacer.Replace("日本的首都是东京"))
}

func TestReplacerReplaceIndefinitely(t *testing.T) {
	mapping := map[string]string{
		"日本的首都": "东京",
		"东京":    "日本的首都",
	}
	assert.NotPanics(t, func() {
		NewReplacer(mapping).Replace("日本的首都是东京")
	})
}

const literal_9058 = "零一二三四五"
