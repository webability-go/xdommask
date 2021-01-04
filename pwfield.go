package xdommask

import (
	"fmt"
	"strconv"

	"github.com/webability-go/wajaf"
)

type PWField struct {
	*TextField
}

func NewPWField(name string) *PWField {
	pf := &PWField{
		TextField: NewTextField(name),
	}
	pf.Type = FIELD
	return pf
}

func (f *PWField) Compile() wajaf.NodeDef {

	t := wajaf.NewTextFieldElement(f.ID)

	t.SetAttribute("style", f.Style)
	t.SetAttribute("classname", f.ClassName)
	t.SetData(f.Title)
	t.SetAttribute("size", f.Size)
	if f.MinLength >= 0 {
		t.SetAttribute("minlength", strconv.Itoa(f.MinLength))
	}
	if f.MaxLength >= 0 {
		t.SetAttribute("maxlength", strconv.Itoa(f.MaxLength))
	}
	if f.MinWords >= 0 {
		t.SetAttribute("minwords", strconv.Itoa(f.MinWords))
	}
	if f.MaxWords >= 0 {
		t.SetAttribute("maxwords", strconv.Itoa(f.MaxWords))
	}
	t.SetAttribute("format", f.FormatJS)

	t.SetAttribute("visible", createModes(f.AuthModes))
	t.SetAttribute("info", createModes(f.ViewModes))
	t.SetAttribute("readonly", createModes(f.ReadOnlyModes))
	t.SetAttribute("notnull", createModes(f.NotNullModes))
	t.SetAttribute("disabled", createModes(f.DisabledModes))
	t.SetAttribute("helpmode", createModes(f.HelpModes))

	t.AddHelp("", "", f.HelpDescription)
	t.AddMessage("defaultvalue", fmt.Sprint(f.DefaultValue))
	t.AddMessage("statusnotnull", f.StatusNotNull)
	t.AddMessage("statusbadformat", f.StatusBadFormat)
	t.AddMessage("statustooshort", f.StatusTooShort)
	t.AddMessage("statustoolong", f.StatusTooLong)
	t.AddMessage("statustoofewwords", f.StatusTooFewWords)
	t.AddMessage("statustoomanywords", f.StatusTooManyWords)
	t.AddMessage("statuscheck", f.StatusCheck)

	t.AddEvent("keyup", f.KeyUpJS)
	t.AddEvent("blur", f.BlurJS)
	t.AddEvent("focus", f.FocusJS)

	return t
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
