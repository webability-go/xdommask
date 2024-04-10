package xdommask

import (
//	"github.com/webability-go/xamboo/cms/context"
//
// import "github.com/webability-go/wajaf"
)

type VideoField struct {
	*FileField
}

func NewVideoField(name string) *VideoField {
	vf := &VideoField{FileField: NewFileField(name)}
	vf.Type = FIELD
	vf.ExtensionsAuth = map[string]bool{"mp4": true, "mov": true}
	vf.MimesAuth = map[string]string{
		"video/mp4": "filemp4.png",
		"video/mov": "filemov.png",
	}
	return vf
}

/*
func (f *VideoField) Compile() wajaf.NodeDef {

	b := f.FileField.Compile()

	return b
}
*/
/*
class DomMaskFieldVideo extends DomMaskField
{
  public $AcceptExternal = false;  // true if third parties videos are accepted (will show the value field and it is editable), like facebook links, instagram, google, etc
  public $Path = '';
  public $DeleteButton = '[Delete]';

  public $ExtensionsAuth = array(".mp4", ".mov", ".avi", ".wmv");   // mandatory
  public $OriginDir = null;   // dir and path on where is the video preview (if needed)
  public $OriginPath = null;
  public $DestinationDir = null;
  public $DestinationPath = null;
  public $createName = null;  // function to call to create the name of the file
                              // params: ($record_key, $record_array, $field_new_value, $field_old_value)
                              // if insert mode,

  private $videoname = null;
  private $ext = null;
  private $newvideo = false;

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
    //print('|flag-3|');die;
    $this->videoname = $name;
    if (!$this->videoname)
    {
      if (isset($_POST[$this->urlvariable]))
      {
        $this->videoname = $_POST[$this->urlvariable];
      }
    }
    if (!$this->videoname)
      return;

    if (substr($this->videoname, 0, 6) == 'http:/' || substr($this->videoname, 0, 6) == 'https:')
      return;

    $this->ext = null;
    $this->newvideo = false;
    if (substr($this->videoname, 0, 5) == 'temp:')
    {
      $this->newvideo = true;
      $this->videoname = substr($this->videoname, 5);
    }
    foreach ($this->ExtensionsAuth as $ext)
    {
      if (strtolower(substr($this->videoname, -strlen($ext))) == $ext)
      {
        $this->ext = strtolower(substr($this->videoname, -strlen($ext)));
        break;
      }
    }
    if (!$this->ext)
      throw new \throwables\InsertError("The extension of the file is not recognized.");

    if (isset($_POST['idbrightcove'])) // custom file name
    {
      $this->videoname = $_POST['idbrightcove'] . $this->ext;
    }
  }

  public function preInsert($newrecord)
  {
    //print('|flag-4|');die;
    $this->checkExtension($newrecord->{$this->urlvariable});
    $newrecord->{$this->urlvariable} = $this->videoname;
  }

  public function preUpdate($key, $newrecord, $oldrecord)
  {
    //print('|flag-5|');die;
    $this->checkExtension($newrecord->{$this->urlvariable});
    $newrecord->{$this->urlvariable} = $this->videoname;
  }
  public function getext()
  {
    //print('|flag-6|');die;
    return $this->ext;
  }

  private function saveVideo($key, $newrecord, $newval, $oldval)
  {
    //print('|flag-7|');die;
    $newkey = $this->getKEY($key, $newrecord);

    if (!$this->ext) // no temporal new video
      return null;

    if ($this->newvideo)
    {
      if ($this->createName)
      {
        $fct = $this->createName;
        if ($this->caller)
          $File_name_store = $this->caller->$fct($newkey, $newrecord, $newval, $oldval, $this->ext);
        else
          $File_name_store = $fct($newkey, $newrecord, $newval, $oldval, $this->ext);
      }
      else
        $File_name_store = $newkey.$this->ext;

    // si hay temporal, moverlo
      $DPATH = str_replace('{KEY}', $newkey, $this->DestinationPath);
      if (isset($newrecord['DIRECTORY']))
        $DPATH = str_replace('{DIRECTORY}', $newrecord['DIRECTORY'], $DPATH);

      // print("DestinationPath");var_dump($this->DestinationPath);
      //  print("DPATH");var_dump($DPATH);
      //  die;
      // \core\WAFile::createDirectory($this->DestinationDir, $DPATH);

      $tmp = $this->OriginDir.$this->OriginPath.$this->videoname;

      $shellcommand = "sudo /home/sites/kiwi4.kiwilimon.com/application/shell/copyvideo.php -i $tmp -od {$DPATH}{$newkey}/ -of $File_name_store";
      // print("command: ");var_dump($shellcommand);die;
      print nl2br(shell_exec($shellcommand));

      // copy($tmp, $this->DestinationDir . $DPATH . $File_name_store);

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
    //print('|flag-8|');die;
    //$newkey = $this->getKEY($key, $newrecord);
    return $this->saveVideo($key, $newrecord, $newrecord->{$this->urlvariable}, null);
  }

  public function postUpdate($key, $newrecord, $oldrecord)
  {
    //print('|flag-9|');die;
    //$newkey = $this->getKEY($key, $newrecord);

    // print('|key|');var_dump($key);
    // print('|newrecord|');var_dump($newrecord);
    // print('|oldrecord|');var_dump($oldrecord);
    // die;

    $newname = $oldname = null;
    if ($newrecord instanceof \dominion\DB_Record)
      $newname = $newrecord->{$this->urlvariable};
    if ($oldrecord instanceof \dominion\DB_Record)
      $oldname = $oldrecord->{$this->urlvariable};

    if (!$newname && $oldname)
    {
      // delete file !
      $DPATH = str_replace('{KEY}', $newkey, $this->DestinationPath);
      if (isset($newrecord['DIRECTORY']))
        $DPATH = str_replace('{DIRECTORY}', $newrecord['DIRECTORY'], $DPATH);
      if (file_exists($this->DestinationDir . $DPATH . $oldname))
        unlink($this->DestinationDir . $DPATH . $oldname);
      $newnamerecord->{$this->urlvariable} = null;
    }
    return $this->saveVideo($key, $newrecord, $newname, $oldname);
  }

  public function postDelete($key, $oldrecord)
  {
    //print('|flag-10|');die;
    if ($oldrecord->{$this->urlvariable})
    {
      $newkey = $this->getKEY($key, $oldrecord);
      // delete file !
      $DPATH = str_replace('{KEY}', $newkey, $this->DestinationPath);
      if (isset($oldrecord['DIRECTORY']))
        $DPATH = str_replace('{DIRECTORY}', $oldrecord['DIRECTORY'], $DPATH);
      unlink($this->DestinationDir . $DPATH . $oldrecord->{$this->urlvariable});
    }
  }

  public function getKEY($key, $newrecord)
  {
    return (isset($newrecord['idbrightcove']) && !empty($newrecord['idbrightcove'])) ? $newrecord['idbrightcove']: $key;
  }

  //public function prepareVideo()
  public function prepareImage()
  {
    //return array('quecosatapasando' => 'nose');
    try
    {
      $tempname = $this->base->createKey(10);
      $this->checkExtension(strtolower($_FILES['images']['name'][0]));

      if ($this->ext)
      {
        // 1. save the video in a temporary public directory
        \core\WAFile::createDirectory($this->OriginDir, $this->OriginPath);
        move_uploaded_file($_FILES['images']['tmp_name'][0], $this->OriginDir . $this->OriginPath . $tempname . $this->ext);

        return array('status' => 'OK', 'tempname' => $this->OriginPath . $tempname . $this->ext, 'name' => $tempname . $this->ext);
      }
      else
      {
        return array('status' => 'error', 'message' => 'Error: el archivo que subi� no es un video.');
      }
    } catch (Exception $e)
    {
      return array('status' => 'error', 'message' => ''.$e);
    }
  }

}

*/
