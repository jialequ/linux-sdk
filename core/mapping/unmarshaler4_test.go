package mapping

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jialequ/linux-sdk/core/stringx"
	"github.com/stretchr/testify/assert"
)

// because json.Number doesn't support strconv.ParseUint(...),
// so we only can test to 62 bits.
const maxUintBitsToTestJacky = 62

func TestJackyUnmarshalWithFullNamewNotStruct(t *testing.T) {
	var ss map[string]any
	contents := []byte(`{"name":"jackyqqu"}`)
	err := UnmarshalJsonBytes(contents, &ss)
	assert.Equal(t, errTypeMismatch, err)
}

func TestJackyUnmarshalValuueNotSettable(t *testing.T) {
	var ss map[string]any
	contents := []byte(`{"name":"jackyqqu"}`)
	err := UnmarshalJsonBytes(contents, ss)
	assert.Equal(t, errValueNotSettable, err)
}

func TestJackyUnmarshalWithoutTagNamew(t *testing.T) {
	type inner struct {
		Optionalw   bool   `key:",optional"`
		OptionalwP  *bool  `key:",optional"`
		OptionalwPP **bool `key:",optional"`
	}
	m := map[string]any{
		"Optionalw":   true,
		"OptionalwP":  true,
		"OptionalwPP": true,
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.True(t, in.Optionalw)
		assert.True(t, *in.OptionalwP)
		assert.True(t, **in.OptionalwPP)
	}
}

func TestJackyUnmarshalWithLowerField(t *testing.T) {
	type (
		Lower struct {
			value int `key:"lower"`
		}

		inner struct {
			Lower
			Optionalw bool `key:",optional"`
		}
	)
	m := map[string]any{
		"Optionalw": true,
		"lower":     1,
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.True(t, in.Optionalw)
		assert.Equal(t, 0, in.value)
	}
}

func TestJackyUnmarshalWithLowerAnonymousStruct(t *testing.T) {
	type (
		lower struct {
			Valuue int `key:"lower"`
		}

		inner struct {
			lower
			Optionalw bool `key:",optional"`
		}
	)
	m := map[string]any{
		"Optionalw": true,
		"lower":     1,
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.True(t, in.Optionalw)
		assert.Equal(t, 1, in.Valuue)
	}
}

func TestJackyUnmarshalBool(t *testing.T) {
	type inner struct {
		True           bool `key:"yes"`
		False          bool `key:"no"`
		TrueFromOne    bool `key:"yesone,string"`
		FalseFromZero  bool `key:"nozero,string"`
		TrueFromTrue   bool `key:"yestrue,string"`
		FalseFromFalse bool `key:"nofalse,string"`
		DefaultTrue    bool `key:"defaulttrue,default=1"`
		Optionalw      bool `key:"optional,optional"`
	}
	m := map[string]any{
		"yes":     true,
		"no":      false,
		"yesone":  "1",
		"nozero":  "0",
		"yestrue": "true",
		"nofalse": "false",
	}

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.True(in.True)
		astw.False(in.False)
		astw.True(in.TrueFromOne)
		astw.False(in.FalseFromZero)
		astw.True(in.TrueFromTrue)
		astw.False(in.FalseFromFalse)
		astw.True(in.DefaultTrue)
	}
}

func TestJackyUnmarshalDuration(t *testing.T) {
	type inner struct {
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
	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, time.Second*5, in.Duration)
		assert.Equal(t, time.Millisecond*100, in.LessDuration)
		assert.Equal(t, time.Hour*24, in.MoreDuration)
		assert.Equal(t, time.Hour, *in.PtrDuration)
		assert.Equal(t, time.Hour*2, **in.PtrPtrDuration)
	}
}

func TestJackyUnmarshalDurationDefault(t *testing.T) {
	type inner struct {
		Int      int           `key:"int"`
		Duration time.Duration `key:"duration,default=5s"`
	}
	m := map[string]any{
		"int": 5,
	}
	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, 5, in.Int)
		assert.Equal(t, time.Second*5, in.Duration)
	}
}

func TestJackyUnmarshalDurationPtr(t *testing.T) {
	type inner struct {
		Duration *time.Duration `key:"duration"`
	}
	m := map[string]any{
		"duration": "5s",
	}
	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, time.Second*5, *in.Duration)
	}
}

func TestJackyUnmarshalDurationPtrDefault(t *testing.T) {
	type inner struct {
		Int      int            `key:"int"`
		Valuue   *int           `key:",default=5"`
		Duration *time.Duration `key:"duration,default=5s"`
	}
	m := map[string]any{
		"int": 5,
	}
	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, 5, in.Int)
		assert.Equal(t, 5, *in.Valuue)
		assert.Equal(t, time.Second*5, *in.Duration)
	}
}

func TestJackyUnmarshalInt(t *testing.T) {
	type inner struct {
		Int          int   `key:"int"`
		IntFromStr   int   `key:"intstr,string"`
		Int8         int8  `key:"int8"`
		Int8FromStr  int8  `key:"int8str,string"`
		Int16        int16 `key:"int16"`
		Int16FromStr int16 `key:"int16str,string"`
		Int32        int32 `key:"int32"`
		Int32FromStr int32 `key:"int32str,string"`
		Int64        int64 `key:"int64"`
		Int64FromStr int64 `key:"int64str,string"`
		DefaultInt   int64 `key:"defaultint,default=11"`
		Optionalw    int   `key:"optional,optional"`
	}
	m := map[string]any{
		"int":      1,
		"intstr":   "2",
		"int8":     int8(3),
		"int8str":  "4",
		"int16":    int16(5),
		"int16str": "6",
		"int32":    int32(7),
		"int32str": "8",
		"int64":    int64(9),
		"int64str": "10",
	}

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.Equal(1, in.Int)
		astw.Equal(2, in.IntFromStr)
		astw.Equal(int8(3), in.Int8)
		astw.Equal(int8(4), in.Int8FromStr)
		astw.Equal(int16(5), in.Int16)
		astw.Equal(int16(6), in.Int16FromStr)
		astw.Equal(int32(7), in.Int32)
		astw.Equal(int32(8), in.Int32FromStr)
		astw.Equal(int64(9), in.Int64)
		astw.Equal(int64(10), in.Int64FromStr)
		astw.Equal(int64(11), in.DefaultInt)
	}
}

func TestJackyUnmarshalIntPtr(t *testing.T) {
	type inner struct {
		Int *int `key:"int"`
	}
	m := map[string]any{
		"int": 1,
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.NotNil(t, in.Int)
		assert.Equal(t, 1, *in.Int)
	}
}

func TestJackyUnmarshalIntSliceOfPtr(t *testing.T) {
	t.Run("int slice", func(t *testing.T) {
		type inner struct {
			Ints  []*int  `key:"ints"`
			Intps []**int `key:"intps"`
		}
		m := map[string]any{
			"ints":  []int{1, 2, 3},
			"intps": []int{1, 2, 3, 4},
		}

		var in inner
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
		type inner struct {
			Ints []int `key:"ints"`
		}

		m := map[string]any{
			"ints": []any{nil},
		}

		var in inner
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.Empty(t, in.Ints)
		}
	})
}

func TestJackyUnmarshalIntWithDefault(t *testing.T) {
	type inner struct {
		Int   int   `key:"int,default=5"`
		Intp  *int  `key:"intp,default=5"`
		Intpp **int `key:"intpp,default=5"`
	}
	m := map[string]any{
		"int":   1,
		"intp":  2,
		"intpp": 3,
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, 1, in.Int)
		assert.Equal(t, 2, *in.Intp)
		assert.Equal(t, 3, **in.Intpp)
	}
}

