package xdommask

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/webability-go/wajaf"
	"github.com/webability-go/xamboo/cms/context"
)

type IntegerField struct {
	*TextField
	DefaultValue  int
	Min           int
	Max           int
	StatusTooLow  string
	StatusTooHigh string
}

func NewIntegerField(name string) *IntegerField {
	inf := &IntegerField{
		TextField: NewTextField(name),
	}
	inf.TextType = TEXTTYPE_INTEGER
	return inf
}

func (f *IntegerField) Compile() wajaf.NodeDef {

	t := f.TextField.Compile()
	t.SetAttribute("format", "^-{0,1}[0-9]{1,}$")
	if f.Min != 0 {
		t.SetAttribute("min", strconv.Itoa(f.Min))
		t.AddMessage("statustoolow", f.StatusTooLow)
	}
	if f.Max != 0 {
		t.SetAttribute("max", strconv.Itoa(f.Max))
		t.AddMessage("statustoohigh", f.StatusTooHigh)
	}
	return t
}

func (f *IntegerField) GetValue(ctx *context.Context, mode Mode) (interface{}, bool, error) {
	val, ignored, err := f.DataField.GetValue(ctx, mode)
	if err != nil {
		return val, ignored, err
	}
	fmt.Printf("GetValue: %v %T\n", val, val)
	newval, err := f.ConvertValue(val)
	fmt.Printf("New value: %v %T\n", newval, newval)
	if newval == nil && mode == INSERT {
		newval = f.DefaultValue
	}
	return newval, ignored, err
}

func (f *IntegerField) ConvertValue(value interface{}) (interface{}, error) {
	ival, ok := value.(int)
	if ok {
		return ival, nil
	}
	i8val, ok := value.(int8)
	if ok {
		return int(i8val), nil
	}
	i16val, ok := value.(int16)
	if ok {
		return int(i16val), nil
	}
	i32val, ok := value.(int32)
	if ok {
		return int(i32val), nil
	}
	i64val, ok := value.(int64)
	if ok {
		return int(i64val), nil
	}
	sval, ok := value.(string)
	if ok {
		return strconv.Atoi(sval)
	}
	return value, errors.New("Cannot convert value")
}
