package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/jimmyjames85/eflag"
)

func main() {

	//mike := "Mike 'The Default' Pike"

	ss := &simpleSettings{
		OneFishName: "Wally 'The Default' Walleye",
		TwoFishName: "Tod 'The Default' Cod",
		//ThreeFishName: &mike,  //This is do able, Right now I don't see the need for it though
	}

	eflag.StructVar(ss)
	flag.Usage = eflag.POSIXStyle
	flag.Parse()
	fmt.Println(ss)
}

type simpleSettings struct {
	OneFishName   string  `flag:"one,1,theonlyone" desc:"name of one fish"`
	TwoFishName   string  `flag:"2,two"`                             // MUST HAVE flag TAG, but no description is necessary
	ThreeFishName *string `flag:"3,three" desc:"name of three fish"` // pointers default to nil (unless specified on the CMD line)
	IsAFish       *bool   `flag:"f,isafish"`                         // Boolean types can be specified  -f OR -f=true
	WormCount     int64   f`flag:"w,wormcount" desc:"number of worms you have"`

	Boool      bool   `flag:"b,B" desc:"this is a boolean value"`
	hiddenFish string `flag:"h"` // Only exported fields can be parsed
}

func (ss *simpleSettings) String() string {
	b, err := json.Marshal(ss)
	if err != nil {
		return "null"
	}
	pp := &bytes.Buffer{}
	json.Indent(pp, []byte(string(b)), " ", "\t")
	return pp.String()
}
