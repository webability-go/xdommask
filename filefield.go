package xdommask

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/webability-go/wajaf"
	"github.com/webability-go/xamboo/cms/context"
	"github.com/webability-go/xdominion"
)

type Ifilefield interface {
	PrepareFile(ctx *context.Context) error
	GetOriginalFileName() string
	GetTemporalFileName() string
	GetFinalFileName() string
	GetMime() string
	GetIconName() string
	GetTemporalImage() string
	GetMultiFile() bool
	GetOriginalFileNames() []string
	GetTemporalFileNames() []string
	GetIconNames() []string
	GetTemporalImages() []string
}

type FileField struct {
	*DataField

	AcceptExternal bool // true if third parties images are accepted (will show the value field and it is editable), like facebook links, instagram, google, etc
	DeleteButton   string
	Loading        string
	MultiFile      bool // true accept multifiles
	MaxSize        int64

	ExtensionsAuth  map[string]bool
	MimesAuth       map[string]string // format:   mime => icon. If icon == "*" will use the file image itself (if it's an image authorized, png, gif, jpeg only)
	OriginDir       string
	OriginPath      string
	DestinationDir  string
	DestinationPath string
	DestinationName string

	TemporalFileName  string   // Set only when a new upload happens
	TemporalFileNames []string // when it's multifile
	OriginalFileName  string   // Set only when a new upload happens
	OriginalFileNames []string // Set only when a new upload happens
	FinalFileName     string   // Official name of file, not modified if no upload, or temporal value if new upload
	IconPath          string
	IconName          string
	IconNames         []string
	Ext               string
	Mime              string
	Image             string   // only if filename is an official image (jpeg, png)
	Images            []string // only if filename is an official image (jpeg, png)

	changed bool
}

func NewFileField(name string) *FileField {
	ff := &FileField{DataField: NewDataField(name)}
	ff.Type = FIELD
	ff.DeleteButton = "[Delete file]"
	ff.ExtensionsAuth = map[string]bool{"txt": true, "csv": true, "json": true, "xml": true, "pdf": true, "xls": true, "doc": true, "ppt": true, "xlsx": true, "docx": true, "pptx": true}
	ff.MimesAuth = map[string]string{
		"text/plain":                    "filetext.png",
		"text/csv":                      "filecsv.png",
		"application/json":              "filejson.png",
		"application/xml":               "filexml.png",
		"application/pdf":               "filepdf.png",
		"application/vnd.ms-excel":      "filexls.png",
		"application/msword":            "filedoc.png",
		"application/vnd.ms-powerpoint": "fileppt.png",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         "filexls.png",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   "filedoc.png",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": "fileppt.png",
		"default": "filefile.png",
	}
	return ff
}

func (f *FileField) GetFinalFileName() string {
	return f.FinalFileName
}

func (f *FileField) GetTemporalFileName() string {
	return f.TemporalFileName
}

func (f *FileField) GetOriginalFileName() string {
	return f.OriginalFileName
}

func (f *FileField) GetTemporalImage() string {
	return f.OriginPath + f.TemporalFileName
}

func (f *FileField) GetIconName() string {
	return f.IconName
}

func (f *FileField) GetMime() string {
	return f.Mime
}

func (f *FileField) GetMultiFile() bool {
	return f.MultiFile
}

func (f *FileField) GetTemporalFileNames() []string {
	return f.TemporalFileNames
}

func (f *FileField) GetOriginalFileNames() []string {
	return f.OriginalFileNames
}

func (f *FileField) GetIconNames() []string {
	return f.IconNames
}

func (f *FileField) GetTemporalImages() []string {
	tmp := []string{}
	for _, v := range f.TemporalFileNames {
		tmp = append(tmp, f.OriginPath+v)
	}
	return tmp
}

