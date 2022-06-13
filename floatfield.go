package xdommask

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/webability-go/wajaf"
	"github.com/webability-go/xamboo/cms/context"
)

type FloatField struct {
	*TextField
	DefaultValue  float64
	Min           float64
	Max           float64
	StatusTooLow  string
	StatusTooHigh string
}

func NewFloatField(name string) *FloatField {
	ff := &FloatField{
		TextField: NewTextField(name),
	}
	ff.TextType = TEXTTYPE_FLOAT
	return ff
}

func (f *FloatField) Compile() wajaf.NodeDef {

	//	t := wajaf.NewTextFieldElement(f.ID)
	t := f.TextField.Compile()
	t.SetAttribute("format", "^-{0,1}[0-9.,]{1,}([eE]-{0,1}[0-9,.]{1,}){0,}$")
	if f.Min != 0 {
		t.SetAttribute("min", fmt.Sprintf("%f", f.Min))
		t.AddMessage("statustoolow", f.StatusTooLow)
	}
	if f.Max != 0 {
		t.SetAttribute("max", fmt.Sprintf("%f", f.Max))
		t.AddMessage("statustoohigh", f.StatusTooHigh)
	}
	return t
}

func (f *FloatField) GetValue(ctx *context.Context, mode Mode) (interface{}, bool, error) {
	val, ignored, err := f.DataField.GetValue(ctx, mode)
	fmt.Println("Float value = ", val, ignored, err)
	if err != nil {
		return val, ignored, err
	}
	newval, err := f.ConvertValue(val)
	return newval, ignored, err
}

func (f *FloatField) ConvertValue(value interface{}) (interface{}, error) {
	fval, ok := value.(float64)
	if ok {
		return fval, nil
	}
	f32val, ok := value.(float32)
	if ok {
		return float64(f32val), nil
	}
	sval, ok := value.(string)
	if ok {
		return strconv.ParseFloat(sval, 64)
	}
	return value, errors.New("Cannot convert value")
}
