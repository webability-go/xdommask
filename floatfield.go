package xdommask

import "github.com/webability-go/wajaf"

type FloatField struct {
	*TextField
	Default float64
	Min     float64
	Max     float64
}

func NewFloatField(name string) *FloatField {
	ff := &FloatField{TextField: NewTextField(name)}
	ff.Type = FIELD
	return ff
}

func (f *FloatField) Compile() wajaf.NodeDef {

	b := wajaf.NewTextFieldElement(f.ID)

	return b
}

/*
class DomMaskRealField extends DomMaskField
{
  function __construct($name = '', $iftable = false)
  {
    parent::__construct($name, $iftable);
    $this->type = 'real';
  }

  public function create()
  {
    $f = new \wajaf\textfieldElement($this->name);

    $f->setSize($this->size);
    $f->setMinlength($this->minlength);
    $f->setMaxlength($this->maxlength);
    $f->setMinwords($this->minwords);
    $f->setMaxwords($this->maxwords);
    $f->setFormat($this->formatjs);

    $f->setVisible($this->DomMask->createModes($this->authmodes));
    $f->setInfo($this->DomMask->createModes($this->viewmodes));
    $f->setReadonly($this->DomMask->createModes($this->readonlymodes));
    $f->setNotnull($this->DomMask->createModes($this->notnullmodes));
    $f->setDisabled('');
    $f->setHelpmode('12');
//    $f->setTabindex($this->tabindex);

    $f->setData($this->title);

    $f->setMessage('defaultvalue', $this->default);
    $f->setMessage('helpdescription', $this->helpdescription);
    $f->setMessage('statusnotnull', $this->statusnotnull);
    $f->setMessage('statusbadformat', $this->statusbadformat);
    $f->setMessage('statustooshort', $this->statustooshort);
    $f->setMessage('statustoolong', $this->statustoolong);
    $f->setMessage('statustoofewwords', $this->statustoofewwords);
    $f->setMessage('statustoomanywords', $this->statustoomanywords);
    $f->setMessage('statuscheck', $this->statuscheck);

    return $f;
  }

}
*/
