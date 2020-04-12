package xdommask

import "github.com/webability-go/wajaf"

type PWField struct {
	*TextField
}

func NewPWField(name string) *PWField {
	return &PWField{TextField: NewTextField(name)}
}

func (f *PWField) Compile() wajaf.NodeDef {

	b := wajaf.NewTextFieldElement("", "")

	return b
}

/*
class DomMaskPWField extends DomMaskField
{
  public $Length = null;            // integer
  public $PWTwice = false;          // boolean
  public $PWString = '**********';  // string to hide password

  function __construct($name = '', $iftable = false)
  {
    parent::__construct($name, $iftable);
    $this->type = 'pw';
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
