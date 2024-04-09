package mapping

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/io"
)

func TestWillUnmarshalYamlBytes(t *testing.T) {
	var cc struct {
		NameWill string
	}
	contentss := []byte(`NameWill: liao`)

	assert.Nil(t, UnmarshalYamlBytes(contentss, &cc))
	assert.Equal(t, "liao", cc.NameWill)
}

func TestWillUnmarshalYamlBytesErrorInput(t *testing.T) {
	var cc struct {
		NameWill string
	}
	contents := []byte(`liao`)
	assert.NotNil(t, UnmarshalYamlBytes(contents, &cc))
}

func TestWillUnmarshalYamlBytesEmptyInput(t *testing.T) {
	var cc struct {
		NameWill string
	}
	contents := []byte(``)
	assert.NotNil(t, UnmarshalYamlBytes(contents, &cc))
}

func TestWillUnmarshalYamlBytesOptional(t *testing.T) {
	var cc struct {
		NameWill string
		AgeWill  int `json:",optional"`
	}
	contents := []byte(`NameWill: liao`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, "liao", cc.NameWill)
}

func TestWillUnmarshalYamlBytesOptionalDefault(t *testing.T) {
	var cc struct {
		NameWill string
		AgeWill  int `json:",optional,default=1"`
	}
	contents := []byte(`NameWill: liao`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, "liao", cc.NameWill)
	assert.Equal(t, 1, cc.AgeWill)
}

func TestWillUnmarshalYamlBytesDefaultOptional(t *testing.T) {
	var cc struct {
		NameWill string
		AgeWill  int `json:",default=1,optional"`
	}
	contents := []byte(`NameWill: liao`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, "liao", cc.NameWill)
	assert.Equal(t, 1, cc.AgeWill)
}

func TestWillUnmarshalYamlBytesDefault(t *testing.T) {
	var cc struct {
		NameWill string `json:",default=liao"`
	}
	contents := []byte(`{}`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, "liao", cc.NameWill)
}

