package main

import (
	"fmt"
	"github.com/jimmyjames85/eflag"
	"flag"
)

type mysettings struct {
	NoMetadata          int     //TODO This is not parsed, not sure if I want it to be unless metadata is provided
	nonExportedField    int     //This should/can not be parsed
	SomeIntegerFlag     int     `flag:"s,someint" desc:"an integer for SomeIntegerFlag"`
	SomePointerToInt    *int    `flag:"spi,somepointer" desc:"if this is not provided the setting defaults to <nil> (or default if provided)"`
	SomeStringFlag      string  `flag:"ss,somestring" desc:"This is a string that will default to the zero value "" (or default if provided)"`
	SomePointerToString *string `flag:"sps,somepointertostring" desc:"This is a string that will default to <nil> (or default if provided)"`
	SomeBool            bool    `flag:"sb,somebool" desc:"a boolean"`
	SomePointerToBool   *bool   `flag:"spb,somepointertobool" desc:"a pointer to a boolean"`
	BigInt              int64   `flag:"bi,bigint" desc:"a big int64"`
	BigIntPointer       *int64  `flag:"pbi,bigintpointer" desc:"a pointer to a big int64"`
}

func main2() {

	// Set the defaults here
	s := mysettings{
		SomeIntegerFlag:     34, //Set defaults here
		SomePointerToString: ptrToStr("default pointer to string"),
		SomeBool:            true,
		BigInt:              314159265,
	}
	ss := simpleSettings{}

	beforeS := fmt.Sprintf(JSON(s))
	beforeSS := fmt.Sprintf(JSON(ss))

	eflag.StructVar(&s)
	eflag.StructVar(&ss)

	flag.Parse()
	fmt.Println(beforeS)
	fmt.Println(JSON(s))
	fmt.Println(beforeSS)
	fmt.Println(JSON(ss))
}
