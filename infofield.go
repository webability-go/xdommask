package xdommask

import (
	"github.com/webability-go/wajaf"
)

type InfoField struct {
	*Field
	Data string
}

func NewInfoField(name string, data string) *InfoField {
	inf := &InfoField{
		Field: NewField(name),
		Data:  data,
	}
	inf.Type = INFO
	return inf
}

func (f *InfoField) Compile() wajaf.NodeDef {

	t := wajaf.NewHTMLElement(f.ID, f.Data)
	t.SetAttribute("display", "block")
	return t
}
