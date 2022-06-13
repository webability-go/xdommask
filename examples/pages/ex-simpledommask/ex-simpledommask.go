package main

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"
	"github.com/webability-go/xdommask"

	"github.com/webability-go/xamboo/cms/context"
)

var base, _ = verifyDB()

func Run(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	mode := ctx.Request.Form.Get("mode")
	dbmode := "Use the form without Database"
	if base != nil {
		dbmode = "Use the form on a Database"
	}

	//	bridge.EntityLog_LogStat(ctx)
	params := &xcore.XDataset{
		"DB":   dbmode,
		"FORM": createXMLMask("formaccess", mode, ctx),
		"#":    language,
	}

	return template.Execute(params)
}

// Formaccess is the listenere of the groupContainer (XDomMask)
func Formaccess(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	xdommask.DEBUG = true
	mask, _ := createMask("formaccess", ctx)
	data, _ := mask.Run(ctx)
	return data
}

func Createdb(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	err := createDB()
	if err != nil {
		return err.Error()
	}
	base, err = verifyDB()
	if err != nil {
		return err.Error()
	}
	return "OK"
}

func createMask(id string, ctx *context.Context) (*xdommask.Mask, error) {

	hooks := xdommask.MaskHooks{
		Build: build,
	}
	if base == nil {
		hooks.GetRecord = getrecord
		hooks.Insert = insert
		hooks.Update = update
		hooks.Delete = delete
	}
	return xdommask.NewMask(id, hooks, ctx)
}

func createXMLMask(id string, mode string, ctx *context.Context) string {
	mask, _ := createMask(id, ctx)
	cmask := mask.Compile(mode, ctx)
	xmlmask, _ := xml.Marshal(cmask)
	return string(xmlmask)
}

