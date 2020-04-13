package xdommask

import "github.com/webability-go/wajaf"

type LOVField struct {
	*LOOField
}

func NewLOVField(name string) *LOVField {
	return &LOVField{LOOField: NewLOOField(name)}
}

func (f *LOVField) Compile() wajaf.NodeDef {

	b := wajaf.NewLOVFieldElement(f.ID)

	return b
}

/*
class DomMaskLOVField extends DomMaskField
{
  public $RadioButton = false;      // boolean
  public $MultiSelect = false;      // boolean

  public $ListTable = null;         // List DB_Table object
  public $ListKey = null;           // the key field on list
  public $ListName = null;          // the name field on list
  public $ListOrder = null;         // field order by on list
  public $ListWhere = null;         // where object on list
  public $ListSeparator = ' / ';    // separator if the list name is an array
  public $ListEncoded = false;      // boolean true if the list result is encoded
  public $ListEntities = false;     // boolean true if the list result has entities
  public $Controlling = null;       // this LOV controls another LOV (Id of the FIELD)
  public $ControllingOptions = null; // The special options (array( father => array(childs => childs) ) )
  public $ControllingIndex = null;   // the tabindex of the controlled field, used to actualize field validity
  public $OnEvent = null;            // If a DB_MaskField::LOO or DB_MaskField::LOV have a javascript event

  function __construct($MF = null)
  {
    $name = ""; $iftable = false;
    if ($MF)
    {
      $name = $MF->Name;
      $iftable = $MF->IfTable;
      if ($MF instanceof DomMaskLOVField)
      {
        $this->RadioButton = $MF->RadioButton;
        $this->MultiSelect = $MF->MultiSelect;
        $this->ListTable = $MF->ListTable;
        $this->ListKey = $MF->ListKey;
        $this->ListName = $MF->ListName;
        $this->ListEncoded = $MF->ListEncoded;
        $this->ListEntities = $MF->ListEntities;
        $this->ListWhere = $MF->ListWhere;
        $this->ListOrder = $MF->ListOrder;
        $this->ListSeparator = $MF->ListSeparator;
        $this->Controlling = $MF->Controlling;
        $this->ControllingOptions = $MF->ControllingOptions;
        $this->ControllingIndex = $MF->ControllingIndex;
        $this->OnEvent = $MF->OnEvent;
      }
    }
    parent::__construct($name, $iftable, DB_MaskField::LOV, $MF);
  }

  public function create()
  {
    $f = new \wajaf\dommasklovfieldElement();
    $f->setLink($this->DomMask->RealMaskId);
    $f->setMessage('title', $this->Title);
    $f->setMessage('helpsummary', $this->HelpSummary);
    $f->setMessage('helptitle', $this->HelpTitle);
    $f->setMessage('helpdescription', $this->HelpDescription);
    $f->setMessage('statusok', $this->StatusOK);
    $f->setMessage('statusnotnull', $this->StatusNotNull);
    $f->setMessage('statuscheck', $this->StatusCheck);
    $f->setNotnull($this->NotNull);

    return $f;
  }

}

*/
