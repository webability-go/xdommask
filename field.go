package xdommask

import "github.com/webability-go/wajaf"

const (
	CONTROL = "control"
	FIELD   = "field"
	INFO    = "info"
	HIDDEN  = "hidden"
)

type FieldDef interface {
	Compile() wajaf.NodeDef
	GetType() string
}

type Field struct {
	Name string
	Type string
	ID   string

	Style     string
	ClassName string

	AuthModes     Mode
	ViewModes     Mode
	ReadOnlyModes Mode
	NotNullModes  Mode
	DisabledModes Mode
	HelpModes     Mode

	HelpToolTip     string
	HelpDescription string
	HelpTitle       string
	TabIndex        int
	Size            string

	JS    string
	Focus string
	Blur  string
}

func NewField(name string) *Field {
	return &Field{Name: name, ID: name}
}

func (f *Field) Compile() wajaf.NodeDef {
	return nil
}

func (f *Field) GetType() string {
	return f.Type
}

type DataField struct {
	*Field
	InRecord    bool
	URLVariable string

	Title    string // title for editable, in button for control
	Auto     bool
	Encoded  bool
	Entities bool

	AutoMessage string

	EmptyNotNull bool
	NullOnEmpty  bool
	NullValue    string
	MD5Encrypted bool

	StatusOK      string
	StatusNotNull string
	StatusCheck   string
}

func NewDataField(name string) *DataField {
	return &DataField{Field: NewField(name)}
}

func (f *DataField) Compile() wajaf.NodeDef {
	return nil
}

type ControlField struct {
	*Field
	TitleInsert string
	TitleUpdate string
	TitleDelete string
	TitleView   string
	Control     bool
}

func NewControlField(name string) *ControlField {
	return &ControlField{Field: NewField(name), Control: true}
}

func (f *ControlField) Compile() wajaf.NodeDef {
	return nil
}

/*
/*
class DomMaskField extends \core\WAClass
{
  protected $DomMask;                  // the DomMask that owns us
  protected $caller;
  public $calcfunction = null;       // function to calculate the value(s) of the field based on the value from table


  public function __construct($name, $inrecord)
  {
    parent::__construct();
    $this->name = $name;
    $this->inrecord = $inrecord;
    $this->urlvariable = $name; // by default, name of field
  }

  public function fix($DomMask, $caller)
  {
    $this->DomMask = $DomMask;
    $this->caller = $caller;
  }

  public function needEncType()
  {
    // replace this method in your extended field to get back true if enctype is needed in the form: file, image, etc.
    return false;
  }


  public function gettitle()
  {
    return $this->title;
  }

  public function getType()
  {
    return $this->type;
  }

  public function getMessages()
  {
    return $this->helpdescription;
  }

  // VARIABLES functions
  // gets the value from DATABASE (record array) or DEFAULT or FUNCTION
  public function getValue($record)
  {
    if ($this->DomMask->realmode != DomMask::INSERT && $this->name && $this->inrecord)
    {
      if (@isset($record[$this->name]))
      {
        $val = $record[$this->name];
      }
      else
      {
        $val = '';
      }
      if ($this->encoded)
        $val = rawurldecode($val);
      if ($this->entities)
      {
        $trans = get_html_translation_table (HTML_ENTITIES);
        $trans = array_flip ($trans);
        $val = strtr($val, $trans);
      }
    }
    elseif ($this->DomMask->realmode == DomMask::INSERT && $this->auto)
    {
      $val = $this->automessage;
    }
    elseif ($this->DomMask->realmode == DomMask::INSERT && $this->default)
    {
      $val = $this->default;
    }
    elseif (!$this->inrecord)
    {
      $val = $this->default;
    }
    else
      $val = "";

    if (($this->DomMask->realmode == DomMask::DELETE || $this->DomMask->realmode == DomMask::VIEW) && is_null($val))
      $val = "";

    if ($this->calcfunction)
    {
      $Fct = $this->calcfunction;
      $val = $this->DomMask->$Fct($val, $record);
    }
    return $val;
  }

  // format the value from URL
  protected function formatField($val)
  {
    if ($this->encoded)
      $val = rawurldecode($val);
    if ($this->entities)
      $val = nl2br(htmlentities($val, ENT_COMPAT, $this->DomMask->CharSet));
    else
    {
      $val = nl2br($val);
    }
    return $val;
  }

  protected function filterParameter($val)
  {
    if ($val !== null && $val!="")
    {
      if ($this->entities)
        $val = htmlentities($val, ENT_COMPAT, $this->DomMask->CharSet);
      if ($this->encoded)
        $val = rawurlencode($val);
      if ($this->md5encrypted && $this->DomMask->realmode == DomMask::INSERT)  // MD5 only apply on INSERT (cannot change on update ?)
        $val = MD5($val);
      if ($this->maxlength)
        $val=substr($val, 0, $this->maxlength);
    }
    return $val;
  }

  // get and filter the parameter from URL
  public function getParameter()
  {
    $val = $this->DomMask->getParameter($this->urlvariable);
    if ($this->format)
    {
      if (!preg_match($this->format, $val))
        $val = null;
    }
    if ($this->nullonempty && !$val)
    { $val = null; }
    if ($this->emptyonnull && !$val)
    { $val = ''; }
    return $this->filterParameter($val);
  }

  // CONTROL functions, do nothing by default, are called by DomMask
  public function preInsert($record) // record is a DB_Record, by ref
  {
  }

  public function postInsert($key, $record)
  {
  }

  public function preUpdate($key, $record, $oldrecord) // record is a DB_Record, by ref
  {
  }

  public function postUpdate($key, $record, $oldrecord)
  {
  }

  public function preDelete($key, $oldrecord)
  {
  }

  public function postDelete($key, $oldrecord)
  {
  }

  // 4GL method to create a wajaffield object
  public function create()
  {
    return null;
  }

  public function loadDefinition($data)
  {
    foreach($data as $p => $v)
    {
      if ($p == 'type')
        continue;
      if (is_array($v))
      {
        $t = array();
        foreach($v as $p1 => $v1)
          $t[constant('DomMask::'. $p1)] = $v1;
        $this->$p = $t;
      }
      else
        $this->$p = $v;
    }
  }

}
*/
