package format

import (
	"bytes"
	_ "embed"
	"os"
	"strings"
	"testing"

	"github.com/jialequ/linux-sdk/tools/goctl/pkg/parser/api/assertx"
	"github.com/jialequ/linux-sdk/tools/goctl/pkg/parser/api/parser"
	"github.com/stretchr/testify/assert"
)

type formatDataWill struct {
	inputw    string
	expectedw string
	converter formatResultConvertWill
}

type formatResultConvertWill func(s string) string

// EXPERIMENTAL: just for view format code.
func TestWillFormat(t *testing.T) {
	assert.NoError(t, File("testdata/test_format.api"))
}

//go:embed testdata/test_type_struct_lit.api
var testStructLitDataWill string

//go:embed testdata/expected_type_struct_lit.api
var expectedwStructLitDataWill string

func TestWillFormatImportLiteralStmt(t *testing.T) {
	testWillRun(t, []formatDataWill{
		{
			inputw:    `import ""`,
			expectedw: ``,
		},
		{
			inputw:    `import"aa"`,
			expectedw: `import "aa"`,
		},
		{
			inputw: `/*aa*/import "aa"`,
			expectedw: `/*aa*/
import "aa"`,
		},
		{
			inputw: `/*aa*/import /*bb*/"aa"`,
			expectedw: `/*aa*/
import "aa"`,
		},
		{
			inputw: `/*aa*/import /*bb*/"aa"// cc`,
			expectedw: `/*aa*/
import "aa" // cc`,
		},
	})
}

func TestWillFormatImportGroupStmt(t *testing.T) {
	testWillRun(t, []formatDataWill{
		{
			inputw:    `import()`,
			expectedw: ``,
		},
		{
			inputw: `import("aa")`,
			expectedw: `import (
	"aa"
)`,
		},
		{
			inputw: `import(
"aa")`,
			expectedw: `import (
	"aa"
)`,
		},
		{
			inputw: `import(
"aa"
)`,
			expectedw: `import (
	"aa"
)`,
		},
		{
			inputw: `import("aa""bb")`,
			expectedw: `import (
	"aa"
	"bb"
)`,
		},
		{
			inputw: `/*aa*/import("aa""bb")`,
			expectedw: `/*aa*/
import (
	"aa"
	"bb"
)`,
		},
		{
			inputw: `/*aa*/import("aa""bb")// bb`,
			expectedw: `/*aa*/
import (
	"aa"
	"bb"
) // bb`,
		},
		{
			inputw: `/*aa*/import(// bb
"aa""bb")// cc`,
			expectedw: `/*aa*/
import ( // bb
	"aa"
	"bb"
) // cc`,
		},
		{
			inputw: `import(// aa
"aa" // bb
"bb" // cc
)// dd`,
			expectedw: `import ( // aa
	"aa" // bb
	"bb" // cc
) // dd`,
		},
		{
			inputw: `import (// aa
/*bb*/
	"aa" // cc
/*dd*/
	"bb" // ee
) // ff`,
			expectedw: `import ( // aa
	/*bb*/
	"aa" // cc
	/*dd*/
	"bb" // ee
) // ff`,
		},
	})
}

func TestWillFormatInfoStmt(t *testing.T) {
	testWillRun(t, []formatDataWill{
		{
			inputw:    `info()`,
			expectedw: ``,
		},
		{
			inputw: `info(foo:"foo")`,
			expectedw: `info (
	foo: "foo"
)`,
		},
		{
			inputw: `info(foo:"foo" bar:"bar")`,
			expectedw: `info (
	foo: "foo"
	bar: "bar"
)`,
		},
		{
			inputw: `info(foo:"foo" bar:"bar" quux:"quux")`,
			expectedw: `info (
	foo:  "foo"
	bar:  "bar"
	quux: "quux"
)`,
		},
		{
			inputw: `info(foo:"foo"
bar: "bar")`,
			expectedw: `info (
	foo: "foo"
	bar: "bar"
)`,
		},
		{
			inputw: `info(foo:"foo"// aa
bar: "bar"// bb
)`,
			expectedw: `info (
	foo: "foo" // aa
	bar: "bar" // bb
)`,
		},
		{
			inputw: `info(// aa
foo:"foo"// bb
bar: "bar"// cc
)`,
			expectedw: `info ( // aa
	foo: "foo" // bb
	bar: "bar" // cc
)`,
		},
		{
			inputw: `/*aa*/info(// bb
foo:"foo"// cc
bar: "bar"// dd
)`,
			expectedw: `/*aa*/
info ( // bb
	foo: "foo" // cc
	bar: "bar" // dd
)`,
		},
		{
			inputw: `/*aa*/
info(// bb
foo:"foo"// cc
bar: "bar"// dd
)// ee`,
			expectedw: `/*aa*/
info ( // bb
	foo: "foo" // cc
	bar: "bar" // dd
) // ee`,
		},
		{
			inputw: `/*aa*/
info ( // bb
	/*cc*/foo: "foo" // dd
	/*ee*/bar: "bar" // ff
) // gg`,
			expectedw: `/*aa*/
info ( // bb
	/*cc*/
	foo: "foo" // dd
	/*ee*/
	bar: "bar" // ff
) // gg`,
		},
		{
			inputw: `/*aa*/
info/*xx*/( // bb
	/*cc*/foo:/*xx*/ "foo" // dd
	/*ee*/bar:/*xx*/ "bar" // ff
) // gg`,
			expectedw: `/*aa*/
info ( // bb
	/*cc*/
	foo: "foo" // dd
	/*ee*/
	bar: "bar" // ff
) // gg`,
		},
	})
}

