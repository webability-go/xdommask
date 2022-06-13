package xdommask

import (
	"fmt"

	"github.com/webability-go/wajaf"
)

// List Of Options field:
// it is dynamic. when the user write the field, a list of possible options displays under the field and can be selected.
type LOOField struct {
	*DataField
	DefaultValue string
	Options      map[string]string
	MultiSelect  bool

	FocusJS string
	BlurJS  string
}

func NewLOOField(name string) *LOOField {
	lf := &LOOField{DataField: NewDataField(name)}
	lf.Type = FIELD
	return lf
}

func (f *LOOField) Compile() wajaf.NodeDef {

	l := wajaf.NewLOVFieldElement(f.ID)

	l.SetAttribute("style", f.Style)
	l.SetAttribute("classname", f.ClassName)
	l.SetData(f.Title)
	l.SetAttribute("size", f.Size)

	l.SetAttribute("visible", convertModes(f.AuthModes))
	l.SetAttribute("info", convertModes(f.ViewModes))
	l.SetAttribute("readonly", convertModes(f.ReadOnlyModes))
	l.SetAttribute("notnull", convertModes(f.NotNullModes))
	l.SetAttribute("disabled", convertModes(f.DisabledModes))
	l.SetAttribute("helpmode", convertModes(f.HelpModes))

	l.AddHelp("", "", f.HelpDescription)
	l.AddMessage("defaultvalue", fmt.Sprint(f.DefaultValue))
	l.AddMessage("statusnotnull", f.StatusNotNull)
	l.AddMessage("statuscheck", f.StatusCheck)

	l.AddEvent("blur", f.BlurJS)
	l.AddEvent("focus", f.FocusJS)

	ms := "yes"
	if !f.MultiSelect {
		ms = "no"
	}
	l.SetAttribute("multiselect", ms)
	opts := wajaf.NewOptions()
	for p, v := range f.Options {
		opts.AddChild(wajaf.NewOption(p, v))
	}
	l.AddChild(opts)

	return l
}
