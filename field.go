package xdommask

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"

	"github.com/webability-go/wajaf"
	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/cms/context"
)

const (
	CONTROL = "control"
	FIELD   = "field"
	INFO    = "info"
	HIDDEN  = "hidden"
	GROUP   = "group"
)

type FieldDef interface {
	Compile() wajaf.NodeDef
	GetType() string
	GetName() string
	GetInRecord() bool
	GetValue(ctx *context.Context, mode Mode) (interface{}, bool, error)
	ConvertValue(value interface{}) (interface{}, error)
	PostGet(ctx *context.Context, key interface{}, rec *xdominion.XRecord) error
	PreInsert(ctx *context.Context, rec *xdominion.XRecord) error
	PostInsert(ctx *context.Context, key interface{}, rec *xdominion.XRecord) (bool, error) // bool = field changed true/false, need an update
	PreUpdate(ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, newrec *xdominion.XRecord) error
	PostUpdate(ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, newrec *xdominion.XRecord) (bool, error)
	PreDelete(ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, rec *xdominion.XRecord) error
	PostDelete(ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, rec *xdominion.XRecord) error
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

	CheckJS string
	Focus   string
	Blur    string
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

func (f *Field) GetName() string {
	return f.Name
}

func (f *Field) GetInRecord() bool {
	return false
}

func (f *Field) GetValue(ctx *context.Context, mode Mode) (interface{}, bool, error) {
	if DEBUG {
		fmt.Println("xdominion.Field::GetValue", mode, "field ignored by construct")
	}
	return nil, true, nil
}

// ConvertValue will convert the entry value to the correct type for this field
func (f *Field) ConvertValue(value interface{}) (interface{}, error) {
	if DEBUG {
		fmt.Println("xdominion.Field::ConvertValue", value)
	}
	return value, nil
}

func (f *Field) PreInsert(ctx *context.Context, rec *xdominion.XRecord) error {
	return nil
}

func (f *Field) PostGet(ctx *context.Context, key interface{}, rec *xdominion.XRecord) error {
	return nil
}

func (f *Field) PostInsert(ctx *context.Context, key interface{}, rec *xdominion.XRecord) (bool, error) {
	return false, nil
}

func (f *Field) PreUpdate(ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, newrec *xdominion.XRecord) error {
	return nil
}

func (f *Field) PostUpdate(ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, newrec *xdominion.XRecord) (bool, error) {
	return false, nil
}

func (f *Field) PreDelete(ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, rec *xdominion.XRecord) error {
	return nil
}

func (f *Field) PostDelete(ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, rec *xdominion.XRecord) error {
	return nil
}

type DataField struct {
	*Field
	InRecord    bool
	URLVariable string

	Title        string
	DefaultValue string

	Auto        bool
	AutoMessage string

	Encoded  bool
	Entities bool

	NullOnEmpty  bool
	NullValue    string
	MD5Encrypted bool

	StatusNotNull string
	StatusCheck   string
}

func NewDataField(name string) *DataField {
	return &DataField{Field: NewField(name),
		URLVariable: name,
	}
}

func (f *DataField) GetInRecord() bool {
	return f.InRecord
}

func (f *DataField) Compile() wajaf.NodeDef {
	return nil
}

func (f *DataField) PostGet(ctx *context.Context, key interface{}, rec *xdominion.XRecord) error {
	if f.MD5Encrypted {
		rec.Del(f.Name) // does not show the field on client side
	}
	return nil
}

// GetValue to get the value from the field when needed. return value, ignored bool (true = ignored by construct)
// return ONLY string or nil
func (f *DataField) GetValue(ctx *context.Context, mode Mode) (interface{}, bool, error) {

	if DEBUG {
		fmt.Println("xdominion.DataField::GetValue", f.Name, mode)
	}

	// CAN BE USED?
	//	if !f.InRecord {
	//		return nil, true, nil
	//	}
	if f.AuthModes&mode == 0 {
		return nil, true, nil
	}
	if f.ViewModes&mode != 0 {
		return nil, true, nil
	}

	// BUILD VALUE
	sval := ""
	if (mode&INSERT != 0) && f.Auto {
		sval = f.DefaultValue
	} else {
		sval = ctx.Request.Form.Get(f.URLVariable)
	}
	// NOT NULL ?
	if sval == "" && (f.NotNullModes&mode != 0) {
		return nil, false, errors.New(f.StatusNotNull)
	}

	if f.NullOnEmpty && sval == "" {
		return nil, false, nil
	}

	// FILTER VALUE
	if sval != "" {
		if f.Entities {
			sval = html.EscapeString(sval)
		}
		if f.Encoded {
			sval = strings.Replace(url.QueryEscape(sval), "+", "%20", -1)
		}
		if f.MD5Encrypted {
			data := []byte(sval)
			sval = fmt.Sprintf("%x", md5.Sum(data))
		}
	} else if f.MD5Encrypted { // ignore MD5 is there is no value to modify
		return nil, true, nil
	}
	return sval, false, nil
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

*/