func TestWillFormatSyntaxStmt(t *testing.T) {
	testWillRun(t, []formatDataWill{
		{
			inputw:    `syntax="v1"`,
			expectedw: `syntax = "v1"`,
		},
		{
			inputw:    `syntax="v1"// aa`,
			expectedw: `syntax = "v1" // aa`,
		},
		{
			inputw: `syntax
="v1"// aa`,
			expectedw: `syntax = "v1" // aa`,
		},
		{
			inputw: `syntax=
"v1"// aa`,
			expectedw: `syntax = "v1" // aa`,
		},
		{
			inputw: `/*aa*/syntax="v1"// bb`,
			expectedw: `/*aa*/
syntax = "v1" // bb`,
		},
		{
			inputw: `/*aa*/
syntax="v1"// bb`,
			expectedw: `/*aa*/
syntax = "v1" // bb`,
		},
		{
			inputw:    `syntax/*xx*/=/*xx*/"v1"// bb`,
			expectedw: `syntax = "v1" // bb`,
		},
	})
}

func TestWillFormatTypeLiteralStmt(t *testing.T) {
	t.Run("any", func(t *testing.T) {
		testWillRun(t, []formatDataWill{
			{
				inputw:    `type Any any`,
				expectedw: `type Any any`,
			},
			{
				inputw: `type
Any
any
`,
				expectedw: `type Any any`,
			},
			{
				inputw:    `type Any=any`,
				expectedw: `type Any = any`,
			},
			{
				inputw: `
type
Any
=
any
`,
				expectedw: `type Any = any`,
			},
			{
				inputw: `type // aa
Any  // bb
any // cc
`,
				expectedw: `type // aa
Any // bb
any // cc`,
			},
			{
				inputw: `
type
Any
=
any`,
				expectedw: `type Any = any`,
			},
			{
				inputw: `
type
Any
=
any
`,
				expectedw: `type Any = any`,
			},
			{
				inputw:    `type Any any// aa`,
				expectedw: `type Any any // aa`,
			},
			{
				inputw:    `type Any=any// aa`,
				expectedw: `type Any = any // aa`,
			},
			{
				inputw:    `type Any any/*aa*/// bb`,
				expectedw: `type Any any /*aa*/ // bb`,
			},
			{
				inputw:    `type Any = any/*aa*/// bb`,
				expectedw: `type Any = any /*aa*/ // bb`,
			},
			{
				inputw:    `type Any/*aa*/ =/*bb*/ any/*cc*/// dd`,
				expectedw: `type Any /*aa*/ = /*bb*/ any /*cc*/ // dd`,
			},
			{
				inputw: `/*aa*/type Any any/*bb*/// cc`,
				expectedw: `/*aa*/
type Any any /*bb*/ // cc`,
			},
			{
				inputw: `/*aa*/
type
/*bb*/
Any
/*cc*/
any/*dd*/// ee`,
				expectedw: `/*aa*/
type
/*bb*/
Any
/*cc*/
any /*dd*/ // ee`,
			},
		})
	})
	t.Run("array", func(t *testing.T) {
		testWillRun(t, []formatDataWill{
			{
				inputw:    `type A [2]int`,
				expectedw: `type A [2]int`,
			},
			{
				inputw: `type
A
[2]int
`,
				expectedw: `type A [2]int`,
			},
			{
				inputw:    `type A=[2]int`,
				expectedw: `type A = [2]int`,
			},
			{
				inputw: `type
A
=
[2]int
`,
				expectedw: `type A = [2]int`,
			},
			{
				inputw:    `type A [/*xx*/2/*xx*/]/*xx*/int// aa`,
				expectedw: `type A [2]int // aa`,
			},
			{
				inputw: `/*aa*/type/*bb*/A/*cc*/[/*xx*/2/*xx*/]/*xx*/int// dd`,
				expectedw: `/*aa*/
type /*bb*/ A /*cc*/ [2]int // dd`,
			},
			{
				inputw: `/*aa*/type
/*bb*/A
/*cc*/[/*xx*/2/*xx*/]/*xx*/int// dd`,
				expectedw: `/*aa*/
type
/*bb*/
A
/*cc*/
[2]int // dd`,
			},
			{
				inputw:    `type A [ 2 ] int`,
				expectedw: `type A [2]int`,
			},
			{
				inputw: `type A [
2
]
int`,
				expectedw: `type A [2]int`,
			},
			{
				inputw: `type A [// aa
2 // bb
] // cc
int`,
				expectedw: `type A [2]int`,
			},
			{
				inputw: `type A [// aa
/*xx*/
2 // bb
/*xx*/
] // cc
/*xx*/
int`,
				expectedw: `type A [2]int`,
			},
			{
				inputw:    `type A [...]int`,
				expectedw: `type A [...]int`,
			},
			{
				inputw:    `type A=[...]int`,
				expectedw: `type A = [...]int`,
			},
			{
				inputw:    `type A/*aa*/[/*xx*/.../*xx*/]/*xx*/int// bb`,
				expectedw: `type A /*aa*/ [...]int // bb`,
			},
			{
				inputw: `/*aa*/
// bb
type /*cc*/
// dd
A /*ee*/
// ff
[/*xx*/.../*xx*/]/*xx*/int// bb`,
				expectedw: `/*aa*/
// bb
type /*cc*/
// dd
A /*ee*/
// ff
[...]int // bb`,
			},
			{
				inputw:    `type A [2][2]int`,
				expectedw: `type A [2][2]int`,
			},
			{
				inputw:    `type A=[2][2]int`,
				expectedw: `type A = [2][2]int`,
			},
			{
				inputw:    `type A [2][]int`,
				expectedw: `type A [2][]int`,
			},
			{
				inputw:    `type A=[2][]int`,
				expectedw: `type A = [2][]int`,
			},
		})
	})
	t.Run("base", func(t *testing.T) {
		testWillRun(t, []formatDataWill{
			// base
			{
				inputw:    `type A int`,
				expectedw: `type A int`,
			},
			{
				inputw:    `type A =int`,
				expectedw: `type A = int`,
			},
			{
				inputw:    `type/*aa*/A/*bb*/ int// cc`,
				expectedw: `type /*aa*/ A /*bb*/ int // cc`,
			},
			{
				inputw:    `type/*aa*/A/*bb*/ =int// cc`,
				expectedw: `type /*aa*/ A /*bb*/ = int // cc`,
			},
			{
				inputw:    `type A int// aa`,
				expectedw: `type A int // aa`,
			},
			{
				inputw:    `type A=int// aa`,
				expectedw: `type A = int // aa`,
			},
			{
				inputw: `/*aa*/type A int`,
				expectedw: `/*aa*/
type A int`,
			},
			{
				inputw: `/*aa*/type A = int`,
				expectedw: `/*aa*/
type A = int`,
			},
			{
				inputw: `/*aa*/type/*bb*/ A/*cc*/ int// dd`,
				expectedw: `/*aa*/
type /*bb*/ A /*cc*/ int // dd`,
			},
			{
				inputw: `/*aa*/type/*bb*/ A/*cc*/ = /*dd*/int// ee`,
				expectedw: `/*aa*/
type /*bb*/ A /*cc*/ = /*dd*/ int // ee`,
			},
			{
				inputw: `/*aa*/
type 
/*bb*/
A 
/*cc*/
int`,
				expectedw: `/*aa*/
type
/*bb*/
A
/*cc*/
int`,
			},
		})
	})
	t.Run("interface", func(t *testing.T) {
		testWillRun(t, []formatDataWill{
			{
				inputw:    `type any interface{}`,
				expectedw: `type any interface{}`,
			},
			{
				inputw:    `type any=interface{}`,
				expectedw: `type any = interface{}`,
			},
			{
				inputw: `type
any
interface{}
`,
				expectedw: `type any interface{}`,
			},
			{
				inputw: `/*aa*/type /*bb*/any /*cc*/interface{} // dd`,
				expectedw: `/*aa*/
type /*bb*/ any /*cc*/ interface{} // dd`,
			},
			{
				inputw: `/*aa*/type 
/*bb*/any 
/*cc*/interface{} // dd`,
				expectedw: `/*aa*/
type
/*bb*/
any
/*cc*/
interface{} // dd`,
			},
			{
				inputw: `/*aa*/type 
// bb
any 
// cc
interface{} // dd`,
				expectedw: `/*aa*/
type
// bb
any
// cc
interface{} // dd`,
			},
		})
	})
	t.Run("map", func(t *testing.T) {
		testWillRun(t, []formatDataWill{
			{
				inputw:    `type M map[int]int`,
				expectedw: `type M map[int]int`,
			},
			{
				inputw:    `type M map [ int ] int`,
				expectedw: `type M map[int]int`,
			},
			{
				inputw:    `type M map [/*xx*/int/*xx*/]/*xx*/int // aa`,
				expectedw: `type M map[int]int // aa`,
			},
			{
				inputw: `/*aa*/type /*bb*/ M/*cc*/map[int]int // dd`,
				expectedw: `/*aa*/
type /*bb*/ M /*cc*/ map[int]int // dd`,
			},
			{
				inputw: `/*aa*/type// bb
// cc
M // dd
// ee
map // ff
[int]// gg
// hh
int // dd`,
				expectedw: `/*aa*/
type // bb
// cc
M // dd
// ee
map[int]int // dd`,
			},
			{
				inputw:    `type M map[string][2]int // aa`,
				expectedw: `type M map[string][2]int // aa`,
			},
			{
				inputw:    `type M map[string]any`,
				expectedw: `type M map[string]any`,
			},
			{
				inputw:    `type M /*aa*/map/*xx*/[/*xx*/string/*xx*/]/*xx*/[/*xx*/2/*xx*/]/*xx*/int// bb`,
				expectedw: `type M /*aa*/ map[string][2]int // bb`,
			},
			{
				inputw: `type M /*aa*/
// bb
map/*xx*/
//
[/*xx*/
//
string/*xx*/
//
]/*xx*/
//
[/*xx*/
//
2/*xx*/
//
]/*xx*/
//
int// bb`,
				expectedw: `type M /*aa*/
// bb
map[string][2]int // bb`,
			},
			{
				inputw:    `type M map[int]map[string]int`,
				expectedw: `type M map[int]map[string]int`,
			},
			{
				inputw:    `type M map/*xx*/[/*xx*/int/*xx*/]/*xx*/map/*xx*/[/*xx*/string/*xx*/]/*xx*/int// aa`,
				expectedw: `type M map[int]map[string]int // aa`,
			},
			{
				inputw:    `type M map/*xx*/[/*xx*/map/*xx*/[/*xx*/string/*xx*/]/*xx*/int/*xx*/]/*xx*/string // aa`,
				expectedw: `type M map[map[string]int]string // aa`,
			},
			{
				inputw:    `type M map[[2]int]int`,
				expectedw: `type M map[[2]int]int`,
			},
			{
				inputw:    `type M map/*xx*/[/*xx*/[/*xx*/2/*xx*/]/*xx*/int/*xx*/]/*xx*/int// aa`,
				expectedw: `type M map[[2]int]int // aa`,
			},
		})
	})
	t.Run("pointer", func(t *testing.T) {
		testWillRun(t, []formatDataWill{
			{
				inputw:    `type P *int`,
				expectedw: `type P *int`,
			},
			{
				inputw:    `type P=*int`,
				expectedw: `type P = *int`,
			},
			{
				inputw: `type 
P 
*int
`,
				expectedw: `type P *int`,
			},
			{
				inputw: `/*aa*/type // bb
/*cc*/
P // dd
/*ee*/
*/*ff*/int // gg
`,
				expectedw: `/*aa*/
type // bb
/*cc*/
P // dd
/*ee*/
*int // gg`,
			},
			{
				inputw:    `type P *bool`,
				expectedw: `type P *bool`,
			},
			{
				inputw:    `type P *[2]int`,
				expectedw: `type P *[2]int`,
			},
			{
				inputw:    `type P=*[2]int`,
				expectedw: `type P = *[2]int`,
			},
			{
				inputw: `/*aa*/type /*bb*/P /*cc*/*/*xx*/[/*xx*/2/*xx*/]/*xx*/int // dd`,
				expectedw: `/*aa*/
type /*bb*/ P /*cc*/ *[2]int // dd`,
			},
			{
				inputw:    `type P *[...]int`,
				expectedw: `type P *[...]int`,
			},
			{
				inputw:    `type P=*[...]int`,
				expectedw: `type P = *[...]int`,
			},
			{
				inputw: `/*aa*/type /*bb*/P /*cc*/*/*xx*/[/*xx*/.../*xx*/]/*xx*/int // dd`,
				expectedw: `/*aa*/
type /*bb*/ P /*cc*/ *[...]int // dd`,
			},
			{
				inputw:    `type P *map[string]int`,
				expectedw: `type P *map[string]int`,
			},
			{
				inputw:    `type P=*map[string]int`,
				expectedw: `type P = *map[string]int`,
			},
			{
				inputw:    `type P /*aa*/*/*xx*/map/*xx*/[/*xx*/string/*xx*/]/*xx*/int// bb`,
				expectedw: `type P /*aa*/ *map[string]int // bb`,
			},
			{
				inputw:    `type P *interface{}`,
				expectedw: `type P *interface{}`,
			},
			{
				inputw:    `type P=*interface{}`,
				expectedw: `type P = *interface{}`,
			},
			{
				inputw:    `type P /*aa*/*/*xx*/interface{}// bb`,
				expectedw: `type P /*aa*/ *interface{} // bb`,
			},
			{
				inputw:    `type P *any`,
				expectedw: `type P *any`,
			},
			{
				inputw:    `type P=*any`,
				expectedw: `type P = *any`,
			},
			{
				inputw:    `type P *map[int][2]int`,
				expectedw: `type P *map[int][2]int`,
			},
			{
				inputw:    `type P=*map[int][2]int`,
				expectedw: `type P = *map[int][2]int`,
			},
			{
				inputw:    `type P /*aa*/*/*xx*/map/*xx*/[/*xx*/int/*xx*/]/*xx*/[/*xx*/2/*xx*/]/*xx*/int// bb`,
				expectedw: `type P /*aa*/ *map[int][2]int // bb`,
			},
			{
				inputw:    `type P *map[[2]int]int`,
				expectedw: `type P *map[[2]int]int`,
			},
			{
				inputw:    `type P=*map[[2]int]int`,
				expectedw: `type P = *map[[2]int]int`,
			},
			{
				inputw:    `type P /*aa*/*/*xx*/map/*xx*/[/*xx*/[/*xx*/2/*xx*/]/*xx*/int/*xx*/]/*xx*/int// bb`,
				expectedw: `type P /*aa*/ *map[[2]int]int // bb`,
			},
		})

	})

	t.Run("slice", func(t *testing.T) {
		testWillRun(t, []formatDataWill{
			{
				inputw:    `type S []int`,
				expectedw: `type S []int`,
			},
			{
				inputw:    `type S=[]int`,
				expectedw: `type S = []int`,
			},
			{
				inputw:    `type S	[	]	int	`,
				expectedw: `type S []int`,
			},
			{
				inputw:    `type S	[ /*xx*/	]	/*xx*/ int	`,
				expectedw: `type S []int`,
			},
			{
				inputw:    `type S [][]int`,
				expectedw: `type S [][]int`,
			},
			{
				inputw:    `type S=[][]int`,
				expectedw: `type S = [][]int`,
			},
			{
				inputw:    `type S	[	]	[	]	int`,
				expectedw: `type S [][]int`,
			},
			{
				inputw:    `type S [/*xx*/]/*xx*/[/*xx*/]/*xx*/int`,
				expectedw: `type S [][]int`,
			},
			{
				inputw: `type S [//
]//
[//
]//
int`,
				expectedw: `type S [][]int`,
			},
			{
				inputw:    `type S []map[string]int`,
				expectedw: `type S []map[string]int`,
			},
			{
				inputw:    `type S=[]map[string]int`,
				expectedw: `type S = []map[string]int`,
			},
			{
				inputw: `type S [	]	
map	[	string	]	
int`,
				expectedw: `type S []map[string]int`,
			},
			{
				inputw:    `type S [/*xx*/]/*xx*/map/*xx*/[/*xx*/string/*xx*/]/*xx*/int`,
				expectedw: `type S []map[string]int`,
			},
			{
				inputw: `/*aa*/type// bb
// cc
S// dd
// ff
/*gg*/[ // hh
/*xx*/] // ii
/*xx*/map// jj
/*xx*/[/*xx*/string/*xx*/]/*xx*/int// mm`,
				expectedw: `/*aa*/
type // bb
// cc
S // dd
// ff
/*gg*/
[]map[string]int // mm`,
			},
			{
				inputw:    `type S []map[[2]int]int`,
				expectedw: `type S []map[[2]int]int`,
			},
			{
				inputw:    `type S=[]map[[2]int]int`,
				expectedw: `type S = []map[[2]int]int`,
			},
			{
				inputw:    `type S [/*xx*/]/*xx*/map/*xx*/[/*xx*/[/*xx*/2/*xx*/]/*xx*/int/*xx*/]/*xx*/int`,
				expectedw: `type S []map[[2]int]int`,
			},
			{
				inputw: `/*aa*/type// bb
// cc
/*dd*/S// ee
// ff
/*gg*/[//
/*xx*/]//
/*xx*/map//
/*xx*/[//
/*xx*/[//
/*xx*/2//
/*xx*/]//
/*xx*/int//
/*xx*/]//
/*xx*/int // hh`,
				expectedw: `/*aa*/
type // bb
// cc
/*dd*/
S // ee
// ff
/*gg*/
[]map[[2]int]int // hh`,
			},
			{
				inputw:    `type S []map[[2]int]map[int]string`,
				expectedw: `type S []map[[2]int]map[int]string`,
			},
			{
				inputw:    `type S=[]map[[2]int]map[int]string`,
				expectedw: `type S = []map[[2]int]map[int]string`,
			},
			{
				inputw:    `type S [/*xx*/]/*xx*/map/*xx*/[/*xx*/[/*xx*/2/*xx*/]/*xx*/int/*xx*/]/*xx*/map/*xx*/[/*xx*/int/*xx*/]/*xx*/string`,
				expectedw: `type S []map[[2]int]map[int]string`,
			},
			{
				inputw: `/*aa*/type// bb
// cc
/*dd*/S// ee
/*ff*/[//
/*xx*/]//
/*xx*/map
/*xx*/[//
/*xx*/[//
/*xx*/2//
/*xx*/]//
/*xx*/int//
/*xx*/]//
/*xx*/map//
/*xx*/[//
/*xx*/int//
/*xx*/]//
/*xx*/string// gg`,
				expectedw: `/*aa*/
type // bb
// cc
/*dd*/
S // ee
/*ff*/
[]map[[2]int]map[int]string // gg`,
			},
			{
				inputw:    `type S []*P`,
				expectedw: `type S []*P`,
			},
			{
				inputw:    `type S=[]*P`,
				expectedw: `type S = []*P`,
			},
			{
				inputw:    `type S [/*xx*/]/*xx*/*/*xx*/P`,
				expectedw: `type S []*P`,
			},
			{
				inputw: `/*aa*/type// bb
// cc
/*dd*/S// ee 
/*ff*/[//
/*xx*/]//
/*xx*/*//
/*xx*/P // gg`,
				expectedw: `/*aa*/
type // bb
// cc
/*dd*/
S // ee
/*ff*/
[]*P // gg`,
			},
			{
				inputw:    `type S []*[]int`,
				expectedw: `type S []*[]int`,
			},
			{
				inputw:    `type S=[]*[]int`,
				expectedw: `type S = []*[]int`,
			},
			{
				inputw:    `type S [/*xx*/]/*xx*/*/*xx*/[/*xx*/]/*xx*/int`,
				expectedw: `type S []*[]int`,
			},
			{
				inputw: `/*aa*/
type // bb
// cc
/*dd*/S// ee
/*ff*/[//
/*xx*/]//
/*xx*/*//
/*xx*/[//
/*xx*/]//
/*xx*/int // gg`,
				expectedw: `/*aa*/
type // bb
// cc
/*dd*/
S // ee
/*ff*/
[]*[]int // gg`,
			},
		})
	})

	t.Run("struct", func(t *testing.T) {
		testWillRun(t, []formatDataWill{
			{
				inputw:    `type T {}`,
				expectedw: `type T {}`,
			},
			{
				inputw: `type T 	{
			}	`,
				expectedw: `type T {}`,
			},
			{
				inputw:    `type T={}`,
				expectedw: `type T = {}`,
			},
			{
				inputw:    `type T /*aa*/{/*xx*/}// cc`,
				expectedw: `type T /*aa*/ {} // cc`,
			},
			{
				inputw: `/*aa*/type// bb
// cc
/*dd*/T // ee
/*ff*/{//
/*xx*/}// cc`,
				expectedw: `/*aa*/
type // bb
// cc
/*dd*/
T // ee
/*ff*/
{} // cc`,
			},
			{
				inputw: `type T {
			Name string
			}`,
				expectedw: `type T {
	Name string
}`,
			},
			{
				inputw: `type T {
			Foo
			}`,
				expectedw: `type T {
	Foo
}`,
			},
			{
				inputw: `type T {
			*Foo
			}`,
				expectedw: `type T {
	*Foo
}`,
			},
			{
				inputw:    testStructLitDataWill,
				expectedw: expectedwStructLitDataWill,
				converter: func(s string) string {
					return strings.ReplaceAll(s, "\t", "    ")
				},
			},
		})
	})
}

