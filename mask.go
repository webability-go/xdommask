package xdommask

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/webability-go/wajaf"
	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/cms/context"
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

type MaskHooks struct {
	Build func(m *Mask, ctx *context.Context) error

	GetRecord func(m *Mask, ctx *context.Context, key interface{}, mode int) (string, *xdominion.XRecord, error)

	PreInsert  func(m *Mask, ctx *context.Context, rec *xdominion.XRecord) error
	Insert     func(m *Mask, ctx *context.Context, rec *xdominion.XRecord) (interface{}, error)
	PostInsert func(m *Mask, ctx *context.Context, key interface{}, rec *xdominion.XRecord) error

	PreUpdate  func(m *Mask, ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, newrec *xdominion.XRecord) error
	Update     func(m *Mask, ctx *context.Context, key interface{}, newrec *xdominion.XRecord) error
	PostUpdate func(m *Mask, ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, newrec *xdominion.XRecord) error

	PreDelete  func(m *Mask, ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, rec *xdominion.XRecord) error
	Delete     func(m *Mask, ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, rec *xdominion.XRecord) error
	PostDelete func(m *Mask, ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, rec *xdominion.XRecord) error
}

type Mask struct {
	Hooks MaskHooks

	ID string

	Display       string
	Style         string
	ClassName     string
	ClassNameZone string
	VarMode       string
	VarOrder      string
	VarKey        string
	VarField      string
	Maingroup     string
	Template      string

	Mode        Mode
	AuthModes   Mode
	KeyField    string
	Key         interface{}
	InsertedKey interface{}

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

	Table      *xdominion.XTable
	Order      *xdominion.XOrder
	Conditions *xdominion.XConditions
	FieldSet   *xdominion.XFieldSet
}

func NewMask(id string, hooks MaskHooks, ctx *context.Context) (*Mask, error) {
	if DEBUG {
		fmt.Println("xdominion.NewMask", id)
	}
	m := &Mask{
		ID:        id,
		Hooks:     hooks,
		Variables: map[string]string{},
		Fields:    []FieldDef{},
		VarMode:   "Mode",
		VarOrder:  "Order",
		VarKey:    "Key",
		VarField:  "groupfield",
		Counter:   1,
	}
	var err error
	if m.Hooks.Build != nil {
		err = m.Hooks.Build(m, ctx)
	}
	return m, err
}

func (m *Mask) AddField(f FieldDef) {
	if DEBUG {
		fmt.Println("xdominion.Mask::AddField", f.GetName())
	}
	m.Fields = append(m.Fields, f)
}

func (m *Mask) Compile(mode string, ctx *context.Context) wajaf.NodeDef {

	if DEBUG {
		fmt.Println("xdominion.Mask::Compile", mode)
	}
	mode = verifyMode(mode)
	group := wajaf.NewGroupContainer(m.ID)
	group.SetAttribute("display", m.Display)
	group.SetAttribute("style", m.Style)
	group.SetAttribute("classname", m.ClassName)
	group.SetAttribute("classnamezone", m.ClassNameZone)
	group.SetAttribute("haslistener", "yes")
	group.SetAttribute("maingroup", m.Maingroup)
	if m.Template != "" {
		group.AddMessage("template", m.Template)
	}

	group.SetAttribute("varmode", m.VarMode)
	group.SetAttribute("varorder", m.VarOrder)
	group.SetAttribute("varkey", m.VarKey)
	group.SetAttribute("varfield", m.VarField)

	group.SetAttribute("authmodes", convertModes(m.AuthModes))
	group.SetAttribute("mode", mode)
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
		if f.GetType() == CONTROL || f.GetType() == INFO {
			zcontrol.AddChild(f.Compile())
			continue
		}
		z := wajaf.NewGroupZone(f.GetType(), f.GetName()+"_zone")
		z.AddChild(f.Compile())
		group.AddChild(z)
	}
	group.AddChild(zcontrol)

	zhidden := wajaf.NewGroupZone("hidden", "")

	for id, val := range m.Variables {
		h := wajaf.NewHiddenFieldElement(id)
		h.SetData(val)
		zhidden.AddChild(h)
	}
	group.AddChild(zhidden)

	// Original dataset
	var rec *xdominion.XRecord
	if m.Hooks.GetRecord != nil {
		_, rec, _ = m.Hooks.GetRecord(m, ctx, m.Key, 0)
	} else {
		_, rec, _ = m.getrecord(m.Key, 0)
	}
	if rec != nil {
		// rec must be a JSON
		jsonrec, _ := json.Marshal(rec)
		zdata := wajaf.NewGroupDataset("", string(jsonrec))
		group.AddChild(zdata)
	}
	return group
}