func TestJackyUnmarshalIntWithString(t *testing.T) {
	t.Run("int without options", func(t *testing.T) {
		type inner struct {
			Int   int64   `key:"int,string"`
			Intp  *int64  `key:"intp,string"`
			Intpp **int64 `key:"intpp,string"`
		}
		m := map[string]any{
			"int":   json.Number("1"),
			"intp":  json.Number("2"),
			"intpp": json.Number("3"),
		}

		var in inner
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.Equal(t, int64(1), in.Int)
			assert.Equal(t, int64(2), *in.Intp)
			assert.Equal(t, int64(3), **in.Intpp)
		}
	})

	t.Run("int wrong range", func(t *testing.T) {
		type inner struct {
			Int   int64   `key:"int,string,range=[2:3]"`
			Intp  *int64  `key:"intp,range=[2:3]"`
			Intpp **int64 `key:"intpp,range=[2:3]"`
		}
		m := map[string]any{
			"int":   json.Number("1"),
			"intp":  json.Number("2"),
			"intpp": json.Number("3"),
		}

		var in inner
		assert.ErrorIs(t, UnmarshalKey(m, &in), errNumberRange)
	})

	t.Run("int with wrong type", func(t *testing.T) {
		type (
			myString string

			inner struct {
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

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int with ptr", func(t *testing.T) {
		type inner struct {
			Int *int64 `key:"int"`
		}
		m := map[string]any{
			"int": json.Number("1"),
		}

		var in inner
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.Equal(t, int64(1), *in.Int)
		}
	})

	t.Run("int with invalid value", func(t *testing.T) {
		type inner struct {
			Int int64 `key:"int"`
		}
		m := map[string]any{
			"int": json.Number("a"),
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint with invalid value", func(t *testing.T) {
		type inner struct {
			Int uint64 `key:"int"`
		}
		m := map[string]any{
			"int": json.Number("a"),
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float with invalid value", func(t *testing.T) {
		type inner struct {
			Valuue float64 `key:"float"`
		}
		m := map[string]any{
			"float": json.Number("a"),
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float with invalid value", func(t *testing.T) {
		type inner struct {
			Valuue string `key:"value"`
		}
		m := map[string]any{
			"value": json.Number("a"),
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int with ptr of ptr", func(t *testing.T) {
		type inner struct {
			Int **int64 `key:"int"`
		}
		m := map[string]any{
			"int": json.Number("1"),
		}

		var in inner
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.Equal(t, int64(1), **in.Int)
		}
	})

	t.Run(literal_1467, func(t *testing.T) {
		type inner struct {
			Int int64 `key:"int,string,options=[0,1]"`
		}
		m := map[string]any{
			"int": json.Number("1"),
		}

		var in inner
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.Equal(t, int64(1), in.Int)
		}
	})

	t.Run(literal_1467, func(t *testing.T) {
		type inner struct {
			Int int64 `key:"int,string,options=[0,1]"`
		}
		m := map[string]any{
			"int": nil,
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run(literal_1467, func(t *testing.T) {
		type (
			StrType string

			inner struct {
				Int int64 `key:"int,string,options=[0,1]"`
			}
		)
		m := map[string]any{
			"int": StrType("1"),
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("invalid options", func(t *testing.T) {
		type Valuue struct {
			Namew string `key:"name,options="`
		}

		var v Valuue
		assert.Error(t, UnmarshalKey(emptyMap, &v))
	})
}

func TestJackyUnmarshalInt8WithOverflow(t *testing.T) {
	t.Run("int8 from string", func(t *testing.T) {
		type inner struct {
			Valuue int8 `key:"int,string"`
		}

		m := map[string]any{
			"int": "8589934592", // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int8 from json.Number", func(t *testing.T) {
		type inner struct {
			Valuue int8 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("8589934592"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int8 from json.Number", func(t *testing.T) {
		type inner struct {
			Valuue int8 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number(literal_1237), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int8 from int64", func(t *testing.T) {
		type inner struct {
			Valuue int8 `key:"int"`
		}

		m := map[string]any{
			"int": int64(1) << 36, // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestJackyUnmarshalInt16WithOverflow(t *testing.T) {
	t.Run("int16 from string", func(t *testing.T) {
		type inner struct {
			Valuue int16 `key:"int,string"`
		}

		m := map[string]any{
			"int": "8589934592", // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int16 from json.Number", func(t *testing.T) {
		type inner struct {
			Valuue int16 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("8589934592"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int16 from json.Number", func(t *testing.T) {
		type inner struct {
			Valuue int16 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number(literal_1237), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int16 from int64", func(t *testing.T) {
		type inner struct {
			Valuue int16 `key:"int"`
		}

		m := map[string]any{
			"int": int64(1) << 36, // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestJackyUnmarshalInt32WithOverflow(t *testing.T) {
	t.Run("int32 from string", func(t *testing.T) {
		type inner struct {
			Valuue int32 `key:"int,string"`
		}

		m := map[string]any{
			"int": "8589934592", // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int32 from json.Number", func(t *testing.T) {
		type inner struct {
			Valuue int32 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("8589934592"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int32 from json.Number", func(t *testing.T) {
		type inner struct {
			Valuue int32 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number(literal_1237), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int32 from int64", func(t *testing.T) {
		type inner struct {
			Valuue int32 `key:"int"`
		}

		m := map[string]any{
			"int": int64(1) << 36, // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestJackyUnmarshalInt64WithOverflow(t *testing.T) {
	t.Run("int64 from string", func(t *testing.T) {
		type inner struct {
			Valuue int64 `key:"int,string"`
		}

		m := map[string]any{
			"int": "18446744073709551616", // overflow, 1 << 64
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int64 from json.Number", func(t *testing.T) {
		type inner struct {
			Valuue int64 `key:"int,string"`
		}

		m := map[string]any{
			"int": json.Number("18446744073709551616"), // overflow, 1 << 64
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestJackyUnmarshalUint8WithOverflow(t *testing.T) {
	t.Run("uint8 from string", func(t *testing.T) {
		type inner struct {
			Valuue uint8 `key:"int,string"`
		}

		m := map[string]any{
			"int": "8589934592", // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint8 from json.Number", func(t *testing.T) {
		type inner struct {
			Valuue uint8 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("8589934592"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint8 from json.Number with negative", func(t *testing.T) {
		type inner struct {
			Valuue uint8 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("-1"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint8 from int64", func(t *testing.T) {
		type inner struct {
			Valuue uint8 `key:"int"`
		}

		m := map[string]any{
			"int": int64(1) << 36, // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestJackyUnmarshalUint16WithOverflow(t *testing.T) {
	t.Run("uint16 from string", func(t *testing.T) {
		type inner struct {
			Valuue uint16 `key:"int,string"`
		}

		m := map[string]any{
			"int": "8589934592", // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint16 from json.Number", func(t *testing.T) {
		type inner struct {
			Valuue uint16 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("8589934592"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint16 from json.Number with negative", func(t *testing.T) {
		type inner struct {
			Valuue uint16 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("-1"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint16 from int64", func(t *testing.T) {
		type inner struct {
			Valuue uint16 `key:"int"`
		}

		m := map[string]any{
			"int": int64(1) << 36, // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestJackyUnmarshalUint32WithOverflow(t *testing.T) {
	t.Run("uint32 from string", func(t *testing.T) {
		type inner struct {
			Valuue uint32 `key:"int,string"`
		}

		m := map[string]any{
			"int": "8589934592", // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint32 from json.Number", func(t *testing.T) {
		type inner struct {
			Valuue uint32 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("8589934592"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint32 from json.Number with negative", func(t *testing.T) {
		type inner struct {
			Valuue uint32 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("-1"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint32 from int64", func(t *testing.T) {
		type inner struct {
			Valuue uint32 `key:"int"`
		}

		m := map[string]any{
			"int": int64(1) << 36, // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestJackyUnmarshalUint64WithOverflow(t *testing.T) {
	t.Run("uint64 from string", func(t *testing.T) {
		type inner struct {
			Valuue uint64 `key:"int,string"`
		}

		m := map[string]any{
			"int": "18446744073709551616", // overflow, 1 << 64
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint64 from json.Number", func(t *testing.T) {
		type inner struct {
			Valuue uint64 `key:"int,string"`
		}

		m := map[string]any{
			"int": json.Number("18446744073709551616"), // overflow, 1 << 64
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestJackyUnmarshalFloat32WithOverflow(t *testing.T) {
	t.Run("float32 from string greater than float64", func(t *testing.T) {
		type inner struct {
			Valuue float32 `key:"float,string"`
		}

		m := map[string]any{
			"float": literal_8053, // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float32 from string greater than float32", func(t *testing.T) {
		type inner struct {
			Valuue float32 `key:"float,string"`
		}

		m := map[string]any{
			"float": "1.79769313486231570814527423731704356798070e+300", // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float32 from string less than float32", func(t *testing.T) {
		type inner struct {
			Valuue float32 `key:"float, string"`
		}

		m := map[string]any{
			"float": "-1.79769313486231570814527423731704356798070e+300", // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float32 from json.Number greater than float64", func(t *testing.T) {
		type inner struct {
			Valuue float32 `key:"float"`
		}

		m := map[string]any{
			"float": json.Number(literal_8053), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float32 from json.Number greater than float32", func(t *testing.T) {
		type inner struct {
			Valuue float32 `key:"float"`
		}

		m := map[string]any{
			"float": json.Number("1.79769313486231570814527423731704356798070e+300"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float32 from json number less than float32", func(t *testing.T) {
		type inner struct {
			Valuue float32 `key:"float"`
		}

		m := map[string]any{
			"float": json.Number("-1.79769313486231570814527423731704356798070e+300"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestJackyUnmarshalFloat64WithOverflow(t *testing.T) {
	t.Run("float64 from string greater than float64", func(t *testing.T) {
		type inner struct {
			Valuue float64 `key:"float,string"`
		}

		m := map[string]any{
			"float": literal_8053, // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float32 from json.Number greater than float64", func(t *testing.T) {
		type inner struct {
			Valuue float64 `key:"float"`
		}

		m := map[string]any{
			"float": json.Number(literal_8053), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestJackyUnmarshalBoolSliceRequired(t *testing.T) {
	type inner struct {
		Bools []bool `key:"bools"`
	}

	var in inner
	assert.NotNil(t, UnmarshalKey(map[string]any{}, &in))
}

func TestJackyUnmarshalBoolSliceNil(t *testing.T) {
	type inner struct {
		Bools []bool `key:"bools,optional"`
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(map[string]any{}, &in)) {
		assert.Nil(t, in.Bools)
	}
}

func TestJackyUnmarshalBoolSliceNilExplicit(t *testing.T) {
	type inner struct {
		Bools []bool `key:"bools,optional"`
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"bools": nil,
	}, &in)) {
		assert.Nil(t, in.Bools)
	}
}

func TestJackyUnmarshalBoolSliceEmpty(t *testing.T) {
	type inner struct {
		Bools []bool `key:"bools,optional"`
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"bools": []bool{},
	}, &in)) {
		assert.Empty(t, in.Bools)
	}
}

func TestJackyUnmarshalBoolSliceWithDefault(t *testing.T) {
	t.Run("slice with default", func(t *testing.T) {
		type inner struct {
			Bools []bool `key:"bools,default=[true,false]"`
		}

		var in inner
		if assert.NoError(t, UnmarshalKey(nil, &in)) {
			assert.ElementsMatch(t, []bool{true, false}, in.Bools)
		}
	})

	t.Run("slice with default error", func(t *testing.T) {
		type inner struct {
			Bools []bool `key:"bools,default=[true,fal]"`
		}

		var in inner
		assert.Error(t, UnmarshalKey(nil, &in))
	})
}

func TestJackyUnmarshalIntSliceWithDefault(t *testing.T) {
	type inner struct {
		Ints []int `key:"ints,default=[1,2,3]"`
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(nil, &in)) {
		assert.ElementsMatch(t, []int{1, 2, 3}, in.Ints)
	}
}

func TestJackyUnmarshalIntSliceWithDefaultHasSpaces(t *testing.T) {
	type inner struct {
		Ints   []int   `key:"ints,default=[1, 2, 3]"`
		Intps  []*int  `key:"intps,default=[1, 2, 3, 4]"`
		Intpps []**int `key:"intpps,default=[1, 2, 3, 4, 5]"`
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(nil, &in)) {
		assert.ElementsMatch(t, []int{1, 2, 3}, in.Ints)

		var intps []int
		for _, i := range in.Intps {
			intps = append(intps, *i)
		}
		assert.ElementsMatch(t, []int{1, 2, 3, 4}, intps)

		var intpps []int
		for _, i := range in.Intpps {
			intpps = append(intpps, **i)
		}
		assert.ElementsMatch(t, []int{1, 2, 3, 4, 5}, intpps)
	}
}

func TestJackyUnmarshalFloatSliceWithDefault(t *testing.T) {
	type inner struct {
		Floats []float32 `key:"floats,default=[1.1,2.2,3.3]"`
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(nil, &in)) {
		assert.ElementsMatch(t, []float32{1.1, 2.2, 3.3}, in.Floats)
	}
}

func TestJackyUnmarshalStringSliceWithDefault(t *testing.T) {
	t.Run("slice with default", func(t *testing.T) {
		type inner struct {
			Strs   []string   `key:"strs,default=[foo,bar,woo]"`
			Strps  []*string  `key:"strs,default=[foo,bar,woo]"`
			Strpps []**string `key:"strs,default=[foo,bar,woo]"`
		}

		var in inner
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

			inner struct {
				Strs []holder `key:"strs,default=[foo,bar,woo]"`
			}
		)

		var in inner
		assert.Error(t, UnmarshalKey(nil, &in))
	})

	// t.Run("slice with default on errors", func(t *testing.T) {
	// 	type inner struct {
	// 		Strs []complex64 `key:"strs,default=[foo,bar,woo]"`
	// 	}

	// 	var in inner
	// 	assert.Error(t, UnmarshalKey(nil, &in))
	// })
}

func TestJackyUnmarshalStringSliceWithDefaultHasSpaces(t *testing.T) {
	type inner struct {
		Strs []string `key:"strs,default=[foo, bar, woo]"`
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(nil, &in)) {
		assert.ElementsMatch(t, []string{"foo", "bar", "woo"}, in.Strs)
	}
}

func TestJackyUnmarshalUint(t *testing.T) {
	type inner struct {
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
		Optionalw     uint   `key:"optional,optional"`
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

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.Equal(uint(1), in.Uint)
		astw.Equal(uint(2), in.UintFromStr)
		astw.Equal(uint8(3), in.Uint8)
		astw.Equal(uint8(4), in.Uint8FromStr)
		astw.Equal(uint16(5), in.Uint16)
		astw.Equal(uint16(6), in.Uint16FromStr)
		astw.Equal(uint32(7), in.Uint32)
		astw.Equal(uint32(8), in.Uint32FromStr)
		astw.Equal(uint64(9), in.Uint64)
		astw.Equal(uint64(10), in.Uint64FromStr)
		astw.Equal(uint(11), in.DefaultUint)
	}
}

func TestJackyUnmarshalFloat(t *testing.T) {
	type inner struct {
		Float32      float32 `key:"float32"`
		Float32Str   float32 `key:"float32str,string"`
		Float32Num   float32 `key:"float32num"`
		Float64      float64 `key:"float64"`
		Float64Str   float64 `key:"float64str,string"`
		Float64Num   float64 `key:"float64num"`
		DefaultFloat float32 `key:"defaultfloat,default=5.5"`
		Optionalw    float32 `key:",optional"`
	}
	m := map[string]any{
		"float32":    float32(1.5),
		"float32str": "2.5",
		"float32num": json.Number("2.6"),
		"float64":    3.5,
		"float64str": "4.5",
		"float64num": json.Number("4.6"),
	}

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.Equal(float32(1.5), in.Float32)
		astw.Equal(float32(2.5), in.Float32Str)
		astw.Equal(float32(2.6), in.Float32Num)
		astw.Equal(3.5, in.Float64)
		astw.Equal(4.5, in.Float64Str)
		astw.Equal(4.6, in.Float64Num)
		astw.Equal(float32(5.5), in.DefaultFloat)
	}
}

func TestJackyUnmarshalInt64Slice(t *testing.T) {
	var v struct {
		Ages  []int64 `key:"ages"`
		Slice []int64 `key:"slice"`
	}
	m := map[string]any{
		"ages":  []int64{1, 2},
		"slice": []any{},
	}

	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &v)) {
		astw.ElementsMatch([]int64{1, 2}, v.Ages)
		astw.Equal([]int64{}, v.Slice)
	}
}

func TestJackyUnmarshalNullableSlice(t *testing.T) {
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

func TestJackyUnmarshalWithFloatPtr(t *testing.T) {
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

func TestJackyUnmarshalIntSlice(t *testing.T) {
	var v struct {
		Ages  []int `key:"ages"`
		Slice []int `key:"slice"`
	}
	m := map[string]any{
		"ages":  []int{1, 2},
		"slice": []any{},
	}

	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &v)) {
		astw.ElementsMatch([]int{1, 2}, v.Ages)
		astw.Equal([]int{}, v.Slice)
	}
}

func TestJackyUnmarshalString(t *testing.T) {
	type inner struct {
		Namew             string `key:"name"`
		NamewStr          string `key:"namestr,string"`
		NotPresent        string `key:",optional"`
		NotPresentWithTag string `key:"notpresent,optional"`
		DefaultString     string `key:"defaultstring,default=hello"`
		Optionalw         string `key:",optional"`
	}
	m := map[string]any{
		"name":    "kevin",
		"namestr": "namewithstring",
	}

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.Equal("kevin", in.Namew)
		astw.Equal("namewithstring", in.NamewStr)
		astw.Empty(in.NotPresent)
		astw.Empty(in.NotPresentWithTag)
		astw.Equal("hello", in.DefaultString)
	}
}

func TestJackyUnmarshalStringWithMissing(t *testing.T) {
	type inner struct {
		Namew string `key:"name"`
	}
	m := map[string]any{}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestJackyUnmarshalStringSliceFromString(t *testing.T) {
	t.Run("slice from string", func(t *testing.T) {
		var v struct {
			Namews []string `key:"names"`
		}
		m := map[string]any{
			"names": `["first", "second"]`,
		}

		astw := assert.New(t)
		if astw.NoError(UnmarshalKey(m, &v)) {
			astw.Equal(2, len(v.Namews))
			astw.Equal("first", v.Namews[0])
			astw.Equal("second", v.Namews[1])
		}
	})

	t.Run("slice from string with slice error", func(t *testing.T) {
		var v struct {
			Namews []int `key:"names"`
		}
		m := map[string]any{
			"names": `["first", 1]`,
		}

		assert.Error(t, UnmarshalKey(m, &v))
	})

	t.Run("slice from string with error", func(t *testing.T) {
		type myString string

		var v struct {
			Namews []string `key:"names"`
		}
		m := map[string]any{
			"names": myString("not a slice"),
		}

		assert.Error(t, UnmarshalKey(m, &v))
	})
}

func TestJackyUnmarshalIntSliceFromString(t *testing.T) {
	var v struct {
		Valuues []int `key:"values"`
	}
	m := map[string]any{
		"values": `[1, 2]`,
	}

	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &v)) {
		astw.Equal(2, len(v.Valuues))
		astw.Equal(1, v.Valuues[0])
		astw.Equal(2, v.Valuues[1])
	}
}

func TestJackyUnmarshalIntMapFromString(t *testing.T) {
	var v struct {
		Sort map[string]int `key:"sort"`
	}
	m := map[string]any{
		"sort": `{"value":12345,"zeroVal":0,"nullVal":null}`,
	}

	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &v)) {
		astw.Equal(3, len(v.Sort))
		astw.Equal(12345, v.Sort["value"])
		astw.Equal(0, v.Sort["zeroVal"])
		astw.Equal(0, v.Sort["nullVal"])
	}
}

func TestJackyUnmarshalBoolMapFromString(t *testing.T) {
	var v struct {
		Sort map[string]bool `key:"sort"`
	}
	m := map[string]any{
		"sort": `{"value":true,"zeroVal":false,"nullVal":null}`,
	}

	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &v)) {
		astw.Equal(3, len(v.Sort))
		astw.Equal(true, v.Sort["value"])
		astw.Equal(false, v.Sort["zeroVal"])
		astw.Equal(false, v.Sort["nullVal"])
	}
}

type CustomStringerJacky string

type UnsupportedStringerJacky string

func (c CustomStringerJacky) String() string {
	return fmt.Sprintf("{%s}", string(c))
}

func TestJackyUnmarshalStringMapFromStringer(t *testing.T) {
	t.Run("CustomStringerJacky", func(t *testing.T) {
		var v struct {
			Sort map[string]string `key:"sort"`
		}
		m := map[string]any{
			"sort": CustomStringerJacky(`"value":"ascend","emptyStr":""`),
		}

		astw := assert.New(t)
		if astw.NoError(UnmarshalKey(m, &v)) {
			astw.Equal(2, len(v.Sort))
			astw.Equal("ascend", v.Sort["value"])
			astw.Equal("", v.Sort["emptyStr"])
		}
	})

	t.Run("CustomStringerJacky incorrect", func(t *testing.T) {
		var v struct {
			Sort map[string]string `key:"sort"`
		}
		m := map[string]any{
			"sort": CustomStringerJacky(`"value"`),
		}

		assert.Error(t, UnmarshalKey(m, &v))
	})
}

func TestJackyUnmarshalStringMapFromUnsupportedType(t *testing.T) {
	var v struct {
		Sort map[string]string `key:"sort"`
	}
	m := map[string]any{
		"sort": UnsupportedStringerJacky(`{"value":"ascend","emptyStr":""}`),
	}

	astw := assert.New(t)
	astw.Error(UnmarshalKey(m, &v))
}

func TestJackyUnmarshalStringMapFromNotSettableValuue(t *testing.T) {
	var v struct {
		sort  map[string]string  `key:"sort"`
		psort *map[string]string `key:"psort"`
	}
	m := map[string]any{
		"sort":  `{"value":"ascend","emptyStr":""}`,
		"psort": `{"value":"ascend","emptyStr":""}`,
	}

	astw := assert.New(t)
	astw.NoError(UnmarshalKey(m, &v))
	assert.Empty(t, v.sort)
	assert.Nil(t, v.psort)
}

func TestJackyUnmarshalStringMapFromString(t *testing.T) {
	var v struct {
		Sort map[string]string `key:"sort"`
	}
	m := map[string]any{
		"sort": `{"value":"ascend","emptyStr":""}`,
	}

	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &v)) {
		astw.Equal(2, len(v.Sort))
		astw.Equal("ascend", v.Sort["value"])
		astw.Equal("", v.Sort["emptyStr"])
	}
}

func TestJackyUnmarshalStructMapFromString(t *testing.T) {
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

	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &v)) {
		astw.Equal(1, len(v.Filter))
		astw.NotNil(v.Filter["obj"])
		astw.Equal(true, v.Filter["obj"].Field1)
		astw.Equal(int64(1573570455447539712), v.Filter["obj"].Field2)
		astw.Equal("this is a string", v.Filter["obj"].Field3)
		astw.Equal("this is a string pointer", *v.Filter["obj"].Field4)
		astw.ElementsMatch([]string{"str1", "str2"}, v.Filter["obj"].Field5)
	}
}

func TestJackyUnmarshalStringSliceMapFromString(t *testing.T) {
	var v struct {
		Filter map[string][]string `key:"filter"`
	}
	m := map[string]any{
		"filter": `{"assignType":null,"status":["process","comment"],"rate":[]}`,
	}

	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &v)) {
		astw.Equal(3, len(v.Filter))
		astw.Equal([]string(nil), v.Filter["assignType"])
		astw.Equal(2, len(v.Filter["status"]))
		astw.Equal("process", v.Filter["status"][0])
		astw.Equal("comment", v.Filter["status"][1])
		astw.Equal(0, len(v.Filter["rate"]))
	}
}

func TestJackyUnmarshalStruct(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		type address struct {
			City          string `key:"city"`
			ZipCode       int    `key:"zipcode,string"`
			DefaultString string `key:"defaultstring,default=hello"`
			Optionalw     string `key:",optional"`
		}
		type inner struct {
			Namew     string    `key:"name"`
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

		var in inner
		astw := assert.New(t)
		if astw.NoError(UnmarshalKey(m, &in)) {
			astw.Equal("kevin", in.Namew)
			astw.Equal("shanghai", in.Address.City)
			astw.Equal(200000, in.Address.ZipCode)
			astw.Equal("hello", in.AddressP.DefaultString)
			astw.Equal("beijing", in.AddressP.City)
			astw.Equal(300000, in.AddressP.ZipCode)
			astw.Equal("hello", in.AddressP.DefaultString)
			astw.Equal("guangzhou", (*in.AddressPP).City)
			astw.Equal(400000, (*in.AddressPP).ZipCode)
			astw.Equal("hello", (*in.AddressPP).DefaultString)
		}
	})

	t.Run("struct with error", func(t *testing.T) {
		type address struct {
			City          string `key:"city"`
			ZipCode       int    `key:"zipcode,string"`
			DefaultString string `key:"defaultstring,default=hello"`
			Optionalw     string `key:",optional"`
		}
		type inner struct {
			Namew     string    `key:"name"`
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

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestJackyUnmarshalStructOptionalwDepends(t *testing.T) {
	type address struct {
		City             string `key:"city"`
		Optionalw        string `key:",optional"`
		OptionalwDepends string `key:",optional=Optionalw"`
	}
	type inner struct {
		Namew   string  `key:"name"`
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
				"OptionalwDepends": "b",
			},
			pass: false,
		},
		{
			input: map[string]string{
				"Optionalw": "a",
			},
			pass: false,
		},
		{
			input: map[string]string{
				"Optionalw":        "a",
				"OptionalwDepends": "b",
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

			var in inner
			astw := assert.New(t)
			if test.pass {
				if astw.NoError(UnmarshalKey(m, &in)) {
					astw.Equal("kevin", in.Namew)
					astw.Equal("shanghai", in.Address.City)
					astw.Equal(test.input["Optionalw"], in.Address.Optionalw)
					astw.Equal(test.input["OptionalwDepends"], in.Address.OptionalwDepends)
				}
			} else {
				astw.Error(UnmarshalKey(m, &in))
			}
		})
	}
}

func TestJackyUnmarshalStructOptionalwDependsNot(t *testing.T) {
	type address struct {
		City             string `key:"city"`
		Optionalw        string `key:",optional"`
		OptionalwDepends string `key:",optional=!Optionalw"`
	}
	type inner struct {
		Namew   string  `key:"name"`
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
				"Optionalw":        "a",
				"OptionalwDepends": "b",
			},
			pass: false,
		},
		{
			input: map[string]string{
				"Optionalw": "a",
			},
			pass: true,
		},
		{
			input: map[string]string{
				"OptionalwDepends": "b",
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

			var in inner
			astw := assert.New(t)
			if test.pass {
				if astw.NoError(UnmarshalKey(m, &in)) {
					astw.Equal("kevin", in.Namew)
					astw.Equal("shanghai", in.Address.City)
					astw.Equal(test.input["Optionalw"], in.Address.Optionalw)
					astw.Equal(test.input["OptionalwDepends"], in.Address.OptionalwDepends)
				}
			} else {
				astw.Error(UnmarshalKey(m, &in))
			}
		})
	}
}

func TestJackyUnmarshalStructOptionalwDependsNotErrorDetails(t *testing.T) {
	t.Run("mutal optionals", func(t *testing.T) {
		type address struct {
			Optionalw        string `key:",optional"`
			OptionalwDepends string `key:",optional=!Optionalw"`
		}
		type inner struct {
			Namew   string  `key:"name"`
			Address address `key:"address"`
		}

		m := map[string]any{
			"name": "kevin",
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("with default", func(t *testing.T) {
		type address struct {
			Optionalw        string `key:",optional"`
			OptionalwDepends string `key:",default=value,optional"`
		}
		type inner struct {
			Namew   string  `key:"name"`
			Address address `key:"address"`
		}

		m := map[string]any{
			"name": "kevin",
		}

		var in inner
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.Equal(t, "kevin", in.Namew)
			assert.Equal(t, "value", in.Address.OptionalwDepends)
		}
	})
}

func TestJackyUnmarshalStructOptionalwDependsNotNested(t *testing.T) {
	t.Run("mutal optionals", func(t *testing.T) {
		type address struct {
			Optionalw        string `key:",optional"`
			OptionalwDepends string `key:",optional=!Optionalw"`
		}
		type combo struct {
			Namew   string  `key:"name,optional"`
			Address address `key:"address"`
		}
		type inner struct {
			Namew string `key:"name"`
			Combo combo  `key:"combo"`
		}

		m := map[string]any{
			"name": "kevin",
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("bad format", func(t *testing.T) {
		type address struct {
			Optionalw        string `key:",optional"`
			OptionalwDepends string `key:",optional=!Optionalw=abcd"`
		}
		type combo struct {
			Namew   string  `key:"name,optional"`
			Address address `key:"address"`
		}
		type inner struct {
			Namew string `key:"name"`
			Combo combo  `key:"combo"`
		}

		m := map[string]any{
			"name": "kevin",
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("invalid option", func(t *testing.T) {
		type address struct {
			Optionalw        string `key:",optional"`
			OptionalwDepends string `key:",opt=abcd"`
		}
		type combo struct {
			Namew   string  `key:"name,optional"`
			Address address `key:"address"`
		}
		type inner struct {
			Namew string `key:"name"`
			Combo combo  `key:"combo"`
		}

		m := map[string]any{
			"name": "kevin",
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestJackyUnmarshalStructOptionalwNestedDifferentKey(t *testing.T) {
	type address struct {
		Optionalw        string `dkey:",optional"`
		OptionalwDepends string `key:",optional"`
	}
	type combo struct {
		Namew   string  `key:"name,optional"`
		Address address `key:"address"`
	}
	type inner struct {
		Namew string `key:"name"`
		Combo combo  `key:"combo"`
	}

	m := map[string]any{
		"name": "kevin",
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestJackyUnmarshalStructOptionalwDependsNotEnoughValuue(t *testing.T) {
	type address struct {
		Optionalw        string `key:",optional"`
		OptionalwDepends string `key:",optional=!"`
	}
	type inner struct {
		Namew   string  `key:"name"`
		Address address `key:"address"`
	}

	m := map[string]any{
		"name":    "kevin",
		"address": map[string]any{},
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestJackyUnmarshalStructOptionalwDependsMoreValuues(t *testing.T) {
	type address struct {
		Optionalw        string `key:",optional"`
		OptionalwDepends string `key:",optional=a=b"`
	}
	type inner struct {
		Namew   string  `key:"name"`
		Address address `key:"address"`
	}

	m := map[string]any{
		"name":    "kevin",
		"address": map[string]any{},
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestJackyUnmarshalStructMissing(t *testing.T) {
	type address struct {
		Optionalw        string `key:",optional"`
		OptionalwDepends string `key:",optional=a=b"`
	}
	type inner struct {
		Namew   string  `key:"name"`
		Address address `key:"address"`
	}

	m := map[string]any{
		"name": "kevin",
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestJackyUnmarshalNestedStructMissing(t *testing.T) {
	type mostInner struct {
		Namew string `key:"name"`
	}
	type address struct {
		Optionalw        string `key:",optional"`
		OptionalwDepends string `key:",optional=a=b"`
		MostInner        mostInner
	}
	type inner struct {
		Namew   string  `key:"name"`
		Address address `key:"address"`
	}

	m := map[string]any{
		"name":    "kevin",
		"address": map[string]any{},
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestJackyUnmarshalAnonymousStructOptionalwDepends(t *testing.T) {
	type AnonAddress struct {
		City             string `key:"city"`
		Optionalw        string `key:",optional"`
		OptionalwDepends string `key:",optional=Optionalw"`
	}
	type inner struct {
		Namew string `key:"name"`
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
				"OptionalwDepends": "b",
			},
			pass: false,
		},
		{
			input: map[string]string{
				"Optionalw": "a",
			},
			pass: false,
		},
		{
			input: map[string]string{
				"Optionalw":        "a",
				"OptionalwDepends": "b",
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

			var in inner
			astw := assert.New(t)
			if test.pass {
				if astw.NoError(UnmarshalKey(m, &in)) {
					astw.Equal("kevin", in.Namew)
					astw.Equal("shanghai", in.City)
					astw.Equal(test.input["Optionalw"], in.Optionalw)
					astw.Equal(test.input["OptionalwDepends"], in.OptionalwDepends)
				}
			} else {
				astw.Error(UnmarshalKey(m, &in))
			}
		})
	}
}

func TestJackyUnmarshalStructPtr(t *testing.T) {
	type address struct {
		City          string `key:"city"`
		ZipCode       int    `key:"zipcode,string"`
		DefaultString string `key:"defaultstring,default=hello"`
		Optionalw     string `key:",optional"`
	}
	type inner struct {
		Namew   string   `key:"name"`
		Address *address `key:"address"`
	}
	m := map[string]any{
		"name": "kevin",
		"address": map[string]any{
			"city":    "shanghai",
			"zipcode": "200000",
		},
	}

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.Equal("kevin", in.Namew)
		astw.Equal("shanghai", in.Address.City)
		astw.Equal(200000, in.Address.ZipCode)
		astw.Equal("hello", in.Address.DefaultString)
	}
}

func TestJackyUnmarshalWithStringIgnored(t *testing.T) {
	type inner struct {
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

	var in inner
	um := NewUnmarshaler("key", WithStringValues())
	astw := assert.New(t)
	if astw.NoError(um.Unmarshal(m, &in)) {
		astw.True(in.True)
		astw.False(in.False)
		astw.Equal(1, in.Int)
		astw.Equal(int8(3), in.Int8)
		astw.Equal(int16(5), in.Int16)
		astw.Equal(int32(7), in.Int32)
		astw.Equal(int64(9), in.Int64)
		astw.Equal(uint(1), in.Uint)
		astw.Equal(uint8(3), in.Uint8)
		astw.Equal(uint16(5), in.Uint16)
		astw.Equal(uint32(7), in.Uint32)
		astw.Equal(uint64(9), in.Uint64)
		astw.Equal(float32(1.5), in.Float32)
		astw.Equal(3.5, in.Float64)
	}
}

func TestJackyUnmarshalJsonNumberInt64(t *testing.T) {
	for i := 0; i <= maxUintBitsToTestJacky; i++ {
		var intValuue int64 = 1 << uint(i)
		strValuue := strconv.FormatInt(intValuue, 10)
		number := json.Number(strValuue)
		m := map[string]any{
			"ID": number,
		}
		var v struct {
			ID int64
		}
		if assert.NoError(t, UnmarshalKey(m, &v)) {
			assert.Equal(t, intValuue, v.ID)
		}
	}
}

func TestJackyUnmarshalJsonNumberUint64(t *testing.T) {
	for i := 0; i <= maxUintBitsToTestJacky; i++ {
		var intValuue uint64 = 1 << uint(i)
		strValuue := strconv.FormatUint(intValuue, 10)
		number := json.Number(strValuue)
		m := map[string]any{
			"ID": number,
		}
		var v struct {
			ID uint64
		}
		if assert.NoError(t, UnmarshalKey(m, &v)) {
			assert.Equal(t, intValuue, v.ID)
		}
	}
}

func TestJackyUnmarshalJsonNumberUint64Ptr(t *testing.T) {
	for i := 0; i <= maxUintBitsToTestJacky; i++ {
		var intValuue uint64 = 1 << uint(i)
		strValuue := strconv.FormatUint(intValuue, 10)
		number := json.Number(strValuue)
		m := map[string]any{
			"ID": number,
		}
		var v struct {
			ID *uint64
		}
		astw := assert.New(t)
		if astw.NoError(UnmarshalKey(m, &v)) {
			astw.NotNil(v.ID)
			astw.Equal(intValuue, *v.ID)
		}
	}
}

func TestJackyUnmarshalMapOfInt(t *testing.T) {
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

func TestJackyUnmarshalMapOfStruct(t *testing.T) {
	t.Run("map of struct", func(t *testing.T) {
		m := map[string]any{
			"Ids": map[string]any{
				"foo": map[string]any{"Namew": "foo"},
			},
		}
		var v struct {
			Ids map[string]struct {
				Namew string
			}
		}
		if assert.NoError(t, UnmarshalKey(m, &v)) {
			assert.Equal(t, "foo", v.Ids["foo"].Namew)
		}
	})
}

func TestJackyUnmarshalSlice(t *testing.T) {
	t.Run("slice of string", func(t *testing.T) {
		m := map[string]any{
			"Ids": []any{"first", "second"},
		}
		var v struct {
			Ids []string
		}
		astw := assert.New(t)
		if astw.NoError(UnmarshalKey(m, &v)) {
			astw.Equal(2, len(v.Ids))
			astw.Equal("first", v.Ids[0])
			astw.Equal("second", v.Ids[1])
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
		astw := assert.New(t)
		if astw.NoError(NewUnmarshaler(jsonTagKey).Unmarshal([]any{1, 2}, &v)) {
			astw.Equal(2, len(v))
			astw.Equal(1, v[0])
			astw.Equal(2, v[1])
		}
	})

	t.Run("slice with unsupported type", func(t *testing.T) {
		var v int
		assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(1, &v))
	})
}

func TestJackyUnmarshalSliceOfStruct(t *testing.T) {
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
		astw := assert.New(t)
		if astw.NoError(UnmarshalKey(m, &v)) {
			astw.Equal(1, len(v.Ids))
			astw.Equal(1, v.Ids[0].First)
			astw.Equal(2, v.Ids[0].Second)
		}
	})
}

func TestJackyUnmarshalWithStringOptionsCorrect(t *testing.T) {
	type inner struct {
		Valuue  string `key:"value,options=first|second"`
		Foo     string `key:"foo,options=[bar,baz]"`
		Correct string `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "first",
		"foo":     "bar",
		"correct": "2",
	}

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.Equal("first", in.Valuue)
		astw.Equal("bar", in.Foo)
		astw.Equal("2", in.Correct)
	}
}

func TestJackyUnmarshalOptionsOptionalw(t *testing.T) {
	type inner struct {
		Valuue        string `key:"value,options=first|second,optional"`
		OptionalwVaue string `key:"optional_value,options=first|second,optional"`
		Foo           string `key:"foo,options=[bar,baz]"`
		Correct       string `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "first",
		"foo":     "bar",
		"correct": "2",
	}

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.Equal("first", in.Valuue)
		astw.Equal("", in.OptionalwVaue)
		astw.Equal("bar", in.Foo)
		astw.Equal("2", in.Correct)
	}
}

func TestJackyUnmarshalOptionsMissingValuues(t *testing.T) {
	type inner struct {
		Valuue string `key:"value,options"`
	}
	m := map[string]any{
		"value": "first",
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestJackyUnmarshalStringOptionsWithStringOptionsNotString(t *testing.T) {
	type inner struct {
		Valuue  string `key:"value,options=first|second"`
		Correct string `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "first",
		"correct": 2,
	}

	var in inner
	unmarshaler := NewUnmarshaler(defaultKeyName, WithStringValues())
	assert.Error(t, unmarshaler.Unmarshal(m, &in))
}

func TestJackyUnmarshalStringOptionsWithStringOptions(t *testing.T) {
	type inner struct {
		Valuue  string `key:"value,options=first|second"`
		Correct string `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "first",
		"correct": "2",
	}

	var in inner
	unmarshaler := NewUnmarshaler(defaultKeyName, WithStringValues())
	astw := assert.New(t)
	if astw.NoError(unmarshaler.Unmarshal(m, &in)) {
		astw.Equal("first", in.Valuue)
		astw.Equal("2", in.Correct)
	}
}

func TestJackyUnmarshalStringOptionsWithStringOptionsPtr(t *testing.T) {
	type inner struct {
		Valuue  *string  `key:"value,options=first|second"`
		ValuueP **string `key:"valuep,options=first|second"`
		Correct *int     `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "first",
		"valuep":  "second",
		"correct": "2",
	}

	var in inner
	unmarshaler := NewUnmarshaler(defaultKeyName, WithStringValues())
	astw := assert.New(t)
	if astw.NoError(unmarshaler.Unmarshal(m, &in)) {
		astw.True(*in.Valuue == "first")
		astw.True(**in.ValuueP == "second")
		astw.True(*in.Correct == 2)
	}
}

func TestJackyUnmarshalStringOptionsWithStringOptionsIncorrect(t *testing.T) {
	type inner struct {
		Valuue  string `key:"value,options=first|second"`
		Correct string `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "third",
		"correct": "2",
	}

	var in inner
	unmarshaler := NewUnmarshaler(defaultKeyName, WithStringValues())
	assert.Error(t, unmarshaler.Unmarshal(m, &in))
}

func TestJackyUnmarshalStringOptionsWithStringOptionsIncorrectGrouped(t *testing.T) {
	type inner struct {
		Valuue  string `key:"value,options=[first,second]"`
		Correct string `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "third",
		"correct": "2",
	}

	var in inner
	unmarshaler := NewUnmarshaler(defaultKeyName, WithStringValues())
	assert.Error(t, unmarshaler.Unmarshal(m, &in))
}

func TestJackyUnmarshalWithStringOptionsIncorrect(t *testing.T) {
	type inner struct {
		Valuue    string `key:"value,options=first|second"`
		Incorrect string `key:"incorrect,options=1|2"`
	}
	m := map[string]any{
		"value":     "first",
		"incorrect": "3",
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestJackyUnmarshalWithIntOptionsCorrect(t *testing.T) {
	type inner struct {
		Valuue string `key:"value,options=first|second"`
		Number int    `key:"number,options=1|2"`
	}
	m := map[string]any{
		"value":  "first",
		"number": 2,
	}

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.Equal("first", in.Valuue)
		astw.Equal(2, in.Number)
	}
}

func TestJackyUnmarshalWithIntOptionsCorrectPtr(t *testing.T) {
	type inner struct {
		Valuue *string `key:"value,options=first|second"`
		Number *int    `key:"number,options=1|2"`
	}
	m := map[string]any{
		"value":  "first",
		"number": 2,
	}

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.True(*in.Valuue == "first")
		astw.True(*in.Number == 2)
	}
}

func TestJackyUnmarshalWithIntOptionsIncorrect(t *testing.T) {
	type inner struct {
		Valuue    string `key:"value,options=first|second"`
		Incorrect int    `key:"incorrect,options=1|2"`
	}
	m := map[string]any{
		"value":     "first",
		"incorrect": 3,
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestJackyUnmarshalWithJsonNumberOptionsIncorrect(t *testing.T) {
	type inner struct {
		Valuue    string `key:"value,options=first|second"`
		Incorrect int    `key:"incorrect,options=1|2"`
	}
	m := map[string]any{
		"value":     "first",
		"incorrect": json.Number("3"),
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestJackyUnmarshalerUnmarshalIntOptions(t *testing.T) {
	var val struct {
		Sex int `json:"sex,options=0|1"`
	}
	input := []byte(`{"sex": 2}`)
	assert.Error(t, UnmarshalJsonBytes(input, &val))
}

func TestJackyUnmarshalWithUintOptionsCorrect(t *testing.T) {
	type inner struct {
		Valuue string `key:"value,options=first|second"`
		Number uint   `key:"number,options=1|2"`
	}
	m := map[string]any{
		"value":  "first",
		"number": uint(2),
	}

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.Equal("first", in.Valuue)
		astw.Equal(uint(2), in.Number)
	}
}

func TestJackyUnmarshalWithUintOptionsIncorrect(t *testing.T) {
	type inner struct {
		Valuue    string `key:"value,options=first|second"`
		Incorrect uint   `key:"incorrect,options=1|2"`
	}
	m := map[string]any{
		"value":     "first",
		"incorrect": uint(3),
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestJackyUnmarshalWithOptionsAndDefault(t *testing.T) {
	type inner struct {
		Valuue string `key:"value,options=first|second|third,default=second"`
	}
	m := map[string]any{}

	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, "second", in.Valuue)
	}
}

func TestJackyUnmarshalWithOptionsAndSet(t *testing.T) {
	type inner struct {
		Valuue string `key:"value,options=first|second|third,default=second"`
	}
	m := map[string]any{
		"value": "first",
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, "first", in.Valuue)
	}
}

func TestJackyUnmarshalNestedKey(t *testing.T) {
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

func TestJackyUnmarhsalNestedKeyArray(t *testing.T) {
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

func TestJackyUnmarshalAnonymousOptionalwRequiredProvided(t *testing.T) {
	type (
		Foo struct {
			Valuue string `json:"v"`
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
		assert.Equal(t, "anything", b.Valuue)
	}
}

func TestJackyUnmarshalAnonymousOptionalwRequiredMissed(t *testing.T) {
	type (
		Foo struct {
			Valuue string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Valuue) == 0)
	}
}

func TestJackyUnmarshalAnonymousOptionalwOptionalwProvided(t *testing.T) {
	type (
		Foo struct {
			Valuue string `json:"v,optional"`
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
		assert.Equal(t, "anything", b.Valuue)
	}
}

func TestJackyUnmarshalAnonymousOptionalwOptionalwMissed(t *testing.T) {
	type (
		Foo struct {
			Valuue string `json:"v,optional"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Valuue) == 0)
	}
}

func TestJackyUnmarshalAnonymousOptionalwRequiredBothProvided(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n"`
			Valuue string `json:"v"`
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
		assert.Equal(t, "kevin", b.Namew)
		assert.Equal(t, "anything", b.Valuue)
	}
}

func TestJackyUnmarshalAnonymousOptionalwRequiredOneProvidedOneMissed(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n"`
			Valuue string `json:"v"`
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

func TestJackyUnmarshalAnonymousOptionalwRequiredBothMissed(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n"`
			Valuue string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Namew) == 0)
		assert.True(t, len(b.Valuue) == 0)
	}
}

func TestJackyUnmarshalAnonymousOptionalwOneRequiredOneOptionalwBothProvided(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n,optional"`
			Valuue string `json:"v"`
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
		assert.Equal(t, "kevin", b.Namew)
		assert.Equal(t, "anything", b.Valuue)
	}
}

func TestJackyUnmarshalAnonymousOptionalwOneRequiredOneOptionalwBothMissed(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n,optional"`
			Valuue string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Namew) == 0)
		assert.True(t, len(b.Valuue) == 0)
	}
}

func TestJackyUnmarshalAnonymousOptionalwOneRequiredOneOptionalwRequiredProvidedOptionalwMissed(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n,optional"`
			Valuue string `json:"v"`
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
		assert.True(t, len(b.Namew) == 0)
		assert.Equal(t, "anything", b.Valuue)
	}
}

func TestJackyUnmarshalAnonymousOptionalwOneRequiredOneOptionalwRequiredMissedOptionalwProvided(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n,optional"`
			Valuue string `json:"v"`
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

func TestJackyUnmarshalAnonymousOptionalwBothOptionalwBothProvided(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n,optional"`
			Valuue string `json:"v,optional"`
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
		assert.Equal(t, "kevin", b.Namew)
		assert.Equal(t, "anything", b.Valuue)
	}
}

func TestJackyUnmarshalAnonymousOptionalwBothOptionalwOneProvidedOneMissed(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n,optional"`
			Valuue string `json:"v,optional"`
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
		assert.True(t, len(b.Namew) == 0)
		assert.Equal(t, "anything", b.Valuue)
	}
}

func TestJackyUnmarshalAnonymousOptionalwBothOptionalwBothMissed(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n,optional"`
			Valuue string `json:"v,optional"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Namew) == 0)
		assert.True(t, len(b.Valuue) == 0)
	}
}

func TestJackyUnmarshalAnonymousRequiredProvided(t *testing.T) {
	type (
		Foo struct {
			Valuue string `json:"v"`
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
		assert.Equal(t, "anything", b.Valuue)
	}
}

func TestJackyUnmarshalAnonymousRequiredMissed(t *testing.T) {
	type (
		Foo struct {
			Valuue string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{}

	var b Bar
	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b))
}

func TestJackyUnmarshalAnonymousOptionalwProvided(t *testing.T) {
	type (
		Foo struct {
			Valuue string `json:"v,optional"`
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
		assert.Equal(t, "anything", b.Valuue)
	}
}

func TestJackyUnmarshalAnonymousOptionalwMissed(t *testing.T) {
	type (
		Foo struct {
			Valuue string `json:"v,optional"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Valuue) == 0)
	}
}

func TestJackyUnmarshalAnonymousRequiredBothProvided(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n"`
			Valuue string `json:"v"`
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
		assert.Equal(t, "kevin", b.Namew)
		assert.Equal(t, "anything", b.Valuue)
	}
}

func TestJackyUnmarshalAnonymousRequiredOneProvidedOneMissed(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n"`
			Valuue string `json:"v"`
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

func TestJackyUnmarshalAnonymousRequiredBothMissed(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n"`
			Valuue string `json:"v"`
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

func TestJackyUnmarshalAnonymousOneRequiredOneOptionalwBothProvided(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n,optional"`
			Valuue string `json:"v"`
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
		assert.Equal(t, "kevin", b.Namew)
		assert.Equal(t, "anything", b.Valuue)
	}
}

func TestJackyUnmarshalAnonymousOneRequiredOneOptionalwBothMissed(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n,optional"`
			Valuue string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{}

	var b Bar
	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b))
}

func TestJackyUnmarshalAnonymousOneRequiredOneOptionalwRequiredProvidedOptionalwMissed(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n,optional"`
			Valuue string `json:"v"`
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
		assert.True(t, len(b.Namew) == 0)
		assert.Equal(t, "anything", b.Valuue)
	}
}

func TestJackyUnmarshalAnonymousOneRequiredOneOptionalwRequiredMissedOptionalwProvided(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n,optional"`
			Valuue string `json:"v"`
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

func TestJackyUnmarshalAnonymousBothOptionalwBothProvided(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n,optional"`
			Valuue string `json:"v,optional"`
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
		assert.Equal(t, "kevin", b.Namew)
		assert.Equal(t, "anything", b.Valuue)
	}
}

func TestJackyUnmarshalAnonymousBothOptionalwOneProvidedOneMissed(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n,optional"`
			Valuue string `json:"v,optional"`
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
		assert.True(t, len(b.Namew) == 0)
		assert.Equal(t, "anything", b.Valuue)
	}
}

func TestJackyUnmarshalAnonymousBothOptionalwBothMissed(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n,optional"`
			Valuue string `json:"v,optional"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Namew) == 0)
		assert.True(t, len(b.Valuue) == 0)
	}
}

func TestJackyUnmarshalAnonymousWrappedToMuch(t *testing.T) {
	type (
		Foo struct {
			Namew  string `json:"n"`
			Valuue string `json:"v"`
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

func TestJackyUnmarshalWrappedObject(t *testing.T) {
	type (
		Foo struct {
			Valuue string `json:"v"`
		}

		Bar struct {
			Inner Foo
		}
	)
	m := map[string]any{
		"Inner": map[string]any{
			"v": "anything",
		},
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "anything", b.Inner.Valuue)
	}
}

func TestJackyUnmarshalWrappedObjectOptionalw(t *testing.T) {
	type (
		Foo struct {
			Hosts []string
			Key   string
		}

		Bar struct {
			Inner Foo `json:",optional"`
			Namew string
		}
	)
	m := map[string]any{
		"Namew": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "anything", b.Namew)
	}
}

func TestJackyUnmarshalWrappedObjectOptionalwFilled(t *testing.T) {
	type (
		Foo struct {
			Hosts []string
			Key   string
		}

		Bar struct {
			Inner Foo `json:",optional"`
			Namew string
		}
	)
	hosts := []string{"1", "2"}
	m := map[string]any{
		"Inner": map[string]any{
			"Hosts": hosts,
			"Key":   "key",
		},
		"Namew": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.EqualValues(t, hosts, b.Inner.Hosts)
		assert.Equal(t, "key", b.Inner.Key)
		assert.Equal(t, "anything", b.Namew)
	}
}

func TestJackyUnmarshalWrappedNamewdObjectOptionalw(t *testing.T) {
	type (
		Foo struct {
			Host string
			Key  string
		}

		Bar struct {
			Inner Foo `json:",optional"`
			Namew string
		}
	)
	m := map[string]any{
		"Inner": map[string]any{
			"Host": "thehost",
			"Key":  "thekey",
		},
		"Namew": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "thehost", b.Inner.Host)
		assert.Equal(t, "thekey", b.Inner.Key)
		assert.Equal(t, "anything", b.Namew)
	}
}

func TestJackyUnmarshalWrappedObjectNamewdPtr(t *testing.T) {
	type (
		Foo struct {
			Valuue string `json:"v"`
		}

		Bar struct {
			Inner *Foo `json:"foo,optional"`
		}
	)
	m := map[string]any{
		"foo": map[string]any{
			"v": "anything",
		},
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "anything", b.Inner.Valuue)
	}
}

func TestJackyUnmarshalWrappedObjectPtr(t *testing.T) {
	type (
		Foo struct {
			Valuue string `json:"v"`
		}

		Bar struct {
			Inner *Foo
		}
	)
	m := map[string]any{
		"Inner": map[string]any{
			"v": "anything",
		},
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "anything", b.Inner.Valuue)
	}
}

func TestJackyUnmarshalInt2String(t *testing.T) {
	type inner struct {
		Int string `key:"int"`
	}
	m := map[string]any{
		"int": 123,
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestJackyUnmarshalZeroValuues(t *testing.T) {
	type inner struct {
		False  bool   `key:"no"`
		Int    int    `key:"int"`
		String string `key:"string"`
	}
	m := map[string]any{
		"no":     false,
		"int":    0,
		"string": "",
	}

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.False(in.False)
		astw.Equal(0, in.Int)
		astw.Equal("", in.String)
	}
}

func TestJackyUnmarshalUsingDifferentKeys(t *testing.T) {
	type inner struct {
		False  bool   `key:"no"`
		Int    int    `key:"int"`
		String string `bson:"string"`
	}
	m := map[string]any{
		"no":     false,
		"int":    9,
		"string": "value",
	}

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.False(in.False)
		astw.Equal(9, in.Int)
		astw.True(len(in.String) == 0)
	}
}

func TestJackyUnmarshalNumberRangeInt(t *testing.T) {
	type inner struct {
		Valuue1  int    `key:"value1,range=[1:]"`
		Valuue2  int8   `key:"value2,range=[1:5]"`
		Valuue3  int16  `key:"value3,range=[1:5]"`
		Valuue4  int32  `key:"value4,range=[1:5]"`
		Valuue5  int64  `key:"value5,range=[1:5]"`
		Valuue6  uint   `key:"value6,range=[:5]"`
		Valuue8  uint8  `key:"value8,range=[1:5],string"`
		Valuue9  uint16 `key:"value9,range=[1:5],string"`
		Valuue10 uint32 `key:"value10,range=[1:5],string"`
		Valuue11 uint64 `key:"value11,range=[1:5],string"`
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

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.Equal(10, in.Valuue1)
		astw.Equal(int8(1), in.Valuue2)
		astw.Equal(int16(2), in.Valuue3)
		astw.Equal(int32(4), in.Valuue4)
		astw.Equal(int64(5), in.Valuue5)
		astw.Equal(uint(0), in.Valuue6)
		astw.Equal(uint8(1), in.Valuue8)
		astw.Equal(uint16(2), in.Valuue9)
		astw.Equal(uint32(4), in.Valuue10)
		astw.Equal(uint64(5), in.Valuue11)
	}
}

func TestJackyUnmarshalNumberRangeJsonNumber(t *testing.T) {
	type inner struct {
		Valuue3 uint   `key:"value3,range=(1:5]"`
		Valuue4 uint8  `key:"value4,range=(1:5]"`
		Valuue5 uint16 `key:"value5,range=(1:5]"`
	}
	m := map[string]any{
		"value3": json.Number("2"),
		"value4": json.Number("4"),
		"value5": json.Number("5"),
	}

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.Equal(uint(2), in.Valuue3)
		astw.Equal(uint8(4), in.Valuue4)
		astw.Equal(uint16(5), in.Valuue5)
	}

	type inner1 struct {
		Valuue int `key:"value,range=(1:5]"`
	}
	m = map[string]any{
		"value": json.Number("a"),
	}

	var in1 inner1
	astw.Error(UnmarshalKey(m, &in1))
}

func TestJackyUnmarshalNumberRangeIntLeftExclude(t *testing.T) {
	type inner struct {
		Valuue3  uint   `key:"value3,range=(1:5]"`
		Valuue4  uint32 `key:"value4,default=4,range=(1:5]"`
		Valuue5  uint64 `key:"value5,range=(1:5]"`
		Valuue9  int    `key:"value9,range=(1:5],string"`
		Valuue10 int    `key:"value10,range=(1:5],string"`
		Valuue11 int    `key:"value11,range=(1:5],string"`
	}
	m := map[string]any{
		"value3":  uint(2),
		"value4":  uint32(4),
		"value5":  uint64(5),
		"value9":  "2",
		"value10": "4",
		"value11": "5",
	}

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.Equal(uint(2), in.Valuue3)
		astw.Equal(uint32(4), in.Valuue4)
		astw.Equal(uint64(5), in.Valuue5)
		astw.Equal(2, in.Valuue9)
		astw.Equal(4, in.Valuue10)
		astw.Equal(5, in.Valuue11)
	}
}

func TestJackyUnmarshalNumberRangeIntRightExclude(t *testing.T) {
	type inner struct {
		Valuue2  uint   `key:"value2,range=[1:5)"`
		Valuue3  uint8  `key:"value3,range=[1:5)"`
		Valuue4  uint16 `key:"value4,range=[1:5)"`
		Valuue8  int    `key:"value8,range=[1:5),string"`
		Valuue9  int    `key:"value9,range=[1:5),string"`
		Valuue10 int    `key:"value10,range=[1:5),string"`
	}
	m := map[string]any{
		"value2":  uint(1),
		"value3":  uint8(2),
		"value4":  uint16(4),
		"value8":  "1",
		"value9":  "2",
		"value10": "4",
	}

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.Equal(uint(1), in.Valuue2)
		astw.Equal(uint8(2), in.Valuue3)
		astw.Equal(uint16(4), in.Valuue4)
		astw.Equal(1, in.Valuue8)
		astw.Equal(2, in.Valuue9)
		astw.Equal(4, in.Valuue10)
	}
}

func TestJackyUnmarshalNumberRangeIntExclude(t *testing.T) {
	type inner struct {
		Valuue3  int `key:"value3,range=(1:5)"`
		Valuue4  int `key:"value4,range=(1:5)"`
		Valuue9  int `key:"value9,range=(1:5),string"`
		Valuue10 int `key:"value10,range=(1:5),string"`
	}
	m := map[string]any{
		"value3":  2,
		"value4":  4,
		"value9":  "2",
		"value10": "4",
	}

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.Equal(2, in.Valuue3)
		astw.Equal(4, in.Valuue4)
		astw.Equal(2, in.Valuue9)
		astw.Equal(4, in.Valuue10)
	}
}

func TestJackyUnmarshalNumberRangeIntOutOfRange(t *testing.T) {
	type inner1 struct {
		Valuue int64 `key:"value,default=3,range=(1:5)"`
	}

	var in1 inner1
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

	type inner2 struct {
		Valuue int64 `key:"value,optional,range=[1:5)"`
	}

	var in2 inner2
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(0),
	}, &in2))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(5),
	}, &in2))

	type inner3 struct {
		Valuue int64 `key:"value,range=(1:5]"`
	}

	var in3 inner3
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(1),
	}, &in3))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(6),
	}, &in3))

	type inner4 struct {
		Valuue int64 `key:"value,range=[1:5]"`
	}

	var in4 inner4
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(0),
	}, &in4))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(6),
	}, &in4))
}

func TestJackyUnmarshalNumberRangeFloat(t *testing.T) {
	type inner struct {
		Valuue2  float32 `key:"value2,range=[1:5]"`
		Valuue3  float32 `key:"value3,range=[1:5]"`
		Valuue4  float64 `key:"value4,range=[1:5]"`
		Valuue5  float64 `key:"value5,range=[1:5]"`
		Valuue8  float64 `key:"value8,range=[1:5],string"`
		Valuue9  float64 `key:"value9,range=[1:5],string"`
		Valuue10 float64 `key:"value10,range=[1:5],string"`
		Valuue11 float64 `key:"value11,range=[1:5],string"`
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

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.Equal(float32(1), in.Valuue2)
		astw.Equal(float32(2), in.Valuue3)
		astw.Equal(float64(4), in.Valuue4)
		astw.Equal(float64(5), in.Valuue5)
		astw.Equal(float64(1), in.Valuue8)
		astw.Equal(float64(2), in.Valuue9)
		astw.Equal(float64(4), in.Valuue10)
		astw.Equal(float64(5), in.Valuue11)
	}
}

func TestJackyUnmarshalNumberRangeFloatLeftExclude(t *testing.T) {
	type inner struct {
		Valuue3  float64 `key:"value3,range=(1:5]"`
		Valuue4  float64 `key:"value4,range=(1:5]"`
		Valuue5  float64 `key:"value5,range=(1:5]"`
		Valuue9  float64 `key:"value9,range=(1:5],string"`
		Valuue10 float64 `key:"value10,range=(1:5],string"`
		Valuue11 float64 `key:"value11,range=(1:5],string"`
	}
	m := map[string]any{
		"value3":  float64(2),
		"value4":  float64(4),
		"value5":  float64(5),
		"value9":  "2",
		"value10": "4",
		"value11": "5",
	}

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.Equal(float64(2), in.Valuue3)
		astw.Equal(float64(4), in.Valuue4)
		astw.Equal(float64(5), in.Valuue5)
		astw.Equal(float64(2), in.Valuue9)
		astw.Equal(float64(4), in.Valuue10)
		astw.Equal(float64(5), in.Valuue11)
	}
}

func TestJackyUnmarshalNumberRangeFloatRightExclude(t *testing.T) {
	type inner struct {
		Valuue2  float64 `key:"value2,range=[1:5)"`
		Valuue3  float64 `key:"value3,range=[1:5)"`
		Valuue4  float64 `key:"value4,range=[1:5)"`
		Valuue8  float64 `key:"value8,range=[1:5),string"`
		Valuue9  float64 `key:"value9,range=[1:5),string"`
		Valuue10 float64 `key:"value10,range=[1:5),string"`
	}
	m := map[string]any{
		"value2":  float64(1),
		"value3":  float64(2),
		"value4":  float64(4),
		"value8":  "1",
		"value9":  "2",
		"value10": "4",
	}

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.Equal(float64(1), in.Valuue2)
		astw.Equal(float64(2), in.Valuue3)
		astw.Equal(float64(4), in.Valuue4)
		astw.Equal(float64(1), in.Valuue8)
		astw.Equal(float64(2), in.Valuue9)
		astw.Equal(float64(4), in.Valuue10)
	}
}

func TestJackyUnmarshalNumberRangeFloatExclude(t *testing.T) {
	type inner struct {
		Valuue3  float64 `key:"value3,range=(1:5)"`
		Valuue4  float64 `key:"value4,range=(1:5)"`
		Valuue9  float64 `key:"value9,range=(1:5),string"`
		Valuue10 float64 `key:"value10,range=(1:5),string"`
	}
	m := map[string]any{
		"value3":  float64(2),
		"value4":  float64(4),
		"value9":  "2",
		"value10": "4",
	}

	var in inner
	astw := assert.New(t)
	if astw.NoError(UnmarshalKey(m, &in)) {
		astw.Equal(float64(2), in.Valuue3)
		astw.Equal(float64(4), in.Valuue4)
		astw.Equal(float64(2), in.Valuue9)
		astw.Equal(float64(4), in.Valuue10)
	}
}

func TestJackyUnmarshalNumberRangeFloatOutOfRange(t *testing.T) {
	type inner1 struct {
		Valuue float64 `key:"value,range=(1:5)"`
	}

	var in1 inner1
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

	type inner2 struct {
		Valuue float64 `key:"value,range=[1:5)"`
	}

	var in2 inner2
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(0),
	}, &in2))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(5),
	}, &in2))

	type inner3 struct {
		Valuue float64 `key:"value,range=(1:5]"`
	}

	var in3 inner3
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(1),
	}, &in3))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(6),
	}, &in3))

	type inner4 struct {
		Valuue float64 `key:"value,range=[1:5]"`
	}

	var in4 inner4
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(0),
	}, &in4))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(6),
	}, &in4))
}

func TestJackyUnmarshalNestedMap(t *testing.T) {
	t.Run("nested map", func(t *testing.T) {
		var c struct {
			Anything map[string]map[string]string `json:"anything"`
		}
		m := map[string]any{
			"anything": map[string]map[string]any{
				"inner": {
					"id":   "1",
					"name": "any",
				},
			},
		}

		if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &c)) {
			assert.Equal(t, "1", c.Anything["inner"]["id"])
		}
	})

	t.Run("nested map with slice element", func(t *testing.T) {
		var c struct {
			Anything map[string][]string `json:"anything"`
		}
		m := map[string]any{
			"anything": map[string][]any{
				"inner": {
					"id",
					"name",
				},
			},
		}

		if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &c)) {
			assert.Equal(t, []string{"id", "name"}, c.Anything["inner"])
		}
	})

	// t.Run("nested map with slice element error", func(t *testing.T) {
	// 	var c struct {
	// 		Anything map[string][]string `json:"anything"`
	// 	}
	// 	m := map[string]any{
	// 		"anything": map[string][]any{
	// 			"inner": {
	// 				"id",
	// 				1,
	// 			},
	// 		},
	// 	}

	// 	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &c))
	// })
}

// func TestJackyUnmarshalNestedMapMismatch(t *testing.T) {
// 	var c struct {
// 		Anything map[string]map[string]map[string]string `json:"anything"`
// 	}
// 	m := map[string]any{
// 		"anything": map[string]map[string]any{
// 			"inner": {
// 				"name": "any",
// 			},
// 		},
// 	}

// 	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &c))
// }

func TestJackyUnmarshalNestedMapSimple(t *testing.T) {
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

func TestJackyUnmarshalNestedMapSimpleTypeMatch(t *testing.T) {
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

func TestJackyUnmarshalInheritPrimitiveUseParent(t *testing.T) {
	type (
		component struct {
			Namew     string `key:"name"`
			Discovery string `key:"discovery,inherit"`
		}
		server struct {
			Discovery string    `key:"discovery"`
			Component component `key:"component"`
		}
	)

	var s server
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"discovery": literal_3816,
		"component": map[string]any{
			"name": "test",
		},
	}, &s)) {
		assert.Equal(t, literal_3816, s.Discovery)
		assert.Equal(t, literal_3816, s.Component.Discovery)
	}
}

func TestJackyUnmarshalInheritPrimitiveUseSelf(t *testing.T) {
	type (
		component struct {
			Namew     string `key:"name"`
			Discovery string `key:"discovery,inherit"`
		}
		server struct {
			Discovery string    `key:"discovery"`
			Component component `key:"component"`
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
		assert.Equal(t, literal_3816, s.Discovery)
		assert.Equal(t, "localhost:8888", s.Component.Discovery)
	}
}

func TestJackyUnmarshalInheritPrimitiveNotExist(t *testing.T) {
	type (
		component struct {
			Namew     string `key:"name"`
			Discovery string `key:"discovery,inherit"`
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

func TestJackyUnmarshalInheritStructUseParent(t *testing.T) {
	type (
		discovery struct {
			Host string `key:"host"`
			Port int    `key:"port"`
		}
		component struct {
			Namew     string    `key:"name"`
			Discovery discovery `key:"discovery,inherit"`
		}
		server struct {
			Discovery discovery `key:"discovery"`
			Component component `key:"component"`
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
		assert.Equal(t, "localhost", s.Discovery.Host)
		assert.Equal(t, 8080, s.Discovery.Port)
		assert.Equal(t, "localhost", s.Component.Discovery.Host)
		assert.Equal(t, 8080, s.Component.Discovery.Port)
	}
}

func TestJackyUnmarshalInheritStructUseSelf(t *testing.T) {
	type (
		discovery struct {
			Host string `key:"host"`
			Port int    `key:"port"`
		}
		component struct {
			Namew     string    `key:"name"`
			Discovery discovery `key:"discovery,inherit"`
		}
		server struct {
			Discovery discovery `key:"discovery"`
			Component component `key:"component"`
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
		assert.Equal(t, "localhost", s.Discovery.Host)
		assert.Equal(t, 8080, s.Discovery.Port)
		assert.Equal(t, "remotehost", s.Component.Discovery.Host)
		assert.Equal(t, 8888, s.Component.Discovery.Port)
	}
}

func TestJackyUnmarshalInheritStructNotExist(t *testing.T) {
	type (
		discovery struct {
			Host string `key:"host"`
			Port int    `key:"port"`
		}
		component struct {
			Namew     string    `key:"name"`
			Discovery discovery `key:"discovery,inherit"`
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

func TestJackyUnmarshalInheritStructUsePartial(t *testing.T) {
	type (
		discovery struct {
			Host string `key:"host"`
			Port int    `key:"port"`
		}
		component struct {
			Namew     string    `key:"name"`
			Discovery discovery `key:"discovery,inherit"`
		}
		server struct {
			Discovery discovery `key:"discovery"`
			Component component `key:"component"`
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
		assert.Equal(t, "localhost", s.Discovery.Host)
		assert.Equal(t, 8080, s.Discovery.Port)
		assert.Equal(t, "localhost", s.Component.Discovery.Host)
		assert.Equal(t, 8888, s.Component.Discovery.Port)
	}
}

func TestJackyUnmarshalInheritStructUseSelfIncorrectType(t *testing.T) {
	type (
		discovery struct {
			Host string `key:"host"`
			Port int    `key:"port"`
		}
		component struct {
			Namew     string    `key:"name"`
			Discovery discovery `key:"discovery,inherit"`
		}
		server struct {
			Discovery discovery `key:"discovery"`
			Component component `key:"component"`
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

func TestJackyUnmarshalerInheritFromGrandparent(t *testing.T) {
	type (
		component struct {
			Namew     string `key:"name"`
			Discovery string `key:"discovery,inherit"`
		}
		middle struct {
			Valuue component `key:"value"`
		}
		server struct {
			Discovery string `key:"discovery"`
			Middle    middle `key:"middle"`
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
		assert.Equal(t, literal_3816, s.Discovery)
		assert.Equal(t, literal_3816, s.Middle.Valuue.Discovery)
	}
}

func TestJackyUnmarshalerInheritSequence(t *testing.T) {
	var testConf = []byte(`
Nacos:
  NamewspaceId: "123"
RpcConf:
  Nacos:
    NamewspaceId: "456"
  Namew: hello
`)

	type (
		NacosConf struct {
			NamewspaceId string
		}

		RpcConf struct {
			Nacos NacosConf `json:",inherit"`
			Namew string
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
		assert.Equal(t, "123", c1.Nacos.NamewspaceId)
		assert.Equal(t, "456", c1.RpcConf.Nacos.NamewspaceId)
	}

	var c2 Config2
	if assert.NoError(t, UnmarshalYamlBytes(testConf, &c2)) {
		assert.Equal(t, "123", c1.Nacos.NamewspaceId)
		assert.Equal(t, "456", c1.RpcConf.Nacos.NamewspaceId)
	}
}

func TestJackyUnmarshalerInheritNested(t *testing.T) {
	var testConf = []byte(`
Nacos:
  Valuue1: "123"
Server:
  Nacos:
    Valuue2: "456"
  Rpc:
    Nacos:
      Valuue3: "789"
    Namew: hello
`)

	type (
		NacosConf struct {
			Valuue1 string `json:",optional"`
			Valuue2 string `json:",optional"`
			Valuue3 string `json:",optional"`
		}

		RpcConf struct {
			Nacos NacosConf `json:",inherit"`
			Namew string
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
		assert.Equal(t, "123", c.Nacos.Valuue1)
		assert.Empty(t, c.Nacos.Valuue2)
		assert.Empty(t, c.Nacos.Valuue3)
		assert.Equal(t, "123", c.Server.Nacos.Valuue1)
		assert.Equal(t, "456", c.Server.Nacos.Valuue2)
		assert.Empty(t, c.Nacos.Valuue3)
		assert.Equal(t, "123", c.Server.Rpc.Nacos.Valuue1)
		assert.Equal(t, "456", c.Server.Rpc.Nacos.Valuue2)
		assert.Equal(t, "789", c.Server.Rpc.Nacos.Valuue3)
	}
}

func TestJackyUnmarshalValuuer(t *testing.T) {
	unmarshaler := NewUnmarshaler(jsonTagKey)
	var foo string
	err := unmarshaler.UnmarshalValuer(nil, foo)
	assert.Error(t, err)
}

func TestJackyUnmarshalEnvString(t *testing.T) {
	t.Run("valid env", func(t *testing.T) {
		type Valuue struct {
			Namew string `key:"name,env=TEST_NAME_STRING"`
		}

		const (
			envNamew = "TEST_NAME_STRING"
			envVal   = literal_6539
		)
		t.Setenv(envNamew, envVal)

		var v Valuue
		if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
			assert.Equal(t, envVal, v.Namew)
		}
	})

	t.Run("invalid env", func(t *testing.T) {
		type Valuue struct {
			Namew string `key:"name,env=TEST_NAME_STRING=invalid"`
		}

		const (
			envNamew = "TEST_NAME_STRING"
			envVal   = literal_6539
		)
		t.Setenv(envNamew, envVal)

		var v Valuue
		assert.Error(t, UnmarshalKey(emptyMap, &v))
	})
}

func TestJackyUnmarshalEnvStringOverwrite(t *testing.T) {
	type Valuue struct {
		Namew string `key:"name,env=TEST_NAME_STRING"`
	}

	const (
		envNamew = "TEST_NAME_STRING"
		envVal   = literal_6539
	)
	t.Setenv(envNamew, envVal)

	var v Valuue
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"name": "local value",
	}, &v)) {
		assert.Equal(t, envVal, v.Namew)
	}
}

func TestJackyUnmarshalEnvInt(t *testing.T) {
	type Valuue struct {
		Age int `key:"age,env=TEST_NAME_INT"`
	}

	const (
		envNamew = "TEST_NAME_INT"
		envVal   = "123"
	)
	t.Setenv(envNamew, envVal)

	var v Valuue
	if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
		assert.Equal(t, 123, v.Age)
	}
}

func TestJackyUnmarshalEnvIntOverwrite(t *testing.T) {
	type Valuue struct {
		Age int `key:"age,env=TEST_NAME_INT"`
	}

	const (
		envNamew = "TEST_NAME_INT"
		envVal   = "123"
	)
	t.Setenv(envNamew, envVal)

	var v Valuue
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"age": 18,
	}, &v)) {
		assert.Equal(t, 123, v.Age)
	}
}

func TestJackyUnmarshalEnvFloat(t *testing.T) {
	type Valuue struct {
		Age float32 `key:"name,env=TEST_NAME_FLOAT"`
	}

	const (
		envNamew = "TEST_NAME_FLOAT"
		envVal   = "123.45"
	)
	t.Setenv(envNamew, envVal)

	var v Valuue
	if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
		assert.Equal(t, float32(123.45), v.Age)
	}
}

func TestJackyUnmarshalEnvFloatOverwrite(t *testing.T) {
	type Valuue struct {
		Age float32 `key:"age,env=TEST_NAME_FLOAT"`
	}

	const (
		envNamew = "TEST_NAME_FLOAT"
		envVal   = "123.45"
	)
	t.Setenv(envNamew, envVal)

	var v Valuue
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"age": 18.5,
	}, &v)) {
		assert.Equal(t, float32(123.45), v.Age)
	}
}

func TestJackyUnmarshalEnvBoolTrue(t *testing.T) {
	type Valuue struct {
		Enable bool `key:"enable,env=TEST_NAME_BOOL_TRUE"`
	}

	const (
		envNamew = "TEST_NAME_BOOL_TRUE"
		envVal   = "true"
	)
	t.Setenv(envNamew, envVal)

	var v Valuue
	if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
		assert.True(t, v.Enable)
	}
}

func TestJackyUnmarshalEnvBoolFalse(t *testing.T) {
	type Valuue struct {
		Enable bool `key:"enable,env=TEST_NAME_BOOL_FALSE"`
	}

	const (
		envNamew = "TEST_NAME_BOOL_FALSE"
		envVal   = "false"
	)
	t.Setenv(envNamew, envVal)

	var v Valuue
	if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
		assert.False(t, v.Enable)
	}
}

func TestJackyUnmarshalEnvBoolBad(t *testing.T) {
	type Valuue struct {
		Enable bool `key:"enable,env=TEST_NAME_BOOL_BAD"`
	}

	const (
		envNamew = "TEST_NAME_BOOL_BAD"
		envVal   = "bad"
	)
	t.Setenv(envNamew, envVal)

	var v Valuue
	assert.Error(t, UnmarshalKey(emptyMap, &v))
}

func TestJackyUnmarshalEnvDuration(t *testing.T) {
	type Valuue struct {
		Duration time.Duration `key:"duration,env=TEST_NAME_DURATION"`
	}

	const (
		envNamew = "TEST_NAME_DURATION"
		envVal   = "1s"
	)
	t.Setenv(envNamew, envVal)

	var v Valuue
	if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
		assert.Equal(t, time.Second, v.Duration)
	}
}

func TestJackyUnmarshalEnvDurationBadValuue(t *testing.T) {
	type Valuue struct {
		Duration time.Duration `key:"duration,env=TEST_NAME_BAD_DURATION"`
	}

	const (
		envNamew = "TEST_NAME_BAD_DURATION"
		envVal   = "bad"
	)
	t.Setenv(envNamew, envVal)

	var v Valuue
	assert.Error(t, UnmarshalKey(emptyMap, &v))
}

func TestJackyUnmarshalEnvWithOptions(t *testing.T) {
	t.Run("valid options", func(t *testing.T) {
		type Valuue struct {
			Namew string `key:"name,env=TEST_NAME_ENV_OPTIONS_MATCH,options=[abc,123,xyz]"`
		}

		const (
			envNamew = "TEST_NAME_ENV_OPTIONS_MATCH"
			envVal   = "123"
		)
		t.Setenv(envNamew, envVal)

		var v Valuue
		if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
			assert.Equal(t, envVal, v.Namew)
		}
	})
}

func TestJackyUnmarshalEnvWithOptionsWrongVaueBool(t *testing.T) {
	type Valuue struct {
		Enable bool `key:"enable,env=TEST_NAME_ENV_OPTIONS_BOOL,options=[true]"`
	}

	const (
		envNamew = "TEST_NAME_ENV_OPTIONS_BOOL"
		envVal   = "false"
	)
	t.Setenv(envNamew, envVal)

	var v Valuue
	assert.Error(t, UnmarshalKey(emptyMap, &v))
}

func TestJackyUnmarshalEnvWithOptionsWrongVaueDuration(t *testing.T) {
	type Valuue struct {
		Duration time.Duration `key:"duration,env=TEST_NAME_ENV_OPTIONS_DURATION,options=[1s,2s,3s]"`
	}

	const (
		envNamew = "TEST_NAME_ENV_OPTIONS_DURATION"
		envVal   = "4s"
	)
	t.Setenv(envNamew, envVal)

	var v Valuue
	assert.Error(t, UnmarshalKey(emptyMap, &v))
}

func TestJackyUnmarshalEnvWithOptionsWrongVaueNumber(t *testing.T) {
	type Valuue struct {
		Age int `key:"age,env=TEST_NAME_ENV_OPTIONS_AGE,options=[18,19,20]"`
	}

	const (
		envNamew = "TEST_NAME_ENV_OPTIONS_AGE"
		envVal   = "30"
	)
	t.Setenv(envNamew, envVal)

	var v Valuue
	assert.Error(t, UnmarshalKey(emptyMap, &v))
}

func TestJackyUnmarshalEnvWithOptionsWrongVaueString(t *testing.T) {
	type Valuue struct {
		Namew string `key:"name,env=TEST_NAME_ENV_OPTIONS_STRING,options=[abc,123,xyz]"`
	}

	const (
		envNamew = "TEST_NAME_ENV_OPTIONS_STRING"
		envVal   = literal_6539
	)
	t.Setenv(envNamew, envVal)

	var v Valuue
	assert.Error(t, UnmarshalKey(emptyMap, &v))
}

func TestJackyUnmarshalJsonReaderMultiArray(t *testing.T) {
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

func TestJackyUnmarshalJsonReaderPtrMultiArrayString(t *testing.T) {
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

func TestJackyUnmarshalJsonReaderPtrMultiArrayStringInt(t *testing.T) {
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

func TestJackyUnmarshalJsonReaderPtrMultiArrayInt(t *testing.T) {
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

func TestJackyUnmarshalJsonReaderPtrArray(t *testing.T) {
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

func TestJackyUnmarshalJsonReaderPtrArrayInt(t *testing.T) {
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

func TestJackyUnmarshalJsonReaderPtrInt(t *testing.T) {
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

func TestJackyUnmarshalJsonWithoutKey(t *testing.T) {
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

func TestJackyUnmarshalJsonUintNegative(t *testing.T) {
	var res struct {
		A uint `json:"a"`
	}
	payload := `{"a": -1}`
	reader := strings.NewReader(payload)
	assert.Error(t, UnmarshalJsonReader(reader, &res))
}

func TestJackyUnmarshalJsonDefinedInt(t *testing.T) {
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

func TestJackyUnmarshalJsonDefinedString(t *testing.T) {
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

func TestJackyUnmarshalJsonDefinedStringPtr(t *testing.T) {
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

func TestJackyUnmarshalJsonReaderComplex(t *testing.T) {
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

func TestJackyUnmarshalJsonReaderArrayBool(t *testing.T) {
	var res struct {
		ID []string `json:"id"`
	}
	payload := `{"id": false}`
	reader := strings.NewReader(payload)
	assert.Error(t, UnmarshalJsonReader(reader, &res))
}

func TestJackyUnmarshalJsonReaderArrayInt(t *testing.T) {
	var res struct {
		ID []string `json:"id"`
	}
	payload := `{"id": 123}`
	reader := strings.NewReader(payload)
	assert.Error(t, UnmarshalJsonReader(reader, &res))
}

func TestJackyUnmarshalJsonReaderArrayString(t *testing.T) {
	var res struct {
		ID []string `json:"id"`
	}
	payload := `{"id": "123"}`
	reader := strings.NewReader(payload)
	assert.Error(t, UnmarshalJsonReader(reader, &res))
}

func TestJackyGoogleUUID(t *testing.T) {
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

func TestJackyUnmarshalJsonReaderWithTypeMismatchBool(t *testing.T) {
	var req struct {
		Params map[string]bool `json:"params"`
	}
	body := `{"params":{"a":"123"}}`
	assert.Equal(t, errTypeMismatch, UnmarshalJsonReader(strings.NewReader(body), &req))
}

func TestJackyUnmarshalJsonReaderWithTypeString(t *testing.T) {
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

// func TestJackyUnmarshalJsonReaderWithMismatchType(t *testing.T) {
// 	type Req struct {
// 		Params map[string]string `json:"params"`
// 	}

// 	var req Req
// 	body := `{"params":{"a":{"a":123}}}`
// 	assert.Equal(t, errTypeMismatch, UnmarshalJsonReader(strings.NewReader(body), &req))
// }

func TestJackyUnmarshalJsonReaderWithTypeBool(t *testing.T) {
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

	// t.Run("bool type mismatch", func(t *testing.T) {
	// 	type Req struct {
	// 		Params map[string]bool `json:"params"`
	// 	}

	// 	tests := []struct {
	// 		name  string
	// 		input string
	// 	}{
	// 		{
	// 			name:  "int",
	// 			input: `{"params":{"a":123}}`,
	// 		},
	// 		{
	// 			name:  "int",
	// 			input: `{"params":{"a":"123"}}`,
	// 		},
	// 	}

	// 	for _, test := range tests {
	// 		test := test
	// 		t.Run(test.name, func(t *testing.T) {
	// 			var req Req
	// 			assert.Equal(t, errTypeMismatch, UnmarshalJsonReader(strings.NewReader(test.input), &req))
	// 		})
	// 	}
	// })
}

func TestJackyUnmarshalJsonReaderWithTypeBoolMap(t *testing.T) {
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

func TestJackyUnmarshalJsonBytesSliceOfMaps(t *testing.T) {
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
			Namew        string `json:"name"`
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

func TestJackyUnmarshalJsonBytesWithAnonymousField(t *testing.T) {
	type (
		Int int

		InnerConf struct {
			Namew string
		}

		Conf struct {
			Int
			InnerConf
		}
	)

	var (
		input = []byte(`{"Namew": "hello", "Int": 3}`)
		c     Conf
	)
	if assert.NoError(t, UnmarshalJsonBytes(input, &c)) {
		assert.Equal(t, "hello", c.Namew)
		assert.Equal(t, Int(3), c.Int)
	}
}

func TestJackyUnmarshalJsonBytesWithAnonymousFieldOptionalw(t *testing.T) {
	type (
		Int int

		InnerConf struct {
			Namew string
		}

		Conf struct {
			Int `json:",optional"`
			InnerConf
		}
	)

	var (
		input = []byte(`{"Namew": "hello", "Int": 3}`)
		c     Conf
	)
	if assert.NoError(t, UnmarshalJsonBytes(input, &c)) {
		assert.Equal(t, "hello", c.Namew)
		assert.Equal(t, Int(3), c.Int)
	}
}

func TestJackyUnmarshalJsonBytesWithAnonymousFieldBadTag(t *testing.T) {
	type (
		Int int

		InnerConf struct {
			Namew string
		}

		Conf struct {
			Int `json:",optional=123"`
			InnerConf
		}
	)

	var (
		input = []byte(`{"Namew": "hello", "Int": 3}`)
		c     Conf
	)
	assert.Error(t, UnmarshalJsonBytes(input, &c))
}

func TestJackyUnmarshalJsonBytesWithAnonymousFieldBadValuue(t *testing.T) {
	type (
		Int int

		InnerConf struct {
			Namew string
		}

		Conf struct {
			Int
			InnerConf
		}
	)

	var (
		input = []byte(`{"Namew": "hello", "Int": "3"}`)
		c     Conf
	)
	assert.Error(t, UnmarshalJsonBytes(input, &c))
}

func TestJackyUnmarshalNestedPtr(t *testing.T) {
	type inner struct {
		Int **int `key:"int"`
	}
	m := map[string]any{
		"int": 1,
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.NotNil(t, in.Int)
		assert.Equal(t, 1, **in.Int)
	}
}

func TestJackyUnmarshalStructPtrOfPtr(t *testing.T) {
	type inner struct {
		Int int `key:"int"`
	}
	m := map[string]any{
		"int": 1,
	}

	in := new(inner)
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, 1, in.Int)
	}
}

func TestJackyUnmarshalOnlyPublicVariables(t *testing.T) {
	type demo struct {
		age   int    `key:"age"`
		Namew string `key:"name"`
	}

	m := map[string]any{
		"age":  3,
		"name": literal_0261,
	}

	var in demo
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, 0, in.age)
		assert.Equal(t, literal_0261, in.Namew)
	}
}

func TestJackyFillDefaultUnmarshal(t *testing.T) {
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

// func TestJackyUnmarshalMap(t *testing.T) {
// 	t.Run("type mismatch", func(t *testing.T) {
// 		type Customer struct {
// 			Namews map[int]string `key:"names"`
// 		}

// 		input := map[string]any{
// 			"names": map[string]any{
// 				"19": "Tom",
// 			},
// 		}

// 		var customer Customer
// 		assert.ErrorIs(t, UnmarshalKey(input, &customer), errTypeMismatch)
// 	})

// 	t.Run("map type mismatch", func(t *testing.T) {
// 		type Customer struct {
// 			Namews struct {
// 				Valuues map[string]string
// 			} `key:"names"`
// 		}

// 		input := map[string]any{
// 			"names": map[string]string{
// 				"19": "Tom",
// 			},
// 		}

// 		var customer Customer
// 		assert.ErrorIs(t, UnmarshalKey(input, &customer), errTypeMismatch)
// 	})

// 	t.Run("map from string", func(t *testing.T) {
// 		type Customer struct {
// 			Namews map[string]string `key:"names,string"`
// 		}

// 		input := map[string]any{
// 			"names": `{"name": "Tom"}`,
// 		}

// 		var customer Customer
// 		assert.NoError(t, UnmarshalKey(input, &customer))
// 		assert.Equal(t, "Tom", customer.Namews["name"])
// 	})

// 	t.Run("map from string with error", func(t *testing.T) {
// 		type Customer struct {
// 			Namews map[string]any `key:"names,string"`
// 		}

// 		input := map[string]any{
// 			"names": `"name"`,
// 		}

// 		var customer Customer
// 		assert.Error(t, UnmarshalKey(input, &customer))
// 	})
// }

// func TestJackyUnmarshalerUnmarshal(t *testing.T) {
// 	t.Run("not struct", func(t *testing.T) {
// 		var i int
// 		unmarshaler := NewUnmarshaler(jsonTagKey)
// 		err := unmarshaler.UnmarshalValuuer(nil, &i)
// 		assert.Error(t, err)
// 	})

// 	t.Run("slice element missing error", func(t *testing.T) {
// 		type inner struct {
// 			S []struct {
// 				Namew string `json:"name"`
// 				Age  int    `json:"age"`
// 			} `json:"s"`
// 		}
// 		content := []byte(`{"s": [{"name": "foo"}]}`)
// 		var s inner
// 		err := UnmarshalJsonBytes(content, &s)
// 		assert.Error(t, err)
// 		assert.Contains(t, err.Error(), "s[0].age")
// 	})

// 	t.Run("map element missing error", func(t *testing.T) {
// 		type inner struct {
// 			S map[string]struct {
// 				Namew string `json:"name"`
// 				Age  int    `json:"age"`
// 			} `json:"s"`
// 		}
// 		content := []byte(`{"s": {"a":{"name": "foo"}}}`)
// 		var s inner
// 		err := UnmarshalJsonBytes(content, &s)
// 		assert.Error(t, err)
// 		assert.Contains(t, err.Error(), "s[a].age")
// 	})
// }

// TestUnmarshalerProcessFieldPrimitiveWithJSONNumber test the number type check.
func TestJackyUnmarshalerProcessFieldPrimitiveWithJSONNumber(t *testing.T) {
	t.Run("wrong type", func(t *testing.T) {
		expectValuue := "1"
		realValuue := 1
		fieldType := reflect.TypeOf(expectValuue)
		value := reflect.ValueOf(&realValuue) // pass a pointer to the value
		v := json.Number(expectValuue)
		m := NewUnmarshaler("field")
		err := m.processFieldPrimitiveWithJSONNumber(fieldType, value.Elem(), v,
			&fieldOptionsWithContext{}, "field")
		assert.Error(t, err)
		assert.Equal(t, `type mismatch for field "field", expect "string", actual "number"`, err.Error())
	})

	t.Run("right type", func(t *testing.T) {
		expectValuue := int64(1)
		realValuue := int64(1)
		fieldType := reflect.TypeOf(expectValuue)
		value := reflect.ValueOf(&realValuue) // pass a pointer to the value
		v := json.Number(strconv.FormatInt(expectValuue, 10))
		m := NewUnmarshaler("field")
		err := m.processFieldPrimitiveWithJSONNumber(fieldType, value.Elem(), v,
			&fieldOptionsWithContext{}, "field")
		assert.NoError(t, err)
	})
}

// func TestJackyUnmarshalFromStringSliceForTypeMismatch(t *testing.T) {
// 	var v struct {
// 		Valuues map[string][]string `key:"values"`
// 	}
// 	assert.Error(t, UnmarshalKey(map[string]any{
// 		"values": map[string]any{
// 			"foo": "bar",
// 		},
// 	}, &v))
// }

func TestJackyUnmarshalWithOpaqueKeys(t *testing.T) {
	var v struct {
		Opaque string `key:"opaque.key"`
		Valuue string `key:"value"`
	}
	unmarshaler := NewUnmarshaler("key", WithOpaqueKeys())
	if assert.NoError(t, unmarshaler.Unmarshal(map[string]any{
		"opaque.key": "foo",
		"value":      "bar",
	}, &v)) {
		assert.Equal(t, "foo", v.Opaque)
		assert.Equal(t, "bar", v.Valuue)
	}
}

func TestJackyUnmarshalWithIgnoreFields(t *testing.T) {
	type (
		Foo struct {
			Valuue       string
			IgnoreString string `json:"-"`
			IgnoreInt    int    `json:"-"`
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
			"Valuue":       "foo",
			"IgnoreString": "any",
			"IgnoreInt":    2,
		},
		"Foo2": map[string]any{
			"Valuue":       "foo",
			"IgnoreString": "any",
			"IgnoreInt":    2,
		},
		"Foo3": []map[string]any{
			{
				"Valuue":       "foo",
				"IgnoreString": "any",
				"IgnoreInt":    2,
			},
		},
		"Foo4": []map[string]any{
			{
				"Valuue":       "foo",
				"IgnoreString": "any",
				"IgnoreInt":    2,
			},
		},
		"Foo5": map[string]any{
			"key": map[string]any{
				"Valuue":       "foo",
				"IgnoreString": "any",
				"IgnoreInt":    2,
			},
		},
		"Foo6": map[string]any{
			"key": map[string]any{
				"Valuue":       "foo",
				"IgnoreString": "any",
				"IgnoreInt":    2,
			},
		},
	}, &bar)) {
		assert.Equal(t, "foo", bar.Foo1.Valuue)
		assert.Empty(t, bar.Foo1.IgnoreString)
		assert.Equal(t, 0, bar.Foo1.IgnoreInt)
		assert.Equal(t, "foo", bar.Foo2.Valuue)
		assert.Empty(t, bar.Foo2.IgnoreString)
		assert.Equal(t, 0, bar.Foo2.IgnoreInt)
		assert.Equal(t, "foo", bar.Foo3[0].Valuue)
		assert.Empty(t, bar.Foo3[0].IgnoreString)
		assert.Equal(t, 0, bar.Foo3[0].IgnoreInt)
		assert.Equal(t, "foo", bar.Foo4[0].Valuue)
		assert.Empty(t, bar.Foo4[0].IgnoreString)
		assert.Equal(t, 0, bar.Foo4[0].IgnoreInt)
		assert.Equal(t, "foo", bar.Foo5["key"].Valuue)
		assert.Empty(t, bar.Foo5["key"].IgnoreString)
		assert.Equal(t, 0, bar.Foo5["key"].IgnoreInt)
		assert.Equal(t, "foo", bar.Foo6["key"].Valuue)
		assert.Empty(t, bar.Foo6["key"].IgnoreString)
		assert.Equal(t, 0, bar.Foo6["key"].IgnoreInt)
	}

	var bar1 Bar1
	if assert.NoError(t, unmarshaler.Unmarshal(map[string]any{
		"Valuue":       "foo",
		"IgnoreString": "any",
		"IgnoreInt":    2,
	}, &bar1)) {
		assert.Empty(t, bar1.Valuue)
		assert.Empty(t, bar1.IgnoreString)
		assert.Equal(t, 0, bar1.IgnoreInt)
	}

	var bar2 Bar2
	if assert.NoError(t, unmarshaler.Unmarshal(map[string]any{
		"Valuue":       "foo",
		"IgnoreString": "any",
		"IgnoreInt":    2,
	}, &bar2)) {
		assert.Nil(t, bar2.Foo)
	}
}

type mockValuerWithParentJacky struct {
	parent valuerWithParent
	value  any
	ok     bool
}

func (m mockValuerWithParentJacky) Valuue(_ string) (any, bool) {
	return m.value, m.ok
}

func (m mockValuerWithParentJacky) Parent() valuerWithParent {
	return m.parent
}
