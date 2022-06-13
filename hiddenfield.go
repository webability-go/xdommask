package xdommask

import "github.com/webability-go/wajaf"

// hiddenfield is a hidden field to pass information between field group and server.
// It should have a default value and is setable by code. Its value will be send to the server in any mode if authorized
type HiddenField struct {
	*DataField
	DefaultValue string
}

func NewHiddenField(name string) *HiddenField {
	hf := &HiddenField{DataField: NewDataField(name)}
	hf.Type = HIDDEN
	return hf
}

func (f *HiddenField) Compile() wajaf.NodeDef {

	h := wajaf.NewHiddenFieldElement(f.ID)
	h.SetData(f.DefaultValue)

	return h
}

/*
class DomMaskHiddenField extends DomMaskField
{
  function __construct($name = '', $iftable = false)
  {
    parent::__construct($name, $iftable);
    $this->type = 'hidden';
  }

  public function create()
  {
    $f = new \wajaf\hiddenfieldElement($this->name);

    $f->setData($this->default);

    return $f;
  }

}

*/
