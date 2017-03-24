package eflag

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type Arg struct {
	Names  []string
	Usage  string
	DefVal string
	Type   string
	//TODO Type??
}

// GetDeclaredArgs
// This will return a slice of Arg pointers that represent
// which flags have been flagged from either this library
// or the `flag` library
func GetDeclaredArgs() []*Arg {
	destValues := make(map[uintptr][]*flag.Flag)
	destTypes := make(map[uintptr]string)
	flag.VisitAll(func(f *flag.Flag) {

		val := reflect.ValueOf(f.Value)
		var addr uintptr

		if val.Type() == reflect.TypeOf(&rValue{}) {
			// rValue is just a reflect.Value
			val = (reflect.Value)(val.Elem().Interface().(rValue))
			addr = val.Addr().Pointer()
		} else {
			val = val.Elem()
			if val.CanAddr() {
				addr = val.Addr().Pointer()
			} else {
				addr = val.Pointer()
			}

		}

		destValues[addr] = append(destValues[addr], f)
		destTypes[addr] = val.Type().String()
	})

	args := make([]*Arg, 0)

	for addr, v := range destValues {
		var usage, defVal string
		var names []string

		for _, f := range v {
			names = append(names, f.Name)
			usage = f.Usage
			defVal = f.DefValue
		}
		sort.Strings(names)
		args = append(args, &Arg{Names: names, Usage: usage, DefVal: defVal, Type: destTypes[addr]})
	}

	sort.Sort((argSorter)(args))
	return args
}

type argSorter []*Arg

func (r argSorter) Len() int {
	return len(r)
}
func (r argSorter) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
func (r argSorter) Less(i, j int) bool {
	return r[i].Names[0] < r[j].Names[0]
}

func POSIXStyle() {

	fmt.Printf("Usage: %s\n\n", os.Args[0])
	args := GetDeclaredArgs()

	for _, arg := range args {
		fmt.Printf("\t")
		for _, name := range arg.Names {
			fmt.Printf("-%s ", name)
		}
		defVal := arg.DefVal
		if len(defVal) > 0 {
			defVal = fmt.Sprintf("(default: %s)", defVal)
		}
		fmt.Printf("\n\t\t%s %s [%s]\n\n", arg.Usage, defVal, arg.Type)
	}
}

type rValue reflect.Value

// flag lib will allow boolean flags without values if this method returns true
// see flag.Value interface
func (r *rValue) IsBoolFlag() bool {

	v := (*reflect.Value)(r)

	fmt.Printf("isbool? %s\n", v.Type())

	if _, ok := v.Interface().(*bool); !ok {

		if _, ok2 := v.Interface().(bool); !ok2 {

			return false
		}
	}
	return true
}

// Because the flag library only displays the default value in the help message, if the default
// value's .String() is different from the type's zero value .String() we have to
// TODO This explanation is way too verbose and confusing...what to do
// rvalue must be a ptr
//
// The flag library looks at the zero value to determine
// whether or not to display (default 'the default') in the help message
//
// Because the flag lib is unaware of type, it compares the element's .String()
// result to the zero value's .String() result
//
// If they are equal then no (default 'the default') is displayed in the message
// otherwise the default is displayed in the help message
//
// See flag.isZeroVal method

// This is used to display default values
func (r *rValue) String() string {

	v := (*reflect.Value)(r)

	// rvalue nil check is because of flag.isZeroVal method
	if v == nil || v.Kind() != reflect.Ptr || v.IsNil() {
		return ""
	}
	val := v.Elem()

	switch val.Interface().(type) {

	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val.Int())
	case string:
		return fmt.Sprintf("\"%s\"", val.String())
	case bool:
		return fmt.Sprintf("%t", val.Bool())
	default:
		return fmt.Sprintf("%v", val.String())
	}

}

func (r *rValue) Set(s string) error {

	v := (*reflect.Value)(r)

	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fmt.Errorf("_unable to parse %s: %v", s, err) // TODO how to make DRY (repeated below)
		}
		v.SetInt(i)
		return nil
	//case reflect.Bool:
	//	b, err := strconv.ParseBool(s)
	//	if err != nil {
	//		return fmt.Errorf("_unable to parse %s: %v", s, err) // TODO how to make DRY (repeated below)
	//	}
	//	v.SetBool(b)
	//	return nil
	case reflect.Ptr:
		switch v.Interface().(type) {

		case *int: //, *int8, *int16, *int32, *int64:
			i, err := strconv.Atoi(s) // strconv.ParseInt(s,10,64)
			if err != nil {
				return fmt.Errorf("_unable to parse %s: %v", s, err) // TODO how to make DRY (repeated below)
			}
			v.Set(reflect.ValueOf(&i))

		case *string:
			v.Set(reflect.ValueOf(&s))
		case *bool:
			fmt.Printf("Parsing %s\n", s)
			b, err := strconv.ParseBool(s)
			if err != nil {
				return fmt.Errorf("_unable to parse %s: %v", s, err) // TODO how to make DRY (repeated above)
			}
			v.Set(reflect.ValueOf(&b))
		default:
			return fmt.Errorf("__unsuported type %s", v.Type().String())
		}
	default:
		return fmt.Errorf(" ___ unsuported type %s", v.Type().String())
	}

	//if v.Kind() != reflect.Ptr {
	//
	//	reflect.Int8
	//	v.SetInt(12)
	//	return nil
	//	// then we can't set it
	//	return fmt.Errorf(" ___ unsuported type %s", v.Type().String())
	//}

	return nil
}

func StructVar(ps interface{}) {

	val := reflect.ValueOf(ps)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		log.Printf("Un-assignable: StructVar is not a pointer") //todo return err? This would break convention...
		return
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		log.Printf("ps is not a pointer to a struct") //todo return err? This would break convention...
		return
	}

	for i := 0; i < val.NumField(); i++ {

		tag := val.Type().Field(i).Tag // tag is the metadata in the back-ticks
		field := val.Field(i)          // returns the i-th field of the struct

		tags := tag.Get("flag")

		// TODO I need/want to get a --flag option (two dashes)
		// TODO also add support for no dashes...like "git help clone"
		eflags := strings.Split(tags, ",")

		for _, eflag := range eflags {
			if eflag == "" {
				continue
			}
			desc := tag.Get("desc")
			if field.Kind() == reflect.Ptr {
				dtype := field.Type().String()[1:] // this removes the *

				// back-ticks in the desc are used to denote the type
				// in the help message. See flag.UnquoteUsage
				desc = fmt.Sprintf("`%s` %s", dtype, desc)
			}

			if field.CanSet() {

				//fmt.Printf("CanSet Address: %v\n", field.Addr().Interface())

				switch field.Kind() {
				case reflect.Ptr:
					//ptrVal := &pValue{rvalue: &field}
					//flag.Var(ptrVal, eflag, desc)
					rVal := (rValue)(field)
					flag.Var(&rVal, eflag, desc)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					rVal := (rValue)(field)
					flag.Var(&rVal, eflag, desc)
				case reflect.String:
					rVal := (rValue)(field)
					flag.Var(&rVal, eflag, desc)
				//flag.StringVar(pfield.(*string), eflag, field.String(), desc)
				case reflect.Bool:
					rVal := (rValue)(field)
					flag.Var(&rVal, eflag, desc)
				//pfield := field.Addr().Interface()
				//flag.BoolVar(pfield.(*bool), eflag, field.Bool(), desc)
				default:
					//ptrVal := (rValue)(field)
					//flag.Var(&ptrVal, eflag, desc)
					log.Fatal("unsuported StructVar field: TODO") //TODO support it! ...
				}
			}
		}
	}
}