func (f *FileField) Compile() wajaf.NodeDef {

	b := wajaf.NewMMCFieldElement(f.ID)

	b.SetAttribute("style", f.Style)
	b.SetAttribute("classname", f.ClassName)
	b.SetAttribute("defaultvalue", fmt.Sprint(f.DefaultValue))
	b.SetData(f.Title)

	b.SetAttribute("size", f.Size)
	mf := "no"
	if f.MultiFile {
		mf = "yes"
	}
	b.SetAttribute("multifile", mf)

	b.SetAttribute("visible", convertModes(f.AuthModes))
	b.SetAttribute("info", convertModes(f.ViewModes))
	b.SetAttribute("readonly", convertModes(f.ReadOnlyModes))
	b.SetAttribute("notnull", convertModes(f.NotNullModes))
	b.SetAttribute("disabled", convertModes(f.DisabledModes))
	b.SetAttribute("helpmode", convertModes(f.HelpModes))

	b.AddHelp("", "", f.HelpDescription)
	b.AddMessage("defaultvalue", fmt.Sprint(f.DefaultValue))
	b.AddMessage("automessage", f.AutoMessage)
	b.AddMessage("statusnotnull", f.StatusNotNull)

	b.SetAttribute("deletebutton", f.DeleteButton)
	b.SetAttribute("external", fmt.Sprint(f.AcceptExternal))
	b.SetAttribute("loading", f.Loading)
	// create "accept"
	if len(f.MimesAuth) > 0 {
		acc := ""
		for t := range f.MimesAuth {
			if acc != "" {
				acc += ","
			}
			acc += t
		}
		b.SetAttribute("accept", acc)
	}

	return b
}

func (f *FileField) PostGet(ctx *context.Context, key interface{}, rec *xdominion.XRecord) error {

	fmt.Println("PostGet filefield", f.ID, key)
	v, _ := rec.GetString(f.ID)
	if v == "" {
		return nil
	}
	// check validity of the file, mime type on HD
	nv := map[string]string{
		"value":    v,
		"filename": v,
		"iconname": "",
		"mime":     "",
	}
	localFileName := f.buildPath(f.DestinationDir+f.DestinationPath+v, rec)
	// load mime
	file, err := os.Open(localFileName)
	if err != nil {
		// notify error in message
		nv["value"] = "Error, file not accessible: " + err.Error()
		nv["iconname"] = f.IconPath + "fileerror.png"
		rec.Set(f.ID, nv)
		return nil
	}
	defer file.Close()
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		nv["value"] = "Error, file not accessible: " + err.Error()
		nv["iconname"] = f.IconPath + "fileerror.png"
		rec.Set(f.ID, nv)
		return nil
	}
	mime := http.DetectContentType(buff)
	f.IconName = f.MimesAuth[mime]
	if f.IconName == "" {
		f.IconName = f.MimesAuth["default"]
	}
	nv["iconname"] = f.IconPath + f.IconName
	nv["mime"] = mime
	if mime == "image/jpeg" || mime == "image/png" {
		nv["image"] = f.buildPath(f.Image, rec)
	}

	rec.Set(f.ID, nv)
	return nil
}

func (f *FileField) CreateTempName() string {
	rand.Seed(time.Now().UnixNano())
	result := ""
	for i := 0; i < 10; i++ {
		key := 97 + rand.Intn(26)
		result += string(key)
	}
	return result
}

func (f *FileField) CheckExtension(name string, mime string) error {

	if name == "" {
		return nil
	}

	// 1. separate the extension based on last "."
	xnames := strings.Split(name, ".")
	ext := xnames[len(xnames)-1]
	// 2. verify auth ext
	if ext == "" {
		return errors.New("Error, no extension")
	}
	if !f.ExtensionsAuth[ext] {
		return errors.New("Error, extension not authorized")
	}
	if mime != "" {
		// 3. verify auth mime
		if f.MimesAuth[mime] == "" {
			if f.MimesAuth["default"] == "" {
				return errors.New("Error, file mime not authorized")
			}
		}
	}
	// 4. set vars
	f.Ext = ext
	f.Mime = mime
	f.IconName = ""
	if f.MimesAuth[mime] != "" {
		f.IconName = f.IconPath + f.MimesAuth[mime]
	} else if f.MimesAuth["default"] != "" {
		f.IconName = f.IconPath + f.MimesAuth["default"]
	}
	if f.MultiFile {
		f.IconNames = append(f.IconNames, f.IconName)
	}

	return nil
}

