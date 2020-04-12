package xdommask

import "github.com/webability-go/wajaf"

type HiddenField struct {
	*DataField
}

func NewHiddenField(name string) *HiddenField {
	return &HiddenField{DataField: NewDataField(name)}
}

func (f *HiddenField) Compile() wajaf.NodeDef {

	b := wajaf.NewHiddenFieldElement("", "")

	return b
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
