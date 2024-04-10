package xdommask

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/webability-go/wajaf"
	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/cms/context"
)

const (
	SEARCHABLETEXTTYPE_TEXT    = "text"
	SEARCHABLETEXTTYPE_MASKED  = "masked"
	SEARCHABLETEXTTYPE_INTEGER = "integer"
	SEARCHABLETEXTTYPE_FLOAT   = "float"
	SEARCHABLETEXTTYPE_EMAIL   = "email"
)

type SearchableTextField struct {
	*DataField
	TextType string

	Format   string
	FormatJS string

	MinLength int
	MaxLength int
	MinWords  int
	MaxWords  int

	StatusBadFormat    string
	StatusTooShort     string
	StatusTooLong      string
	StatusTooFewWords  string
	StatusTooManyWords string

	KeyUpJS string
	FocusJS string
	BlurJS  string

	Check func(ctx *context.Context, mode Mode, value interface{}) error
	Calc  func(ctx *context.Context, mode Mode, rec *xdominion.XRecord) (interface{}, error)
}

func NewSearchableTextField(name string) *SearchableTextField {
	tf := &SearchableTextField{
		DataField: NewDataField(name),
		TextType:  TEXTTYPE_TEXT,
		MinLength: -1,
		MaxLength: -1,
		MinWords:  -1,
		MaxWords:  -1,
	}
	tf.Type = FIELD
	return tf
}

func (f *SearchableTextField) Compile() wajaf.NodeDef {

	t := wajaf.NewSearchableTextFieldElement(f.ID)

	t.SetAttribute("texttype", f.TextType)
	t.SetAttribute("style", f.Style)
	t.SetAttribute("classname", f.ClassName)
	t.SetAttribute("defaultvalue", fmt.Sprint(f.DefaultValue))
	t.SetData(f.Title)

	t.SetAttribute("size", f.Size)
	if f.MinLength >= 0 {
		t.SetAttribute("minlength", strconv.Itoa(f.MinLength))
	}
	if f.MaxLength >= 0 {
		t.SetAttribute("maxlength", strconv.Itoa(f.MaxLength))
	}
	if f.MinWords >= 0 {
		t.SetAttribute("minwords", strconv.Itoa(f.MinWords))
	}
	if f.MaxWords >= 0 {
		t.SetAttribute("maxwords", strconv.Itoa(f.MaxWords))
	}
	t.SetAttribute("format", f.FormatJS)

	t.SetAttribute("visible", convertModes(f.AuthModes))
	t.SetAttribute("info", convertModes(f.ViewModes))
	t.SetAttribute("readonly", convertModes(f.ReadOnlyModes))
	t.SetAttribute("notnull", convertModes(f.NotNullModes))
	t.SetAttribute("disabled", convertModes(f.DisabledModes))
	t.SetAttribute("helpmode", convertModes(f.HelpModes))
	if f.Auto {
		t.SetAttribute("auto", "yes")
	} else {
		t.SetAttribute("auto", "no")
	}

	t.AddHelp("", "", f.HelpDescription)
	t.AddMessage("defaultvalue", fmt.Sprint(f.DefaultValue))
	t.AddMessage("statusnotnull", f.StatusNotNull)
	t.AddMessage("statusbadformat", f.StatusBadFormat)
	t.AddMessage("statustooshort", f.StatusTooShort)
	t.AddMessage("statustoolong", f.StatusTooLong)
	t.AddMessage("statustoofewwords", f.StatusTooFewWords)
	t.AddMessage("statustoomanywords", f.StatusTooManyWords)
	t.AddMessage("statuscheck", f.StatusCheck)
	t.AddMessage("automessage", f.AutoMessage)

	t.AddEvent("keyup", f.KeyUpJS)
	t.AddEvent("blur", f.BlurJS)
	t.AddEvent("focus", f.FocusJS)

	return t
}

// GetValue to get the value from the field when needed. return value, ignored bool (true = ignored by construct)
func (f *SearchableTextField) GetValue(ctx *context.Context, mode Mode) (interface{}, bool, error) {

	if DEBUG {
		fmt.Println("xdominion.SearchableTextField::GetValue", f.Name, mode)
	}

	val, ignore, err := f.DataField.GetValue(ctx, mode)
	if err != nil || ignore {
		return val, ignore, err
	}

	// FILTER VALUE:
	// Lengths, formats
	if val != nil {
		sval := val.(string)
		if f.MinLength > 0 && len(sval) < f.MinLength {
			return val, ignore, errors.New(f.StatusTooShort)
		}
		if f.MaxLength > 0 && len(sval) > f.MaxLength {
			return val, ignore, errors.New(f.StatusTooLong)
		}
		nw := countWords(sval)
		if f.MinWords > 0 && nw < f.MinWords {
			return val, ignore, errors.New(f.StatusTooFewWords)
		}
		if f.MaxWords > 0 && nw > f.MaxWords {
			return val, ignore, errors.New(f.StatusTooManyWords)
		}
		if f.Format != "" {
			matched, err := regexp.MatchString(f.Format, sval)
			if err != nil {
				return val, ignore, err
			}
			if !matched {
				return val, ignore, errors.New(f.StatusBadFormat)
			}
		}
	}

	// extra code check
	if f.Check != nil {
		err = f.Check(ctx, mode, val)
	}

	return val, ignore, err
}