func (f *FileField) PrepareFile(ctx *context.Context) error {

	f.TemporalFileNames = []string{}
	f.OriginalFileNames = []string{}
	f.IconNames = []string{}

	files := ctx.Request.MultipartForm.File["file"]
	for _, fileHeader := range files {
		if fileHeader.Size > f.MaxSize {
			fmt.Println("ERROR PREPARE FILE 1: TOO BIG", fileHeader.Filename, fileHeader.Size)
			return nil
		}

		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			fmt.Println("ERROR PREPARE FILE 2: COULD NOT OPEN", err)
			return err
		}

		defer file.Close()
		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			fmt.Println("ERROR PREPARE FILE 3: COULD NOT READ 512 HEADER BUFFER", err)
			return err
		}
		filetype := http.DetectContentType(buff)

		err = f.CheckExtension(fileHeader.Filename, filetype)
		if err != nil {
			fmt.Println("ERROR PREPARE FILE 4: NOT A CORRECT FILE TYPE", err)
			return err
		}
		f.OriginalFileName = fileHeader.Filename
		if f.MultiFile {
			f.OriginalFileNames = append(f.OriginalFileNames, fileHeader.Filename)
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			fmt.Println("ERROR PREPARE FILE 5:", err)
			return err
		}

		tempname := f.CreateTempName()
		f.TemporalFileName = tempname + "." + f.Ext
		if f.MultiFile {
			f.TemporalFileNames = append(f.TemporalFileNames, tempname+"."+f.Ext)
		}
		localFileName := f.OriginDir + f.OriginPath + f.TemporalFileName

		if _, err = os.Stat(f.OriginDir); err != nil {
			fmt.Println("ERROR PREPARE FILE ORIGINDIR:", err)
			return err
		}

		err = os.MkdirAll(f.OriginDir+f.OriginPath, 0771)
		if err != nil {
			fmt.Println("ERROR PREPARE FILE ORIGINPATH:", err)
			return err
		}

		out, err := os.Create(localFileName)
		if err != nil {
			fmt.Printf("ERROR PREPARE FILE 6: failed to open the file %s for writing", localFileName)
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			fmt.Printf("ERROR PREPARE FILE 7: copy file err:%s\n", err)
			return err
		}
		fmt.Printf("file %s uploaded to %s ok\n", fileHeader.Filename, localFileName)
		if !f.MultiFile { // Only one accepted if not multifile
			break
		}
	}

	return nil
}

func (f *FileField) PreInsert(ctx *context.Context, rec *xdominion.XRecord) error {

	fmt.Println("FileField::PreInsert")

	// Keep the temporal name in memory for Post (copy the file)
	name := ctx.Request.Form.Get(f.URLVariable + "[temporal]")
	if name == "" {
		// if NOT NULL => error
		return nil
	}
	f.TemporalFileName = name
	f.OriginalFileName = ctx.Request.Form.Get(f.URLVariable + "[filename]")
	f.changed = true

	// open temporal, verify temporal, get MIME
	// ******************

	// Verify extension, mime
	f.CheckExtension(name, "")

	rec.Set(f.URLVariable, name)
	return nil
}

func (f *FileField) PreUpdate(ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, newrec *xdominion.XRecord) error {

	// Keep the temporal name in memory for Post (copy the file)
	f.changed = true
	name := ctx.Request.Form.Get(f.URLVariable + "[temporal]")
	if name == "" {
		f.changed = false
		// do we have an old value?
		name = ctx.Request.Form.Get(f.URLVariable + "[filename]")
		newrec.Set(f.Name, name)

		// if NOT NULL => error
		return nil
	}
	f.TemporalFileName = name
	f.OriginalFileName = ctx.Request.Form.Get(f.URLVariable + "[filename]")

	// open temporal, verify temporal, get MIME
	// ******************

	// Verify extension, mime
	f.CheckExtension(name, "")

	newrec.Set(f.Name, name)
	return nil
}

