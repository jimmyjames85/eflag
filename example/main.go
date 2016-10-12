package main

import (
	"encoding/json"
	"flag"

	"github.com/jimmyjames85/eflag"
	"fmt"
)

func main() {

	ss := simpleSettings{
		OneFish: "This is a default",
	}

	eflag.StructVar(&ss)
	flag.Parse()

	fmt.Printf(JSON(ss))
}

type simpleSettings struct {

	// Default values are set to the type's zero value
	OneFish    string `flag:"1,one" desc:"name of one fish"`
	OneFishAge int    `flag:"1a,oneage" desc:"age of one fish"`

	// These will be nil pointers, unless specified on the command line
	TwoFish    *string `flag:"2,two" desc:"name of two fish"`
	TwoFishAge *int    `flag:"2a,twoage" desc:"age of two fish"`

	// MUST HAVE flag TAG, but no description is necessary
	ThreeFish string `flag:"3"`

	// Boolean types can be specified with no value -f
	IsAFish *bool `flag:"f"`

	// Only exported fields can be parsed
	hiddenFish string `flag:"h"`
}

func JSON(o interface{}) string {
	b, err := json.Marshal(o)
	if err != nil {
		return "null"
	}
	return string(b)
}

func ptrToStr(s string) *string {
	return &s
}
