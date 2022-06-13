package xdommask

import (
	"fmt"

	"github.com/webability-go/wajaf"
)

type ColorField struct {
	*TextField
}

func NewColorField(name string) *ColorField {
	df := &ColorField{
		TextField: NewTextField(name),
	}
	df.Type = FIELD
	return df
}

func (f *ColorField) Compile() wajaf.NodeDef {

	t := wajaf.NewColorFieldElement(f.ID)
	t.SetAttribute("style", f.Style)
	t.SetAttribute("classname", f.ClassName)
	t.SetData(f.Title)

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
	t.AddMessage("statuscheck", f.StatusCheck)

	t.AddEvent("keyup", f.KeyUpJS)
	t.AddEvent("blur", f.BlurJS)
	t.AddEvent("focus", f.FocusJS)
	return t
}
