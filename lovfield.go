package xdommask

import (
	"fmt"

	"github.com/webability-go/wajaf"
	"github.com/webability-go/xdominion"
)

// List Of Values field:
// it is static. The list comes as tags. Can come from table or array of values, can be reloaded with a listener (sub-list)
type LOVField struct {
	*DataField
	DefaultValue string
	Options      map[string]string
	MultiSelect  bool
	RadioButton  bool

	Table      *xdominion.XTable
	Order      *xdominion.XOrder
	Conditions *xdominion.XConditions
	FieldSet   *xdominion.XFieldSet

	FocusJS  string
	BlurJS   string
	ChangeJS string
}

func NewLOVField(name string) *LOVField {
	lf := &LOVField{DataField: NewDataField(name)}
	lf.Type = FIELD
	return lf
}

func (f *LOVField) Compile() wajaf.NodeDef {

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

	if f.BlurJS != "" {
		l.AddEvent("blur", f.BlurJS)
	}
	if f.FocusJS != "" {
		l.AddEvent("focus", f.FocusJS)
	}
	if f.ChangeJS != "" {
		l.AddEvent("change", f.ChangeJS)
	}
	if f.CheckJS != "" {
		l.AddChild(wajaf.NewCodeNode("", "check", f.CheckJS))
	}

	ms := "yes"
	if !f.MultiSelect {
		ms = "no"
	}
	l.SetAttribute("multiselect", ms)
	rb := "yes"
	if !f.RadioButton {
		rb = "no"
	}
	l.SetAttribute("radiobutton", rb)
	opts := wajaf.NewOptions()
	// Table ?
	if f.Table != nil {
		recs, _ := f.Table.SelectAll(f.Conditions, f.Order, f.FieldSet)
		if recs != nil {
			for _, rec := range *recs {
				p, _ := rec.GetString((*f.FieldSet)[0])
				v, _ := rec.GetString((*f.FieldSet)[1])
				opts.AddChild(wajaf.NewOption(p, p+" / "+v))
			}
		}
	} else if f.Options != nil {
		for p, v := range f.Options {
			opts.AddChild(wajaf.NewOption(p, v))
		}
	}
	l.AddChild(opts)

	return l
}
