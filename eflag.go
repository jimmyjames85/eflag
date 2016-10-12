package eflag

import (
	"flag"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)


// This represents a pointer to any golang data type
type pValue struct {

	// TODO This is way too verbose and confusing...what to do
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
	rvalue *reflect.Value
}

// flag lib will allow boolean flags without values if this method returns true
// see flag.Value interface
func (p *pValue) IsBoolFlag() bool {
	if _, ok := p.rvalue.Interface().(*bool); !ok {
		return false
	}
	return true
}

// This is used to display default values
func (p *pValue) String() string {

	// rvalue nil check is because of flag.isZeroVal method
	if p.rvalue == nil || p.rvalue.Kind() != reflect.Ptr || p.rvalue.IsNil() {
		return ""
	}
	val := p.rvalue.Elem()

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

func (p *pValue) Set(s string) error {
	if p.rvalue.Kind() != reflect.Ptr {
		return fmt.Errorf("unsuported type %s", p.rvalue.Type().String())
	}

	switch p.rvalue.Interface().(type) {

	case *int:
		i, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("unable to parse %s: %v", s, err) // TODO repeated below
		}
		p.rvalue.Set(reflect.ValueOf(&i))

	case *string:
		p.rvalue.Set(reflect.ValueOf(&s))
	case *bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return fmt.Errorf("unable to parse %s: %v", s, err) // TODO repeated above (how to make DRY)
		}
		p.rvalue.Set(reflect.ValueOf(&b))
	default:
		return fmt.Errorf("unsuported type %s", p.rvalue.Type().String())
	}

	return nil
}

func StructVar(ps interface{}) {

	val := reflect.ValueOf(ps)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		log.Printf("StructVar is not a pointer") //todo return err? This would break convention...
		return
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		log.Printf("ps is not a pointer to a struct") //todo return err? This would break convention...
		return
	}

	for i := 0; i < val.NumField(); i++ {

		tag := val.Type().Field(i).Tag // tag is the metadata in the back-ticks
		field := val.Field(i)          // returns the fields of the struct

		tags := tag.Get("flag")

		//TODO I need/want to get a --flag option (two dashes)
		eflags := strings.Split(tags, ",")
		for _, eflag := range eflags {
			if eflag == "" {
				continue
			}
			desc := tag.Get("desc")
			if field.Kind() == reflect.Ptr {
				//back-ticks around dtype are used in flag.UnquoteUsage to determine type
				dtype := field.Type().String()
				desc = fmt.Sprintf("`%s` %s", dtype, desc)
			}
			if field.CanSet() {

				pfield := field.Addr().Interface()
				switch field.Kind() {
				case reflect.Ptr:
					pint := &pValue{rvalue: &field}
					flag.Var(pint, eflag, desc)
				case reflect.Int:
					flag.IntVar(pfield.(*int), eflag, int(field.Int()), desc)
				case reflect.Int64:
					flag.Int64Var(pfield.(*int64), eflag, field.Int(), desc)
				case reflect.String:
					flag.StringVar(pfield.(*string), eflag, field.String(), desc)
				case reflect.Bool:
					flag.BoolVar(pfield.(*bool), eflag, field.Bool(), desc)
				default:
					log.Fatal("unsuported: TODO") //TODO support it!
				}
			}

		}

	}
}