func build(mask *xdommask.Mask, ctx *context.Context) error {
	if base != nil {
		mask.Table = getTableDef(base)
		mask.Order = &xdominion.XOrder{xdominion.NewXOrderBy("int", xdominion.ASC)}
	}

	mask.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	mask.KeyField = "int"
	mask.Key = 1

	mask.AlertMessage = "##mask.errormessage##"
	mask.ServerMessage = "##mask.servermessage##"
	mask.InsertTitle = "##mask.inserttitle##"
	mask.UpdateTitle = "##mask.updatetitle##"
	mask.DeleteTitle = "##mask.deletetitle##"
	mask.ViewTitle = "##mask.viewtitle##"
	mask.FailureJS = "function(params) { this.icontainer.setMessages(params); }"

	// Text field, integer, PRIMARY KEY
	f10 := xdommask.NewIntegerField("int")
	f10.Title = "##int.title##"
	f10.HelpDescription = "##int.help.description##"
	f10.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f10.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f10.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f10.ViewModes = xdommask.DELETE | xdommask.VIEW
	f10.StatusNotNull = "##int.status.notnull##"
	f10.StatusBadFormat = "##int.status.badformat##"
	f10.StatusTooLow = "##int.status.toolow##"
	f10.StatusTooHigh = "##int.status.toohigh##"
	f10.Min = -10
	f10.Max = 222
	f10.Size = "400"
	f10.URLVariable = "int"
	f10.DefaultValue = -1
	mask.AddField(f10)

	// textfield
	f11 := xdommask.NewTextField("text")
	f11.Title = "##text.title##"
	f11.HelpDescription = "##text.help.description##"
	f11.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f11.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f11.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f11.ViewModes = xdommask.DELETE | xdommask.VIEW
	f11.StatusNotNull = "##text.status.notnull##"
	f11.MaxLength = 40 // chars
	f11.MinLength = 3  // chars
	f11.StatusTooShort = "##text.status.tooshort##"
	f11.StatusTooLong = "##text.status.toolong##"
	f11.Size = "400"
	f11.URLVariable = "text"
	f11.DefaultValue = "write something"
	mask.AddField(f11)

	// password field
	f12 := xdommask.NewMaskedField("password")
	f12.Title = "##password.title##"
	f12.HelpDescription = "##password.help.description##"
	f12.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f12.AuthModes = xdommask.INSERT | xdommask.UPDATE
	f12.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f12.StatusNotNull = "##password.status.notnull##"
	f12.MaxLength = 40 // chars
	f12.MinLength = 4  // chars
	f12.StatusTooShort = "##password.status.tooshort##"
	f12.StatusTooLong = "##password.status.toolong##"
	f12.Size = "400"
	f12.URLVariable = "password"
	mask.AddField(f12)

	// floatfield
	f13 := xdommask.NewFloatField("float")
	f13.Title = "##float.title##"
	f13.HelpDescription = "##float.help.description##"
	f13.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f13.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f13.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f13.ViewModes = xdommask.DELETE | xdommask.VIEW
	f13.StatusNotNull = "##float.status.notnull##"
	f13.StatusTooLow = "##float.status.toolow##"
	f13.StatusTooHigh = "##float.status.toohigh##"
	f13.StatusBadFormat = "##float.status.badformat##"
	f13.Min = -123.456
	f13.Max = 234.567
	f13.Size = "400"
	f13.URLVariable = "float"
	f13.DefaultValue = 3.1415927
	mask.AddField(f13)

	// mailfield
	f14 := xdommask.NewMailField("mail")
	f14.Title = "##mail.title##"
	f14.HelpDescription = "##mail.help.description##"
	f14.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f14.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f14.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f14.ViewModes = xdommask.DELETE | xdommask.VIEW
	f14.StatusNotNull = "##mail.status.notnull##"
	f14.StatusBadFormat = "##mail.status.badformat##"
	f14.Size = "400"
	f14.URLVariable = "mail"
	mask.AddField(f14)

	// textareafield
	f15 := xdommask.NewTextAreaField("textarea")
	f15.Title = "##textarea.title##"
	f15.HelpDescription = "##textarea.help.description##"
	f15.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f15.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f15.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f15.ViewModes = xdommask.DELETE | xdommask.VIEW
	f15.StatusNotNull = "##textarea.status.notnull##"
	f15.Width = 500
	f15.Height = 120
	f15.URLVariable = "textarea"
	mask.AddField(f15)

	// infofield
	f16 := xdommask.NewInfoField("info", "##info.title##")
	f16.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	mask.AddField(f16)

	// loofield
	f17 := xdommask.NewLOOField("loo")
	f17.Title = "##loo.title##"
	f17.HelpDescription = "##loo.help.description##"
	f17.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f17.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f17.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f17.ViewModes = xdommask.DELETE | xdommask.VIEW
	f17.StatusNotNull = "##loo.status.notnull##"
	f17.Size = "400"
	f17.URLVariable = "loo"
	f17.Options = map[string]string{
		"v1": "Value 1",
		"v2": "Value 2",
		"v3": "Value 3",
	}
	mask.AddField(f17)

	// lovfield
	f18 := xdommask.NewLOVField("lov")
	f18.Title = "##lov.title##"
	f18.HelpDescription = "##lov.help.description##"
	f18.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f18.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f18.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f18.ViewModes = xdommask.DELETE | xdommask.VIEW
	f18.StatusNotNull = "##lov.status.notnull##"
	f18.Size = "400"
	f18.URLVariable = "lov"
	f18.Options = map[string]string{
		"v1": "Value 1",
		"v2": "Value 2",
		"v3": "Value 3",
	}
	mask.AddField(f18)

	// datefield
	f19 := xdommask.NewColorField("color")
	f19.Title = "##color.title##"
	f19.HelpDescription = "##color.help.description##"
	f19.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f19.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f19.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f19.ViewModes = xdommask.DELETE | xdommask.VIEW
	f19.StatusNotNull = "##color.status.notnull##"
	f19.URLVariable = "color"
	f19.DefaultValue = ""
	mask.AddField(f19)

	// datefield
	f20 := xdommask.NewDateField("date")
	f20.Title = "##date.title##"
	f20.HelpDescription = "##date.help.description##"
	f20.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f20.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f20.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f20.ViewModes = xdommask.DELETE | xdommask.VIEW
	f20.StatusNotNull = "##date.status.notnull##"
	f20.Size = "400"
	f20.URLVariable = "date"
	f20.DefaultValue = ""
	mask.AddField(f20)

	// imagefield
	f21 := xdommask.NewImageField("image")
	f21.Title = "##image.title##"
	f21.HelpDescription = "##image.help.description##"
	f21.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f21.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f21.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f21.ViewModes = xdommask.DELETE | xdommask.VIEW
	f21.StatusNotNull = "##image.status.notnull##"
	f21.Size = "400"
	f21.URLVariable = "image"
	mask.AddField(f21)

	// videofield
	f22 := xdommask.NewVideoField("video")
	f22.Title = "##video.title##"
	f22.HelpDescription = "##video.help.description##"
	f22.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f22.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f22.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f22.ViewModes = xdommask.DELETE | xdommask.VIEW
	f22.StatusNotNull = "##video.status.notnull##"
	f22.Size = "400"
	f22.URLVariable = "video"
	mask.AddField(f22)

	// filefield
	f23 := xdommask.NewFileField("file")
	f23.Title = "##file.title##"
	f23.HelpDescription = "##file.help.description##"
	f23.NotNullModes = xdommask.INSERT | xdommask.UPDATE
	f23.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE | xdommask.VIEW
	f23.HelpModes = xdommask.INSERT | xdommask.UPDATE
	f23.ViewModes = xdommask.DELETE | xdommask.VIEW
	f23.StatusNotNull = "##file.status.notnull##"
	f23.Size = "400"
	f23.URLVariable = "file"
	mask.AddField(f23)

	// hiddenfield
	f24 := xdommask.NewHiddenField("hidden")
	f24.Title = "##hidden.title##"
	f24.AuthModes = xdommask.INSERT | xdommask.UPDATE
	f24.URLVariable = "hidden"
	//	f24.DefaultValue = "field.hidden"
	mask.AddField(f24)

	//buttons

	// first record
	f1 := xdommask.NewButtonField("", "first")
	f1.AuthModes = xdommask.VIEW
	f1.TitleView = "##mask.buttonfirst.titleview##"
	mask.AddField(f1)

	// previous record
	f2 := xdommask.NewButtonField("", "previous")
	f2.AuthModes = xdommask.VIEW
	f2.TitleView = "##mask.buttonprevious.titleview##"
	mask.AddField(f2)

	// next record
	f3 := xdommask.NewButtonField("", "next")
	f3.AuthModes = xdommask.VIEW
	f3.TitleView = "##mask.buttonnext.titleview##"
	mask.AddField(f3)

	// last record
	f4 := xdommask.NewButtonField("", "last")
	f4.AuthModes = xdommask.VIEW
	f4.TitleView = "##mask.buttonlast.titleview##"
	mask.AddField(f4)

	// Insert
	f5 := xdommask.NewButtonField("", "insert")
	f5.AuthModes = xdommask.VIEW
	f5.TitleView = "##mask.buttoninsert.titleview##"
	mask.AddField(f5)

	// Update
	f6 := xdommask.NewButtonField("", "update")
	f6.AuthModes = xdommask.VIEW
	f6.TitleView = "##mask.buttonupdate.titleview##"
	mask.AddField(f6)

	// Delete
	f7 := xdommask.NewButtonField("", "delete")
	f7.AuthModes = xdommask.VIEW
	f7.TitleView = "##mask.buttondelete.titleview##"
	mask.AddField(f7)

	// Submit
	f8 := xdommask.NewButtonField("", "submit")
	f8.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE
	f8.TitleInsert = "##mask.buttonsubmit.titleinsert##"
	f8.TitleUpdate = "##mask.buttonsubmit.titleupdate##"
	f8.TitleDelete = "##mask.buttonsubmit.titledelete##"
	mask.AddField(f8)

	// Reset
	f9 := xdommask.NewButtonField("", "reset")
	f9.AuthModes = xdommask.INSERT | xdommask.UPDATE
	f9.TitleInsert = "##mask.buttonreset.titleinsert##"
	f9.TitleUpdate = "##mask.buttonreset.titleupdate##"
	mask.AddField(f9)

	// View
	f91 := xdommask.NewButtonField("", "view")
	f91.AuthModes = xdommask.INSERT | xdommask.UPDATE | xdommask.DELETE
	f91.TitleInsert = "##mask.buttonview.titleinsert##"
	f91.TitleUpdate = "##mask.buttonview.titleupdate##"
	f91.TitleDelete = "##mask.buttonview.titledelete##"
	mask.AddField(f91)

	return nil
}

