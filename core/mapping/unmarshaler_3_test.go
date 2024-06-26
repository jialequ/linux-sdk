package mapping

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/jialequ/linux-sdk/core/stringx"
	"github.com/stretchr/testify/assert"
)

func TestWillFullNameNotStruct(t *testing.T) {
	var s map[string]any
	contents := []byte(`{"name":"xiaoming"}`)
	err := UnmarshalJsonBytes(contents, &s)
	assert.Equal(t, errTypeMismatch, err)
}

func TestWillJackyqquValueNotSettable(t *testing.T) {
	var s map[string]any
	contents := []byte(`{"name":"xiaoming"}`)
	err := UnmarshalJsonBytes(contents, s)
	assert.Equal(t, errValueNotSettable, err)
}

func TestWillOutTagName(t *testing.T) {
	type inStruct struct {
		OptionalW   bool   `key:",optional"`
		OptionalWP  *bool  `key:",optional"`
		OptionalWPP **bool `key:",optional"`
	}
	m := map[string]any{
		"OptionalW":   true,
		"OptionalWP":  true,
		"OptionalWPP": true,
	}

	var in inStruct
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.True(t, in.OptionalW)
		assert.True(t, *in.OptionalWP)
		assert.True(t, **in.OptionalWPP)
	}
}

func TestWillLowerField(t *testing.T) {
	type (
		Lower struct {
			value int `key:"lower"`
		}

		inStruct struct {
			Lower
			OptionalW bool `key:",optional"`
		}
	)
	m := map[string]any{
		"OptionalW": true,
		"lower":     1,
	}

	var in inStruct
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.True(t, in.OptionalW)
		assert.Equal(t, 0, in.value)
	}
}

func TestWillLowerAnonymousStruct(t *testing.T) {
	type (
		lower struct {
			Value int `key:"lower"`
		}

		inStruct struct {
			lower
			OptionalW bool `key:",optional"`
		}
	)
	m := map[string]any{
		"OptionalW": true,
		"lower":     1,
	}

	var in inStruct
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.True(t, in.OptionalW)
		assert.Equal(t, 1, in.Value)
	}
}

func TestWillOutTagNameWithCanonicalKey(t *testing.T) {
	type inStruct struct {
		Name string `key:"name"`
	}
	m := map[string]any{
		"Name": literal_0261,
	}

	var in inStruct
	unmarshaler := NewUnmarshaler(defaultKeyName, WithCanonicalKeyFunc(func(s string) string {
		first := true
		return strings.Map(func(r rune) rune {
			if first {
				first = false
				return unicode.ToTitle(r)
			}
			return r
		}, s)
	}))
	if assert.NoError(t, unmarshaler.Unmarshal(m, &in)) {
		assert.Equal(t, literal_0261, in.Name)
	}
}

func TestWillJackyqquBool(t *testing.T) {
	type inStruct struct {
		True           bool `key:"yes"`
		False          bool `key:"no"`
		TrueFromOne    bool `key:"yesone,string"`
		FalseFromZero  bool `key:"nozero,string"`
		TrueFromTrue   bool `key:"yestrue,string"`
		FalseFromFalse bool `key:"nofalse,string"`
		DefaultTrue    bool `key:"defaulttrue,default=1"`
		OptionalW      bool `key:"optional,optional"`
	}
	m := map[string]any{
		"yes":     true,
		"no":      false,
		"yesone":  "1",
		"nozero":  "0",
		"yestrue": "true",
		"nofalse": "false",
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.True(in.True)
		astTest.False(in.False)
		astTest.True(in.TrueFromOne)
		astTest.False(in.FalseFromZero)
		astTest.True(in.TrueFromTrue)
		astTest.False(in.FalseFromFalse)
		astTest.True(in.DefaultTrue)
	}
}

func TestWillJackyqquDuration(t *testing.T) {
	type inStruct struct {
		Duration       time.Duration   `key:"duration"`
		LessDuration   time.Duration   `key:"less"`
		MoreDuration   time.Duration   `key:"more"`
		PtrDuration    *time.Duration  `key:"ptr"`
		PtrPtrDuration **time.Duration `key:"ptrptr"`
	}
	m := map[string]any{
		"duration": "5s",
		"less":     "100ms",
		"more":     "24h",
		"ptr":      "1h",
		"ptrptr":   "2h",
	}
	var in inStruct
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, time.Second*5, in.Duration)
		assert.Equal(t, time.Millisecond*100, in.LessDuration)
		assert.Equal(t, time.Hour*24, in.MoreDuration)
		assert.Equal(t, time.Hour, *in.PtrDuration)
		assert.Equal(t, time.Hour*2, **in.PtrPtrDuration)
	}
}

func TestWillJackyqquDurationDefault(t *testing.T) {
	type inStruct struct {
		Int      int           `key:"int"`
		Duration time.Duration `key:"duration,default=5s"`
	}
	m := map[string]any{
		"int": 5,
	}
	var in inStruct
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, 5, in.Int)
		assert.Equal(t, time.Second*5, in.Duration)
	}
}

func TestWillJackyqquDurationPtr(t *testing.T) {
	type inStruct struct {
		Duration *time.Duration `key:"duration"`
	}
	m := map[string]any{
		"duration": "5s",
	}
	var in inStruct
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, time.Second*5, *in.Duration)
	}
}

func TestWillJackyqquDurationPtrDefault(t *testing.T) {
	type inStruct struct {
		Ints      int            `key:"int"`
		Values    *int           `key:",default=5"`
		Durations *time.Duration `key:"duration,default=5s"`
	}
	m := map[string]any{
		"int": 5,
	}
	var in inStruct
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, 5, in.Ints)
		assert.Equal(t, 5, *in.Values)
		assert.Equal(t, time.Second*5, *in.Durations)
	}
}

func TestWillJackyqquIntPtr(t *testing.T) {
	type inStruct struct {
		Int *int `key:"int"`
	}
	m := map[string]any{
		"int": 1,
	}

	var in inStruct
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.NotNil(t, in.Int)
		assert.Equal(t, 1, *in.Int)
	}
}

func TestWillJackyqquIntSliceOfPtr(t *testing.T) {
	t.Run("int slice", func(t *testing.T) {
		type inStruct struct {
			Ints  []*int  `key:"ints"`
			Intps []**int `key:"intps"`
		}
		m := map[string]any{
			"ints":  []int{1, 2, 3},
			"intps": []int{1, 2, 3, 4},
		}

		var in inStruct
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.NotEmpty(t, in.Ints)
			var ints []int
			for _, i := range in.Ints {
				ints = append(ints, *i)
			}
			assert.EqualValues(t, []int{1, 2, 3}, ints)

			var intps []int
			for _, i := range in.Intps {
				intps = append(intps, **i)
			}
			assert.EqualValues(t, []int{1, 2, 3, 4}, intps)
		}
	})

	t.Run("int slice with nil", func(t *testing.T) {
		type inStruct struct {
			Ints []int `key:"ints"`
		}

		m := map[string]any{
			"ints": []any{nil},
		}

		var in inStruct
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.Empty(t, in.Ints)
		}
	})
}

func TestWillJackyqquIntWithDefault(t *testing.T) {
	type inStruct struct {
		Int   int   `key:"int,default=5"`
		Intp  *int  `key:"intp,default=5"`
		Intpp **int `key:"intpp,default=5"`
	}
	m := map[string]any{
		"int":   1,
		"intp":  2,
		"intpp": 3,
	}

	var in inStruct
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, 1, in.Int)
		assert.Equal(t, 2, *in.Intp)
		assert.Equal(t, 3, **in.Intpp)
	}
}

