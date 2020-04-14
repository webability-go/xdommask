package xdommask

import (
	"encoding/json"
	"fmt"

	"github.com/webability-go/wajaf"
	"github.com/webability-go/xdominion"
)

type Mode int

const (
	INSERT Mode = 1 << iota
	UPDATE
	DELETE
	VIEW
	DOINSERT
	DOUPDATE
	DODELETE
	CONFIRMDELETE
)

type Mask struct {
	ID string

	Display       string
	Style         string
	ClassName     string
	ClassNameZone string
	VarMode       string
	VarOrder      string
	VarKey        string

	Mode      Mode
	AuthModes Mode
	Key       interface{}

	Variables map[string]string
	SuccessJS string
	FailureJS string

	AlertMessage         string
	ServerMessage        string
	InsertTitle          string
	UpdateTitle          string
	DeleteTitle          string
	ViewTitle            string
	DoInsertMessage      string
	DoUpdateMessage      string
	DoDeleteMessage      string
	ConfirmDeleteMessage string

	Counter int
	Fields  []FieldDef

	// Hooks
	PreInsert  func(*Mask, *xdominion.XRecord) error
	Insert     func(*Mask, *xdominion.XRecord) error
	PostInsert func(*Mask, *xdominion.XRecord) error
	PreUpdate  func(*Mask, *xdominion.XRecord) error
	Update     func(*Mask, *xdominion.XRecord) error
	PostUpdate func(*Mask, *xdominion.XRecord) error
	PreDelete  func(*Mask, *xdominion.XRecord) error
	Delete     func(*Mask, *xdominion.XRecord) error
	PostDelete func(*Mask, *xdominion.XRecord) error
	GetRecord  func(*Mask) *xdominion.XRecord
}

func NewMask(id string) *Mask {
	return &Mask{
		ID:        id,
		Variables: map[string]string{},
		Fields:    []FieldDef{},
		VarMode:   "Mode",
		VarOrder:  "Order",
		VarKey:    "Key",
		Counter:   1,
	}
}

func (m *Mask) AddField(f FieldDef) {
	m.Fields = append(m.Fields, f)
}

func (m *Mask) Compile() wajaf.NodeDef {

	group := wajaf.NewGroupContainer(m.ID)
	group.SetAttribute("display", m.Display)
	group.SetAttribute("style", m.Style)
	group.SetAttribute("classname", m.ClassName)
	group.SetAttribute("classnamezone", m.ClassNameZone)
	group.SetAttribute("haslistener", "yes")

	group.SetAttribute("varmode", m.VarMode)
	group.SetAttribute("varorder", m.VarOrder)
	group.SetAttribute("varkey", m.VarKey)

	group.SetAttribute("authmodes", createModes(m.AuthModes))
	group.SetAttribute("mode", createModes(m.Mode))
	if m.Key != nil {
		group.SetAttribute("key", fmt.Sprint(m.Key))
	}

	group.AddMessage("alertmessage", m.AlertMessage)
	group.AddMessage("servermessage", m.ServerMessage)
	group.AddMessage("titleinsert", m.InsertTitle)
	group.AddMessage("titleupdate", m.UpdateTitle)
	group.AddMessage("titledelete", m.DeleteTitle)
	group.AddMessage("titleview", m.ViewTitle)
	group.AddMessage("insertok", m.DoInsertMessage)
	group.AddMessage("updateok", m.DoUpdateMessage)
	group.AddMessage("deleteok", m.DoDeleteMessage)
	group.AddMessage("confirmdelete", m.ConfirmDeleteMessage)

	if m.SuccessJS != "" {
		group.AddEvent("success", m.SuccessJS)
	}
	if m.FailureJS != "" {
		group.AddEvent("failure", m.FailureJS)
	}

	zcontrol := wajaf.NewGroupZone("control", "")

	for _, f := range m.Fields {
		if _, ok := f.(*ButtonField); ok {
			zcontrol.AddChild(f.Compile())
			continue
		}
		z := wajaf.NewGroupZone("field", "")
		z.AddChild(f.Compile())
		group.AddChild(z)
	}

	for id, val := range m.Variables {
		h := wajaf.NewHiddenFieldElement(id)
		h.SetData(val)
		zcontrol.AddChild(h)
	}
	group.AddChild(zcontrol)

	// Original dataset
	if m.GetRecord != nil {
		rec := m.GetRecord(m)
		if rec != nil {
			// rec must be a JSON
			jsonrec, _ := json.Marshal(rec)
			zdata := wajaf.NewGroupDataset("", string(jsonrec))
			group.AddChild(zdata)
		}
	}

	return group
}

