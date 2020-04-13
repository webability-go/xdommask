package xdommask

import (
	"fmt"
	"strconv"

	"github.com/webability-go/wajaf"
)

type TextField struct {
	*DataField

	Format   string
	FormatJS string

	MinLength int
	MaxLength int
	MinWords  int
	MaxWords  int
	Min       string
	Max       string

	StatusBadFormat    string
	StatusTooShort     string
	StatusTooLong      string
	StatusTooFewWords  string
	StatusTooManyWords string

	KeyUpJS string
	FocusJS string
	BlurJS  string
}

func NewTextField(name string) *TextField {
	return &TextField{DataField: NewDataField(name)}
}

func (f *TextField) Compile() wajaf.NodeDef {

	t := wajaf.NewTextFieldElement(f.ID)

	t.SetAttribute("style", f.Style)
	t.SetAttribute("classname", f.ClassName)
	t.SetData(f.Title)
	t.SetAttribute("size", f.Size)
	t.SetAttribute("minlength", strconv.Itoa(f.MinLength))
	t.SetAttribute("maxlength", strconv.Itoa(f.MaxLength))
	t.SetAttribute("minwords", strconv.Itoa(f.MinWords))
	t.SetAttribute("maxwords", strconv.Itoa(f.MaxWords))
	t.SetAttribute("format", f.FormatJS)

	t.SetAttribute("visible", createModes(f.AuthModes))
	t.SetAttribute("info", createModes(f.ViewModes))
	t.SetAttribute("readonly", createModes(f.ReadOnlyModes))
	t.SetAttribute("notnull", createModes(f.NotNullModes))
	t.SetAttribute("disabled", createModes(f.DisabledModes))
	t.SetAttribute("helpmode", createModes(f.HelpModes))

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
	t.AddEvent("blue", f.BlurJS)
	t.AddEvent("focus", f.FocusJS)

	return t
}
