package xdommask

import "github.com/webability-go/wajaf"

type LOOField struct {
	*DataField
}

func NewLOOField(name string) *LOOField {
	return &LOOField{DataField: NewDataField(name)}
}

func (f *LOOField) Compile() wajaf.NodeDef {

	b := wajaf.NewLOVFieldElement("", "")

	return b
}

/*

class DomMaskLOOField extends DomMaskField
{
  public $blur = null;
  public $MultiSelect = false;       // boolean

  public $options = null;            // array

  public $ListTable = null;          // List DB_Table object
  public $ListKey = null;            // the key field on list
  public $ListName = null;           // the name field on list
  public $ListOrder = null;          // field order by on list
  public $ListWhere = null;          // where object on list
  public $ListSeparator = ' / ';     // separator if the list name is an array
  public $ListEncoded = false;       // boolean true if the list result is encoded
  public $ListEntities = false;      // boolean true if the list result has entities

  public $Controlling = null;        // this LOO controls another LOO (Id of the FIELD)
  public $ControllingOptions = null; // The special options (array( father => array(childs => childs) ) )
  public $ControllingDefault = null; // The default of the child controlled LOO
  public $ControllingIndex = null;   // the tabindex of the controlled field, used to actualize field validity
  public $OnEvent = null;            // If a DomMaskField::LOO or DomMaskField::LOV have a javascript event
  public $CheckAsBool = false;       // set to true to get the check bos as a boolean instead of an array

  function __construct($name = '', $iftable = false)
  {
    parent::__construct($name, $iftable);
    $this->type = 'loo';
  }

  public function create()
  {
    $f = new \wajaf\lovfieldElement($this->name);

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
    $f->setOptions($this->options);

    return $f;
  }

}

*/