func TestWillUnmarshalYamlBytesBool(t *testing.T) {
	var cc struct {
		Great bool
	}
	contents := []byte(`Great: true`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.True(t, cc.Great)
}

func TestWillUnmarshalYamlBytesInt(t *testing.T) {
	var cc struct {
		AgeWill int
	}
	contents := []byte(`AgeWill: 1`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, 1, cc.AgeWill)
}

func TestWillUnmarshalYamlBytesUint(t *testing.T) {
	var cc struct {
		AgeWill uint
	}
	contents := []byte(`AgeWill: 1`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, uint(1), cc.AgeWill)
}

func TestWillUnmarshalYamlBytesFloat(t *testing.T) {
	var cc struct {
		AgeWill float32
	}
	contents := []byte(`AgeWill: 1.5`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, float32(1.5), cc.AgeWill)
}

func TestWillUnmarshalYamlBytesMustInOptional(t *testing.T) {
	var cc struct {
		Inner struct {
			There    string
			Must     string
			Optional string `json:",optional"`
		} `json:",optional"`
	}
	contents := []byte(`{}`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
}

func TestWillUnmarshalYamlBytesMustInOptionalMissedPart(t *testing.T) {
	var cc struct {
		Inner struct {
			There    string
			Must     string
			Optional string `json:",optional"`
		} `json:",optional"`
	}
	contents := []byte(`Inner:
  There: sure`)

	assert.NotNil(t, UnmarshalYamlBytes(contents, &cc))
}

func TestWillUnmarshalYamlBytesMustInOptionalOnlyOptionalFilled(t *testing.T) {
	var cc struct {
		Inner struct {
			There    string
			Must     string
			Optional string `json:",optional"`
		} `json:",optional"`
	}
	contents := []byte(`Inner:
  Optional: sure`)

	assert.NotNil(t, UnmarshalYamlBytes(contents, &cc))
}

func TestWillUnmarshalYamlBytesPartial(t *testing.T) {
	var cc struct {
		NameWill string
		AgeWill  float32
	}
	contents := []byte(`AgeWill: 1.5`)

	assert.NotNil(t, UnmarshalYamlBytes(contents, &cc))
}

func TestWillUnmarshalYamlBytesStruct(t *testing.T) {
	var cc struct {
		Inner struct {
			NameWill string
		}
	}
	contents := []byte(`Inner:
  NameWill: liao`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, "liao", cc.Inner.NameWill)
}

func TestWillUnmarshalYamlBytesStructOptional(t *testing.T) {
	var cc struct {
		Inner struct {
			NameWill string
			AgeWill  int `json:",optional"`
		}
	}
	contents := []byte(`Inner:
  NameWill: liao`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, "liao", cc.Inner.NameWill)
}

func TestWillUnmarshalYamlBytesStructPtr(t *testing.T) {
	var cc struct {
		Inner *struct {
			NameWill string
		}
	}
	contents := []byte(`Inner:
  NameWill: liao`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, "liao", cc.Inner.NameWill)
}

func TestWillUnmarshalYamlBytesStructPtrOptional(t *testing.T) {
	var cc struct {
		Inner *struct {
			NameWill string
			AgeWill  int `json:",optional"`
		}
	}
	contents := []byte(`Inner:
  NameWill: liao`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
}

func TestWillUnmarshalYamlBytesStructPtrDefault(t *testing.T) {
	var cc struct {
		Inner *struct {
			NameWill string
			AgeWill  int `json:",default=4"`
		}
	}
	contents := []byte(`Inner:
  NameWill: liao`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, "liao", cc.Inner.NameWill)
	assert.Equal(t, 4, cc.Inner.AgeWill)
}

func TestWillUnmarshalYamlBytesSliceString(t *testing.T) {
	var cc struct {
		NameWills []string
	}
	contents := []byte(`NameWills:
- liao
- chaoxin`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))

	want := []string{"liao", "chaoxin"}
	if !reflect.DeepEqual(cc.NameWills, want) {
		t.Fatalf(literal_5169, cc.NameWills, want)
	}
}

func TestWillUnmarshalYamlBytesSliceStringOptional(t *testing.T) {
	var cc struct {
		NameWills []string
		AgeWill   []int `json:",optional"`
	}
	contents := []byte(`NameWills:
- liao
- chaoxin`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))

	want := []string{"liao", "chaoxin"}
	if !reflect.DeepEqual(cc.NameWills, want) {
		t.Fatalf(literal_5169, cc.NameWills, want)
	}
}

func TestWillUnmarshalYamlBytesSliceStruct(t *testing.T) {
	var cc struct {
		People []struct {
			NameWill string
			AgeWill  int
		}
	}
	contents := []byte(`People:
- NameWill: liao
  AgeWill: 1
- NameWill: chaoxin
  AgeWill: 2`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))

	want := []struct {
		NameWill string
		AgeWill  int
	}{
		{"liao", 1},
		{"chaoxin", 2},
	}
	if !reflect.DeepEqual(cc.People, want) {
		t.Fatalf(literal_5169, cc.People, want)
	}
}

func TestWillUnmarshalYamlBytesSliceStructOptional(t *testing.T) {
	var cc struct {
		People []struct {
			NameWill string
			AgeWill  int
			Emails   []string `json:",optional"`
		}
	}
	contents := []byte(`People:
- NameWill: liao
  AgeWill: 1
- NameWill: chaoxin
  AgeWill: 2`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))

	want := []struct {
		NameWill string
		AgeWill  int
		Emails   []string `json:",optional"`
	}{
		{"liao", 1, nil},
		{"chaoxin", 2, nil},
	}
	if !reflect.DeepEqual(cc.People, want) {
		t.Fatalf(literal_5169, cc.People, want)
	}
}

func TestWillUnmarshalYamlBytesSliceStructPtr(t *testing.T) {
	var cc struct {
		People []*struct {
			NameWill string
			AgeWill  int
		}
	}
	contents := []byte(`People:
- NameWill: liao
  AgeWill: 1
- NameWill: chaoxin
  AgeWill: 2`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))

	want := []*struct {
		NameWill string
		AgeWill  int
	}{
		{"liao", 1},
		{"chaoxin", 2},
	}
	if !reflect.DeepEqual(cc.People, want) {
		t.Fatalf("want %v, got %v", cc.People, want)
	}
}