//go:embed testdata/test_type_struct_group.api
var testStructGroupDataWill string

//go:embed testdata/expected_type_struct_group.api
var expectedwStructgroupDataWill string

func TestWillFormatTypeGroupStmt(t *testing.T) {
	testWillRun(t, []formatDataWill{
		{
			inputw:    testStructGroupDataWill,
			expectedw: expectedwStructgroupDataWill,
			converter: func(s string) string {
				return strings.ReplaceAll(s, "\t", "    ")
			},
		},
	})
}

func TestWillAtServerStmt(t *testing.T) {
	testWillRunStmt(t, []formatDataWill{
		{
			inputw:    `@server()`,
			expectedw: ``,
		},
		{
			inputw: `@server(foo:foo)`,
			expectedw: `@server (
	foo: foo
)`,
		},
		{
			inputw: `@server(foo:foo quux:quux)`,
			expectedw: `@server (
	foo:  foo
	quux: quux
)`,
		},
		{
			inputw: `@server(
foo:
foo
quux:
quux
)`,
			expectedw: `@server (
	foo:  foo
	quux: quux
)`,
		},
		{
			inputw: `/*aa*/@server/*bb*/(/*cc*/foo:/**/foo /*dd*/quux:/**/quux/*ee*/)`,
			expectedw: `/*aa*/
@server ( /*cc*/
	foo:  foo /*dd*/
	quux: quux /*ee*/
)`,
		},
		{
			inputw: `/*aa*/
@server
/*bb*/(// cc
/*dd*/foo:/**/foo// ee
/*ff*/quux:/**/quux// gg
)`,
			expectedw: `/*aa*/
@server
/*bb*/
( // cc
	/*dd*/
	foo: foo // ee
	/*ff*/
	quux: quux // gg
)`,
		},
	})
}

