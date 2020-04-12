package xdommask

import "github.com/webability-go/wajaf"

type TextAreaField struct {
	*TextField
}

func NewTextAreaField(name string) *TextAreaField {
	return &TextAreaField{TextField: NewTextField(name)}
}

func (f *TextAreaField) Compile() wajaf.NodeDef {

	b := wajaf.NewTextAreaFieldElement("", "")

	return b
}

/*
class DomMaskTextAreaField extends DomMaskField
{
  public $width = '400';
  public $height = '100';
  public $cols = null;                // anything
  public $lines = null;                // anything

  function __construct($name = '', $iftable = false)
  {
    parent::__construct($name, $iftable);
    $this->type = 'textarea';
  }

  public function create()
  {
    $f = new \wajaf\textareafieldElement($this->name);

    $f->setAreawidth($this->width);
    $f->setAreaheight($this->height);
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
