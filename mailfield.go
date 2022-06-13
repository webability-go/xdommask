package xdommask

import (
	"github.com/webability-go/wajaf"
)

type MailField struct {
	*TextField
}

func NewMailField(name string) *MailField {
	mf := &MailField{
		TextField: NewTextField(name),
	}
	mf.TextType = TEXTTYPE_EMAIL
	return mf
}

func (f *MailField) Compile() wajaf.NodeDef {

	//	t := wajaf.NewTextFieldElement(f.ID)
	t := f.TextField.Compile()
	t.SetAttribute("format", "^[_a-z0-9-+]+(\\.[_a-z0-9-+]+)*@([0-9a-z][0-9a-z-]*[0-9a-z]\\.)+([a-z]{2,})$")
	return t
}
