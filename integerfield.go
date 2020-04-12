package xdommask

import "github.com/webability-go/wajaf"

type IntegerField struct {
	*TextField
	Min int
	Max int
}

func NewIntegerField(name string) *IntegerField {
	return &IntegerField{TextField: NewTextField(name)}
}

func (f *IntegerField) Compile() wajaf.NodeDef {

	b := wajaf.NewTextFieldElement("", "")

	return b
}

/*
namespace dommask;

class DomMaskIntegerField extends DomMaskField
{
  function __construct($name = '', $iftable = false)
  {
    parent::__construct($name, $iftable);
    $this->type = 'integer';
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
    if ($this->auto)
      $f->setAuto('yes');

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
    $f->setMessage('automessage', $this->automessage);

    return $f;
  }

}

*/