func createModes(mode Mode) string {
	f := ""
	if mode&INSERT == INSERT {
		f += "1"
	}
	if mode&UPDATE == UPDATE {
		f += "2"
	}
	if mode&DELETE == DELETE {
		f += "3"
	}
	if mode&VIEW == VIEW {
		f += "4"
	}
	return f
}

/*
class DomMask extends \core\WAClass
{
  // I18N server error messages
  private static $init = false;
  private static $messages = array(
    'dommask.badtable' => 'Error: the first parameter of the constructor is not a DB_table.',
    'dommask.baddescriptor' => 'Error: The descriptor has a bad format.',
    'dommask.badfield' => 'Error: The added field is not a DomMaskfield object.',
    'dommask.badtemplate' => 'Error: You need a valid template to create a Form.'
  );

  public $table = null;                 // the DB_table to use, if no record
  public $maskid = null;                // default id="" in <form> tag. mandatory




  public $realmode;                   // operating mode of state machine, based on mode, authmodes
  private $allmodes = array(1,2,4,8,16,32,64,128);  // all the available modes
  private $execmodes = array(16,32,64);             // only the modes that implies an action on the data (execute called)
  private $datamodes = array(2,4,8,32,64,128);      // only the modes to display already existing data (getRecord called)
  private $move = 0;                  // move the records for view mode


  private $counter = 0;               // fields counter. Do not set or modify
  private $execmessages = '';         // execute result messages, global level


  public $record = null;              // a defined record(s), array, type DB_record or DB_records, if no DB_table
  public $key = null;                 // the default record key
  public $keyfield = null;            // the key field name

  // HTML
  public $action = '';                // the action of the form
  public $method = 'POST';            // the method of the form, can be get/post (only for HTML, since JSON and 4GL post on ajax)
  public $varmode = 'groupmode';
  public $varorder = 'grouporder';
  public $varkey = 'groupkey';
  public $varfield = 'groupfield';

  protected $insertedkey = null;        // last inserted key(s)

  public $charset = "UTF-8";          // the charset to calculate entities if needed

  public $variables = null;           // array, variables to keep present when submitting info

  public $jsonsuccess = null;
  public $jsonfailure = null;
  public $alertmessage = 'Error, please check the fields in red.';
  public $servermessage = 'The server had an unknown internal error.';
  public $titles = array(
    DomMask::INSERT => 'Insert:',
    DomMask::UPDATE => 'Update:',
    DomMask::DELETE => 'Delete:',
    DomMask::VIEW => 'View:',
    DomMask::DOINSERT => 'Insert Result:',
    DomMask::DOUPDATE => 'Update Result:',
    DomMask::CONFIRMDELETE => 'Delete Confirmation:',
    DomMask::DODELETE => 'Delete Result:',
    );

  public $actionmessages = array(
    DomMask::DOINSERT => 'Insert successfull.',
    DomMask::DOUPDATE => 'Update successfull.',
    DomMask::DODELETE => 'Delete successfull.',
    );

  // template is optional in 4GL mode, mandatory into HTML mode
  function __construct($table = null, $descriptor = null, $caller = null)
  {
    if (!self::$init)
    {
      // send messages to  \core\WAMessage
       \core\WAMessage::addmessages(self::$messages);
      self::$init = true;
    }

    parent::__construct();

    if (self::$debug || $this->localdebug)
      $this->doDebug("DomMask::__construct(table, descriptor)", WADebug::SYSTEM);

    if ($table !== null && !($table instanceof \dominion\DB_table))
      throw new \throwables\DomMaskError( \core\WAMessage::getMessage('dommask.badtable'));
    $this->table = $table;

    if (is_array($descriptor))
    {
      $this->loadDefinition($descriptor);
    }
    elseif (is_string($descriptor))
    {
      if (strpos($descriptor, '<?xml') !== false)
        $this->loadDefinition(WASimpleXML::tags($descriptor));
      elseif (strlen($descriptor) < 512 && is_file($descriptor))
        $this->loadDefinition(WASimpleXML::tags(file_get_contents($descriptor)));
    }
    $this->caller = $caller;
  }

  public function addMaskfield($Maskfield)
  {
    if (self::$debug || $this->localdebug)
      $this->doDebug("DomMask::addMaskfield(Maskfield)", WADebug::SYSTEM);

    if (!($Maskfield instanceof DomMaskfield))
      throw new \throwables\DomMaskError( \core\WAMessage::getMessage('dommask.badfield'));

    $Maskfield->fix($this, $this->caller);
    $this->fields[] = $Maskfield;
  }

  public function getUniqueID()
  {
    return $this->counter++;
  }

  public function getInsertedkey()
  {
    return $this->insertedkey;
  }

  // get the record from the selected source
  // position is: -2: first, -1: previous, 0 = this one, 1 = next, 2 = last
  protected function getRecord($key, $position = 0)
  {
    if ($this->caller && isset($this->caller->getRecord))
      return $this->caller->getRecord($key, $position);

    if ($this->realmode == DomMask::INSERT)
      return null;
    if ($this->table)
      return $this->table->doSelect($this->key);
    if ($this->record)
      return $this->record;
    return null;
  }

  protected function preInsert($record)
  {
    if ($this->caller && method_exists($this->caller, 'preInsert'))
      $this->caller->preInsert($record);
  }

  protected function Insert($record)
  {
    if ($this->caller && method_exists($this->caller, 'Insert'))
      return $this->caller->Insert($record);

    if ($this->table)
    {
      $this->table->doInsert($record);
      return $this->table->getInsertedkey();
    }
    return null;
  }

  protected function postInsert($key, $record)
  {
    if ($this->caller && method_exists($this->caller, 'postInsert'))
      $this->caller->postInsert($key, $record);
  }

  protected function preUpdate($key, $record, $oldrecord)
  {
    if ($this->caller && method_exists($this->caller, 'preUpdate'))
      $this->caller->preUpdate($key, $record, $oldrecord);
  }

  protected function Update($key, $record, $oldrecord)
  {
    if ($this->caller && method_exists($this->caller, 'Update'))
      return $this->caller->Update($key, $record, $oldrecord);

    if ($this->table)
      $this->table->doUpdate($key, $record);
  }

  protected function postUpdate($key, $record, $oldrecord)
  {
    if ($this->caller && method_exists($this->caller, 'postUpdate'))
      $this->caller->postUpdate($key, $record, $oldrecord);
  }

  protected function preDelete($key, $oldrecord)
  {
    if ($this->caller && method_exists($this->caller, 'preDelete'))
      $this->caller->preDelete($key, $oldrecord);
  }

  protected function Delete($key, $oldrecord)
  {
    if ($this->caller && method_exists($this->caller, 'Delete'))
      return $this->caller->Delete($key, $oldrecord);

    if ($this->table)
      $this->table->doDelete($key);
  }

  protected function postDelete($key, $oldrecord)
  {
    if ($this->caller && method_exists($this->caller, 'postDelete'))
      $this->caller->postDelete($key, $oldrecord);
  }

  private function __buildrecord()
  {
    // we loop over the fields
    $record = new \dominion\DB_Record();
    foreach ($this->fields as $K => $F)
    {
      if ($F->inrecord) // the field must be part of the table
      {
        if (!($F->authmodes & $this->mode))
          continue;
        if ($F->viewmodes & $this->mode)
          continue;

        if ($F->auto)
        {
          if ($this->mode == DomMask::INSERT)
            $val = $F->default;
          else
            continue;
        }
        else
        {
          $val = $F->getParameter();
        }
        if (is_array($val))
        {
          foreach($val as $rf => $rv)
          {
            $record[$rf] = $rv;
          }
        }
        else
          $record[$F->name] = $val;
      }
    }
    return $record;
  }

  // will return the necesary data based on context for JSON and 4GL
  private function execute()
  {
    switch($this->realmode)
    {
      case DomMask::DOINSERT:
        // BUILD RAW RECORD
        $record = $this->__buildrecord(1);
        // PREINSERT
        $this->preInsert($record);
        foreach ($this->fields as $K => $F)
          $F->preInsert($record);
        // INSERT
        $this->insertedkey = $this->Insert($record);
        // POSTINSERT
        foreach ($this->fields as $K => $F)
          $F->postInsert($this->insertedkey, $record);
        $this->postInsert($this->insertedkey, $record);
        $this->key = $this->insertedkey;
        $this->execmessages = $this->actionmessages[DomMask::DOINSERT];
        break;
      case DomMask::DOUPDATE:
        // GET OLD RECORD
        $oldrecord = $this->getRecord($this->key);
        // BUILD RAW RECORD
        $record = $this->__buildrecord(2);
        // PREUPDATE
        $this->preUpdate($this->key, $record, $oldrecord);
        foreach ($this->fields as $K => $F)
          $F->preUpdate($this->key, $record, $oldrecord);
        // UPDATE
        $this->Update($this->key, $record, $oldrecord);
        // POST UPDATE
        foreach ($this->fields as $K => $F)
          $F->postUpdate($this->key, $record, $oldrecord);
        $this->postUpdate($this->key, $record, $oldrecord);
        $this->execmessages = $this->actionmessages[DomMask::DOUPDATE];
        break;
      case DomMask::DODELETE:
        // GET OLD RECORD
        $oldrecord = $this->getRecord($this->key);
        // PREDELETE
        $this->preDelete($this->key, $oldrecord);
        foreach ($this->fields as $K => $F)
          $F->preDelete($this->key, $oldrecord);
        // DELETE
        $this->Delete($this->key, $oldrecord);
        // POSTDELETE
        foreach ($this->fields as $K => $F)
          $F->postDelete($this->key, $oldrecord);
        $this->postDelete($this->key, $oldrecord);
        $this->execmessages = $this->actionmessages[DomMask::DODELETE];
        break;
    }
  }

  public function createModes($modes)
  {
    return (($modes & DomMask::INSERT)?'1':'') . (($modes & DomMask::UPDATE)?'2':'') . (($modes & DomMask::DELETE)?'3':'') . (($modes & DomMask::VIEW)?'4':'');
  }

  // creates the WAJAF objects and gets back the container
  // type can be json, 4gl
  public function run($type = DomMask::JSON)
  {
    $order = $this->getParameter($this->varorder);
    $this->mode = $this->getParameter($this->varmode);
    $this->realmode = $this->mode;
    $data = array();
    switch($order)
    {
      case 'start':
        $data = $this->createJSON($data);
        break;
      case 'next':
        $key = $this->getParameter($this->varkey);
        $data = $this->getRecord($key, 1);
        break;
      case 'previous':
        $key = $this->getParameter($this->varkey);
        $data = $this->getRecord($key, -1);
        break;
      case 'last':
        $key = $this->getParameter($this->varkey);
        $data = $this->getRecord($key, 2);
        break;
      case 'first':
        $key = $this->getParameter($this->varkey);
        $data = $this->getRecord($key, -2);
        break;
      case 'image':
        $key = $this->getParameter($this->varkey);
        $field = $this->getParameter($this->varfield);
        // gives the control to the field
        foreach($this->fields as $f)
          if ($f->name == $field)
            $data = $f->prepareImage();
        break;
      case 'submit':
        if ($this->mode == 1)
          $this->realmode = DomMask::DOINSERT;
        elseif ($this->mode == 2)
        {
          $this->realmode = DomMask::DOUPDATE;
          $this->key = $this->getParameter($this->varkey);
        }
        elseif ($this->mode == 3)
        {
          $this->realmode = DomMask::DODELETE;
          $this->key = $this->getParameter($this->varkey);
        }
        $this->execute();
        $data = array('success' => true, 'messages' => array('text' => 'Exito'));

//*************** NOTE: SHOULD GETS BACK THE FINAL RECORD FOR IF THERE ARE SOME CALCULATED FIELDS (LIKE IMAGES NAMES; SUM FIELDS; LINKS, PATHS, ETC)



        break;
    }

    return $data;
//    return json_encode($data);
  }

  public function code()
  {
    $container = new \wajaf\groupContainer($this->maskid);
    $container->setStyle($this->style);
    $container->setAuthmodes( $this->createModes($this->authmodes) );
    $container->setMode( $this->createModes($this->mode) );
    $container->setKey($this->key);
    $container->setVarmode($this->varmode);
    $container->setVarorder($this->varorder);
    $container->setVarkey($this->varkey);
    $container->setVarfield($this->varfield);

    $container->setMessage('alertmessage', $this->alertmessage);
    $container->setMessage('servermessage', $this->servermessage);
    $container->setMessage('titleinsert', $this->titles[DomMask::INSERT]);
    $container->setMessage('titleupdate', $this->titles[DomMask::UPDATE]);
    $container->setMessage('titledelete', $this->titles[DomMask::DELETE]);
    $container->setMessage('titleview', $this->titles[DomMask::VIEW]);
    $container->setMessage('insertok', $this->actionmessages[DomMask::DOINSERT]);
    $container->setMessage('updateok', $this->actionmessages[DomMask::DOUPDATE]);
    $container->setMessage('deleteok', $this->actionmessages[DomMask::DODELETE]);
    if ($this->jsonsuccess)
      $container->setEvent('success', $this->jsonsuccess);
    if ($this->jsonfailure)
      $container->setEvent('failure', $this->jsonfailure);

    $controlzone = new \wajaf\groupZone();
    $controlzone->setType('control');

    // fill the zones with fields
    foreach ($this->fields as $K => $F)
    {
      $f = $F->create();
      if ($F->type == 'button')
      {
        $controlzone->add($f);
        continue;
      }
      // creates the zone
      $zone = new \wajaf\groupZone();
      $zone->setType('field');
      $container->add($zone);
      $zone->add($f);
    }

    if ($this->variables)
    {
      foreach ($this->variables as $K => $V)
      {
        $f = new \wajaf\hiddenfieldElement($K);
        $f->setData($V);
        $controlzone->add($f);
      }
    }
    $container->add($controlzone);

    // any template ?

    // creates the dataset if mode > 1 ? (just to save a hit and time)
    // creates the original dataset
    $ds = new \wajaf\groupDataset(json_encode(array($this->key => $this->getRecord($this->key))));
    $container->add($ds);

    return $container;
  }

  private function loadDefinition($data)
  {
    if ($data['def'])
    {
      foreach($data['def'] as $p => $v)
      {
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

    if (!isset($data['fields']))
      return; // throw new DomMaskError( \core\WAMessage::getMessage('dommask.baddescriptor'));

    foreach($data['fields'] as $n => $f)
    {
      if (!isset($f['type']))
        throw new \throwables\DomMaskError( \core\WAMessage::getMessage('dommask.baddescriptor'));
      $type = 'DomMask' . $f['type'] . 'Field';
      $F = new $type($n, isset($f['inrecord'])?$f['inrecord']:false);
      $F->loadDefinition($f);
      $this->addMaskfield($F);
    }
  }

  // can be replaced by user to have its own parameter sources and validations
  public function getParameter($V, $source = 'all')
  {
    // get the variable from the client, first check POST (PRIORITY) then GET
    // if the variable doesnt exists, returns NULL
    if (isset($_POST[$V]) && in_array($source, array('all', 'post')))
    {
      return $_POST[$V];
    }
    if (isset($_GET[$V]) && in_array($source, array('all', 'get')))
    {
      return $_GET[$V];
    }
    return null;
  }

} // class DomMask
*/