func (f *FileField) SaveImage(key interface{}, rec *xdominion.XRecord) error {

	// Create the name of the file for DB (need the key of the record), and then update
	realname := f.buildPath(f.DestinationName, rec)

	fmt.Println("FileField::SaveImage::realname", realname)

	originpath := f.OriginDir + f.OriginPath + f.TemporalFileName
	realpath := f.buildPath(f.DestinationDir+f.DestinationPath, rec) + realname
	fmt.Println("FileField::SaveImage::paths", originpath, realpath)

	_, err := os.Stat(f.buildPath(f.DestinationDir, rec))
	if err != nil {
		fmt.Println("ERROR SAVEIMAGE DESTINATIONDIR:", err)
		return err
	}

	err = os.MkdirAll(f.buildPath(f.DestinationDir+f.DestinationPath, rec), 0771)
	if err != nil {
		fmt.Println("ERROR SAVEIMAGE DESTINATIONPATH:", err)
		return err
	}

	// We copy the remove origin, both files may be in different filesystems (we cannot use Rename)
	sourceFile, err := os.Open(originpath)
	if err != nil {
		fmt.Println("ERROR SAVEIMAGE OPEN ORIGIN:", err)
		return err
	}
	defer sourceFile.Close()

	newFile, err := os.Create(realpath)
	if err != nil {
		fmt.Println("ERROR SAVEIMAGE OPEN DESTINATION:", err)
		return err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, sourceFile)
	if err != nil {
		fmt.Println("ERROR SAVEIMAGE COPY FILE:", err)
		return err
	}

	err = os.Remove(originpath)
	if err != nil {
		fmt.Println("ERROR SAVEIMAGE REMOVE ORIGIN FILE:", err)
		return err
	}

	// save file in DB, we need the mask
	rec.Set(f.Name, realname)

	/*
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
	*/
	return nil
}

func (f *FileField) PostInsert(ctx *context.Context, key interface{}, rec *xdominion.XRecord) (bool, error) {

	if !f.changed {
		return false, nil
	}

	err := f.SaveImage(key, rec)

	return true, err
}

func (f *FileField) PostUpdate(ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, newrec *xdominion.XRecord) (bool, error) {

	if !f.changed { // no temporal file, verify if we delete something
		oldv, _ := oldrec.Get(f.ID)
		v, _ := newrec.Get(f.ID)
		fmt.Printf("FileField::PostUpdate %#v %#v \n", oldv, v)
		if v == oldv {
			// Nothing to do, same values
			return false, nil
		}
		if v == nil { // no image
			// verify NOT NULL

			if oldv != nil {
				// DELETE old images

			}
			return false, nil
		}
	}

	err := f.SaveImage(key, newrec)
	return true, err

	/*
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
	*/

}

func (f *FileField) PostDelete(ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, rec *xdominion.XRecord) error {

	fmt.Println("FileField::PostDelete")
	/*
	   public function postDelete($key, $oldrecord)
	   {
	     if ($oldrecord->{$this->urlvariable})
	     {
	       // delete file !
	       $DPATH = str_replace('{KEY}', $key, $this->DestinationPath);
	       unlink($this->DestinationDir . $DPATH . $oldrecord->{$this->urlvariable});
	     }
	   }
	*/
	return nil
}

func (f *FileField) buildPath(path string, rec *xdominion.XRecord) string {

	path = strings.ReplaceAll(path, "{_ext}", f.Ext)

	// regexp {} to find fields
	code := `(\{)(.*?)\}` // index based 2, One URL param, index-1 based
	codex := regexp.MustCompile(code)
	matches := codex.FindAllStringSubmatch(path, -1)

	for _, x := range matches {
		v, _ := rec.GetString(x[2])
		path = strings.ReplaceAll(path, x[0], v)
	}

	return path
}

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