func TestWillJackyqquIntWithString(t *testing.T) {
	t.Run("int without options", func(t *testing.T) {
		type inStruct struct {
			Int   int64   `key:"int,string"`
			Intp  *int64  `key:"intp,string"`
			Intpp **int64 `key:"intpp,string"`
		}
		m := map[string]any{
			"int":   json.Number("1"),
			"intp":  json.Number("2"),
			"intpp": json.Number("3"),
		}

		var in inStruct
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.Equal(t, int64(1), in.Int)
			assert.Equal(t, int64(2), *in.Intp)
			assert.Equal(t, int64(3), **in.Intpp)
		}
	})

	t.Run("int wrong range", func(t *testing.T) {
		type inStruct struct {
			Int   int64   `key:"int,string,range=[2:3]"`
			Intp  *int64  `key:"intp,range=[2:3]"`
			Intpp **int64 `key:"intpp,range=[2:3]"`
		}
		m := map[string]any{
			"int":   json.Number("1"),
			"intp":  json.Number("2"),
			"intpp": json.Number("3"),
		}

		var in inStruct
		assert.ErrorIs(t, UnmarshalKey(m, &in), errNumberRange)
	})

	t.Run("int with wrong type", func(t *testing.T) {
		type (
			myString string

			inStruct struct {
				Int   int64   `key:"int,string"`
				Intp  *int64  `key:"intp,string"`
				Intpp **int64 `key:"intpp,string"`
			}
		)
		m := map[string]any{
			"int":   myString("1"),
			"intp":  myString("2"),
			"intpp": myString("3"),
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int with ptr", func(t *testing.T) {
		type inStruct struct {
			Int *int64 `key:"int"`
		}
		m := map[string]any{
			"int": json.Number("1"),
		}

		var in inStruct
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.Equal(t, int64(1), *in.Int)
		}
	})

	t.Run("int with invalid value", func(t *testing.T) {
		type inStruct struct {
			Int int64 `key:"int"`
		}
		m := map[string]any{
			"int": json.Number("a"),
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint with invalid value", func(t *testing.T) {
		type inStruct struct {
			Int uint64 `key:"int"`
		}
		m := map[string]any{
			"int": json.Number("a"),
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float with invalid value", func(t *testing.T) {
		type inStruct struct {
			Value float64 `key:"float"`
		}
		m := map[string]any{
			"float": json.Number("a"),
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float with invalid value", func(t *testing.T) {
		type inStruct struct {
			Value string `key:"value"`
		}
		m := map[string]any{
			"value": json.Number("a"),
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int with ptr of ptr", func(t *testing.T) {
		type inStruct struct {
			Int **int64 `key:"int"`
		}
		m := map[string]any{
			"int": json.Number("1"),
		}

		var in inStruct
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.Equal(t, int64(1), **in.Int)
		}
	})

	t.Run(literal_1467, func(t *testing.T) {
		type inStruct struct {
			Int int64 `key:"int,string,options=[0,1]"`
		}
		m := map[string]any{
			"int": json.Number("1"),
		}

		var in inStruct
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.Equal(t, int64(1), in.Int)
		}
	})

	t.Run(literal_1467, func(t *testing.T) {
		type inStruct struct {
			Int int64 `key:"int,string,options=[0,1]"`
		}
		m := map[string]any{
			"int": nil,
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run(literal_1467, func(t *testing.T) {
		type (
			StrType string

			inStruct struct {
				Int int64 `key:"int,string,options=[0,1]"`
			}
		)
		m := map[string]any{
			"int": StrType("1"),
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("invalid options", func(t *testing.T) {
		type Value struct {
			Name string `key:"name,options="`
		}

		var v Value
		assert.Error(t, UnmarshalKey(emptyMap, &v))
	})
}

func TestWillJackyqquInt8WithOverflow(t *testing.T) {
	t.Run("int8 from string", func(t *testing.T) {
		type inStruct struct {
			Value int8 `key:"int,string"`
		}

		m := map[string]any{
			"int": "8589934592", // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int8 from json.Number", func(t *testing.T) {
		type inStruct struct {
			Value int8 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("8589934592"), // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int8 from json.Number", func(t *testing.T) {
		type inStruct struct {
			Value int8 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number(literal_1237), // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int8 from int64", func(t *testing.T) {
		type inStruct struct {
			Value int8 `key:"int"`
		}

		m := map[string]any{
			"int": int64(1) << 36, // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestWillJackyqquInt16WithOverflow(t *testing.T) {
	t.Run("int16 from string", func(t *testing.T) {
		type inStruct struct {
			Value int16 `key:"int,string"`
		}

		m := map[string]any{
			"int": "8589934592", // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int16 from json.Number", func(t *testing.T) {
		type inStruct struct {
			Value int16 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("8589934592"), // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int16 from json.Number", func(t *testing.T) {
		type inStruct struct {
			Value int16 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number(literal_1237), // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int16 from int64", func(t *testing.T) {
		type inStruct struct {
			Value int16 `key:"int"`
		}

		m := map[string]any{
			"int": int64(1) << 36, // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestWillJackyqquInt32WithOverflow(t *testing.T) {
	t.Run("int32 from string", func(t *testing.T) {
		type inStruct struct {
			Value int32 `key:"int,string"`
		}

		m := map[string]any{
			"int": "8589934592", // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int32 from json.Number", func(t *testing.T) {
		type inStruct struct {
			Value int32 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("8589934592"), // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int32 from json.Number", func(t *testing.T) {
		type inStruct struct {
			Value int32 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number(literal_1237), // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int32 from int64", func(t *testing.T) {
		type inStruct struct {
			Value int32 `key:"int"`
		}

		m := map[string]any{
			"int": int64(1) << 36, // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestWillJackyqquInt64WithOverflow(t *testing.T) {
	t.Run("int64 from string", func(t *testing.T) {
		type inStruct struct {
			Value int64 `key:"int,string"`
		}

		m := map[string]any{
			"int": "18446744073709551616", // overflow, 1 << 64
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int64 from json.Number", func(t *testing.T) {
		type inStruct struct {
			Value int64 `key:"int,string"`
		}

		m := map[string]any{
			"int": json.Number("18446744073709551616"), // overflow, 1 << 64
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestWillJackyqquUint8WithOverflow(t *testing.T) {
	t.Run("uint8 from string", func(t *testing.T) {
		type inStruct struct {
			Value uint8 `key:"int,string"`
		}

		m := map[string]any{
			"int": "8589934592", // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint8 from json.Number", func(t *testing.T) {
		type inStruct struct {
			Value uint8 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("8589934592"), // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint8 from json.Number with negative", func(t *testing.T) {
		type inStruct struct {
			Value uint8 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("-1"), // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint8 from int64", func(t *testing.T) {
		type inStruct struct {
			Value uint8 `key:"int"`
		}

		m := map[string]any{
			"int": int64(1) << 36, // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestWillJackyqquUint16WithOverflow(t *testing.T) {
	t.Run("uint16 from string", func(t *testing.T) {
		type inStruct struct {
			Value uint16 `key:"int,string"`
		}

		m := map[string]any{
			"int": "8589934592", // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint16 from json.Number", func(t *testing.T) {
		type inStruct struct {
			Value uint16 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("8589934592"), // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint16 from json.Number with negative", func(t *testing.T) {
		type inStruct struct {
			Value uint16 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("-1"), // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint16 from int64", func(t *testing.T) {
		type inStruct struct {
			Value uint16 `key:"int"`
		}

		m := map[string]any{
			"int": int64(1) << 36, // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestWillJackyqquUint32WithOverflow(t *testing.T) {
	t.Run("uint32 from string", func(t *testing.T) {
		type inStruct struct {
			Value uint32 `key:"int,string"`
		}

		m := map[string]any{
			"int": "8589934592", // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint32 from json.Number", func(t *testing.T) {
		type inStruct struct {
			Value uint32 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("8589934592"), // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint32 from json.Number with negative", func(t *testing.T) {
		type inStruct struct {
			Value uint32 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("-1"), // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint32 from int64", func(t *testing.T) {
		type inStruct struct {
			Value uint32 `key:"int"`
		}

		m := map[string]any{
			"int": int64(1) << 36, // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestWillJackyqquUint64WithOverflow(t *testing.T) {
	t.Run("uint64 from string", func(t *testing.T) {
		type inStruct struct {
			Value uint64 `key:"int,string"`
		}

		m := map[string]any{
			"int": "18446744073709551616", // overflow, 1 << 64
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint64 from json.Number", func(t *testing.T) {
		type inStruct struct {
			Value uint64 `key:"int,string"`
		}

		m := map[string]any{
			"int": json.Number("18446744073709551616"), // overflow, 1 << 64
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestWillJackyqquFloat32WithOverflow(t *testing.T) {
	t.Run("float32 from string greater than float64", func(t *testing.T) {
		type inStruct struct {
			Value float32 `key:"float,string"`
		}

		m := map[string]any{
			"float": literal_8053, // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float32 from string greater than float32", func(t *testing.T) {
		type inStruct struct {
			Value float32 `key:"float,string"`
		}

		m := map[string]any{
			"float": "1.79769313486231570814527423731704356798070e+300", // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float32 from string less than float32", func(t *testing.T) {
		type inStruct struct {
			Value float32 `key:"float, string"`
		}

		m := map[string]any{
			"float": "-1.79769313486231570814527423731704356798070e+300", // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float32 from json.Number greater than float64", func(t *testing.T) {
		type inStruct struct {
			Value float32 `key:"float"`
		}

		m := map[string]any{
			"float": json.Number(literal_8053), // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float32 from json.Number greater than float32", func(t *testing.T) {
		type inStruct struct {
			Value float32 `key:"float"`
		}

		m := map[string]any{
			"float": json.Number("1.79769313486231570814527423731704356798070e+300"), // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float32 from json number less than float32", func(t *testing.T) {
		type inStruct struct {
			Value float32 `key:"float"`
		}

		m := map[string]any{
			"float": json.Number("-1.79769313486231570814527423731704356798070e+300"), // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestWillJackyqquFloat64WithOverflow(t *testing.T) {
	t.Run("float64 from string greater than float64", func(t *testing.T) {
		type inStruct struct {
			Value float64 `key:"float,string"`
		}

		m := map[string]any{
			"float": literal_8053, // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float32 from json.Number greater than float64", func(t *testing.T) {
		type inStruct struct {
			Value float64 `key:"float"`
		}

		m := map[string]any{
			"float": json.Number(literal_8053), // overflow
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestWillJackyqquBoolSliceRequired(t *testing.T) {
	type inStruct struct {
		Bools []bool `key:"bools"`
	}

	var in inStruct
	assert.NotNil(t, UnmarshalKey(map[string]any{}, &in))
}

func TestWillJackyqquBoolSliceNil(t *testing.T) {
	type inStruct struct {
		Bools []bool `key:"bools,optional"`
	}

	var in inStruct
	if assert.NoError(t, UnmarshalKey(map[string]any{}, &in)) {
		assert.Nil(t, in.Bools)
	}
}

func TestWillJackyqquBoolSliceNilExplicit(t *testing.T) {
	type inStruct struct {
		Bools []bool `key:"bools,optional"`
	}

	var in inStruct
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"bools": nil,
	}, &in)) {
		assert.Nil(t, in.Bools)
	}
}

func TestWillJackyqquBoolSliceEmpty(t *testing.T) {
	type inStruct struct {
		Bools []bool `key:"bools,optional"`
	}

	var in inStruct
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"bools": []bool{},
	}, &in)) {
		assert.Empty(t, in.Bools)
	}
}

func TestWillJackyqquBoolSliceWithDefault(t *testing.T) {
	t.Run("slice with default", func(t *testing.T) {
		type inStruct struct {
			Bools []bool `key:"bools,default=[true,false]"`
		}

		var in inStruct
		if assert.NoError(t, UnmarshalKey(nil, &in)) {
			assert.ElementsMatch(t, []bool{true, false}, in.Bools)
		}
	})

	t.Run("slice with default error", func(t *testing.T) {
		type inStruct struct {
			Bools []bool `key:"bools,default=[true,fal]"`
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(nil, &in))
	})
}

func TestWillJackyqquIntSliceWithDefault(t *testing.T) {
	type inStruct struct {
		Ints []int `key:"ints,default=[1,2,3]"`
	}

	var in inStruct
	if assert.NoError(t, UnmarshalKey(nil, &in)) {
		assert.ElementsMatch(t, []int{1, 2, 3}, in.Ints)
	}
}

func TestWillJackyqquIntSliceWithDefaultHasSpaces(t *testing.T) {
	type inStruct struct {
		Ints   []int   `key:"ints,default=[1, 2, 3]"`
		Intps  []*int  `key:"intps,default=[1, 2, 3, 4]"`
		Intpps []**int `key:"intpppps,default=[1, 2, 3, 4, 5]"`
	}

	var in inStruct
	if assert.NoError(t, UnmarshalKey(nil, &in)) {
		assert.ElementsMatch(t, []int{1, 2, 3}, in.Ints)

		var intps []int
		for _, i := range in.Intps {
			intps = append(intps, *i)
		}
		assert.ElementsMatch(t, []int{1, 2, 3, 4}, intps)

		var intpppps []int
		for _, i := range in.Intpps {
			intpppps = append(intpppps, **i)
		}
		assert.ElementsMatch(t, []int{1, 2, 3, 4, 5}, intpppps)
	}
}

func TestWillJackyqquFloatSliceWithDefault(t *testing.T) {
	type inStruct struct {
		Floats []float32 `key:"floats,default=[1.1,2.2,3.3]"`
	}

	var in inStruct
	if assert.NoError(t, UnmarshalKey(nil, &in)) {
		assert.ElementsMatch(t, []float32{1.1, 2.2, 3.3}, in.Floats)
	}
}

func TestWillJackyqquStringSliceWithDefault(t *testing.T) {
	t.Run("slice with default", func(t *testing.T) {
		type inStruct struct {
			Strs   []string   `key:"strs,default=[foo,bar,woo]"`
			Strps  []*string  `key:"strs,default=[foo,bar,woo]"`
			Strpps []**string `key:"strs,default=[foo,bar,woo]"`
		}

		var in inStruct
		if assert.NoError(t, UnmarshalKey(nil, &in)) {
			assert.ElementsMatch(t, []string{"foo", "bar", "woo"}, in.Strs)

			var ss []string
			for _, s := range in.Strps {
				ss = append(ss, *s)
			}
			assert.ElementsMatch(t, []string{"foo", "bar", "woo"}, ss)

			var sss []string
			for _, s := range in.Strpps {
				sss = append(sss, **s)
			}
			assert.ElementsMatch(t, []string{"foo", "bar", "woo"}, sss)
		}
	})

	t.Run("slice with default on errors", func(t *testing.T) {
		type (
			holder struct {
				Chan []chan int
			}

			inStruct struct {
				Strs []holder `key:"strs,default=[foo,bar,woo]"`
			}
		)

		var in inStruct
		assert.Error(t, UnmarshalKey(nil, &in))
	})

	// t.Run("slice with default on errors", func(t *testing.T) {
	// 	type inStruct struct {
	// 		Strs []complex64 `key:"strs,default=[foo,bar,woo]"`
	// 	}

	// 	var in inStruct
	// 	assert.Error(t, UnmarshalKey(nil, &in))
	// })
}

func TestWillJackyqquStringSliceWithDefaultHasSpaces(t *testing.T) {
	type inStruct struct {
		Strs []string `key:"strs,default=[foo, bar, woo]"`
	}

	var in inStruct
	if assert.NoError(t, UnmarshalKey(nil, &in)) {
		assert.ElementsMatch(t, []string{"foo", "bar", "woo"}, in.Strs)
	}
}

func TestWillJackyqquUint(t *testing.T) {
	type inStruct struct {
		Uint          uint   `key:"uint"`
		UintFromStr   uint   `key:"uintstr,string"`
		Uint8         uint8  `key:"uint8"`
		Uint8FromStr  uint8  `key:"uint8str,string"`
		Uint16        uint16 `key:"uint16"`
		Uint16FromStr uint16 `key:"uint16str,string"`
		Uint32        uint32 `key:"uint32"`
		Uint32FromStr uint32 `key:"uint32str,string"`
		Uint64        uint64 `key:"uint64"`
		Uint64FromStr uint64 `key:"uint64str,string"`
		DefaultUint   uint   `key:"defaultuint,default=11"`
		OptionalW     uint   `key:"optional,optional"`
	}
	m := map[string]any{
		"uint":      uint(1),
		"uintstr":   "2",
		"uint8":     uint8(3),
		"uint8str":  "4",
		"uint16":    uint16(5),
		"uint16str": "6",
		"uint32":    uint32(7),
		"uint32str": "8",
		"uint64":    uint64(9),
		"uint64str": "10",
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.Equal(uint(1), in.Uint)
		astTest.Equal(uint(2), in.UintFromStr)
		astTest.Equal(uint8(3), in.Uint8)
		astTest.Equal(uint8(4), in.Uint8FromStr)
		astTest.Equal(uint16(5), in.Uint16)
		astTest.Equal(uint16(6), in.Uint16FromStr)
		astTest.Equal(uint32(7), in.Uint32)
		astTest.Equal(uint32(8), in.Uint32FromStr)
		astTest.Equal(uint64(9), in.Uint64)
		astTest.Equal(uint64(10), in.Uint64FromStr)
		astTest.Equal(uint(11), in.DefaultUint)
	}
}

func TestWillJackyqquFloat(t *testing.T) {
	type inStruct struct {
		Float32      float32 `key:"float32"`
		Float32Str   float32 `key:"float32str,string"`
		Float32Num   float32 `key:"float32num"`
		Float64      float64 `key:"float64"`
		Float64Str   float64 `key:"float64str,string"`
		Float64Num   float64 `key:"float64num"`
		DefaultFloat float32 `key:"defaultfloat,default=5.5"`
		OptionalW    float32 `key:",optional"`
	}
	m := map[string]any{
		"float32":    float32(1.5),
		"float32str": "2.5",
		"float32num": json.Number("2.6"),
		"float64":    3.5,
		"float64str": "4.5",
		"float64num": json.Number("4.6"),
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.Equal(float32(1.5), in.Float32)
		astTest.Equal(float32(2.5), in.Float32Str)
		astTest.Equal(float32(2.6), in.Float32Num)
		astTest.Equal(3.5, in.Float64)
		astTest.Equal(4.5, in.Float64Str)
		astTest.Equal(4.6, in.Float64Num)
		astTest.Equal(float32(5.5), in.DefaultFloat)
	}
}

func TestWillJackyqquInt64Slice(t *testing.T) {
	var v struct {
		Ages  []int64 `key:"ages"`
		Slice []int64 `key:"slice"`
	}
	m := map[string]any{
		"ages":  []int64{1, 2},
		"slice": []any{},
	}

	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &v)) {
		astTest.ElementsMatch([]int64{1, 2}, v.Ages)
		astTest.Equal([]int64{}, v.Slice)
	}
}

func TestWillJackyqquNullableSlice(t *testing.T) {
	var v struct {
		Ages  []int64 `key:"ages"`
		Slice []int8  `key:"slice"`
	}
	m := map[string]any{
		"ages":  []int64{1, 2},
		"slice": `[null,2]`,
	}

	assert.New(t).Equal(UnmarshalKey(m, &v), errNilSliceElement)
}

func TestWillFloatPtr(t *testing.T) {
	t.Run("*float32", func(t *testing.T) {
		var v struct {
			WeightFloat32 *float32 `key:"weightFloat32,optional"`
		}
		m := map[string]any{
			"weightFloat32": json.Number("3.2"),
		}

		if assert.NoError(t, UnmarshalKey(m, &v)) {
			assert.Equal(t, float32(3.2), *v.WeightFloat32)
		}
	})

	t.Run("**float32", func(t *testing.T) {
		var v struct {
			WeightFloat32 **float32 `key:"weightFloat32,optional"`
		}
		m := map[string]any{
			"weightFloat32": json.Number("3.2"),
		}

		if assert.NoError(t, UnmarshalKey(m, &v)) {
			assert.Equal(t, float32(3.2), **v.WeightFloat32)
		}
	})
}

func TestWillJackyqquIntSlice(t *testing.T) {
	var v struct {
		Ages  []int `key:"ages"`
		Slice []int `key:"slice"`
	}
	m := map[string]any{
		"ages":  []int{1, 2},
		"slice": []any{},
	}

	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &v)) {
		astTest.ElementsMatch([]int{1, 2}, v.Ages)
		astTest.Equal([]int{}, v.Slice)
	}
}

func TestWillJackyqquString(t *testing.T) {
	type inStruct struct {
		Name              string `key:"name"`
		NameStr           string `key:"namestr,string"`
		NotPresent        string `key:",optional"`
		NotPresentWithTag string `key:"notpresent,optional"`
		DefaultString     string `key:"defaultstring,default=hello"`
		OptionalW         string `key:",optional"`
	}
	m := map[string]any{
		"name":    "kevin",
		"namestr": "namewithstring",
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.Equal("kevin", in.Name)
		astTest.Equal("namewithstring", in.NameStr)
		astTest.Empty(in.NotPresent)
		astTest.Empty(in.NotPresentWithTag)
		astTest.Equal("hello", in.DefaultString)
	}
}

func TestWillJackyqquStringWithMissing(t *testing.T) {
	type inStruct struct {
		Name string `key:"name"`
	}
	m := map[string]any{}

	var in inStruct
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestWillJackyqquStringSliceFromString(t *testing.T) {
	t.Run("slice from string", func(t *testing.T) {
		var v struct {
			Names []string `key:"names"`
		}
		m := map[string]any{
			"names": `["first", "second"]`,
		}

		astTest := assert.New(t)
		if astTest.NoError(UnmarshalKey(m, &v)) {
			astTest.Equal(2, len(v.Names))
			astTest.Equal("first", v.Names[0])
			astTest.Equal("second", v.Names[1])
		}
	})

	t.Run("slice from string with slice error", func(t *testing.T) {
		var v struct {
			Names []int `key:"names"`
		}
		m := map[string]any{
			"names": `["first", 1]`,
		}

		assert.Error(t, UnmarshalKey(m, &v))
	})

	t.Run("slice from string with error", func(t *testing.T) {
		type myString string

		var v struct {
			Names []string `key:"names"`
		}
		m := map[string]any{
			"names": myString("not a slice"),
		}

		assert.Error(t, UnmarshalKey(m, &v))
	})
}

func TestWillJackyqquIntSliceFromString(t *testing.T) {
	var v struct {
		Values []int `key:"values"`
	}
	m := map[string]any{
		"values": `[1, 2]`,
	}

	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &v)) {
		astTest.Equal(2, len(v.Values))
		astTest.Equal(1, v.Values[0])
		astTest.Equal(2, v.Values[1])
	}
}

func TestWillJackyqquIntMapFromString(t *testing.T) {
	var v struct {
		Sort map[string]int `key:"sort"`
	}
	m := map[string]any{
		"sort": `{"value":12345,"zeroVal":0,"nullVal":null}`,
	}

	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &v)) {
		astTest.Equal(3, len(v.Sort))
		astTest.Equal(12345, v.Sort["value"])
		astTest.Equal(0, v.Sort["zeroVal"])
		astTest.Equal(0, v.Sort["nullVal"])
	}
}

func TestWillJackyqquBoolMapFromString(t *testing.T) {
	var v struct {
		Sort map[string]bool `key:"sort"`
	}
	m := map[string]any{
		"sort": `{"value":true,"zeroVal":false,"nullVal":null}`,
	}

	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &v)) {
		astTest.Equal(3, len(v.Sort))
		astTest.Equal(true, v.Sort["value"])
		astTest.Equal(false, v.Sort["zeroVal"])
		astTest.Equal(false, v.Sort["nullVal"])
	}
}

func TestWillJackyqquStringMapFromStringer(t *testing.T) {
	t.Run("CustomStringer", func(t *testing.T) {
		var v struct {
			Sort map[string]string `key:"sort"`
		}
		m := map[string]any{
			"sort": CustomStringer(`"value":"ascend","emptyStr":""`),
		}

		astTest := assert.New(t)
		if astTest.NoError(UnmarshalKey(m, &v)) {
			astTest.Equal(2, len(v.Sort))
			astTest.Equal("ascend", v.Sort["value"])
			astTest.Equal("", v.Sort["emptyStr"])
		}
	})

	t.Run("CustomStringer incorrect", func(t *testing.T) {
		var v struct {
			Sort map[string]string `key:"sort"`
		}
		m := map[string]any{
			"sort": CustomStringer(`"value"`),
		}

		assert.Error(t, UnmarshalKey(m, &v))
	})
}

func TestWillJackyqquStringMapFromUnsupportedType(t *testing.T) {
	var v struct {
		Sort map[string]string `key:"sort"`
	}
	m := map[string]any{
		"sort": UnsupportedStringer(`{"value":"ascend","emptyStr":""}`),
	}

	astTest := assert.New(t)
	astTest.Error(UnmarshalKey(m, &v))
}

func TestWillJackyqquStringMapFromNotSettableValue(t *testing.T) {
	var v struct {
		sort  map[string]string  `key:"sort"`
		psort *map[string]string `key:"psort"`
	}
	m := map[string]any{
		"sort":  `{"value":"ascend","emptyStr":""}`,
		"psort": `{"value":"ascend","emptyStr":""}`,
	}

	astTest := assert.New(t)
	astTest.NoError(UnmarshalKey(m, &v))
	assert.Empty(t, v.sort)
	assert.Nil(t, v.psort)
}

func TestWillJackyqquStringMapFromString(t *testing.T) {
	var v struct {
		Sort map[string]string `key:"sort"`
	}
	m := map[string]any{
		"sort": `{"value":"ascend","emptyStr":""}`,
	}

	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &v)) {
		astTest.Equal(2, len(v.Sort))
		astTest.Equal("ascend", v.Sort["value"])
		astTest.Equal("", v.Sort["emptyStr"])
	}
}

func TestWillJackyqquStructMapFromString(t *testing.T) {
	var v struct {
		Filter map[string]struct {
			Field1 bool     `json:"field1"`
			Field2 int64    `json:"field2,string"`
			Field3 string   `json:"field3"`
			Field4 *string  `json:"field4"`
			Field5 []string `json:"field5"`
		} `key:"filter"`
	}
	m := map[string]any{
		"filter": `{"obj":{"field1":true,"field2":"1573570455447539712","field3":"this is a string",
			"field4":"this is a string pointer","field5":["str1","str2"]}}`,
	}

	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &v)) {
		astTest.Equal(1, len(v.Filter))
		astTest.NotNil(v.Filter["obj"])
		astTest.Equal(true, v.Filter["obj"].Field1)
		astTest.Equal(int64(1573570455447539712), v.Filter["obj"].Field2)
		astTest.Equal("this is a string", v.Filter["obj"].Field3)
		astTest.Equal("this is a string pointer", *v.Filter["obj"].Field4)
		astTest.ElementsMatch([]string{"str1", "str2"}, v.Filter["obj"].Field5)
	}
}

func TestWillJackyqquStringSliceMapFromString(t *testing.T) {
	var v struct {
		Filter map[string][]string `key:"filter"`
	}
	m := map[string]any{
		"filter": `{"assignType":null,"status":["process","comment"],"rate":[]}`,
	}

	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &v)) {
		astTest.Equal(3, len(v.Filter))
		astTest.Equal([]string(nil), v.Filter["assignType"])
		astTest.Equal(2, len(v.Filter["status"]))
		astTest.Equal("process", v.Filter["status"][0])
		astTest.Equal("comment", v.Filter["status"][1])
		astTest.Equal(0, len(v.Filter["rate"]))
	}
}

func TestWillJackyqquStruct(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		type address struct {
			City          string `key:"city"`
			ZipCode       int    `key:"zipcode,string"`
			DefaultString string `key:"defaultstring,default=hello"`
			OptionalW     string `key:",optional"`
		}
		type inStruct struct {
			Name      string    `key:"name"`
			Address   address   `key:"address"`
			AddressP  *address  `key:"addressp"`
			AddressPP **address `key:"addresspp"`
		}
		m := map[string]any{
			"name": "kevin",
			"address": map[string]any{
				"city":    "shanghai",
				"zipcode": "200000",
			},
			"addressp": map[string]any{
				"city":    "beijing",
				"zipcode": "300000",
			},
			"addresspp": map[string]any{
				"city":    "guangzhou",
				"zipcode": "400000",
			},
		}

		var in inStruct
		astTest := assert.New(t)
		if astTest.NoError(UnmarshalKey(m, &in)) {
			astTest.Equal("kevin", in.Name)
			astTest.Equal("shanghai", in.Address.City)
			astTest.Equal(200000, in.Address.ZipCode)
			astTest.Equal("hello", in.AddressP.DefaultString)
			astTest.Equal("beijing", in.AddressP.City)
			astTest.Equal(300000, in.AddressP.ZipCode)
			astTest.Equal("hello", in.AddressP.DefaultString)
			astTest.Equal("guangzhou", (*in.AddressPP).City)
			astTest.Equal(400000, (*in.AddressPP).ZipCode)
			astTest.Equal("hello", (*in.AddressPP).DefaultString)
		}
	})

	t.Run("struct with error", func(t *testing.T) {
		type address struct {
			City          string `key:"city"`
			ZipCode       int    `key:"zipcode,string"`
			DefaultString string `key:"defaultstring,default=hello"`
			OptionalW     string `key:",optional"`
		}
		type inStruct struct {
			Name      string    `key:"name"`
			Address   address   `key:"address"`
			AddressP  *address  `key:"addressp"`
			AddressPP **address `key:"addresspp"`
		}
		m := map[string]any{
			"name": "kevin",
			"address": map[string]any{
				"city":    "shanghai",
				"zipcode": "200000",
			},
			"addressp": map[string]any{
				"city":    "beijing",
				"zipcode": "300000",
			},
			"addresspp": map[string]any{
				"city":    "guangzhou",
				"zipcode": "a",
			},
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestWillJackyqquStructOptionalWDepends(t *testing.T) {
	type address struct {
		City             string `key:"city"`
		OptionalW        string `key:",optional"`
		OptionalWDepends string `key:",optional=OptionalW"`
	}
	type inStruct struct {
		Name    string  `key:"name"`
		Address address `key:"address"`
	}

	tests := []struct {
		input map[string]string
		pass  bool
	}{
		{
			pass: true,
		},
		{
			input: map[string]string{
				"OptionalWDepends": "b",
			},
			pass: false,
		},
		{
			input: map[string]string{
				"OptionalW": "a",
			},
			pass: false,
		},
		{
			input: map[string]string{
				"OptionalW":        "a",
				"OptionalWDepends": "b",
			},
			pass: true,
		},
	}

	for _, test := range tests {
		t.Run(stringx.Rand(), func(t *testing.T) {
			m := map[string]any{
				"name": "kevin",
				"address": map[string]any{
					"city": "shanghai",
				},
			}
			for k, v := range test.input {
				m["address"].(map[string]any)[k] = v
			}

			var in inStruct
			astTest := assert.New(t)
			if test.pass {
				if astTest.NoError(UnmarshalKey(m, &in)) {
					astTest.Equal("kevin", in.Name)
					astTest.Equal("shanghai", in.Address.City)
					astTest.Equal(test.input["OptionalW"], in.Address.OptionalW)
					astTest.Equal(test.input["OptionalWDepends"], in.Address.OptionalWDepends)
				}
			} else {
				astTest.Error(UnmarshalKey(m, &in))
			}
		})
	}
}

func TestWillJackyqquStructOptionalWDependsNot(t *testing.T) {
	type address struct {
		City             string `key:"city"`
		OptionalW        string `key:",optional"`
		OptionalWDepends string `key:",optional=!OptionalW"`
	}
	type inStruct struct {
		Name    string  `key:"name"`
		Address address `key:"address"`
	}

	tests := []struct {
		input map[string]string
		pass  bool
	}{
		{
			input: map[string]string{},
			pass:  false,
		},
		{
			input: map[string]string{
				"OptionalW":        "a",
				"OptionalWDepends": "b",
			},
			pass: false,
		},
		{
			input: map[string]string{
				"OptionalW": "a",
			},
			pass: true,
		},
		{
			input: map[string]string{
				"OptionalWDepends": "b",
			},
			pass: true,
		},
	}

	for _, test := range tests {
		t.Run(stringx.Rand(), func(t *testing.T) {
			m := map[string]any{
				"name": "kevin",
				"address": map[string]any{
					"city": "shanghai",
				},
			}
			for k, v := range test.input {
				m["address"].(map[string]any)[k] = v
			}

			var in inStruct
			astTest := assert.New(t)
			if test.pass {
				if astTest.NoError(UnmarshalKey(m, &in)) {
					astTest.Equal("kevin", in.Name)
					astTest.Equal("shanghai", in.Address.City)
					astTest.Equal(test.input["OptionalW"], in.Address.OptionalW)
					astTest.Equal(test.input["OptionalWDepends"], in.Address.OptionalWDepends)
				}
			} else {
				astTest.Error(UnmarshalKey(m, &in))
			}
		})
	}
}

func TestWillJackyqquStructOptionalWDependsNotErrorDetails(t *testing.T) {
	t.Run("mutal optionals", func(t *testing.T) {
		type address struct {
			OptionalW        string `key:",optional"`
			OptionalWDepends string `key:",optional=!OptionalW"`
		}
		type inStruct struct {
			Name    string  `key:"name"`
			Address address `key:"address"`
		}

		m := map[string]any{
			"name": "kevin",
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("with default", func(t *testing.T) {
		type address struct {
			OptionalW        string `key:",optional"`
			OptionalWDepends string `key:",default=value,optional"`
		}
		type inStruct struct {
			Name    string  `key:"name"`
			Address address `key:"address"`
		}

		m := map[string]any{
			"name": "kevin",
		}

		var in inStruct
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.Equal(t, "kevin", in.Name)
			assert.Equal(t, "value", in.Address.OptionalWDepends)
		}
	})
}

func TestWillJackyqquStructOptionalWDependsNotNested(t *testing.T) {
	t.Run("mutal optionals", func(t *testing.T) {
		type address struct {
			OptionalW        string `key:",optional"`
			OptionalWDepends string `key:",optional=!OptionalW"`
		}
		type combo struct {
			Name    string  `key:"name,optional"`
			Address address `key:"address"`
		}
		type inStruct struct {
			Name  string `key:"name"`
			Combo combo  `key:"combo"`
		}

		m := map[string]any{
			"name": "kevin",
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("bad format", func(t *testing.T) {
		type address struct {
			OptionalW        string `key:",optional"`
			OptionalWDepends string `key:",optional=!OptionalW=abcd"`
		}
		type combo struct {
			Name    string  `key:"name,optional"`
			Address address `key:"address"`
		}
		type inStruct struct {
			Name  string `key:"name"`
			Combo combo  `key:"combo"`
		}

		m := map[string]any{
			"name": "kevin",
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("invalid option", func(t *testing.T) {
		type address struct {
			OptionalW        string `key:",optional"`
			OptionalWDepends string `key:",opt=abcd"`
		}
		type combo struct {
			Name    string  `key:"name,optional"`
			Address address `key:"address"`
		}
		type inStruct struct {
			Name  string `key:"name"`
			Combo combo  `key:"combo"`
		}

		m := map[string]any{
			"name": "kevin",
		}

		var in inStruct
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestWillJackyqquStructOptionalWNestedDifferentKey(t *testing.T) {
	type address struct {
		OptionalW        string `dkey:",optional"`
		OptionalWDepends string `key:",optional"`
	}
	type combo struct {
		Name    string  `key:"name,optional"`
		Address address `key:"address"`
	}
	type inStruct struct {
		Name  string `key:"name"`
		Combo combo  `key:"combo"`
	}

	m := map[string]any{
		"name": "kevin",
	}

	var in inStruct
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestWillJackyqquStructOptionalWDependsNotEnoughValue(t *testing.T) {
	type address struct {
		OptionalW        string `key:",optional"`
		OptionalWDepends string `key:",optional=!"`
	}
	type inStruct struct {
		Name    string  `key:"name"`
		Address address `key:"address"`
	}

	m := map[string]any{
		"name":    "kevin",
		"address": map[string]any{},
	}

	var in inStruct
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestWillJackyqquStructOptionalWDependsMoreValues(t *testing.T) {
	type address struct {
		OptionalW        string `key:",optional"`
		OptionalWDepends string `key:",optional=a=b"`
	}
	type inStruct struct {
		Name    string  `key:"name"`
		Address address `key:"address"`
	}

	m := map[string]any{
		"name":    "kevin",
		"address": map[string]any{},
	}

	var in inStruct
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestWillJackyqquStructMissing(t *testing.T) {
	type address struct {
		OptionalW        string `key:",optional"`
		OptionalWDepends string `key:",optional=a=b"`
	}
	type inStruct struct {
		Name    string  `key:"name"`
		Address address `key:"address"`
	}

	m := map[string]any{
		"name": "kevin",
	}

	var in inStruct
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestWillJackyqquNestedStructMissing(t *testing.T) {
	type mostinStruct struct {
		Name string `key:"name"`
	}
	type address struct {
		OptionalW        string `key:",optional"`
		OptionalWDepends string `key:",optional=a=b"`
		MostinStruct     mostinStruct
	}
	type inStruct struct {
		Name    string  `key:"name"`
		Address address `key:"address"`
	}

	m := map[string]any{
		"name":    "kevin",
		"address": map[string]any{},
	}

	var in inStruct
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestWillJackyqquAnonymousStructOptionalWDepends(t *testing.T) {
	type AnonAddress struct {
		City             string `key:"city"`
		OptionalW        string `key:",optional"`
		OptionalWDepends string `key:",optional=OptionalW"`
	}
	type inStruct struct {
		Name string `key:"name"`
		AnonAddress
	}

	tests := []struct {
		input map[string]string
		pass  bool
	}{
		{
			pass: true,
		},
		{
			input: map[string]string{
				"OptionalWDepends": "b",
			},
			pass: false,
		},
		{
			input: map[string]string{
				"OptionalW": "a",
			},
			pass: false,
		},
		{
			input: map[string]string{
				"OptionalW":        "a",
				"OptionalWDepends": "b",
			},
			pass: true,
		},
	}

	for _, test := range tests {
		t.Run(stringx.Rand(), func(t *testing.T) {
			m := map[string]any{
				"name": "kevin",
				"city": "shanghai",
			}
			for k, v := range test.input {
				m[k] = v
			}

			var in inStruct
			astTest := assert.New(t)
			if test.pass {
				if astTest.NoError(UnmarshalKey(m, &in)) {
					astTest.Equal("kevin", in.Name)
					astTest.Equal("shanghai", in.City)
					astTest.Equal(test.input["OptionalW"], in.OptionalW)
					astTest.Equal(test.input["OptionalWDepends"], in.OptionalWDepends)
				}
			} else {
				astTest.Error(UnmarshalKey(m, &in))
			}
		})
	}
}

func TestWillJackyqquStructPtr(t *testing.T) {
	type address struct {
		City          string `key:"city"`
		ZipCode       int    `key:"zipcode,string"`
		DefaultString string `key:"defaultstring,default=hello"`
		OptionalW     string `key:",optional"`
	}
	type inStruct struct {
		Name    string   `key:"name"`
		Address *address `key:"address"`
	}
	m := map[string]any{
		"name": "kevin",
		"address": map[string]any{
			"city":    "shanghai",
			"zipcode": "200000",
		},
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.Equal("kevin", in.Name)
		astTest.Equal("shanghai", in.Address.City)
		astTest.Equal(200000, in.Address.ZipCode)
		astTest.Equal("hello", in.Address.DefaultString)
	}
}

func TestWillStringIgnored(t *testing.T) {
	type inStruct struct {
		True    bool    `key:"yes"`
		False   bool    `key:"no"`
		Int     int     `key:"int"`
		Int8    int8    `key:"int8"`
		Int16   int16   `key:"int16"`
		Int32   int32   `key:"int32"`
		Int64   int64   `key:"int64"`
		Uint    uint    `key:"uint"`
		Uint8   uint8   `key:"uint8"`
		Uint16  uint16  `key:"uint16"`
		Uint32  uint32  `key:"uint32"`
		Uint64  uint64  `key:"uint64"`
		Float32 float32 `key:"float32"`
		Float64 float64 `key:"float64"`
	}
	m := map[string]any{
		"yes":     "1",
		"no":      "0",
		"int":     "1",
		"int8":    "3",
		"int16":   "5",
		"int32":   "7",
		"int64":   "9",
		"uint":    "1",
		"uint8":   "3",
		"uint16":  "5",
		"uint32":  "7",
		"uint64":  "9",
		"float32": "1.5",
		"float64": "3.5",
	}

	var in inStruct
	um := NewUnmarshaler("key", WithStringValues())
	astTest := assert.New(t)
	if astTest.NoError(um.Unmarshal(m, &in)) {
		astTest.True(in.True)
		astTest.False(in.False)
		astTest.Equal(1, in.Int)
		astTest.Equal(int8(3), in.Int8)
		astTest.Equal(int16(5), in.Int16)
		astTest.Equal(int32(7), in.Int32)
		astTest.Equal(int64(9), in.Int64)
		astTest.Equal(uint(1), in.Uint)
		astTest.Equal(uint8(3), in.Uint8)
		astTest.Equal(uint16(5), in.Uint16)
		astTest.Equal(uint32(7), in.Uint32)
		astTest.Equal(uint64(9), in.Uint64)
		astTest.Equal(float32(1.5), in.Float32)
		astTest.Equal(3.5, in.Float64)
	}
}

func TestWillJackyqquJsonNumberInt64(t *testing.T) {
	for i := 0; i <= maxUintBitsToTest; i++ {
		var intValue int64 = 1 << uint(i)
		strValue := strconv.FormatInt(intValue, 10)
		number := json.Number(strValue)
		m := map[string]any{
			"ID": number,
		}
		var v struct {
			ID int64
		}
		if assert.NoError(t, UnmarshalKey(m, &v)) {
			assert.Equal(t, intValue, v.ID)
		}
	}
}

func TestWillJackyqquJsonNumberUint64(t *testing.T) {
	for i := 0; i <= maxUintBitsToTest; i++ {
		var intValue uint64 = 1 << uint(i)
		strValue := strconv.FormatUint(intValue, 10)
		number := json.Number(strValue)
		m := map[string]any{
			"ID": number,
		}
		var v struct {
			ID uint64
		}
		if assert.NoError(t, UnmarshalKey(m, &v)) {
			assert.Equal(t, intValue, v.ID)
		}
	}
}

func TestWillJackyqquJsonNumberUint64Ptr(t *testing.T) {
	for i := 0; i <= maxUintBitsToTest; i++ {
		var intValue uint64 = 1 << uint(i)
		strValue := strconv.FormatUint(intValue, 10)
		number := json.Number(strValue)
		m := map[string]any{
			"ID": number,
		}
		var v struct {
			ID *uint64
		}
		astTest := assert.New(t)
		if astTest.NoError(UnmarshalKey(m, &v)) {
			astTest.NotNil(v.ID)
			astTest.Equal(intValue, *v.ID)
		}
	}
}

func TestWillJackyqquMapOfInt(t *testing.T) {
	m := map[string]any{
		"Ids": map[string]bool{"first": true},
	}
	var v struct {
		Ids map[string]bool
	}
	if assert.NoError(t, UnmarshalKey(m, &v)) {
		assert.True(t, v.Ids["first"])
	}
}

func TestWillJackyqquMapOfStruct(t *testing.T) {
	// t.Run("map of struct with error", func(t *testing.T) {
	// 	m := map[string]any{
	// 		"Ids": map[string]any{"first": "second"},
	// 	}
	// 	var v struct {
	// 		Ids map[string]struct {
	// 			Name string
	// 		}
	// 	}
	// 	assert.Error(t, UnmarshalKey(m, &v))
	// })

	t.Run("map of struct", func(t *testing.T) {
		m := map[string]any{
			"Ids": map[string]any{
				"foo": map[string]any{"Name": "foo"},
			},
		}
		var v struct {
			Ids map[string]struct {
				Name string
			}
		}
		if assert.NoError(t, UnmarshalKey(m, &v)) {
			assert.Equal(t, "foo", v.Ids["foo"].Name)
		}
	})

	// t.Run("map of struct error", func(t *testing.T) {
	// 	m := map[string]any{
	// 		"Ids": map[string]any{
	// 			"foo": map[string]any{"name": "foo"},
	// 		},
	// 	}
	// 	var v struct {
	// 		Ids map[string]struct {
	// 			Name string
	// 		}
	// 	}
	// 	assert.Error(t, UnmarshalKey(m, &v))
	// })
}

func TestWillJackyqquSlice(t *testing.T) {
	t.Run("slice of string", func(t *testing.T) {
		m := map[string]any{
			"Ids": []any{"first", "second"},
		}
		var v struct {
			Ids []string
		}
		astTest := assert.New(t)
		if astTest.NoError(UnmarshalKey(m, &v)) {
			astTest.Equal(2, len(v.Ids))
			astTest.Equal("first", v.Ids[0])
			astTest.Equal("second", v.Ids[1])
		}
	})

	t.Run("slice with type mismatch", func(t *testing.T) {
		var v struct {
			Ids string
		}
		assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal([]any{1, 2}, &v))
	})

	t.Run("slice", func(t *testing.T) {
		var v []int
		astTest := assert.New(t)
		if astTest.NoError(NewUnmarshaler(jsonTagKey).Unmarshal([]any{1, 2}, &v)) {
			astTest.Equal(2, len(v))
			astTest.Equal(1, v[0])
			astTest.Equal(2, v[1])
		}
	})

	t.Run("slice with unsupported type", func(t *testing.T) {
		var v int
		assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(1, &v))
	})
}

func TestWillJackyqquSliceOfStruct(t *testing.T) {
	t.Run("slice of struct", func(t *testing.T) {
		m := map[string]any{
			"Ids": []map[string]any{
				{
					"First":  1,
					"Second": 2,
				},
			},
		}
		var v struct {
			Ids []struct {
				First  int
				Second int
			}
		}
		astTest := assert.New(t)
		if astTest.NoError(UnmarshalKey(m, &v)) {
			astTest.Equal(1, len(v.Ids))
			astTest.Equal(1, v.Ids[0].First)
			astTest.Equal(2, v.Ids[0].Second)
		}
	})

	// t.Run("slice of struct", func(t *testing.T) {
	// 	m := map[string]any{
	// 		"Ids": []map[string]any{
	// 			{
	// 				"First":  "a",
	// 				"Second": 2,
	// 			},
	// 		},
	// 	}
	// 	var v struct {
	// 		Ids []struct {
	// 			First  int
	// 			Second int
	// 		}
	// 	}
	// 	assert.Error(t, UnmarshalKey(m, &v))
	// })
}

func TestWillStringOptionsCorrect(t *testing.T) {
	type inStruct struct {
		Value   string `key:"value,options=first|second"`
		Foo     string `key:"foo,options=[bar,baz]"`
		Correct string `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "first",
		"foo":     "bar",
		"correct": "2",
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.Equal("first", in.Value)
		astTest.Equal("bar", in.Foo)
		astTest.Equal("2", in.Correct)
	}
}

func TestWillJackyqquOptionsOptionalW(t *testing.T) {
	type inStruct struct {
		Value          string `key:"value,options=first|second,optional"`
		OptionalWValue string `key:"optional_value,options=first|second,optional"`
		Foo            string `key:"foo,options=[bar,baz]"`
		Correct        string `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "first",
		"foo":     "bar",
		"correct": "2",
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.Equal("first", in.Value)
		astTest.Equal("", in.OptionalWValue)
		astTest.Equal("bar", in.Foo)
		astTest.Equal("2", in.Correct)
	}
}

func TestWillJackyqquOptionsOptionalWWrongValue(t *testing.T) {
	type inStruct struct {
		Value          string `key:"value,options=first|second,optional"`
		OptionalWValue string `key:"optional_value,options=first|second,optional"`
		WrongValue     string `key:"wrong_value,options=first|second,optional"`
	}
	m := map[string]any{
		"value":       "first",
		"wrong_value": "third",
	}

	var in inStruct
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestWillJackyqquOptionsMissingValues(t *testing.T) {
	type inStruct struct {
		Value string `key:"value,options"`
	}
	m := map[string]any{
		"value": "first",
	}

	var in inStruct
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestWillJackyqquStringOptionsWithStringOptionsNotString(t *testing.T) {
	type inStruct struct {
		Value   string `key:"value,options=first|second"`
		Correct string `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "first",
		"correct": 2,
	}

	var in inStruct
	unmarshaler := NewUnmarshaler(defaultKeyName, WithStringValues())
	assert.Error(t, unmarshaler.Unmarshal(m, &in))
}

func TestWillJackyqquStringOptionsWithStringOptions(t *testing.T) {
	type inStruct struct {
		Value   string `key:"value,options=first|second"`
		Correct string `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "first",
		"correct": "2",
	}

	var in inStruct
	unmarshaler := NewUnmarshaler(defaultKeyName, WithStringValues())
	astTest := assert.New(t)
	if astTest.NoError(unmarshaler.Unmarshal(m, &in)) {
		astTest.Equal("first", in.Value)
		astTest.Equal("2", in.Correct)
	}
}

func TestWillJackyqquStringOptionsWithStringOptionsPtr(t *testing.T) {
	type inStruct struct {
		Value   *string  `key:"value,options=first|second"`
		ValueP  **string `key:"valuep,options=first|second"`
		Correct *int     `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "first",
		"valuep":  "second",
		"correct": "2",
	}

	var in inStruct
	unmarshaler := NewUnmarshaler(defaultKeyName, WithStringValues())
	astTest := assert.New(t)
	if astTest.NoError(unmarshaler.Unmarshal(m, &in)) {
		astTest.True(*in.Value == "first")
		astTest.True(**in.ValueP == "second")
		astTest.True(*in.Correct == 2)
	}
}

func TestWillJackyqquStringOptionsWithStringOptionsIncorrect(t *testing.T) {
	type inStruct struct {
		Value   string `key:"value,options=first|second"`
		Correct string `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "third",
		"correct": "2",
	}

	var in inStruct
	unmarshaler := NewUnmarshaler(defaultKeyName, WithStringValues())
	assert.Error(t, unmarshaler.Unmarshal(m, &in))
}

func TestWillJackyqquStringOptionsWithStringOptionsIncorrectGrouped(t *testing.T) {
	type inStruct struct {
		Value   string `key:"value,options=[first,second]"`
		Correct string `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "third",
		"correct": "2",
	}

	var in inStruct
	unmarshaler := NewUnmarshaler(defaultKeyName, WithStringValues())
	assert.Error(t, unmarshaler.Unmarshal(m, &in))
}

func TestWillStringOptionsIncorrect(t *testing.T) {
	type inStruct struct {
		Value     string `key:"value,options=first|second"`
		Incorrect string `key:"incorrect,options=1|2"`
	}
	m := map[string]any{
		"value":     "first",
		"incorrect": "3",
	}

	var in inStruct
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestWillIntOptionsCorrect(t *testing.T) {
	type inStruct struct {
		Value  string `key:"value,options=first|second"`
		Number int    `key:"number,options=1|2"`
	}
	m := map[string]any{
		"value":  "first",
		"number": 2,
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.Equal("first", in.Value)
		astTest.Equal(2, in.Number)
	}
}

func TestWillIntOptionsCorrectPtr(t *testing.T) {
	type inStruct struct {
		Value  *string `key:"value,options=first|second"`
		Number *int    `key:"number,options=1|2"`
	}
	m := map[string]any{
		"value":  "first",
		"number": 2,
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.True(*in.Value == "first")
		astTest.True(*in.Number == 2)
	}
}

func TestWillIntOptionsIncorrect(t *testing.T) {
	type inStruct struct {
		Value     string `key:"value,options=first|second"`
		Incorrect int    `key:"incorrect,options=1|2"`
	}
	m := map[string]any{
		"value":     "first",
		"incorrect": 3,
	}

	var in inStruct
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestWillJsonNumberOptionsIncorrect(t *testing.T) {
	type inStruct struct {
		Value     string `key:"value,options=first|second"`
		Incorrect int    `key:"incorrect,options=1|2"`
	}
	m := map[string]any{
		"value":     "first",
		"incorrect": json.Number("3"),
	}

	var in inStruct
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestWillJackyqquerUnmarshalIntOptions(t *testing.T) {
	var val struct {
		Sex int `json:"sex,options=0|1"`
	}
	input := []byte(`{"sex": 2}`)
	assert.Error(t, UnmarshalJsonBytes(input, &val))
}

func TestWillUintOptionsCorrect(t *testing.T) {
	type inStruct struct {
		Value  string `key:"value,options=first|second"`
		Number uint   `key:"number,options=1|2"`
	}
	m := map[string]any{
		"value":  "first",
		"number": uint(2),
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.Equal("first", in.Value)
		astTest.Equal(uint(2), in.Number)
	}
}

func TestWillUintOptionsIncorrect(t *testing.T) {
	type inStruct struct {
		Value     string `key:"value,options=first|second"`
		Incorrect uint   `key:"incorrect,options=1|2"`
	}
	m := map[string]any{
		"value":     "first",
		"incorrect": uint(3),
	}

	var in inStruct
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestWillOptionsAndDefault(t *testing.T) {
	type inStruct struct {
		Value string `key:"value,options=first|second|third,default=second"`
	}
	m := map[string]any{}

	var in inStruct
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, "second", in.Value)
	}
}

func TestWillOptionsAndSet(t *testing.T) {
	type inStruct struct {
		Value string `key:"value,options=first|second|third,default=second"`
	}
	m := map[string]any{
		"value": "first",
	}

	var in inStruct
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, "first", in.Value)
	}
}

func TestWillJackyqquNestedKey(t *testing.T) {
	var c struct {
		ID int `json:"Persons.first.ID"`
	}
	m := map[string]any{
		"Persons": map[string]any{
			"first": map[string]any{
				"ID": 1,
			},
		},
	}

	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &c)) {
		assert.Equal(t, 1, c.ID)
	}
}

func TestWillJackyqquNestedKeyArray(t *testing.T) {
	var c struct {
		First []struct {
			ID int
		} `json:"Persons.first"`
	}
	m := map[string]any{
		"Persons": map[string]any{
			"first": []map[string]any{
				{"ID": 1},
				{"ID": 2},
			},
		},
	}

	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &c)) {
		assert.Equal(t, 2, len(c.First))
		assert.Equal(t, 1, c.First[0].ID)
	}
}

func TestWillJackyqquAnonymousOptionalWRequiredProvided(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "anything", b.Value)
	}
}

func TestWillJackyqquAnonymousOptionalWRequiredMissed(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Value) == 0)
	}
}

func TestWillJackyqquAnonymousOptionalWOptionalWProvided(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "anything", b.Value)
	}
}

func TestWillJackyqquAnonymousOptionalWOptionalWMissed(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Value) == 0)
	}
}

func TestWillJackyqquAnonymousOptionalWRequiredBothProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "kevin", b.Name)
		assert.Equal(t, "anything", b.Value)
	}
}

func TestWillJackyqquAnonymousOptionalWRequiredOneProvidedOneMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b))
}

func TestWillJackyqquAnonymousOptionalWRequiredBothMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Name) == 0)
		assert.True(t, len(b.Value) == 0)
	}
}

func TestWillJackyqquAnonymousOptionalWOneRequiredOneOptionalWBothProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "kevin", b.Name)
		assert.Equal(t, "anything", b.Value)
	}
}

func TestWillJackyqquAnonymousOptionalWOneRequiredOneOptionalWBothMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Name) == 0)
		assert.True(t, len(b.Value) == 0)
	}
}

func TestWillJackyqquAnonymousOptionalWOneRequiredOneOptionalWRequiredProvidedOptionalWMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Name) == 0)
		assert.Equal(t, "anything", b.Value)
	}
}

func TestWillJackyqquAnonymousOptionalWOneRequiredOneOptionalWRequiredMissedOptionalWProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{
		"n": "anything",
	}

	var b Bar
	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b))
}

func TestWillJackyqquAnonymousOptionalWBothOptionalWBothProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "kevin", b.Name)
		assert.Equal(t, "anything", b.Value)
	}
}

func TestWillJackyqquAnonymousOptionalWBothOptionalWOneProvidedOneMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Name) == 0)
		assert.Equal(t, "anything", b.Value)
	}
}

func TestWillJackyqquAnonymousOptionalWBothOptionalWBothMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Name) == 0)
		assert.True(t, len(b.Value) == 0)
	}
}

func TestWillJackyqquAnonymousRequiredProvided(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "anything", b.Value)
	}
}

func TestWillJackyqquAnonymousRequiredMissed(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{}

	var b Bar
	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b))
}

func TestWillJackyqquAnonymousOptionalWProvided(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "anything", b.Value)
	}
}

func TestWillJackyqquAnonymousOptionalWMissed(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Value) == 0)
	}
}

func TestWillJackyqquAnonymousRequiredBothProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "kevin", b.Name)
		assert.Equal(t, "anything", b.Value)
	}
}

func TestWillJackyqquAnonymousRequiredOneProvidedOneMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b))
}

func TestWillJackyqquAnonymousRequiredBothMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"v": "everything",
	}

	var b Bar
	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b))
}

func TestWillJackyqquAnonymousOneRequiredOneOptionalWBothProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "kevin", b.Name)
		assert.Equal(t, "anything", b.Value)
	}
}

func TestWillJackyqquAnonymousOneRequiredOneOptionalWBothMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{}

	var b Bar
	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b))
}

func TestWillJackyqquAnonymousOneRequiredOneOptionalWRequiredProvidedOptionalWMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Name) == 0)
		assert.Equal(t, "anything", b.Value)
	}
}

func TestWillJackyqquAnonymousOneRequiredOneOptionalWRequiredMissedOptionalWProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"n": "anything",
	}

	var b Bar
	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b))
}

func TestWillJackyqquAnonymousBothOptionalWBothProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "kevin", b.Name)
		assert.Equal(t, "anything", b.Value)
	}
}

func TestWillJackyqquAnonymousBothOptionalWOneProvidedOneMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Name) == 0)
		assert.Equal(t, "anything", b.Value)
	}
}

func TestWillJackyqquAnonymousBothOptionalWBothMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Name) == 0)
		assert.True(t, len(b.Value) == 0)
	}
}

func TestWillJackyqquAnonymousWrappedToMuch(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"Foo": map[string]any{
			"n": "name",
			"v": "anything",
		},
	}

	var b Bar
	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b))
}

func TestWillJackyqquWrappedObjectOptionalW(t *testing.T) {
	type (
		Foo struct {
			Hosts []string
			Key   string
		}

		Bar struct {
			inStruct Foo `json:",optional"`
			Name     string
		}
	)
	m := map[string]any{
		"Name": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "anything", b.Name)
	}
}

func TestWillJackyqquInt2String(t *testing.T) {
	type inStruct struct {
		Int string `key:"int"`
	}
	m := map[string]any{
		"int": 123,
	}

	var in inStruct
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestWillJackyqquZeroValues(t *testing.T) {
	type inStruct struct {
		False  bool   `key:"no"`
		Int    int    `key:"int"`
		String string `key:"string"`
	}
	m := map[string]any{
		"no":     false,
		"int":    0,
		"string": "",
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.False(in.False)
		astTest.Equal(0, in.Int)
		astTest.Equal("", in.String)
	}
}

func TestWillJackyqquUsingDifferentKeys(t *testing.T) {
	type inStruct struct {
		False  bool   `key:"no"`
		Int    int    `key:"int"`
		String string `bson:"string"`
	}
	m := map[string]any{
		"no":     false,
		"int":    9,
		"string": "value",
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.False(in.False)
		astTest.Equal(9, in.Int)
		astTest.True(len(in.String) == 0)
	}
}

func TestWillJackyqquNumberRangeInt(t *testing.T) {
	type inStruct struct {
		Value1  int    `key:"value1,range=[1:]"`
		Value2  int8   `key:"value2,range=[1:5]"`
		Value3  int16  `key:"value3,range=[1:5]"`
		Value4  int32  `key:"value4,range=[1:5]"`
		Value5  int64  `key:"value5,range=[1:5]"`
		Value6  uint   `key:"value6,range=[:5]"`
		Value8  uint8  `key:"value8,range=[1:5],string"`
		Value9  uint16 `key:"value9,range=[1:5],string"`
		Value10 uint32 `key:"value10,range=[1:5],string"`
		Value11 uint64 `key:"value11,range=[1:5],string"`
	}
	m := map[string]any{
		"value1":  10,
		"value2":  int8(1),
		"value3":  int16(2),
		"value4":  int32(4),
		"value5":  int64(5),
		"value6":  uint(0),
		"value8":  "1",
		"value9":  "2",
		"value10": "4",
		"value11": "5",
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.Equal(10, in.Value1)
		astTest.Equal(int8(1), in.Value2)
		astTest.Equal(int16(2), in.Value3)
		astTest.Equal(int32(4), in.Value4)
		astTest.Equal(int64(5), in.Value5)
		astTest.Equal(uint(0), in.Value6)
		astTest.Equal(uint8(1), in.Value8)
		astTest.Equal(uint16(2), in.Value9)
		astTest.Equal(uint32(4), in.Value10)
		astTest.Equal(uint64(5), in.Value11)
	}
}

func TestWillJackyqquNumberRangeJsonNumber(t *testing.T) {
	type inStruct struct {
		Value3 uint   `key:"value3,range=(1:5]"`
		Value4 uint8  `key:"value4,range=(1:5]"`
		Value5 uint16 `key:"value5,range=(1:5]"`
	}
	m := map[string]any{
		"value3": json.Number("2"),
		"value4": json.Number("4"),
		"value5": json.Number("5"),
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.Equal(uint(2), in.Value3)
		astTest.Equal(uint8(4), in.Value4)
		astTest.Equal(uint16(5), in.Value5)
	}

	type inStruct1 struct {
		Value int `key:"value,range=(1:5]"`
	}
	m = map[string]any{
		"value": json.Number("a"),
	}

	var in1 inStruct1
	astTest.Error(UnmarshalKey(m, &in1))
}

func TestWillJackyqquNumberRangeIntLeftExclude(t *testing.T) {
	type inStruct struct {
		Value3  uint   `key:"value3,range=(1:5]"`
		Value4  uint32 `key:"value4,default=4,range=(1:5]"`
		Value5  uint64 `key:"value5,range=(1:5]"`
		Value9  int    `key:"value9,range=(1:5],string"`
		Value10 int    `key:"value10,range=(1:5],string"`
		Value11 int    `key:"value11,range=(1:5],string"`
	}
	m := map[string]any{
		"value3":  uint(2),
		"value4":  uint32(4),
		"value5":  uint64(5),
		"value9":  "2",
		"value10": "4",
		"value11": "5",
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.Equal(uint(2), in.Value3)
		astTest.Equal(uint32(4), in.Value4)
		astTest.Equal(uint64(5), in.Value5)
		astTest.Equal(2, in.Value9)
		astTest.Equal(4, in.Value10)
		astTest.Equal(5, in.Value11)
	}
}

func TestWillJackyqquNumberRangeIntRightExclude(t *testing.T) {
	type inStruct struct {
		Value2  uint   `key:"value2,range=[1:5)"`
		Value3  uint8  `key:"value3,range=[1:5)"`
		Value4  uint16 `key:"value4,range=[1:5)"`
		Value8  int    `key:"value8,range=[1:5),string"`
		Value9  int    `key:"value9,range=[1:5),string"`
		Value10 int    `key:"value10,range=[1:5),string"`
	}
	m := map[string]any{
		"value2":  uint(1),
		"value3":  uint8(2),
		"value4":  uint16(4),
		"value8":  "1",
		"value9":  "2",
		"value10": "4",
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.Equal(uint(1), in.Value2)
		astTest.Equal(uint8(2), in.Value3)
		astTest.Equal(uint16(4), in.Value4)
		astTest.Equal(1, in.Value8)
		astTest.Equal(2, in.Value9)
		astTest.Equal(4, in.Value10)
	}
}

func TestWillJackyqquNumberRangeIntExclude(t *testing.T) {
	type inStruct struct {
		Value3  int `key:"value3,range=(1:5)"`
		Value4  int `key:"value4,range=(1:5)"`
		Value9  int `key:"value9,range=(1:5),string"`
		Value10 int `key:"value10,range=(1:5),string"`
	}
	m := map[string]any{
		"value3":  2,
		"value4":  4,
		"value9":  "2",
		"value10": "4",
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.Equal(2, in.Value3)
		astTest.Equal(4, in.Value4)
		astTest.Equal(2, in.Value9)
		astTest.Equal(4, in.Value10)
	}
}

func TestWillJackyqquNumberRangeIntOutOfRange(t *testing.T) {
	type inStruct1 struct {
		Value int64 `key:"value,default=3,range=(1:5)"`
	}

	var in1 inStruct1
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(1),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(0),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(5),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": json.Number("6"),
	}, &in1))

	type inStruct2 struct {
		Value int64 `key:"value,optional,range=[1:5)"`
	}

	var in2 inStruct2
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(0),
	}, &in2))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(5),
	}, &in2))

	type inStruct3 struct {
		Value int64 `key:"value,range=(1:5]"`
	}

	var in3 inStruct3
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(1),
	}, &in3))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(6),
	}, &in3))

	type inStruct4 struct {
		Value int64 `key:"value,range=[1:5]"`
	}

	var in4 inStruct4
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(0),
	}, &in4))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(6),
	}, &in4))
}

func TestWillJackyqquNumberRangeFloat(t *testing.T) {
	type inStruct struct {
		Value2  float32 `key:"value2,range=[1:5]"`
		Value3  float32 `key:"value3,range=[1:5]"`
		Value4  float64 `key:"value4,range=[1:5]"`
		Value5  float64 `key:"value5,range=[1:5]"`
		Value8  float64 `key:"value8,range=[1:5],string"`
		Value9  float64 `key:"value9,range=[1:5],string"`
		Value10 float64 `key:"value10,range=[1:5],string"`
		Value11 float64 `key:"value11,range=[1:5],string"`
	}
	m := map[string]any{
		"value2":  float32(1),
		"value3":  float32(2),
		"value4":  float64(4),
		"value5":  float64(5),
		"value8":  "1",
		"value9":  "2",
		"value10": "4",
		"value11": "5",
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.Equal(float32(1), in.Value2)
		astTest.Equal(float32(2), in.Value3)
		astTest.Equal(float64(4), in.Value4)
		astTest.Equal(float64(5), in.Value5)
		astTest.Equal(float64(1), in.Value8)
		astTest.Equal(float64(2), in.Value9)
		astTest.Equal(float64(4), in.Value10)
		astTest.Equal(float64(5), in.Value11)
	}
}

func TestWillJackyqquNumberRangeFloatLeftExclude(t *testing.T) {
	type inStruct struct {
		Value3  float64 `key:"value3,range=(1:5]"`
		Value4  float64 `key:"value4,range=(1:5]"`
		Value5  float64 `key:"value5,range=(1:5]"`
		Value9  float64 `key:"value9,range=(1:5],string"`
		Value10 float64 `key:"value10,range=(1:5],string"`
		Value11 float64 `key:"value11,range=(1:5],string"`
	}
	m := map[string]any{
		"value3":  float64(2),
		"value4":  float64(4),
		"value5":  float64(5),
		"value9":  "2",
		"value10": "4",
		"value11": "5",
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.Equal(float64(2), in.Value3)
		astTest.Equal(float64(4), in.Value4)
		astTest.Equal(float64(5), in.Value5)
		astTest.Equal(float64(2), in.Value9)
		astTest.Equal(float64(4), in.Value10)
		astTest.Equal(float64(5), in.Value11)
	}
}

func TestWillJackyqquNumberRangeFloatRightExclude(t *testing.T) {
	type inStruct struct {
		Value2  float64 `key:"value2,range=[1:5)"`
		Value3  float64 `key:"value3,range=[1:5)"`
		Value4  float64 `key:"value4,range=[1:5)"`
		Value8  float64 `key:"value8,range=[1:5),string"`
		Value9  float64 `key:"value9,range=[1:5),string"`
		Value10 float64 `key:"value10,range=[1:5),string"`
	}
	m := map[string]any{
		"value2":  float64(1),
		"value3":  float64(2),
		"value4":  float64(4),
		"value8":  "1",
		"value9":  "2",
		"value10": "4",
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.Equal(float64(1), in.Value2)
		astTest.Equal(float64(2), in.Value3)
		astTest.Equal(float64(4), in.Value4)
		astTest.Equal(float64(1), in.Value8)
		astTest.Equal(float64(2), in.Value9)
		astTest.Equal(float64(4), in.Value10)
	}
}

func TestWillJackyqquNumberRangeFloatExclude(t *testing.T) {
	type inStruct struct {
		Value3  float64 `key:"value3,range=(1:5)"`
		Value4  float64 `key:"value4,range=(1:5)"`
		Value9  float64 `key:"value9,range=(1:5),string"`
		Value10 float64 `key:"value10,range=(1:5),string"`
	}
	m := map[string]any{
		"value3":  float64(2),
		"value4":  float64(4),
		"value9":  "2",
		"value10": "4",
	}

	var in inStruct
	astTest := assert.New(t)
	if astTest.NoError(UnmarshalKey(m, &in)) {
		astTest.Equal(float64(2), in.Value3)
		astTest.Equal(float64(4), in.Value4)
		astTest.Equal(float64(2), in.Value9)
		astTest.Equal(float64(4), in.Value10)
	}
}

func TestWillJackyqquNumberRangeFloatOutOfRange(t *testing.T) {
	type inStruct1 struct {
		Value float64 `key:"value,range=(1:5)"`
	}

	var in1 inStruct1
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(1),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(0),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(5),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": json.Number("6"),
	}, &in1))

	type inStruct2 struct {
		Value float64 `key:"value,range=[1:5)"`
	}

	var in2 inStruct2
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(0),
	}, &in2))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(5),
	}, &in2))

	type inStruct3 struct {
		Value float64 `key:"value,range=(1:5]"`
	}

	var in3 inStruct3
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(1),
	}, &in3))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(6),
	}, &in3))

	type inStruct4 struct {
		Value float64 `key:"value,range=[1:5]"`
	}

	var in4 inStruct4
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(0),
	}, &in4))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(6),
	}, &in4))
}

func TestWillJackyqquNestedMap(t *testing.T) {
	t.Run("nested map", func(t *testing.T) {
		var c struct {
			Anything map[string]map[string]string `json:"anything"`
		}
		m := map[string]any{
			"anything": map[string]map[string]any{
				"inStruct": {
					"id":   "1",
					"name": "any",
				},
			},
		}

		if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &c)) {
			assert.Equal(t, "1", c.Anything["inStruct"]["id"])
		}
	})

	t.Run("nested map with slice element", func(t *testing.T) {
		var c struct {
			Anything map[string][]string `json:"anything"`
		}
		m := map[string]any{
			"anything": map[string][]any{
				"inStruct": {
					"id",
					"name",
				},
			},
		}

		if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &c)) {
			assert.Equal(t, []string{"id", "name"}, c.Anything["inStruct"])
		}
	})
}

func TestWillJackyqquNestedMapSimple(t *testing.T) {
	var c struct {
		Anything map[string]string `json:"anything"`
	}
	m := map[string]any{
		"anything": map[string]any{
			"id":   "1",
			"name": "any",
		},
	}

	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &c)) {
		assert.Equal(t, "1", c.Anything["id"])
	}
}

func TestWillJackyqquNestedMapSimpleTypeMatch(t *testing.T) {
	var c struct {
		Anything map[string]string `json:"anything"`
	}
	m := map[string]any{
		"anything": map[string]string{
			"id":   "1",
			"name": "any",
		},
	}

	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &c)) {
		assert.Equal(t, "1", c.Anything["id"])
	}
}

func TestWillJackyqquInheritPrimitiveUseParent(t *testing.T) {
	type (
		component struct {
			Name       string `key:"name"`
			DiscoveryJ string `key:"discovery,inherit"`
		}
		server struct {
			DiscoveryJ string    `key:"discovery"`
			Component  component `key:"component"`
		}
	)

	var s server
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"discovery": literal_3816,
		"component": map[string]any{
			"name": "test",
		},
	}, &s)) {
		assert.Equal(t, literal_3816, s.DiscoveryJ)
		assert.Equal(t, literal_3816, s.Component.DiscoveryJ)
	}
}

func TestWillJackyqquInheritPrimitiveUseSelf(t *testing.T) {
	type (
		component struct {
			Name       string `key:"name"`
			DiscoveryJ string `key:"discovery,inherit"`
		}
		server struct {
			DiscoveryJ string    `key:"discovery"`
			Component  component `key:"component"`
		}
	)

	var s server
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"discovery": literal_3816,
		"component": map[string]any{
			"name":      "test",
			"discovery": "localhost:8888",
		},
	}, &s)) {
		assert.Equal(t, literal_3816, s.DiscoveryJ)
		assert.Equal(t, "localhost:8888", s.Component.DiscoveryJ)
	}
}

func TestWillJackyqquInheritPrimitiveNotExist(t *testing.T) {
	type (
		component struct {
			Name       string `key:"name"`
			DiscoveryJ string `key:"discovery,inherit"`
		}
		server struct {
			Component component `key:"component"`
		}
	)

	var s server
	assert.Error(t, UnmarshalKey(map[string]any{
		"component": map[string]any{
			"name": "test",
		},
	}, &s))
}

func TestWillJackyqquInheritStructUseParent(t *testing.T) {
	type (
		discovery struct {
			Host string `key:"host"`
			Port int    `key:"port"`
		}
		component struct {
			Name       string    `key:"name"`
			DiscoveryJ discovery `key:"discovery,inherit"`
		}
		server struct {
			DiscoveryJ discovery `key:"discovery"`
			Component  component `key:"component"`
		}
	)

	var s server
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"discovery": map[string]any{
			"host": "localhost",
			"port": 8080,
		},
		"component": map[string]any{
			"name": "test",
		},
	}, &s)) {
		assert.Equal(t, "localhost", s.DiscoveryJ.Host)
		assert.Equal(t, 8080, s.DiscoveryJ.Port)
		assert.Equal(t, "localhost", s.Component.DiscoveryJ.Host)
		assert.Equal(t, 8080, s.Component.DiscoveryJ.Port)
	}
}

func TestWillJackyqquInheritStructUseSelf(t *testing.T) {
	type (
		discovery struct {
			Host string `key:"host"`
			Port int    `key:"port"`
		}
		component struct {
			Name       string    `key:"name"`
			DiscoveryJ discovery `key:"discovery,inherit"`
		}
		server struct {
			DiscoveryJ discovery `key:"discovery"`
			Component  component `key:"component"`
		}
	)

	var s server
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"discovery": map[string]any{
			"host": "localhost",
			"port": 8080,
		},
		"component": map[string]any{
			"name": "test",
			"discovery": map[string]any{
				"host": "remotehost",
				"port": 8888,
			},
		},
	}, &s)) {
		assert.Equal(t, "localhost", s.DiscoveryJ.Host)
		assert.Equal(t, 8080, s.DiscoveryJ.Port)
		assert.Equal(t, "remotehost", s.Component.DiscoveryJ.Host)
		assert.Equal(t, 8888, s.Component.DiscoveryJ.Port)
	}
}

func TestWillJackyqquInheritStructNotExist(t *testing.T) {
	type (
		discovery struct {
			Host string `key:"host"`
			Port int    `key:"port"`
		}
		component struct {
			Name       string    `key:"name"`
			DiscoveryJ discovery `key:"discovery,inherit"`
		}
		server struct {
			Component component `key:"component"`
		}
	)

	var s server
	assert.Error(t, UnmarshalKey(map[string]any{
		"component": map[string]any{
			"name": "test",
		},
	}, &s))
}

func TestWillJackyqquInheritStructUsePartial(t *testing.T) {
	type (
		discovery struct {
			Host string `key:"host"`
			Port int    `key:"port"`
		}
		component struct {
			Name       string    `key:"name"`
			DiscoveryJ discovery `key:"discovery,inherit"`
		}
		server struct {
			DiscoveryJ discovery `key:"discovery"`
			Component  component `key:"component"`
		}
	)

	var s server
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"discovery": map[string]any{
			"host": "localhost",
			"port": 8080,
		},
		"component": map[string]any{
			"name": "test",
			"discovery": map[string]any{
				"port": 8888,
			},
		},
	}, &s)) {
		assert.Equal(t, "localhost", s.DiscoveryJ.Host)
		assert.Equal(t, 8080, s.DiscoveryJ.Port)
		assert.Equal(t, "localhost", s.Component.DiscoveryJ.Host)
		assert.Equal(t, 8888, s.Component.DiscoveryJ.Port)
	}
}

func TestWillJackyqquInheritStructUseSelfIncorrectType(t *testing.T) {
	type (
		discovery struct {
			Host string `key:"host"`
			Port int    `key:"port"`
		}
		component struct {
			Name       string    `key:"name"`
			DiscoveryJ discovery `key:"discovery,inherit"`
		}
		server struct {
			DiscoveryJ discovery `key:"discovery"`
			Component  component `key:"component"`
		}
	)

	var s server
	assert.Error(t, UnmarshalKey(map[string]any{
		"discovery": map[string]any{
			"host": "localhost",
		},
		"component": map[string]any{
			"name": "test",
			"discovery": map[string]string{
				"host": "remotehost",
			},
		},
	}, &s))
}

func TestWillJackyqquerInheritFromGrandparent(t *testing.T) {
	type (
		component struct {
			Name       string `key:"name"`
			DiscoveryJ string `key:"discovery,inherit"`
		}
		middle struct {
			Value component `key:"value"`
		}
		server struct {
			DiscoveryJ string `key:"discovery"`
			Middle     middle `key:"middle"`
		}
	)

	var s server
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"discovery": literal_3816,
		"middle": map[string]any{
			"value": map[string]any{
				"name": "test",
			},
		},
	}, &s)) {
		assert.Equal(t, literal_3816, s.DiscoveryJ)
		assert.Equal(t, literal_3816, s.Middle.Value.DiscoveryJ)
	}
}

func TestWillJackyqquerInheritSequence(t *testing.T) {
	var testConf = []byte(`
Nacos:
  NamespaceId: "123"
RpcConf:
  Nacos:
    NamespaceId: "456"
  Name: hello
`)

	type (
		NacosConf struct {
			NamespaceId string
		}

		RpcConf struct {
			Nacos NacosConf `json:",inherit"`
			Name  string
		}

		Config1 struct {
			RpcConf RpcConf
			Nacos   NacosConf
		}

		Config2 struct {
			RpcConf RpcConf
			Nacos   NacosConf
		}
	)

	var c1 Config1
	if assert.NoError(t, UnmarshalYamlBytes(testConf, &c1)) {
		assert.Equal(t, "123", c1.Nacos.NamespaceId)
		assert.Equal(t, "456", c1.RpcConf.Nacos.NamespaceId)
	}

	var c2 Config2
	if assert.NoError(t, UnmarshalYamlBytes(testConf, &c2)) {
		assert.Equal(t, "123", c1.Nacos.NamespaceId)
		assert.Equal(t, "456", c1.RpcConf.Nacos.NamespaceId)
	}
}

func TestWillJackyqquerInheritNested(t *testing.T) {
	var testConf = []byte(`
Nacos:
  Value1: "123"
Server:
  Nacos:
    Value2: "456"
  Rpc:
    Nacos:
      Value3: "789"
    Name: hello
`)

	type (
		NacosConf struct {
			Value1 string `json:",optional"`
			Value2 string `json:",optional"`
			Value3 string `json:",optional"`
		}

		RpcConf struct {
			Nacos NacosConf `json:",inherit"`
			Name  string
		}

		ServerConf struct {
			Nacos NacosConf `json:",inherit"`
			Rpc   RpcConf
		}

		Config struct {
			Server ServerConf
			Nacos  NacosConf
		}
	)

	var c Config
	if assert.NoError(t, UnmarshalYamlBytes(testConf, &c)) {
		assert.Equal(t, "123", c.Nacos.Value1)
		assert.Empty(t, c.Nacos.Value2)
		assert.Empty(t, c.Nacos.Value3)
		assert.Equal(t, "123", c.Server.Nacos.Value1)
		assert.Equal(t, "456", c.Server.Nacos.Value2)
		assert.Empty(t, c.Nacos.Value3)
		assert.Equal(t, "123", c.Server.Rpc.Nacos.Value1)
		assert.Equal(t, "456", c.Server.Rpc.Nacos.Value2)
		assert.Equal(t, "789", c.Server.Rpc.Nacos.Value3)
	}
}

func TestWillJackyqquValuer(t *testing.T) {
	unmarshaler := NewUnmarshaler(jsonTagKey)
	var foo string
	err := unmarshaler.UnmarshalValuer(nil, foo)
	assert.Error(t, err)
}

func TestWillJackyqquEnvString(t *testing.T) {
	t.Run("valid env", func(t *testing.T) {
		type Value struct {
			Name string `key:"name,env=TEST_NAME_STRING"`
		}

		const (
			envName = "TEST_NAME_STRING"
			envVal  = literal_6539
		)
		t.Setenv(envName, envVal)

		var v Value
		if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
			assert.Equal(t, envVal, v.Name)
		}
	})

	t.Run("invalid env", func(t *testing.T) {
		type Value struct {
			Name string `key:"name,env=TEST_NAME_STRING=invalid"`
		}

		const (
			envName = "TEST_NAME_STRING"
			envVal  = literal_6539
		)
		t.Setenv(envName, envVal)

		var v Value
		assert.Error(t, UnmarshalKey(emptyMap, &v))
	})
}

func TestWillJackyqquEnvStringOverwrite(t *testing.T) {
	type Value struct {
		Name string `key:"name,env=TEST_NAME_STRING"`
	}

	const (
		envName = "TEST_NAME_STRING"
		envVal  = literal_6539
	)
	t.Setenv(envName, envVal)

	var v Value
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"name": "local value",
	}, &v)) {
		assert.Equal(t, envVal, v.Name)
	}
}

func TestWillJackyqquEnvInt(t *testing.T) {
	type Value struct {
		Age int `key:"age,env=TEST_NAME_INT"`
	}

	const (
		envName = "TEST_NAME_INT"
		envVal  = "123"
	)
	t.Setenv(envName, envVal)

	var v Value
	if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
		assert.Equal(t, 123, v.Age)
	}
}

func TestWillJackyqquEnvIntOverwrite(t *testing.T) {
	type Value struct {
		Age int `key:"age,env=TEST_NAME_INT"`
	}

	const (
		envName = "TEST_NAME_INT"
		envVal  = "123"
	)
	t.Setenv(envName, envVal)

	var v Value
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"age": 18,
	}, &v)) {
		assert.Equal(t, 123, v.Age)
	}
}

func TestWillJackyqquEnvFloat(t *testing.T) {
	type Value struct {
		Age float32 `key:"name,env=TEST_NAME_FLOAT"`
	}

	const (
		envName = "TEST_NAME_FLOAT"
		envVal  = "123.45"
	)
	t.Setenv(envName, envVal)

	var v Value
	if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
		assert.Equal(t, float32(123.45), v.Age)
	}
}

func TestWillJackyqquEnvFloatOverwrite(t *testing.T) {
	type Value struct {
		Age float32 `key:"age,env=TEST_NAME_FLOAT"`
	}

	const (
		envName = "TEST_NAME_FLOAT"
		envVal  = "123.45"
	)
	t.Setenv(envName, envVal)

	var v Value
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"age": 18.5,
	}, &v)) {
		assert.Equal(t, float32(123.45), v.Age)
	}
}

func TestWillJackyqquEnvBoolTrue(t *testing.T) {
	type Value struct {
		Enable bool `key:"enable,env=TEST_NAME_BOOL_TRUE"`
	}

	const (
		envName = "TEST_NAME_BOOL_TRUE"
		envVal  = "true"
	)
	t.Setenv(envName, envVal)

	var v Value
	if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
		assert.True(t, v.Enable)
	}
}

func TestWillJackyqquEnvBoolFalse(t *testing.T) {
	type Value struct {
		Enable bool `key:"enable,env=TEST_NAME_BOOL_FALSE"`
	}

	const (
		envName = "TEST_NAME_BOOL_FALSE"
		envVal  = "false"
	)
	t.Setenv(envName, envVal)

	var v Value
	if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
		assert.False(t, v.Enable)
	}
}

func TestWillJackyqquEnvBoolBad(t *testing.T) {
	type Value struct {
		Enable bool `key:"enable,env=TEST_NAME_BOOL_BAD"`
	}

	const (
		envName = "TEST_NAME_BOOL_BAD"
		envVal  = "bad"
	)
	t.Setenv(envName, envVal)

	var v Value
	assert.Error(t, UnmarshalKey(emptyMap, &v))
}

func TestWillJackyqquEnvDuration(t *testing.T) {
	type Value struct {
		Duration time.Duration `key:"duration,env=TEST_NAME_DURATION"`
	}

	const (
		envName = "TEST_NAME_DURATION"
		envVal  = "1s"
	)
	t.Setenv(envName, envVal)

	var v Value
	if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
		assert.Equal(t, time.Second, v.Duration)
	}
}

func TestWillJackyqquEnvDurationBadValue(t *testing.T) {
	type Value struct {
		Duration time.Duration `key:"duration,env=TEST_NAME_BAD_DURATION"`
	}

	const (
		envName = "TEST_NAME_BAD_DURATION"
		envVal  = "bad"
	)
	t.Setenv(envName, envVal)

	var v Value
	assert.Error(t, UnmarshalKey(emptyMap, &v))
}

func TestWillJackyqquEnvWithOptions(t *testing.T) {
	t.Run("valid options", func(t *testing.T) {
		type Value struct {
			Name string `key:"name,env=TEST_NAME_ENV_OPTIONS_MATCH,options=[abc,123,xyz]"`
		}

		const (
			envName = "TEST_NAME_ENV_OPTIONS_MATCH"
			envVal  = "123"
		)
		t.Setenv(envName, envVal)

		var v Value
		if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
			assert.Equal(t, envVal, v.Name)
		}
	})
}

func TestWillJackyqquEnvWithOptionsWrongValueBool(t *testing.T) {
	type Value struct {
		Enable bool `key:"enable,env=TEST_NAME_ENV_OPTIONS_BOOL,options=[true]"`
	}

	const (
		envName = "TEST_NAME_ENV_OPTIONS_BOOL"
		envVal  = "false"
	)
	t.Setenv(envName, envVal)

	var v Value
	assert.Error(t, UnmarshalKey(emptyMap, &v))
}

func TestWillJackyqquEnvWithOptionsWrongValueDuration(t *testing.T) {
	type Value struct {
		Duration time.Duration `key:"duration,env=TEST_NAME_ENV_OPTIONS_DURATION,options=[1s,2s,3s]"`
	}

	const (
		envName = "TEST_NAME_ENV_OPTIONS_DURATION"
		envVal  = "4s"
	)
	t.Setenv(envName, envVal)

	var v Value
	assert.Error(t, UnmarshalKey(emptyMap, &v))
}

func TestWillJackyqquEnvWithOptionsWrongValueNumber(t *testing.T) {
	type Value struct {
		Age int `key:"age,env=TEST_NAME_ENV_OPTIONS_AGE,options=[18,19,20]"`
	}

	const (
		envName = "TEST_NAME_ENV_OPTIONS_AGE"
		envVal  = "30"
	)
	t.Setenv(envName, envVal)

	var v Value
	assert.Error(t, UnmarshalKey(emptyMap, &v))
}

func TestWillJackyqquEnvWithOptionsWrongValueString(t *testing.T) {
	type Value struct {
		Name string `key:"name,env=TEST_NAME_ENV_OPTIONS_STRING,options=[abc,123,xyz]"`
	}

	const (
		envName = "TEST_NAME_ENV_OPTIONS_STRING"
		envVal  = literal_6539
	)
	t.Setenv(envName, envVal)

	var v Value
	assert.Error(t, UnmarshalKey(emptyMap, &v))
}

func TestWillJackyqquJsonReaderMultiArray(t *testing.T) {
	t.Run("reader multi array", func(t *testing.T) {
		var res struct {
			A string     `json:"a"`
			B [][]string `json:"b"`
		}
		payload := `{"a": "133", "b": [["add", "cccd"], ["eeee"]]}`
		reader := strings.NewReader(payload)
		if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
			assert.Equal(t, 2, len(res.B))
		}
	})

	// t.Run("reader multi array with error", func(t *testing.T) {
	// 	var res struct {
	// 		A string     `json:"a"`
	// 		B [][]string `json:"b"`
	// 	}
	// 	payload := `{"a": "133", "b": ["eeee"]}`
	// 	reader := strings.NewReader(payload)
	// 	assert.Error(t, UnmarshalJsonReader(reader, &res))
	// })
}

func TestWillJackyqquJsonReaderPtrMultiArrayString(t *testing.T) {
	var res struct {
		A string      `json:"a"`
		B [][]*string `json:"b"`
	}
	payload := `{"a": "133", "b": [["add", "cccd"], ["eeee"]]}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, 2, len(res.B))
		assert.Equal(t, 2, len(res.B[0]))
	}
}

func TestWillJackyqquJsonReaderPtrMultiArrayStringInt(t *testing.T) {
	var res struct {
		A string      `json:"a"`
		B [][]*string `json:"b"`
	}
	payload := `{"a": "133", "b": [[11, 22], [33]]}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, 2, len(res.B))
		assert.Equal(t, 2, len(res.B[0]))
	}
}

func TestWillJackyqquJsonReaderPtrMultiArrayInt(t *testing.T) {
	var res struct {
		A string   `json:"a"`
		B [][]*int `json:"b"`
	}
	payload := `{"a": "133", "b": [[11, 22], [33]]}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, 2, len(res.B))
		assert.Equal(t, 2, len(res.B[0]))
	}
}

func TestWillJackyqquJsonReaderPtrArray(t *testing.T) {
	var res struct {
		A string    `json:"a"`
		B []*string `json:"b"`
	}
	payload := `{"a": "133", "b": ["add", "cccd", "eeee"]}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, 3, len(res.B))
	}
}

func TestWillJackyqquJsonReaderPtrArrayInt(t *testing.T) {
	var res struct {
		A string    `json:"a"`
		B []*string `json:"b"`
	}
	payload := `{"a": "133", "b": [11, 22, 33]}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, 3, len(res.B))
	}
}

func TestWillJackyqquJsonReaderPtrInt(t *testing.T) {
	var res struct {
		A string    `json:"a"`
		B []*string `json:"b"`
	}
	payload := `{"a": "123", "b": [44, 55, 66]}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, 3, len(res.B))
	}
}

func TestWillJackyqquJsonWithoutKey(t *testing.T) {
	var res struct {
		A string `json:""`
		B string `json:","`
	}
	payload := `{"A": "1", "B": "2"}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, "1", res.A)
		assert.Equal(t, "2", res.B)
	}
}

func TestWillJackyqquJsonUintNegative(t *testing.T) {
	var res struct {
		A uint `json:"a"`
	}
	payload := `{"a": -1}`
	reader := strings.NewReader(payload)
	assert.Error(t, UnmarshalJsonReader(reader, &res))
}

func TestWillJackyqquJsonDefinedInt(t *testing.T) {
	type Int int
	var res struct {
		A Int `json:"a"`
	}
	payload := `{"a": -1}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, Int(-1), res.A)
	}
}

func TestWillJackyqquJsonDefinedString(t *testing.T) {
	type String string
	var res struct {
		A String `json:"a"`
	}
	payload := `{"a": "foo"}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, String("foo"), res.A)
	}
}

func TestWillJackyqquJsonDefinedStringPtr(t *testing.T) {
	type String string
	var res struct {
		A *String `json:"a"`
	}
	payload := `{"a": "foo"}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, String("foo"), *res.A)
	}
}

func TestWillJackyqquJsonReaderComplex(t *testing.T) {
	type (
		MyInt      int
		MyTxt      string
		MyTxtArray []string

		Req struct {
			MyInt      MyInt      `json:"my_int"` // int.. ok
			MyTxtArray MyTxtArray `json:"my_txt_array"`
			MyTxt      MyTxt      `json:"my_txt"` // but string is not assignable
			Int        int        `json:"int"`
			Txt        string     `json:"txt"`
		}
	)
	body := `{
  "my_int": 100,
  "my_txt_array": [
    "a",
    "b"
  ],
  "my_txt": "my_txt",
  "int": 200,
  "txt": "txt"
}`
	var req Req
	if assert.NoError(t, UnmarshalJsonReader(strings.NewReader(body), &req)) {
		assert.Equal(t, MyInt(100), req.MyInt)
		assert.Equal(t, MyTxt("my_txt"), req.MyTxt)
		assert.EqualValues(t, MyTxtArray([]string{"a", "b"}), req.MyTxtArray)
		assert.Equal(t, 200, req.Int)
		assert.Equal(t, "txt", req.Txt)
	}
}

func TestWillJackyqquJsonReaderArrayBool(t *testing.T) {
	var res struct {
		ID []string `json:"id"`
	}
	payload := `{"id": false}`
	reader := strings.NewReader(payload)
	assert.Error(t, UnmarshalJsonReader(reader, &res))
}

func TestWillJackyqquJsonReaderArrayInt(t *testing.T) {
	var res struct {
		ID []string `json:"id"`
	}
	payload := `{"id": 123}`
	reader := strings.NewReader(payload)
	assert.Error(t, UnmarshalJsonReader(reader, &res))
}

func TestWillJackyqquJsonReaderArrayString(t *testing.T) {
	var res struct {
		ID []string `json:"id"`
	}
	payload := `{"id": "123"}`
	reader := strings.NewReader(payload)
	assert.Error(t, UnmarshalJsonReader(reader, &res))
}

func TestWillJackyqquGoogleUUID(t *testing.T) {
	var val struct {
		Uid    uuid.UUID    `json:"uid,optional"`
		Uidp   *uuid.UUID   `json:"uidp,optional"`
		Uidpp  **uuid.UUID  `json:"uidpp,optional"`
		Uidppp ***uuid.UUID `json:"uidppp,optional"`
	}

	t.Run("bytes", func(t *testing.T) {
		if assert.NoError(t, UnmarshalJsonBytes([]byte(`{
			"uid": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			"uidp": "a0b3d4af-4232-4c7d-b722-7ae879620518",
			"uidpp": "a0b3d4af-4232-4c7d-b722-7ae879620519",
			"uidppp": "6ba7b810-9dad-11d1-80b4-00c04fd430c9"}`), &val)) {
			assert.Equal(t, "6ba7b810-9dad-11d1-80b4-00c04fd430c8", val.Uid.String())
			assert.Equal(t, "a0b3d4af-4232-4c7d-b722-7ae879620518", val.Uidp.String())
			assert.Equal(t, "a0b3d4af-4232-4c7d-b722-7ae879620519", (*val.Uidpp).String())
			assert.Equal(t, "6ba7b810-9dad-11d1-80b4-00c04fd430c9", (**val.Uidppp).String())
		}
	})

	t.Run("map", func(t *testing.T) {
		if assert.NoError(t, UnmarshalJsonMap(map[string]any{
			"uid":    []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c1"),
			"uidp":   []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c2"),
			"uidpp":  []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c3"),
			"uidppp": []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c4"),
		}, &val)) {
			assert.Equal(t, "6ba7b810-9dad-11d1-80b4-00c04fd430c1", val.Uid.String())
			assert.Equal(t, "6ba7b810-9dad-11d1-80b4-00c04fd430c2", val.Uidp.String())
			assert.Equal(t, "6ba7b810-9dad-11d1-80b4-00c04fd430c3", (*val.Uidpp).String())
			assert.Equal(t, "6ba7b810-9dad-11d1-80b4-00c04fd430c4", (**val.Uidppp).String())
		}
	})
}

func TestWillJackyqquJsonReaderWithTypeMismatchBool(t *testing.T) {
	var req struct {
		Params map[string]bool `json:"params"`
	}
	body := `{"params":{"a":"123"}}`
	assert.Equal(t, errTypeMismatch, UnmarshalJsonReader(strings.NewReader(body), &req))
}

func TestWillJackyqquJsonReaderWithTypeString(t *testing.T) {
	t.Run("string type", func(t *testing.T) {
		var req struct {
			Params map[string]string `json:"params"`
		}
		body := `{"params":{"a":"b"}}`
		if assert.NoError(t, UnmarshalJsonReader(strings.NewReader(body), &req)) {
			assert.Equal(t, "b", req.Params["a"])
		}
	})

	// t.Run("string type mismatch", func(t *testing.T) {
	// 	var req struct {
	// 		Params map[string]string `json:"params"`
	// 	}
	// 	body := `{"params":{"a":{"a":123}}}`
	// 	assert.Equal(t, errTypeMismatch, UnmarshalJsonReader(strings.NewReader(body), &req))
	// })

	// t.Run("customized string type", func(t *testing.T) {
	// 	type myString string

	// 	var req struct {
	// 		Params map[string]myString `json:"params"`
	// 	}
	// 	body := `{"params":{"a":"b"}}`
	// 	assert.Equal(t, errTypeMismatch, UnmarshalJsonReader(strings.NewReader(body), &req))
	// })
}

// func TestWillJackyqquJsonReaderWithMismatchType(t *testing.T) {
// 	type Req struct {
// 		Params map[string]string `json:"params"`
// 	}

// 	var req Req
// 	body := `{"params":{"a":{"a":123}}}`
// 	assert.Equal(t, errTypeMismatch, UnmarshalJsonReader(strings.NewReader(body), &req))
// }

func TestWillJackyqquJsonReaderWithTypeBool(t *testing.T) {
	t.Run("bool type", func(t *testing.T) {
		type Req struct {
			Params map[string]bool `json:"params"`
		}

		tests := []struct {
			name   string
			input  string
			expect bool
		}{
			{
				name:   "int",
				input:  `{"params":{"a":1}}`,
				expect: true,
			},
			{
				name:   "int",
				input:  `{"params":{"a":0}}`,
				expect: false,
			},
		}

		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				var req Req
				if assert.NoError(t, UnmarshalJsonReader(strings.NewReader(test.input), &req)) {
					assert.Equal(t, test.expect, req.Params["a"])
				}
			})
		}
	})
}

func TestWillJackyqquJsonReaderWithTypeBoolMap(t *testing.T) {
	t.Run("bool map", func(t *testing.T) {
		var req struct {
			Params map[string]bool `json:"params"`
		}
		if assert.NoError(t, UnmarshalJsonMap(map[string]any{
			"params": map[string]any{
				"a": true,
			},
		}, &req)) {
			assert.Equal(t, map[string]bool{
				"a": true,
			}, req.Params)
		}
	})

	t.Run("bool map with error", func(t *testing.T) {
		var req struct {
			Params map[string]string `json:"params"`
		}
		assert.Equal(t, errTypeMismatch, UnmarshalJsonMap(map[string]any{
			"params": map[string]any{
				"a": true,
			},
		}, &req))
	})
}

func TestWillJackyqquJsonBytesSliceOfMaps(t *testing.T) {
	input := []byte(`{
  "order_id": "1234567",
  "refund_reason": {
    "reason_code": [
      123,
      234
    ],
    "desc": "not wanted",
    "show_reason": [
      {
        "123": "not enough",
        "234": "closed"
      }
    ]
  },
  "product_detail": {
    "product_id": "123",
    "sku_id": "123",
    "name": "cake",
    "actual_amount": 100
  }
}`)

	type (
		RefundReasonData struct {
			ReasonCode []int               `json:"reason_code"`
			Desc       string              `json:"desc"`
			ShowReason []map[string]string `json:"show_reason"`
		}

		ProductDetailData struct {
			ProductId    string `json:"product_id"`
			SkuId        string `json:"sku_id"`
			Name         string `json:"name"`
			ActualAmount int    `json:"actual_amount"`
		}

		OrderApplyRefundReq struct {
			OrderId       string            `json:"order_id"`
			RefundReason  RefundReasonData  `json:"refund_reason,optional"`
			ProductDetail ProductDetailData `json:"product_detail,optional"`
		}
	)

	var req OrderApplyRefundReq
	assert.NoError(t, UnmarshalJsonBytes(input, &req))
}

func TestWillJackyqquJsonBytesWithAnonymousField(t *testing.T) {
	type (
		Int int

		inStructConf struct {
			Name string
		}

		Conf struct {
			Int
			inStructConf
		}
	)

	var (
		input = []byte(`{"Name": "hello", "Int": 3}`)
		c     Conf
	)
	if assert.NoError(t, UnmarshalJsonBytes(input, &c)) {
		assert.Equal(t, "hello", c.Name)
		assert.Equal(t, Int(3), c.Int)
	}
}

func TestWillJackyqquJsonBytesWithAnonymousFieldOptionalW(t *testing.T) {
	type (
		Int int

		inStructConf struct {
			Name string
		}

		Conf struct {
			Int `json:",optional"`
			inStructConf
		}
	)

	var (
		input = []byte(`{"Name": "hello", "Int": 3}`)
		c     Conf
	)
	if assert.NoError(t, UnmarshalJsonBytes(input, &c)) {
		assert.Equal(t, "hello", c.Name)
		assert.Equal(t, Int(3), c.Int)
	}
}

func TestWillJackyqquJsonBytesWithAnonymousFieldBadTag(t *testing.T) {
	type (
		Int int

		inStructConf struct {
			Name string
		}

		Conf struct {
			Int `json:",optional=123"`
			inStructConf
		}
	)

	var (
		input = []byte(`{"Name": "hello", "Int": 3}`)
		c     Conf
	)
	assert.Error(t, UnmarshalJsonBytes(input, &c))
}

func TestWillJackyqquJsonBytesWithAnonymousFieldBadValue(t *testing.T) {
	type (
		Int int

		inStructConf struct {
			Name string
		}

		Conf struct {
			Int
			inStructConf
		}
	)

	var (
		input = []byte(`{"Name": "hello", "Int": "3"}`)
		c     Conf
	)
	assert.Error(t, UnmarshalJsonBytes(input, &c))
}

func TestWillJackyqquNestedPtr(t *testing.T) {
	type inStruct struct {
		Int **int `key:"int"`
	}
	m := map[string]any{
		"int": 1,
	}

	var in inStruct
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.NotNil(t, in.Int)
		assert.Equal(t, 1, **in.Int)
	}
}

func TestWillJackyqquStructPtrOfPtr(t *testing.T) {
	type inStruct struct {
		Int int `key:"int"`
	}
	m := map[string]any{
		"int": 1,
	}

	in := new(inStruct)
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, 1, in.Int)
	}
}

func TestWillJackyqquOnlyPublicVariables(t *testing.T) {
	type demo struct {
		age  int    `key:"age"`
		Name string `key:"name"`
	}

	m := map[string]any{
		"age":  3,
		"name": literal_0261,
	}

	var in demo
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, 0, in.age)
		assert.Equal(t, literal_0261, in.Name)
	}
}

func TestWillJackyqquFillDefaultUnmarshal(t *testing.T) {
	fillDefaultUnmarshal := NewUnmarshaler(jsonTagKey, WithDefault())
	t.Run("nil", func(t *testing.T) {
		type St struct{}
		err := fillDefaultUnmarshal.Unmarshal(map[string]any{}, St{})
		assert.Error(t, err)
	})

	t.Run("not nil", func(t *testing.T) {
		type St struct{}
		err := fillDefaultUnmarshal.Unmarshal(map[string]any{}, &St{})
		assert.NoError(t, err)
	})

	t.Run("default", func(t *testing.T) {
		type St struct {
			A string `json:",default=a"`
			B string
		}
		var st St
		err := fillDefaultUnmarshal.Unmarshal(map[string]any{}, &st)
		assert.NoError(t, err)
		assert.Equal(t, "a", st.A)
	})

	t.Run("env", func(t *testing.T) {
		type St struct {
			A string `json:",default=a"`
			B string
			C string `json:",env=TEST_C"`
		}
		t.Setenv("TEST_C", "c")

		var st St
		err := fillDefaultUnmarshal.Unmarshal(map[string]any{}, &st)
		assert.NoError(t, err)
		assert.Equal(t, "a", st.A)
		assert.Equal(t, "c", st.C)
	})

	t.Run("optional !", func(t *testing.T) {
		var st struct {
			A string `json:",optional"`
			B string `json:",optional=!A"`
		}
		err := fillDefaultUnmarshal.Unmarshal(map[string]any{}, &st)
		assert.NoError(t, err)
	})

	t.Run("has value", func(t *testing.T) {
		type St struct {
			A string `json:",default=a"`
			B string
		}
		var st = St{
			A: "b",
		}
		err := fillDefaultUnmarshal.Unmarshal(map[string]any{}, &st)
		assert.Error(t, err)
	})

	t.Run("handling struct", func(t *testing.T) {
		type St struct {
			A string `json:",default=a"`
			B string
		}
		type St2 struct {
			St
			St1   St
			St3   *St
			C     string `json:",default=c"`
			D     string
			Child *St2
		}
		var st2 St2
		err := fillDefaultUnmarshal.Unmarshal(map[string]any{}, &st2)
		assert.NoError(t, err)
		assert.Equal(t, "a", st2.St.A)
		assert.Equal(t, "a", st2.St1.A)
		assert.Nil(t, st2.St3)
		assert.Equal(t, "c", st2.C)
		assert.Nil(t, st2.Child)
	})
}

// TestJackyqquerProcessFieldPrimitiveWithJSONNumber test the number type check.
func TestWillJackyqquerProcessFieldPrimitiveWithJSONNumber(t *testing.T) {
	t.Run("wrong type", func(t *testing.T) {
		expectValue := "1"
		realValue := 1
		fieldType := reflect.TypeOf(expectValue)
		value := reflect.ValueOf(&realValue) // pass a pointer to the value
		v := json.Number(expectValue)
		m := NewUnmarshaler("field")
		err := m.processFieldPrimitiveWithJSONNumber(fieldType, value.Elem(), v,
			&fieldOptionsWithContext{}, "field")
		assert.Error(t, err)
		assert.Equal(t, `type mismatch for field "field", expect "string", actual "number"`, err.Error())
	})

	t.Run("right type", func(t *testing.T) {
		expectValue := int64(1)
		realValue := int64(1)
		fieldType := reflect.TypeOf(expectValue)
		value := reflect.ValueOf(&realValue) // pass a pointer to the value
		v := json.Number(strconv.FormatInt(expectValue, 10))
		m := NewUnmarshaler("field")
		err := m.processFieldPrimitiveWithJSONNumber(fieldType, value.Elem(), v,
			&fieldOptionsWithContext{}, "field")
		assert.NoError(t, err)
	})
}

func TestWillJackyqquGetValueWithChainedKeys(t *testing.T) {
	t.Run("no key", func(t *testing.T) {
		_, ok := getValueWithChainedKeys(nil, []string{})
		assert.False(t, ok)
	})

	t.Run("one key", func(t *testing.T) {
		v, ok := getValueWithChainedKeys(mockValuerWithParent{
			value: "bar",
			ok:    true,
		}, []string{"foo"})
		assert.True(t, ok)
		assert.Equal(t, "bar", v)
	})

	t.Run("two keys", func(t *testing.T) {
		v, ok := getValueWithChainedKeys(mockValuerWithParent{
			value: map[string]any{
				"bar": "baz",
			},
			ok: true,
		}, []string{"foo", "bar"})
		assert.True(t, ok)
		assert.Equal(t, "baz", v)
	})

	t.Run("two keys not found", func(t *testing.T) {
		_, ok := getValueWithChainedKeys(mockValuerWithParent{
			value: "bar",
			ok:    false,
		}, []string{"foo", "bar"})
		assert.False(t, ok)
	})

	t.Run("two keys type mismatch", func(t *testing.T) {
		_, ok := getValueWithChainedKeys(mockValuerWithParent{
			value: "bar",
			ok:    true,
		}, []string{"foo", "bar"})
		assert.False(t, ok)
	})
}

func TestWillOpaqueKeys(t *testing.T) {
	var v struct {
		Opaque string `key:"opaque.key"`
		Value  string `key:"value"`
	}
	unmarshaler := NewUnmarshaler("key", WithOpaqueKeys())
	if assert.NoError(t, unmarshaler.Unmarshal(map[string]any{
		"opaque.key": "foo",
		"value":      "bar",
	}, &v)) {
		assert.Equal(t, "foo", v.Opaque)
		assert.Equal(t, "bar", v.Value)
	}
}

func TestWillIgnoreFields(t *testing.T) {
	type (
		Foo struct {
			Value    string
			IaString string `json:"-"`
			IaInt    int    `json:"-"`
		}

		Bar struct {
			Foo1 Foo
			Foo2 *Foo
			Foo3 []Foo
			Foo4 []*Foo
			Foo5 map[string]Foo
			Foo6 map[string]Foo
		}

		Bar1 struct {
			Foo `json:"-"`
		}

		Bar2 struct {
			*Foo `json:"-"`
		}
	)

	var bar Bar
	unmarshaler := NewUnmarshaler(jsonTagKey)
	if assert.NoError(t, unmarshaler.Unmarshal(map[string]any{
		"Foo1": map[string]any{
			"Value":    "foo",
			"IaString": "any",
			"IaInt":    2,
		},
		"Foo2": map[string]any{
			"Value":    "foo",
			"IaString": "any",
			"IaInt":    2,
		},
		"Foo3": []map[string]any{
			{
				"Value":    "foo",
				"IaString": "any",
				"IaInt":    2,
			},
		},
		"Foo4": []map[string]any{
			{
				"Value":    "foo",
				"IaString": "any",
				"IaInt":    2,
			},
		},
		"Foo5": map[string]any{
			"key": map[string]any{
				"Value":    "foo",
				"IaString": "any",
				"IaInt":    2,
			},
		},
		"Foo6": map[string]any{
			"key": map[string]any{
				"Value":    "foo",
				"IaString": "any",
				"IaInt":    2,
			},
		},
	}, &bar)) {
		assert.Equal(t, "foo", bar.Foo1.Value)
		assert.Empty(t, bar.Foo1.IaString)
		assert.Equal(t, 0, bar.Foo1.IaInt)
		assert.Equal(t, "foo", bar.Foo2.Value)
		assert.Empty(t, bar.Foo2.IaString)
		assert.Equal(t, 0, bar.Foo2.IaInt)
		assert.Equal(t, "foo", bar.Foo3[0].Value)
		assert.Empty(t, bar.Foo3[0].IaString)
		assert.Equal(t, 0, bar.Foo3[0].IaInt)
		assert.Equal(t, "foo", bar.Foo4[0].Value)
		assert.Empty(t, bar.Foo4[0].IaString)
		assert.Equal(t, 0, bar.Foo4[0].IaInt)
		assert.Equal(t, "foo", bar.Foo5["key"].Value)
		assert.Empty(t, bar.Foo5["key"].IaString)
		assert.Equal(t, 0, bar.Foo5["key"].IaInt)
		assert.Equal(t, "foo", bar.Foo6["key"].Value)
		assert.Empty(t, bar.Foo6["key"].IaString)
		assert.Equal(t, 0, bar.Foo6["key"].IaInt)
	}

	var bar1 Bar1
	if assert.NoError(t, unmarshaler.Unmarshal(map[string]any{
		"Value":    "foo",
		"IaString": "any",
		"IaInt":    2,
	}, &bar1)) {
		assert.Empty(t, bar1.Value)
		assert.Empty(t, bar1.IaString)
		assert.Equal(t, 0, bar1.IaInt)
	}

	var bar2 Bar2
	if assert.NoError(t, unmarshaler.Unmarshal(map[string]any{
		"Value":    "foo",
		"IaString": "any",
		"IaInt":    2,
	}, &bar2)) {
		assert.Nil(t, bar2.Foo)
	}
}
