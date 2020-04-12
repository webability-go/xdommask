package xdommask

import "github.com/webability-go/wajaf"

type TextField struct {
	*DataField

	Format   string
	FormatJS string

	MinLength int
	MaxLength int
	MinWords  int
	MaxWords  int
	Min       string
	Max       string

	StatusBadFormat    string
	StatusTooShort     string
	StatusTooLong      string
	StatusTooFewWords  string
	StatusTooManyWords string
}

func NewTextField(name string) *TextField {
	return &TextField{DataField: NewDataField(name)}
}

func (f *TextField) Compile() wajaf.NodeDef {

	t := wajaf.NewTextFieldElement("", "")

	return t

	//	$f = new \wajaf\textfieldElement($this->name);

}

/*
class DomMaskTextField extends DomMaskField
{
  public $keyup = null;
  public $blur = null;
  public $focus = null;

  function __construct($name = '', $iftable = false)
  {
    parent::__construct($name, $iftable);
    $this->type = 'text';
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
    $f->setDisabled($this->DomMask->createModes($this->disabledmodes));
    $f->setHelpmode($this->DomMask->createModes($this->helpmodes));

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
    if($this->keyup)
      $f->setEvent('keyup', $this->keyup);
    if($this->blur)
      $f->setEvent('blur', $this->blur);
    if($this->focus)
      $f->setEvent('focus', $this->focus);

    return $f;
  }

}

*/