func builRunMessage(data map[string]interface{}, messages map[string]string, err error) (map[string]interface{}, error) {
	if err != nil {
		if messages["text"] != "" {
			messages["text"] += "<br />\n"
		}
		messages["text"] += err.Error()
	}
	data["success"] = false
	data["message"] = messages
	return data, err
}

func (m *Mask) Run(ctx *context.Context) (map[string]interface{}, error) {

	order := ctx.Request.Form.Get(m.VarOrder)
	if DEBUG {
		fmt.Println("xdominion.Mask::Run", order, m.VarOrder, ctx.Request.Form)
	}

	messages := map[string]string{"text": ""}
	data := map[string]interface{}{}
	recdata := map[string]interface{}{}
	var err error

	nkey := ""
	var rec *xdominion.XRecord
	switch order {
	case "getrecord":
		key := m.convertPrimaryKey(ctx.Request.Form.Get(m.VarKey))
		fmt.Printf("GetRecord %#v\n", key)
		if m.Hooks.GetRecord != nil {
			nkey, rec, err = m.Hooks.GetRecord(m, ctx, key, 0)
			if err != nil {
				return builRunMessage(data, messages, err)
			}
		} else {
			nkey, rec, err = m.getrecord(key, 0)
			if err != nil {
				return builRunMessage(data, messages, err)
			}
		}
		if rec != nil {
			err = m.PostGet(ctx, key, rec)
			if err != nil {
				return builRunMessage(data, messages, err)
			}
			recdata[nkey] = rec
			data["data"] = recdata
		}
	case "first":
		if m.Hooks.GetRecord != nil {
			nkey, rec, err = m.Hooks.GetRecord(m, ctx, "", -2)
			if err != nil {
				return builRunMessage(data, messages, err)
			}
		} else {
			nkey, rec, err = m.getrecord("", -2)
			if err != nil {
				return builRunMessage(data, messages, err)
			}
		}
		if rec != nil {
			err = m.PostGet(ctx, nkey, rec)
			if err != nil {
				return builRunMessage(data, messages, err)
			}
			recdata[nkey] = rec
			data["data"] = recdata
		}
	case "previous":
		key := m.convertPrimaryKey(ctx.Request.Form.Get(m.VarKey))
		if m.Hooks.GetRecord != nil {
			nkey, rec, err = m.Hooks.GetRecord(m, ctx, key, -1)
			if err != nil {
				return builRunMessage(data, messages, err)
			}
		} else {
			nkey, rec, err = m.getrecord(key, -1)
			if err != nil {
				return builRunMessage(data, messages, err)
			}
		}
		if rec != nil {
			err = m.PostGet(ctx, nkey, rec)
			if err != nil {
				return builRunMessage(data, messages, err)
			}
			recdata[nkey] = rec
			data["data"] = recdata
		}
	case "next":
		key := m.convertPrimaryKey(ctx.Request.Form.Get(m.VarKey))
		if m.Hooks.GetRecord != nil {
			nkey, rec, err = m.Hooks.GetRecord(m, ctx, key, 1)
			if err != nil {
				return builRunMessage(data, messages, err)
			}
		} else {
			nkey, rec, err = m.getrecord(key, 1)
			if err != nil {
				return builRunMessage(data, messages, err)
			}
		}
		if rec != nil {
			err = m.PostGet(ctx, nkey, rec)
			if err != nil {
				return builRunMessage(data, messages, err)
			}
			recdata[nkey] = rec
			data["data"] = recdata
		}
	case "last":
		if m.Hooks.GetRecord != nil {
			nkey, rec, err = m.Hooks.GetRecord(m, ctx, "", 2)
			if err != nil {
				return builRunMessage(data, messages, err)
			}
		} else {
			nkey, rec, err = m.getrecord("", 2)
			if err != nil {
				return builRunMessage(data, messages, err)
			}
		}
		if rec != nil {
			err = m.PostGet(ctx, nkey, rec)
			if err != nil {
				return builRunMessage(data, messages, err)
			}
			recdata[nkey] = rec
			data["data"] = recdata
		}
	case "submit":
		mode := ctx.Request.Form.Get(m.VarMode)
		realmode, formmode := createMode(mode)

		key := m.convertPrimaryKey(ctx.Request.Form.Get(m.VarKey))
		if DEBUG {
			fmt.Println("xdominion.Mask::Run@submit", key, realmode, formmode)
		}
		fieldmessages := map[string]string{}
		rec, fieldmessages, err = m.execute(ctx, realmode, formmode, key)
		if err != nil {
			return builRunMessage(data, messages, err)
		}
		data["message"] = messages
		data["messages"] = fieldmessages
		data["rec"] = rec
		/*
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

		   S//*************** NOTE: SHOULD GETS BACK THE FINAL RECORD FOR IF THERE ARE SOME CALCULATED FIELDS (LIKE IMAGES NAMES; SUM FIELDS; LINKS, PATHS, ETC)

		*/

	case "file":

		key := m.convertPrimaryKey(ctx.Request.Form.Get(m.VarKey))
		field := ctx.Request.Form.Get(m.VarField)
		if DEBUG {
			fmt.Println("xdominion.Mask::Run@file", key, m.VarField, field)
		}
		for _, f := range m.Fields {
			if f.GetName() == field {
				var nf Ifilefield
				var ok bool
				if nf, ok = f.(*FileField); !ok {
					if nf, ok = f.(*ImageField); !ok {
						if nf, ok = f.(*VideoField); !ok {
							// ERROR GRAVE
						}
					}
				}
				err := nf.PrepareFile(ctx)
				if DEBUG {
					fmt.Println("xdominion.Mask::Run@file::field->PrepareFile", err)
				}
				if err != nil {
					data["message"] = err.Error()
					data["success"] = false
					return data, nil
				}
				// get back original image name, temporal name, icon name
				data["temporal"] = nf.GetTemporalFileName()
				data["filename"] = nf.GetOriginalFileName()
				data["iconname"] = nf.GetIconName()
			}
		}
	}
	data["success"] = true
	return data, nil
}