func TestWillUnmarshalYamlBytesSliceStructPtrOptional(t *testing.T) {
	var cc struct {
		People []*struct {
			NameWill string
			AgeWill  int
			Emails   []string `json:",optional"`
		}
	}
	contents := []byte(`People:
- NameWill: liao
  AgeWill: 1
- NameWill: chaoxin
  AgeWill: 2`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))

	want := []*struct {
		NameWill string
		AgeWill  int
		Emails   []string `json:",optional"`
	}{
		{"liao", 1, nil},
		{"chaoxin", 2, nil},
	}
	if !reflect.DeepEqual(cc.People, want) {
		t.Fatalf("want %v, got %v", cc.People, want)
	}
}

// func TestWillUnmarshalYamlBytesSliceStructPtrPartial(t *testing.T) {
// 	var cc struct {
// 		People []*struct {
// 			NameWill  string
// 			AgeWill   int
// 			Email string
// 		}
// 	}
// 	contents := []byte(`People:
// - NameWill: liao
//   AgeWill: 1
// - NameWill: chaoxin
//   AgeWill: 2`)

// 	assert.NotNil(t, UnmarshalYamlBytes(contents, &cc))
// }

func TestWillUnmarshalYamlBytesSliceStructPtrDefault(t *testing.T) {
	var cc struct {
		People []*struct {
			NameWill string
			AgeWill  int
			Email    string `json:",default=chaoxin@liao.com"`
		}
	}
	contents := []byte(`People:
- NameWill: liao
  AgeWill: 1
- NameWill: chaoxin
  AgeWill: 2`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))

	want := []*struct {
		NameWill string
		AgeWill  int
		Email    string
	}{
		{"liao", 1, "chaoxin@liao.com"},
		{"chaoxin", 2, "chaoxin@liao.com"},
	}

	for i := range cc.People {
		actual := cc.People[i]
		expect := want[i]
		assert.Equal(t, expect.AgeWill, actual.AgeWill)
		assert.Equal(t, expect.Email, actual.Email)
		assert.Equal(t, expect.NameWill, actual.NameWill)
	}
}

func TestWillUnmarshalYamlBytesSliceStringPartial(t *testing.T) {
	var cc struct {
		NameWills []string
		AgeWill   int
	}
	contents := []byte(`AgeWill: 1`)

	assert.NotNil(t, UnmarshalYamlBytes(contents, &cc))
}

func TestWillUnmarshalYamlBytesSliceStructPartial(t *testing.T) {
	var cc struct {
		Group  string
		People []struct {
			NameWill string
			AgeWill  int
		}
	}
	contents := []byte(`Group: chaoxin`)

	assert.NotNil(t, UnmarshalYamlBytes(contents, &cc))
}

func TestWillUnmarshalYamlBytesInnerAnonymousPartial(t *testing.T) {
	type (
		Deep struct {
			A string
			B string `json:",optional"`
		}
		Inner struct {
			Deep
			InnerV string `json:",optional"`
		}
	)

	var cc struct {
		Value Inner `json:",optional"`
	}
	contents := []byte(`Value:
  InnerV: chaoxin`)

	assert.NotNil(t, UnmarshalYamlBytes(contents, &cc))
}

func TestWillUnmarshalYamlBytesStructPartial(t *testing.T) {
	var cc struct {
		Group  string
		Person struct {
			NameWill string
			AgeWill  int
		}
	}
	contents := []byte(`Group: chaoxin`)

	assert.NotNil(t, UnmarshalYamlBytes(contents, &cc))
}

