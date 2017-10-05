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
	"time"
)

type Arg struct {
	Names  []string
	Usage  string
	DefVal string
	Type   string
	//TODO Type??
}

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
			addr = val.UnsafeAddr()

		} else {
			val = val.Elem()
			if val.CanAddr() {
				addr = val.UnsafeAddr()
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

func (r argSorter) Len() int           { return len(r) }
func (r argSorter) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r argSorter) Less(i, j int) bool { return r[i].Names[0] < r[j].Names[0] }

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
	if _, ok := v.Interface().(*bool); !ok {
		return false
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

	case int:
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

	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("_unsuported type %s", v.Type().String())
	}

	switch v.Interface().(type) {

	case *int:
		i, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("_unable to parse %s: %v", s, err) // TODO how to make DRY (repeated below)
		}
		v.Set(reflect.ValueOf(&i))

	case *string:
		v.Set(reflect.ValueOf(&s))
	case *bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return fmt.Errorf("_unable to parse %s: %v", s, err) // TODO how to make DRY (repeated above)
		}
		v.Set(reflect.ValueOf(&b))
	default:
		return fmt.Errorf("_unsuported type %s", v.Type().String())
	}

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
				dtype := field.Type().String()[1:] //remove the *

				// back-ticks in the desc are used to denote the type
				// in the help message. See flag.UnquoteUsage
				desc = fmt.Sprintf("`%s` %s", dtype, desc)
			}

			if field.CanSet() {

				//fmt.Printf("CanSet Address: %v\n", field.Addr().Interface())
				pfield := field.Addr().Interface()

				t, ok := pfield.(*time.Duration)
				if ok {
					flag.DurationVar(t, eflag, *t, desc)
					continue
				}

				switch field.Kind() {
				case reflect.Ptr:
					//ptrVal := &pValue{rvalue: &field}
					//flag.Var(ptrVal, eflag, desc)
					ptrVal := (rValue)(field)
					flag.Var(&ptrVal, eflag, desc)
				case reflect.Int:
					flag.IntVar(pfield.(*int), eflag, int(field.Int()), desc)
				case reflect.Int64:
					flag.Int64Var(pfield.(*int64), eflag, field.Int(), desc)
				case reflect.String:
					flag.StringVar(pfield.(*string), eflag, field.String(), desc)
				case reflect.Bool:
					flag.BoolVar(pfield.(*bool), eflag, field.Bool(), desc)
				default:
					log.Fatal("unsuported: TODO") //TODO support it! ...
				}
			}
		}
	}
}
