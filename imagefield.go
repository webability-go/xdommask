package xdommask

import (
//	"github.com/webability-go/xamboo/cms/context"
//
// import "github.com/webability-go/wajaf"
)

type ImageField struct {
	*FileField
}

func NewImageField(name string) *ImageField {
	imf := &ImageField{FileField: NewFileField(name)}
	imf.Type = FIELD
	imf.ExtensionsAuth = map[string]bool{"jpg": true, "png": true}
	imf.MimesAuth = map[string]string{
		"image/jpeg": "filejpg.png",
		"image/png":  "filepng.png",
	}
	return imf
}

/*
func (f *ImageField) Compile() wajaf.NodeDef {

	b := f.FileField.Compile()

	return b
}
*/
/*
class DomMaskFieldImage extends DomMaskField
{
  public $AcceptExternal = false;  // true if third parties images are accepted (will show the value field and it is editable), like facebook links, instagram, google, etc
  public $Path = '';
  public $DeleteButton = '[Delete]';

  public $ExtensionsAuth = array(".gif", ".jpg", ".jpeg", ".png");   // mandatory
  public $OriginDir = null;   // dir and path on where is the image preview (if needed)
  public $OriginPath = null;
  public $DestinationDir = null;
  public $DestinationPath = null;
  public $createName = null;  // function to call to create the name of the file
                              // params: ($record_key, $record_array, $field_new_value, $field_old_value)
                              // if insert mode,

  private $imagename = null;
  private $ext = null;
  private $newimage = false;

  function __construct($name = '', $iftable = false)
  {
    parent::__construct($name, $iftable);
    $this->type = 'file';
  }

  public function create()
  {
    $f = new \wajaf\mmcfieldElement($this->name);

    $f->setVisible($this->DomMask->createModes($this->authmodes));
    $f->setInfo($this->DomMask->createModes($this->viewmodes));
    $f->setReadonly($this->DomMask->createModes($this->readonlymodes));
    $f->setNotnull($this->DomMask->createModes($this->notnullmodes));
    $f->setDisabled($this->DomMask->createModes($this->disabledmodes));
    $f->setHelpmode($this->DomMask->createModes($this->helpmodes));

    $f->setExternal($this->AcceptExternal);
    $f->setPath($this->Path);
    $f->setDeletebutton($this->DeleteButton);

    $f->setData($this->title);

    $f->setMessage('defaultvalue', $this->default);
    $f->setMessage('helpdescription', $this->helpdescription);
    $f->setMessage('statusnotnull', $this->statusnotnull);

    return $f;
  }

  private function checkExtension($name)
  {
    $this->imagename = $name;
    if (!$this->imagename)
    {
      if (isset($_POST[$this->urlvariable]))
      {
        $this->imagename = $_POST[$this->urlvariable];
      }
    }
    if (!$this->imagename)
      return;

    if (substr($this->imagename, 0, 6) == 'http:/' || substr($this->imagename, 0, 6) == 'https:')
      return;

    $this->ext = null;
    $this->newimage = false;
    if (substr($this->imagename, 0, 5) == 'temp:')
    {
      $this->newimage = true;
      $this->imagename = substr($this->imagename, 5);
    }
    foreach ($this->ExtensionsAuth as $ext)
    {
      if (strtolower(substr($this->imagename, -strlen($ext))) == $ext)
      {
        $this->ext = strtolower(substr($this->imagename, -strlen($ext)));
        break;
      }
    }
    if (!$this->ext)
      throw new \throwables\InsertError("The extension of the file is not recognized.");
  }

  public function preInsert($newrecord)
  {
    $this->checkExtension($newrecord->{$this->urlvariable});
    $newrecord->{$this->urlvariable} = $this->imagename;
  }

  public function preUpdate($key, $newrecord, $oldrecord)
  {
    $this->checkExtension($newrecord->{$this->urlvariable});
    $newrecord->{$this->urlvariable} = $this->imagename;
  }

  private function saveImage($key, $newrecord, $newval, $oldval)
  {
    if (!$this->ext) // no temporal new image
      return null;

    if ($this->newimage)
    {
      if ($this->createName)
      {
        $fct = $this->createName;
        if ($this->caller)
          $File_name_store = $this->caller->$fct($key, $newrecord, $newval, $oldval, $this->ext);
        else
          $File_name_store = $fct($key, $newrecord, $newval, $oldval, $this->ext);
      }
      else
        $File_name_store = $key.$this->ext;

    // si hay temporal, moverlo
      $DPATH = str_replace('{KEY}', $key, $this->DestinationPath);
      \core\WAFile::createDirectory($this->DestinationDir, $DPATH);

      $tmp = $this->OriginDir.$this->OriginPath.$this->imagename;
      copy($tmp, $this->DestinationDir . $DPATH . $File_name_store);

      if ($this->inrecord && $this->DomMask->table)
      {
        $rec = array($this->name => $File_name_store);
        $this->DomMask->table->doUpdate($key, $rec);
      }
      if ($this->inrecord)
      {
        $newrecord->{$this->name} = $File_name_store;
      }

      $newval = $File_name_store;
    }
    return $newval;
  }

  public function postInsert($key, $newrecord)
  {
    return $this->saveImage($key, $newrecord, $newrecord->{$this->urlvariable}, null);
  }

  public function postUpdate($key, $newrecord, $oldrecord)
  {
    $newname = $oldname = null;
    if ($newrecord instanceof \dominion\DB_Record)
      $newname = $newrecord->{$this->urlvariable};
    if ($oldrecord instanceof \dominion\DB_Record)
      $oldname = $oldrecord->{$this->urlvariable};

    if (!$newname && $oldname)
    {
      // delete file !
      $DPATH = str_replace('{KEY}', $key, $this->DestinationPath);
      if (file_exists($this->DestinationDir . $DPATH . $oldname))
        unlink($this->DestinationDir . $DPATH . $oldname);
      $newnamerecord->{$this->urlvariable} = null;
    }
    return $this->saveImage($key, $newrecord, $newname, $oldname);
  }

  public function postDelete($key, $oldrecord)
  {
    if ($oldrecord->{$this->urlvariable})
    {
      // delete file !
      $DPATH = str_replace('{KEY}', $key, $this->DestinationPath);
      unlink($this->DestinationDir . $DPATH . $oldrecord->{$this->urlvariable});
    }
  }

  public function prepareImage()
  {
    $tempname = $this->base->createKey(10);
    $this->checkExtension(strtolower($_FILES['images']['name'][0]));

    if ($this->ext)
    {
      // 1. save the image in a temporary public directory
      \core\WAFile::createDirectory($this->OriginDir, $this->OriginPath);
      move_uploaded_file($_FILES['images']['tmp_name'][0], $this->OriginDir . $this->OriginPath . $tempname . $this->ext);

      return array('status' => 'OK', 'tempname' => $this->OriginPath . $tempname . $this->ext, 'name' => $tempname . $this->ext);
    }
    else
    {
      return array('status' => 'error', 'message' => 'Error: el archivo que subi� no es una imagen.');
    }
  }

}

*/
