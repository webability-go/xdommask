package xdommask

import "github.com/webability-go/wajaf"

type ButtonField struct {
	*ControlField
	Action string
}

func NewButtonField(name string) *ButtonField {
	return &ButtonField{ControlField: NewControlField(name)}
}

func (f *ButtonField) Compile() wajaf.NodeDef {

	b := wajaf.NewButtonElement("")

	b.SetAttribute("style", f.Style)
	b.SetAttribute("classname", f.ClassName)
	b.SetData(f.Title)

	b.SetAttribute("action", f.Action)

	return b

	/*
		$title = is_string($this->title)?$this->title:'';
		$f = new \wajaf\buttonElement($title, $this->name);

		if (is_array($this->title))
		{
			if (isset($this->title[DomMask::INSERT]))
				$f->setMessage('titleinsert', $this->title[DomMask::INSERT]);
			if (isset($this->title[DomMask::UPDATE]))
				$f->setMessage('titleupdate', $this->title[DomMask::UPDATE]);
			if (isset($this->title[DomMask::DELETE]))
				$f->setMessage('titledelete', $this->title[DomMask::DELETE]);
			if (isset($this->title[DomMask::VIEW]))
				$f->setMessage('titleview', $this->title[DomMask::VIEW]);
		}

		$f->setVisible($this->DomMask->createModes($this->authmodes));
		$f->setAction($this->action);

		return $f;
	*/
}

/*
class DomMaskButtonField extends DomMaskField
{
  public $action = 'submit';
  public $OnClick = 'reset();';          // for ButtonFields
  public $ButtonFieldInsert = null;      // string
  public $ButtonFieldUpdate = null;      // string
  public $ButtonFieldDelete = null;      // string
  public $ButtonFieldView = null;        // string
  public $ButtonFieldAsImage = null;     // string link of image
  public $OnEvent = null;

  function __construct($name = '')
  {
    parent::__construct($name, false);
    $this->type = 'button';
  }

  public function getAction()
  {
    return $this->action;
  }

  public function create()
  {
    $title = is_string($this->title)?$this->title:'';
    $f = new \wajaf\buttonElement($title, $this->name);

    if (is_array($this->title))
    {
      if (isset($this->title[DomMask::INSERT]))
        $f->setMessage('titleinsert', $this->title[DomMask::INSERT]);
      if (isset($this->title[DomMask::UPDATE]))
        $f->setMessage('titleupdate', $this->title[DomMask::UPDATE]);
      if (isset($this->title[DomMask::DELETE]))
        $f->setMessage('titledelete', $this->title[DomMask::DELETE]);
      if (isset($this->title[DomMask::VIEW]))
        $f->setMessage('titleview', $this->title[DomMask::VIEW]);
    }

    $f->setVisible($this->DomMask->createModes($this->authmodes));
    $f->setAction($this->action);

    return $f;
  }
}
*/