func (m *Mask) RunMultipart(ctx *context.Context) (map[string]interface{}, error) {

	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 1024*1024*1024) // 1 GB
	ctx.Request.ParseMultipartForm(1024 * 1024)

	// Read the form dinamically
	data := map[string]interface{}{}

	key := m.convertPrimaryKey(ctx.Request.Form.Get(m.VarKey))
	field := ctx.Request.Form.Get(m.VarField)
	if DEBUG {
		fmt.Println("xdominion.Mask::Run@file", key, m.VarField, field)
	}
	for _, f := range m.Fields {
		if f.GetName() == field {
			var nf Ifilefield
			var ok bool
			if nf, ok = f.(*FileField); !ok {
				if nf, ok = f.(*ImageField); !ok {
					if nf, ok = f.(*VideoField); !ok {
						// ERROR GRAVE
					}
				}
			}
			err := nf.PrepareFile(ctx)
			if DEBUG {
				fmt.Println("xdominion.Mask::Run@file::field->PrepareFile", err)
			}
			if err != nil {
				data["message"] = err.Error()
				data["success"] = false
				return data, nil
			}
			// get back original image name, temporal name, icon name
			if !nf.GetMultiFile() {
				data["temporal"] = nf.GetTemporalFileName()
				data["filename"] = nf.GetOriginalFileName()
				data["iconname"] = nf.GetIconName()
			} else {
				data["temporal"] = nf.GetTemporalFileNames()
				data["filename"] = nf.GetOriginalFileNames()
				data["iconname"] = nf.GetIconNames()
			}
			mime := nf.GetMime()
			if len(mime) > 6 && mime[0:6] == "image/" {
				if !nf.GetMultiFile() {
					data["image"] = nf.GetTemporalImage()
				} else {
					data["image"] = nf.GetTemporalImages()
				}
			}
			data["mime"] = mime
		}
	}
	data["success"] = true
	return data, nil
}

