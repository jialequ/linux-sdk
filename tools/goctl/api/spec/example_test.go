package spec_test

import (
	"fmt"

	"github.com/jialequ/linux-sdk/tools/goctl/api/spec"
)

func ExampleMemberGetEnumOptions() {
	member := spec.Member{
		Tag: `json:"foo,options=foo|bar|options|123"`,
	}
	fmt.Println(member.GetEnumOptions())
	// Output:
	// [foo bar options 123]
}