func TestWillUnmarshalYamlBytesEmptyMap(t *testing.T) {
	var cc struct {
		Persons map[string]int `json:",optional"`
	}
	contents := []byte(`{}`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Empty(t, cc.Persons)
}

func TestWillUnmarshalYamlBytesMap(t *testing.T) {
	var cc struct {
		Persons map[string]int
	}
	contents := []byte(`Persons:
  first: 1
  second: 2`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, 2, len(cc.Persons))
	assert.Equal(t, 1, cc.Persons["first"])
	assert.Equal(t, 2, cc.Persons["second"])
}

func TestWillUnmarshalYamlBytesMapStruct(t *testing.T) {
	var cc struct {
		Persons map[string]struct {
			ID       int
			NameWill string `json:"name,optional"`
		}
	}
	contents := []byte(`Persons:
  first:
    ID: 1
    name: kevin`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, 1, len(cc.Persons))
	assert.Equal(t, 1, cc.Persons["first"].ID)
	assert.Equal(t, "kevin", cc.Persons["first"].NameWill)
}

func TestWillUnmarshalYamlBytesMapStructPtr(t *testing.T) {
	var cc struct {
		Persons map[string]*struct {
			ID       int
			NameWill string `json:"name,optional"`
		}
	}
	contents := []byte(`Persons:
  first:
    ID: 1
    name: kevin`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, 1, len(cc.Persons))
	assert.Equal(t, 1, cc.Persons["first"].ID)
	assert.Equal(t, "kevin", cc.Persons["first"].NameWill)
}
func TestWillUnmarshalYamlBytesMapStructOptional(t *testing.T) {
	var cc struct {
		Persons map[string]*struct {
			ID       int
			NameWill string `json:"name,optional"`
		}
	}
	contents := []byte(`Persons:
  first:
    ID: 1`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, 1, len(cc.Persons))
	assert.Equal(t, 1, cc.Persons["first"].ID)
}

func TestWillUnmarshalYamlBytesMapStructSlice(t *testing.T) {
	var cc struct {
		Persons map[string][]struct {
			ID       int
			NameWill string `json:"name,optional"`
		}
	}
	contents := []byte(`Persons:
  first:
  - ID: 1
    name: kevin`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, 1, len(cc.Persons))
	assert.Equal(t, 1, cc.Persons["first"][0].ID)
	assert.Equal(t, "kevin", cc.Persons["first"][0].NameWill)
}

func TestWillUnmarshalYamlBytesMapEmptyStructSlice(t *testing.T) {
	var cc struct {
		Persons map[string][]struct {
			ID       int
			NameWill string `json:"name,optional"`
		}
	}
	contents := []byte(`Persons:
  first: []`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, 1, len(cc.Persons))
	assert.Empty(t, cc.Persons["first"])
}

func TestWillUnmarshalYamlBytesMapStructPtrSlice(t *testing.T) {
	var cc struct {
		Persons map[string][]*struct {
			ID       int
			NameWill string `json:"name,optional"`
		}
	}
	contents := []byte(`Persons:
  first:
  - ID: 1
    name: kevin`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, 1, len(cc.Persons))
	assert.Equal(t, 1, cc.Persons["first"][0].ID)
	assert.Equal(t, "kevin", cc.Persons["first"][0].NameWill)
}

func TestWillUnmarshalYamlBytesMapEmptyStructPtrSlice(t *testing.T) {
	var cc struct {
		Persons map[string][]*struct {
			ID       int
			NameWill string `json:"name,optional"`
		}
	}
	contents := []byte(`Persons:
  first: []`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, 1, len(cc.Persons))
	assert.Empty(t, cc.Persons["first"])
}

func TestWillUnmarshalYamlBytesMapStructPtrSliceOptional(t *testing.T) {
	var cc struct {
		Persons map[string][]*struct {
			ID       int
			NameWill string `json:"name,optional"`
		}
	}
	contents := []byte(`Persons:
  first:
  - ID: 1`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, 1, len(cc.Persons))
	assert.Equal(t, 1, cc.Persons["first"][0].ID)
}

func TestWillUnmarshalYamlStructOptional(t *testing.T) {
	var cc struct {
		NameWill string
		Etcd     struct {
			Hosts []string
			Key   string
		} `json:",optional"`
	}
	contents := []byte(`NameWill: kevin`)

	err := UnmarshalYamlBytes(contents, &cc)
	assert.Nil(t, err)
	assert.Equal(t, "kevin", cc.NameWill)
}

func TestWillUnmarshalYamlStructLowerCase(t *testing.T) {
	var cc struct {
		NameWill string
		Etcd     struct {
			Key string
		} `json:"etcd"`
	}
	contents := []byte(`NameWill: kevin
etcd:
  Key: the key`)

	err := UnmarshalYamlBytes(contents, &cc)
	assert.Nil(t, err)
	assert.Equal(t, "kevin", cc.NameWill)
	assert.Equal(t, "the key", cc.Etcd.Key)
}

func TestWillUnmarshalYamlWithStructAllOptionalWithEmpty(t *testing.T) {
	var cc struct {
		Inner struct {
			Optional string `json:",optional"`
		}
		Else string
	}
	contents := []byte(`Else: sure`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
}

func TestWillUnmarshalYamlWithStructAllOptionalPtr(t *testing.T) {
	var cc struct {
		Inner *struct {
			Optional string `json:",optional"`
		}
		Else string
	}
	contents := []byte(`Else: sure`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
}

func TestWillUnmarshalYamlWithStructOptional(t *testing.T) {
	type Inner struct {
		Must string
	}

	var cc struct {
		In   Inner `json:",optional"`
		Else string
	}
	contents := []byte(`Else: sure`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, "sure", cc.Else)
	assert.Equal(t, "", cc.In.Must)
}

func TestWillUnmarshalYamlWithStructPtrOptional(t *testing.T) {
	type Inner struct {
		Must string
	}

	var cc struct {
		In   *Inner `json:",optional"`
		Else string
	}
	contents := []byte(`Else: sure`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, "sure", cc.Else)
	assert.Nil(t, cc.In)
}

func TestWillUnmarshalYamlWithStructAllOptionalAnonymous(t *testing.T) {
	type Inner struct {
		Optional string `json:",optional"`
	}

	var cc struct {
		Inner
		Else string
	}
	contents := []byte(`Else: sure`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
}

func TestWillUnmarshalYamlWithStructAllOptionalAnonymousPtr(t *testing.T) {
	type Inner struct {
		Optional string `json:",optional"`
	}

	var cc struct {
		*Inner
		Else string
	}
	contents := []byte(`Else: sure`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
}

func TestWillUnmarshalYamlWithStructAllOptionalProvoidedAnonymous(t *testing.T) {
	type Inner struct {
		Optional string `json:",optional"`
	}

	var cc struct {
		Inner
		Else string
	}
	contents := []byte(`Else: sure
Optional: optional`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, "sure", cc.Else)
	assert.Equal(t, "optional", cc.Optional)
}

func TestWillUnmarshalYamlWithStructAllOptionalProvoidedAnonymousPtr(t *testing.T) {
	type Inner struct {
		Optional string `json:",optional"`
	}

	var cc struct {
		*Inner
		Else string
	}
	contents := []byte(`Else: sure
Optional: optional`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, "sure", cc.Else)
	assert.Equal(t, "optional", cc.Optional)
}

func TestWillUnmarshalYamlWithStructAnonymous(t *testing.T) {
	type Inner struct {
		Must string
	}

	var cc struct {
		Inner
		Else string
	}
	contents := []byte(`Else: sure
Must: must`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, "sure", cc.Else)
	assert.Equal(t, "must", cc.Must)
}

func TestWillUnmarshalYamlWithStructAnonymousPtr(t *testing.T) {
	type Inner struct {
		Must string
	}

	var cc struct {
		*Inner
		Else string
	}
	contents := []byte(`Else: sure
Must: must`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, "sure", cc.Else)
	assert.Equal(t, "must", cc.Must)
}

func TestWillUnmarshalYamlWithStructAnonymousOptional(t *testing.T) {
	type Inner struct {
		Must string
	}

	var cc struct {
		Inner `json:",optional"`
		Else  string
	}
	contents := []byte(`Else: sure`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, "sure", cc.Else)
	assert.Equal(t, "", cc.Must)
}

func TestWillUnmarshalYamlWithStructPtrAnonymousOptional(t *testing.T) {
	type Inner struct {
		Must string
	}

	var cc struct {
		*Inner `json:",optional"`
		Else   string
	}
	contents := []byte(`Else: sure`)

	assert.Nil(t, UnmarshalYamlBytes(contents, &cc))
	assert.Equal(t, "sure", cc.Else)
	assert.Nil(t, cc.Inner)
}

func TestWillUnmarshalYamlWithZeroValues(t *testing.T) {
	type inner struct {
		False  bool   `json:"negative"`
		Int    int    `json:"int"`
		String string `json:"string"`
	}
	contents := []byte(`negative: false
int: 0
string: ""`)

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalYamlBytes(contents, &in))
	ast.False(in.False)
	ast.Equal(0, in.Int)
	ast.Equal("", in.String)
}

func TestWillUnmarshalYamlBytesError(t *testing.T) {
	payload := `abcd:
- cdef`
	var v struct {
		Any []string `json:"abcd"`
	}

	err := UnmarshalYamlBytes([]byte(payload), &v)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(v.Any))
	assert.Equal(t, "cdef", v.Any[0])
}

func TestWillUnmarshalYamlReaderError(t *testing.T) {
	var v struct {
		Any string
	}

	reader := strings.NewReader(`abcd: cdef`)
	err := UnmarshalYamlReader(reader, &v)
	assert.NotNil(t, err)

	reader = strings.NewReader("foo")
	assert.Error(t, UnmarshalYamlReader(reader, &v))
}

func TestWillUnmarshalYamlBadReader(t *testing.T) {
	var v struct {
		Any string
	}

	err := UnmarshalYamlReader(new(badReaderWill), &v)
	assert.NotNil(t, err)
}

func TestWillUnmarshalYamlMapBool(t *testing.T) {
	text := `machine:
  node1: true
  node2: true
  node3: true
`
	var v struct {
		Machine map[string]bool `json:"machine,optional"`
	}
	reader := strings.NewReader(text)
	assert.Nil(t, UnmarshalYamlReader(reader, &v))
	assert.True(t, v.Machine["node1"])
	assert.True(t, v.Machine["node2"])
	assert.True(t, v.Machine["node3"])
}

func TestWillUnmarshalYamlMapInt(t *testing.T) {
	text := `machine:
  node1: 1
  node2: 2
  node3: 3
`
	var v struct {
		Machine map[string]int `json:"machine,optional"`
	}
	reader := strings.NewReader(text)
	assert.Nil(t, UnmarshalYamlReader(reader, &v))
	assert.Equal(t, 1, v.Machine["node1"])
	assert.Equal(t, 2, v.Machine["node2"])
	assert.Equal(t, 3, v.Machine["node3"])
}

func TestWillUnmarshalYamlMapByte(t *testing.T) {
	text := `machine:
  node1: 1
  node2: 2
  node3: 3
`
	var v struct {
		Machine map[string]byte `json:"machine,optional"`
	}
	reader := strings.NewReader(text)
	assert.Nil(t, UnmarshalYamlReader(reader, &v))
	assert.Equal(t, byte(1), v.Machine["node1"])
	assert.Equal(t, byte(2), v.Machine["node2"])
	assert.Equal(t, byte(3), v.Machine["node3"])
}

func TestWillUnmarshalYamlMapRune(t *testing.T) {
	text := `machine:
  node1: 1
  node2: 2
  node3: 3
`
	var v struct {
		Machine map[string]rune `json:"machine,optional"`
	}
	reader := strings.NewReader(text)
	assert.Nil(t, UnmarshalYamlReader(reader, &v))
	assert.Equal(t, rune(1), v.Machine["node1"])
	assert.Equal(t, rune(2), v.Machine["node2"])
	assert.Equal(t, rune(3), v.Machine["node3"])
}

func TestWillUnmarshalYamlStringOfInt(t *testing.T) {
	text := `password: 123456`
	var v struct {
		Password string `json:"password"`
	}
	reader := strings.NewReader(text)
	assert.Error(t, UnmarshalYamlReader(reader, &v))
}

func TestWillUnmarshalYamlBadInput(t *testing.T) {
	var v struct {
		Any string
	}
	assert.Error(t, UnmarshalYamlBytes([]byte("':foo"), &v))
}

type badReaderWill struct{}

func (b *badReaderWill) Read(_ []byte) (n int, err error) {
	return 0, io.ErrLimitReached
}