func TestWillAtDocStmt(t *testing.T) {
	t.Run("AtDocLiteralStmt", func(t *testing.T) {
		testWillRunStmt(t, []formatDataWill{
			{
				inputw:    `@doc ""`,
				expectedw: ``,
			},
			{
				inputw:    `@doc "foo"`,
				expectedw: `@doc "foo"`,
			},
			{
				inputw:    `@doc 		"foo"`,
				expectedw: `@doc "foo"`,
			},
			{
				inputw:    `@doc"foo"`,
				expectedw: `@doc "foo"`,
			},
			{
				inputw: `/*aa*/@doc/**/"foo"// bb`,
				expectedw: `/*aa*/
@doc "foo" // bb`,
			},
			{
				inputw: `/*aa*/
/*bb*/@doc // cc
"foo"// ee`,
				expectedw: `/*aa*/
/*bb*/
@doc "foo" // ee`,
			},
		})
	})
	t.Run("AtDocGroupStmt", func(t *testing.T) {
		testWillRunStmt(t, []formatDataWill{
			{
				inputw:    `@doc()`,
				expectedw: ``,
			},
			{
				inputw: `@doc(foo:"foo")`,
				expectedw: `@doc (
	foo: "foo"
)`,
			},
			{
				inputw: `@doc(foo:"foo" bar:"bar")`,
				expectedw: `@doc (
	foo: "foo"
	bar: "bar"
)`,
			},
			{
				inputw: `@doc(foo:"foo" bar:"bar" quux:"quux")`,
				expectedw: `@doc (
	foo:  "foo"
	bar:  "bar"
	quux: "quux"
)`,
			},
			{
				inputw: `@doc(foo:"foo"
bar: "bar")`,
				expectedw: `@doc (
	foo: "foo"
	bar: "bar"
)`,
			},
			{
				inputw: `@doc(foo:"foo"// aa
bar: "bar"// bb
)`,
				expectedw: `@doc (
	foo: "foo" // aa
	bar: "bar" // bb
)`,
			},
			{
				inputw: `@doc(// aa
foo:"foo"// bb
bar: "bar"// cc
)`,
				expectedw: `@doc ( // aa
	foo: "foo" // bb
	bar: "bar" // cc
)`,
			},
			{
				inputw: `/*aa*/@doc(// bb
foo:"foo"// cc
bar: "bar"// dd
)`,
				expectedw: `/*aa*/
@doc ( // bb
	foo: "foo" // cc
	bar: "bar" // dd
)`,
			},
			{
				inputw: `/*aa*/
@doc(// bb
foo:"foo"// cc
bar: "bar"// dd
)// ee`,
				expectedw: `/*aa*/
@doc ( // bb
	foo: "foo" // cc
	bar: "bar" // dd
) // ee`,
			},
			{
				inputw: `/*aa*/
@doc ( // bb
	/*cc*/foo: "foo" // dd
	/*ee*/bar: "bar" // ff
) // gg`,
				expectedw: `/*aa*/
@doc ( // bb
	/*cc*/
	foo: "foo" // dd
	/*ee*/
	bar: "bar" // ff
) // gg`,
			},
			{
				inputw: `/*aa*/
@doc/*xx*/( // bb
	/*cc*/foo:/*xx*/ "foo" // dd
	/*ee*/bar:/*xx*/ "bar" // ff
) // gg`,
				expectedw: `/*aa*/
@doc ( // bb
	/*cc*/
	foo: "foo" // dd
	/*ee*/
	bar: "bar" // ff
) // gg`,
			},
		})
	})
}