func getrecord(mask *xdommask.Mask, ctx *context.Context, key interface{}, mode int) (string, *xdominion.XRecord, error) {

	ikey, _ := key.(int)
	var rec *xdominion.XRecord
	switch mode {
	case 0: // aqui
		rec = getDBRecord(ikey)
	case -1:
		rec = getDBRecord(ikey - 1)
	case -2:
		rec = getDBRecord(1)
	case 1:
		rec = getDBRecord(ikey + 1)
	case 2:
		rec = getDBRecord(10)
	}
	if rec != nil {
		nkey, _ := rec.GetString("int")
		return nkey, rec, nil
	}
	return "", nil, nil
}

func insert(m *xdommask.Mask, ctx *context.Context, rec *xdominion.XRecord) (interface{}, error) {
	key, _ := rec.GetInt("int")
	fmt.Println("Inserting", key, rec)
	if key < 1 || key > 20 {
		return nil, errors.New("Error: the key is not between 1 and 20")
	}
	if key >= 1 && key <= 10 {
		return nil, errors.New("Error: the key already exists into the default dataset")
	}
	return key, nil
}

func update(m *xdommask.Mask, ctx *context.Context, key interface{}, newrec *xdominion.XRecord) error {
	newkey, _ := newrec.GetInt("int")
	fmt.Println("Updating", key, newkey, newrec)
	if newkey < 1 || newkey > 20 {
		return errors.New("Error: the key is not between 1 and 20")
	}
	if newkey >= 1 && newkey <= 10 {
		return errors.New("Error: the key already exists into the default dataset")
	}
	return nil
}

