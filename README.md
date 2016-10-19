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
	OneFishName   string  `flag:"one,1,theonlyone" desc:"name of one fish"`
	TwoFishName   string  `flag:"2,two"`                             // MUST HAVE flag TAG, but no description is necessary
	ThreeFishName *string `flag:"3,three" desc:"name of three fish"` // pointers default to nil (unless specified on the CMD line)
	IsAFish       *bool   `flag:"f,isafish"`                         // Boolean types can be specified  -f OR -f=true
	WormCount     int     `flag:"w,wormcount" desc:"number of worms you have"`

	hiddenFish    string `flag:"h"`                                  // Only exported fields can be parsed
}


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

 ```
 
## Output

```
 
 $ ./example -h
 Usage: ./example
 
         -1 -one -theonlyone
                 name of one fish (default: Wally 'The Default' Walleye) [flag.stringValue]
 
         -2 -two
                  (default: Tod 'The Default' Cod) [flag.stringValue]
 
         -3 -three
                 `string` name of three fish  [*string]
 
         -f -isafish
                 `bool`   [*bool]
 
         -w -wormcount
                 number of worms you have (default: 0) [flag.intValue]
 
 $ ./example
 {
         "OneFishName": "Wally 'The Default' Walleye",
         "TwoFishName": "Tod 'The Default' Cod",
         "ThreeFishName": null,
         "IsAFish": null,
         "WormCount": 0
 }
 
 $ ./example -two "Change Tod's Name"
 {
         "OneFishName": "Wally 'The Default' Walleye",
         "TwoFishName": "Change Tod's Name",
         "ThreeFishName": null,
         "IsAFish": null,
         "WormCount": 0
 }
 
 $ ./example -3 "values are allocated for pointer types"
 {
         "OneFishName": "Wally 'The Default' Walleye",
         "TwoFishName": "Tod 'The Default' Cod",
         "ThreeFishName": "values are allocated for pointer types",
         "IsAFish": null,
         "WormCount": 0
  }
 

 ```