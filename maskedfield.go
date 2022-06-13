package xdommask

import (
//	"github.com/webability-go/wajaf"
)

type MaskedField struct {
	*TextField
}

func NewMaskedField(name string) *MaskedField {
	pf := &MaskedField{
		TextField: NewTextField(name),
	}
	pf.TextType = TEXTTYPE_MASKED
	return pf
}

/*
func (f *MaskedField) Compile() wajaf.NodeDef {

	t := f.TextField.Compile()
	return t
}
*/
