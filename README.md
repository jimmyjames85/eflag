# eflag

The `eflag` package lets you specify your flags in a single 
struct through metadata tags. The goal is to seemlesy work 
with the `flag` standard package. You can have multiple flags 
defined for the same variable. Struct fields of type pointer will
default to `nil` if not specified. 


###### TODO: Eventually I want to auto format the help message.

## Example


 ```

type simpleSettings struct {
    
	// Default values are set to the type's zero value
	OneFish    string   `flag:"1,one" desc:"name of one fish"`
	OneFishAge int      `flag:"1a,oneage" desc:"age of one fish"`

	// These will be nil pointers, unless specified on the command line
	TwoFish    *string  `flag:"2,two" desc:"name of two fish"`
	TwoFishAge *int     `flag:"2a,twoage" desc:"age of two fish"`

	// MUST HAVE flag TAG, but no description is necessary
	ThreeFish  string   `flag:"3"`
	
	// Only exported fields can be parsed
	hiddenFish  string   `flag:"3"`	
}

func main() {

	ss := simpleSettings{
		OneFish:              "This is a default",
	}

	eflag.StructVar(&ss)
	flag.Parse()


    /* ... */
}

 ```