func delete(m *xdommask.Mask, ctx *context.Context, key interface{}, oldrec *xdominion.XRecord, rec *xdominion.XRecord) error {
	fmt.Println("Deleting", oldrec, rec)
	return nil
}

func buildRecord(rec *xdominion.XRecord) map[string]interface{} {
	if rec == nil {
		return nil
	}

	key, _ := rec.GetString("key")
	data := map[string]interface{}{}
	data["key"] = key
	data["name"], _ = rec.GetString("name")
	data["description"], _ = rec.GetString("description")
	data["group"], _ = rec.GetString("group")

	return map[string]interface{}{
		key: data,
	}
}

// MINI DATABASE SIMULATED
func getDBRecord(key int) *xdominion.XRecord {
	if key >= 1 && key <= 10 {
		return &xdominion.XRecord{
			"int":  key,
			"text": fmt.Sprint("Data line %d", key),
		}
	}
	return nil
}

func connectDB() (rbase *xdominion.XBase, rerr error) {

	lbase := &xdominion.XBase{
		DBType:   xdominion.DB_Postgres,
		Username: "username",
		Password: "password",
		Database: "test",
		Host:     xdominion.DB_Localhost,
		SSL:      false,
	}

	var err error
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("connectDB::Error Defered:", r)
			rbase = nil
			rerr = errors.New(fmt.Sprint(r))
		}
	}()

	lbase.Logon()
	fmt.Println("CONNECTDB", err)
	return lbase, err
}

func getTableDef(lbase *xdominion.XBase) *xdominion.XTable {
	t := xdominion.NewXTable("simpledommask", "")
	t.AddField(xdominion.XFieldInteger{Name: "int", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
	}}) // ai, pk
	t.AddField(xdominion.XFieldVarChar{Name: "text", Size: 100, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})
	t.AddField(xdominion.XFieldVarChar{Name: "password", Size: 100, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})
	t.AddField(xdominion.XFieldFloat{Name: "float"})
	t.AddField(xdominion.XFieldVarChar{Name: "mail", Size: 255})
	t.AddField(xdominion.XFieldText{Name: "textarea"})
	t.SetBase(lbase)
	return t
}

func createDB() error {
	lbase, err := connectDB()
	fmt.Println("ConnectDB: ", lbase, err)
	if lbase == nil {
		return err
	}
	tb := getTableDef(lbase)
	tb.Synchronize()
	for key := 1; key < 11; key++ {
		tb.Insert(xdominion.XRecord{"int": key, "text": fmt.Sprintf("Data line %d", key), "password": "qwerty", "float": 123.456, "mail": "test@test.com", "textarea": "Big text without limits\nSecond line"})
	}
	return nil
}

func verifyDB() (*xdominion.XBase, error) {

	// verify it can connect AND table exists
	lbase, err := connectDB()
	if lbase == nil {
		return nil, err
	}
	tb := getTableDef(lbase)
	_, err = tb.Count()
	if err != nil {
		return nil, err
	}
	return lbase, nil
}