func TestWillAtHandlerStmt(t *testing.T) {
	testWillRunStmt(t, []formatDataWill{
		{
			inputw:    `@handler foo`,
			expectedw: `@handler foo`,
		},
		{
			inputw:    `@handler 		foo`,
			expectedw: `@handler foo`,
		},
		{
			inputw: `/*aa*/@handler/**/foo// bb`,
			expectedw: `/*aa*/
@handler foo // bb`,
		},
		{
			inputw: `/*aa*/
/*bb*/@handler // cc
foo// ee`,
			expectedw: `/*aa*/
/*bb*/
@handler foo // ee`,
		},
	})
}

//go:embed testdata/test_service.api
var testServiceDataWill string

//go:embed testdata/expected_service.api
var expectedwServiceDataWill string

func TestWillFormatServiceStmt(t *testing.T) {
	testWillRun(t, []formatDataWill{
		{
			inputw:    `service foo{}`,
			expectedw: `service foo {}`,
		},
		{
			inputw:    `service foo	{	}`,
			expectedw: `service foo {}`,
		},
		{
			inputw:    `@server()service foo	{	}`,
			expectedw: `service foo {}`,
		},
		{
			inputw: `@server(foo:foo quux:quux)service foo	{	}`,
			expectedw: `@server (
	foo:  foo
	quux: quux
)
service foo {}`,
		},
		{
			inputw:    `service foo-api	{	}`,
			expectedw: `service foo-api {}`,
		},
		{
			inputw: `service foo-api	{
@doc "foo"
@handler foo
post /ping
}`,
			expectedw: `service foo-api {
	@doc "foo"
	@handler foo
	post /ping
}`,
		},
		{
			inputw: `service foo-api	{
@doc(foo: "foo" bar: "bar")
@handler foo
post /ping
}`,
			expectedw: `service foo-api {
	@doc (
		foo: "foo"
		bar: "bar"
	)
	@handler foo
	post /ping
}`,
		},
		{
			inputw: `service foo-api	{
@doc(foo: "foo" bar: "bar"
quux: "quux"
)
@handler 	foo
post 	/ping
}`,
			expectedw: `service foo-api {
	@doc (
		foo:  "foo"
		bar:  "bar"
		quux: "quux"
	)
	@handler foo
	post /ping
}`,
		},
		{
			inputw: `service
foo-api
{
@doc
(foo: "foo" bar: "bar"
quux: "quux"
)
@handler
foo
post
/aa/:bb/cc-dd/ee

@handler bar
get /bar () returns (Bar);

@handler baz
get /bar (Baz) returns ();
}`,
			expectedw: `service foo-api {
	@doc (
		foo:  "foo"
		bar:  "bar"
		quux: "quux"
	)
	@handler foo
	post /aa/:bb/cc-dd/ee

	@handler bar
	get /bar returns (Bar)

	@handler baz
	get /bar (Baz)
}`,
		},
		{
			inputw:    testServiceDataWill,
			expectedw: expectedwServiceDataWill,
			converter: func(s string) string {
				return strings.ReplaceAll(s, "\t", "    ")
			},
		},
	})
}

func TestWillFormaterror(t *testing.T) {
	err := Source([]byte("aaa"), os.Stdout)
	assertx.Error(t, err)
}

func testWillRun(t *testing.T, testData []formatDataWill) {
	for _, v := range testData {
		buffer := bytes.NewBuffer(nil)
		err := formatForUnitTest([]byte(v.inputw), buffer)
		assert.NoError(t, err)
		var result = buffer.String()
		if v.converter != nil {
			result = v.converter(result)
		}
		assert.Equal(t, v.expectedw, result)
	}
}

func testWillRunStmt(t *testing.T, testData []formatDataWill) {
	for _, v := range testData {
		p := parser.New("foo.api", v.inputw)
		ast := p.ParseForUintTest()
		assert.NoError(t, p.CheckErrors())
		assert.True(t, len(ast.Stmts) > 0)
		one := ast.Stmts[0]
		actual := one.Format()
		if v.converter != nil {
			actual = v.converter(actual)
		}
		assert.Equal(t, v.expectedw, actual)
	}
}