func (m *Mask) PostGet(ctx *context.Context, key interface{}, rec *xdominion.XRecord) error {
	for _, f := range m.Fields {
		err := f.PostGet(ctx, key, rec)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Mask) getrecord(key interface{}, mode int) (string, *xdominion.XRecord, error) {

	if m.Table == nil {
		return "", nil, nil
	}
	primkey := m.Table.GetPrimaryKey()
	primkeyname := primkey.GetName()
	// build QUERY
	if mode == 0 {
		rec, err := m.Table.SelectOne(key, m.FieldSet)
		rkey := ""
		if rec != nil {
			rkey, _ = rec.GetString(primkeyname)
		}
		return rkey, rec, err
	}
	if mode == -2 {
		rec, err := m.Table.SelectOne(m.Conditions, m.Order, m.FieldSet)
		rkey := ""
		if rec != nil {
			rkey, _ = rec.GetString(primkeyname)
		}
		return rkey, rec, err
	}
	if mode == 2 {
		// revert order
		var neworder xdominion.XOrder
		if m.Order != nil {
			neworder = revertOrder(*m.Order)
		}
		rec, err := m.Table.SelectOne(m.Conditions, neworder, m.FieldSet)
		rkey := ""
		if rec != nil {
			rkey, _ = rec.GetString(primkeyname)
		}
		return rkey, rec, err
	}
	if mode == -1 {
		// revert order
		var neworder xdominion.XOrder
		if m.Order != nil {
			neworder = revertOrder(*m.Order)
		}
		var newconditions xdominion.XConditions
		if m.Conditions != nil {
			newconditions = m.Conditions.Clone()
			if len(newconditions) > 0 {
				newconditions = append(newconditions, xdominion.NewXCondition(primkeyname, "<", key, "and"))
			} else {
				newconditions = append(newconditions, xdominion.NewXCondition(primkeyname, "<", key))
			}
		} else {
			newconditions = xdominion.XConditions{xdominion.NewXCondition(primkeyname, "<", key)}
		}
		rec, err := m.Table.SelectOne(newconditions, neworder, m.FieldSet)
		rkey := ""
		if rec != nil {
			rkey, _ = rec.GetString(primkeyname)
		}
		return rkey, rec, err
	}
	if mode == 1 {
		var newconditions xdominion.XConditions
		if m.Conditions != nil {
			newconditions = m.Conditions.Clone()
			if len(newconditions) > 0 {
				newconditions = append(newconditions, xdominion.NewXCondition(primkeyname, ">", key, "and"))
			} else {
				newconditions = append(newconditions, xdominion.NewXCondition(primkeyname, ">", key))
			}
		} else {
			newconditions = xdominion.XConditions{xdominion.NewXCondition(primkeyname, ">", key)}
		}
		rec, err := m.Table.SelectOne(newconditions, m.Order, m.FieldSet)
		rkey := ""
		if rec != nil {
			rkey, _ = rec.GetString(primkeyname)
		}
		return rkey, rec, err
	}
	return "", nil, nil
}

func (m *Mask) execute(ctx *context.Context, realmode Mode, formmode Mode, key interface{}) (*xdominion.XRecord, map[string]string, error) {

	messages := map[string]string{}
	record := m.buildrecord(ctx, formmode)
	switch realmode {
	case DOINSERT:
		if DEBUG {
			fmt.Println("xdominion.Mask::execute-->DOINSERT", realmode, formmode, key, record)
		}
		var err error
		if m.Hooks.PreInsert != nil {
			err = m.Hooks.PreInsert(m, ctx, record)
			if err != nil {
				return nil, messages, err
			}
		}
		isok := true
		for _, f := range m.Fields {
			err = f.PreInsert(ctx, record)
			if err != nil {
				isok = false
				messages[f.GetName()] = err.Error()
			}
		}
		if !isok {
			return nil, messages, errors.New("Error on pre fields. Please check the form.")
		}
		var key interface{}
		if m.Hooks.Insert != nil {
			key, err = m.Hooks.Insert(m, ctx, record)
			if err != nil {
				return nil, messages, err
			}
		} else {
			key, err = m.insert(ctx, record)
			if err != nil {
				return nil, messages, err
			}
		}
		// assign key to prim key field !!!
		primkey := m.Table.GetPrimaryKey()
		primkeyname := primkey.GetName()
		record.Set(primkeyname, key)

		isok = true
		changedrecord := &xdominion.XRecord{}
		haschanged := false
		for _, f := range m.Fields {
			changed, err := f.PostInsert(ctx, key, record)
			if err != nil {
				isok = false
				messages[f.GetName()] = err.Error()
			}
			if changed && f.GetInRecord() {
				haschanged = true
				changedvalue, _ := record.Get(f.GetName())
				changedrecord.Set(f.GetName(), changedvalue)
			}
		}
		if !isok {
			return nil, messages, errors.New("Error on post fields. Please check the form.")
		}
		if haschanged {
			// do an update of new data
			if m.Hooks.Update != nil {
				err := m.Hooks.Update(m, ctx, key, changedrecord)
				if err != nil {
					return nil, messages, err
				}
			} else {
				err := m.update(ctx, key, record, changedrecord)
				if err != nil {
					return nil, messages, err
				}
			}
		}
		if m.Hooks.PostInsert != nil {
			err = m.Hooks.PostInsert(m, ctx, key, record)
			if err != nil {
				return nil, messages, err
			}
		}
	case DOUPDATE:
		if DEBUG {
			fmt.Println("xdominion.Mask::execute-->DOUPDATE", realmode, formmode, key, record)
		}
		var oldrecord *xdominion.XRecord
		var err error
		if m.Hooks.GetRecord != nil {
			_, oldrecord, err = m.Hooks.GetRecord(m, ctx, key, 0)
			if err != nil {
				return nil, messages, err
			}
		} else {
			_, oldrecord, err = m.getrecord(key, 0)
			if err != nil {
				return nil, messages, err
			}
		}
		if m.Hooks.PreUpdate != nil {
			err := m.Hooks.PreUpdate(m, ctx, key, oldrecord, record)
			if err != nil {
				return nil, messages, err
			}
		}
		isok := true
		for _, f := range m.Fields {
			err = f.PreUpdate(ctx, key, oldrecord, record)
			if err != nil {
				isok = false
				messages[f.GetName()] = err.Error()
			}
		}
		if !isok {
			return nil, messages, errors.New("Error on pre fields. Please check the form.")
		}
		if m.Hooks.Update != nil {
			err := m.Hooks.Update(m, ctx, key, record)
			if err != nil {
				return nil, messages, err
			}
		} else {
			err := m.update(ctx, key, oldrecord, record)
			if err != nil {
				return nil, messages, err
			}
		}
		primkey := m.Table.GetPrimaryKey()
		primkeyname := primkey.GetName()
		record.Set(primkeyname, key)

		isok = true
		changedrecord := &xdominion.XRecord{}
		haschanged := false
		for _, f := range m.Fields {
			changed, err := f.PostUpdate(ctx, key, oldrecord, record)
			if err != nil {
				isok = false
				messages[f.GetName()] = err.Error()
			}
			if changed && f.GetInRecord() {
				haschanged = true
				changedvalue, _ := record.Get(f.GetName())
				changedrecord.Set(f.GetName(), changedvalue)
			}
		}
		if !isok {
			return nil, messages, errors.New("Error on post fields. Please check the form.")
		}
		if haschanged {
			// do an update of new data
			if m.Hooks.Update != nil {
				err := m.Hooks.Update(m, ctx, key, changedrecord)
				if err != nil {
					return nil, messages, err
				}
			} else {
				err := m.update(ctx, key, record, changedrecord)
				if err != nil {
					return nil, messages, err
				}
			}
		}
		if m.Hooks.PostUpdate != nil {
			err := m.Hooks.PostUpdate(m, ctx, key, oldrecord, record)
			if err != nil {
				return nil, messages, err
			}
		}
	case DODELETE:
		if DEBUG {
			fmt.Println("xdominion.Mask::execute-->DODELETE", realmode, formmode, key, record)
		}
		var oldrecord *xdominion.XRecord
		var err error
		if m.Hooks.GetRecord != nil {
			_, oldrecord, err = m.Hooks.GetRecord(m, ctx, key, 0)
			if err != nil {
				return nil, messages, err
			}
		} else {
			_, oldrecord, err = m.getrecord(key, 0)
			if err != nil {
				return nil, messages, err
			}
		}
		if m.Hooks.PreDelete != nil {
			err := m.Hooks.PreDelete(m, ctx, key, oldrecord, record)
			if err != nil {
				return nil, messages, err
			}
		}
		isok := true
		for _, f := range m.Fields {
			err := f.PreDelete(ctx, key, oldrecord, record)
			if err != nil {
				isok = false
				messages[f.GetName()] = err.Error()
			}
		}
		if !isok {
			return nil, messages, errors.New("Error on pre fields. Please check the form.")
		}
		if m.Hooks.Delete != nil {
			err := m.Hooks.Delete(m, ctx, key, oldrecord, record)
			if err != nil {
				return nil, messages, err
			}
		} else {
			err := m.delete(ctx, key, oldrecord, record)
			if err != nil {
				return nil, messages, err
			}
		}
		isok = true
		for _, f := range m.Fields {
			err := f.PostDelete(ctx, key, oldrecord, record)
			if err != nil {
				isok = false
				messages[f.GetName()] = err.Error()
			}
		}
		if !isok {
			return nil, messages, errors.New("Error on post fields. Please check the form.")
		}
		if m.Hooks.PostDelete != nil {
			err := m.Hooks.PostDelete(m, ctx, key, oldrecord, record)
			if err != nil {
				return nil, messages, err
			}
		}
	}
	return nil, messages, nil
}

func (m *Mask) buildrecord(ctx *context.Context, mode Mode) *xdominion.XRecord {
	// we loop over the fields
	record := &xdominion.XRecord{}
	for _, f := range m.Fields {
		if f.GetType() == CONTROL {
			continue
		}
		v, ignore, _ := f.GetValue(ctx, mode)
		if ignore {
			continue
		}
		record.Set(f.GetName(), v)
	}
	return record
}

func (m *Mask) insert(ctx *context.Context, rec *xdominion.XRecord) (interface{}, error) {

	if m.Table == nil {
		return nil, nil
	}
	key, err := m.Table.Insert(rec)
	return key, err
}

func (m *Mask) update(ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, newrec *xdominion.XRecord) error {
	if m.Table == nil {
		return nil
	}
	_, err := m.Table.Update(key, newrec)
	return err
}

func (m *Mask) delete(ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, params *xdominion.XRecord) error {

	fmt.Println("Mask::delete", key, oldrec, params)

	if m.Table == nil {
		return nil
	}
	_, err := m.Table.Delete(key)
	return err
}

func (m *Mask) convertPrimaryKey(value interface{}) interface{} {
	if m.KeyField == "" {
		return value
	}
	for _, f := range m.Fields {
		if f.GetName() == m.KeyField {
			nv, err := f.ConvertValue(value)
			if err == nil {
				return nv
			}
		}
	}
	return value
}

// ==============================================
// TOOLS for mask
// ==============================================
func convertModes(mode Mode) string {
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

func createMode(mode string) (Mode, Mode) {
	fmt.Println("Mask::createMode", mode)
	switch mode {
	case "1":
		return DOINSERT, INSERT
	case "2":
		return DOUPDATE, UPDATE
	case "3":
		return DODELETE, DELETE
	default:
		return VIEW, VIEW
	}
}

func verifyMode(mode string) string {
	if mode != "1" && mode != "2" && mode != "3" && mode != "4" {
		mode = "4"
	}
	return mode
}

func revertOrder(order xdominion.XOrder) xdominion.XOrder {

	norder := order.Clone()
	for i, o := range norder {
		if o.Operator == xdominion.ASC {
			norder[i].Operator = xdominion.DESC
		} else {
			norder[i].Operator = xdominion.ASC
		}
	}
	return norder
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

  public function createModes($modes)
  {
    return (($modes & DomMask::INSERT)?'1':'') . (($modes & DomMask::UPDATE)?'2':'') . (($modes & DomMask::DELETE)?'3':'') . (($modes & DomMask::VIEW)?'4':'');
  }

  // creates the WAJAF objects and gets back the container
  // type can be json, 4gl
  public function run($type = DomMask::JSON)
  {

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
