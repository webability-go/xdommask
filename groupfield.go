package xdommask

import (
	"github.com/webability-go/wajaf"
	"github.com/webability-go/xamboo/cms/context"
)

type GroupField struct {
	*Field
	Mask *Mask
	Ctx  *context.Context
}

func NewGroupField(name string, mask *Mask) *GroupField {
	inf := &GroupField{
		Field: NewField(name),
		Mask:  mask,
	}
	inf.Type = GROUP
	return inf
}

func (f *GroupField) Compile() wajaf.NodeDef {

	mask := f.Mask.Compile("1", f.Ctx)
	return mask
}
