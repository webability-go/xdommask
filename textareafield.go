package xdommask

import (
	"fmt"
	"strconv"

	"github.com/webability-go/wajaf"
)

type TextAreaField struct {
	*TextField
	Width  int
	Height int
}

func NewTextAreaField(name string) *TextAreaField {
	tf := &TextAreaField{
		TextField: NewTextField(name),
	}
	tf.Type = FIELD
	tf.Width = 400
	tf.Height = 100
	return tf
}

func (f *TextAreaField) Compile() wajaf.NodeDef {

	t := wajaf.NewTextAreaFieldElement(f.ID)

	t.SetAttribute("style", f.Style)
	t.SetAttribute("classname", f.ClassName)
	t.SetData(f.Title)
	t.SetAttribute("areawidth", strconv.Itoa(f.Width))
	t.SetAttribute("areaheight", strconv.Itoa(f.Height))
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

	t.AddHelp("", "", f.HelpDescription)
	t.AddMessage("defaultvalue", fmt.Sprint(f.DefaultValue))
	t.AddMessage("statusnotnull", f.StatusNotNull)
	t.AddMessage("statusbadformat", f.StatusBadFormat)
	t.AddMessage("statustooshort", f.StatusTooShort)
	t.AddMessage("statustoolong", f.StatusTooLong)
	t.AddMessage("statustoofewwords", f.StatusTooFewWords)
	t.AddMessage("statustoomanywords", f.StatusTooManyWords)
	t.AddMessage("statuscheck", f.StatusCheck)

	t.AddEvent("keyup", f.KeyUpJS)
	t.AddEvent("blur", f.BlurJS)
	t.AddEvent("focus", f.FocusJS)

	return t
}
