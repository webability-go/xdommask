package xdommask

import (
	"fmt"
	"testing"
	//	"encoding/xml"
	"encoding/json"

	"github.com/webability-go/wajaf"
	"github.com/webability-go/xdominion"
)

func GetRecord(m *Mask) *xdominion.XRecord {
	fmt.Println("GET RECORD DE MYMASK")
	rec := xdominion.NewXRecord()
	rec.Set("Field1", "Value1")
	fmt.Println(rec)
	return rec
}

func TestNewMask(t *testing.T) {

	mask := NewMask("formaccount")
	mask.Mode = INSERT
	mask.Variables["COUNTRY"] = "MX"
	mask.Variables["LANGUAGE"] = "es"

	mask.AlertMessage = "##mask.errormessage##"
	mask.ServerMessage = "##mask.servermessage##"
	mask.InsertTitle = "##mask.titleinsert##"
	mask.ViewTitle = "##mask.titleview##"

	mask.SuccessJS = `
function(params)
{
WA.$N('titleform').hide();
WA.$N('titleconfirmation').show();
WA.$N('continue').show();
WA.toDOM('mastermaininstall|single|step1').className = 'installstepdone';
WA.toDOM('mastermaininstall|single|step1').onclick = null;
WA.toDOM('mastermaininstall|single|step2').className = 'installstepdone';
WA.toDOM('mastermaininstall|single|step3').className = 'installstepactual';
}
`
	// serial
	f1 := NewTextField("serial")
	f1.Title = "##serial.title##"
	f1.HelpDescription = "##serial.help.description##"
	f1.NotNullModes = INSERT
	f1.StatusNotNull = "##serial.status.notnull##"
	f1.Size = "200"
	f1.MinLength = 20
	f1.MaxLength = 20
	f1.URLVariable = "serial"
	f1.Format = "^[a-z|A-Z|0-9]{20}$"
	f1.FormatJS = "^[a-z|A-Z|0-9]{20}$"
	f1.StatusBadFormat = "##serial.status.badformat##"
	f1.DefaultValue = "00000000000000000000"
	mask.AddField(f1)

	f2 := NewTextField("locale")
	f2.Title = "##locale.title##"
	f2.HelpDescription = "##locale.help.description##"
	f2.NotNullModes = INSERT
	f2.StatusNotNull = "##locale.status.notnull##"
	f2.Size = "200"
	f2.MaxLength = 200
	f2.URLVariable = "locale"
	f2.DefaultValue = "es_MX"
	mask.AddField(f2)

	f3 := NewTextField("timezone")
	f3.Title = "##timezone.title##"
	f3.HelpDescription = "##timezone.help.description##"
	f3.NotNullModes = INSERT
	f3.StatusNotNull = "##timezone.status.notnull##"
	f3.URLVariable = "timezone"
	//	f3.options =
	f3.DefaultValue = ""
	mask.AddField(f3)

	// username
	f4 := NewTextField("username")
	f4.Title = "##username.title##"
	f4.HelpDescription = "##username.help.description##"
	f4.NotNullModes = INSERT
	f4.StatusNotNull = "##username.status.notnull##"
	f4.MinLength = 4
	f4.MaxLength = 80
	f4.StatusTooShort = "##username.status.tooshort##"
	f4.URLVariable = "username"
	mask.AddField(f4)

	// password
	f5 := NewPWField("password")
	f5.Title = "##password.title##"
	f5.HelpDescription = "##password.help.description##"
	f5.NotNullModes = INSERT
	f5.StatusNotNull = "##password.status.notnull##"
	f5.MaxLength = 80
	f5.MinLength = 4
	f5.StatusTooShort = "##password.status.tooshort##"
	f5.URLVariable = "password"
	mask.AddField(f5)

	// email
	f6 := NewMailField("email")
	f6.Title = "##email.title##"
	f6.HelpDescription = "##email.help.description##"
	f6.NotNullModes = INSERT
	f6.StatusNotNull = "##email.status.notnull##"
	f6.MaxLength = 200 // chars
	f6.URLVariable = "email"
	mask.AddField(f6)

	f7 := NewButtonField("", "submit")
	//	f7.Action = "submit"
	f7.AuthModes = INSERT // insert
	f7.TitleInsert = "##form.continue##"
	mask.AddField(f7)

	f8 := NewButtonField("", "reset")
	//	f8.Action = "reset"
	f8.AuthModes = INSERT // insert + view
	f8.TitleInsert = "##form.reset##"
	mask.AddField(f8)

	// from xdommask structure to wajaf object structures
	stmask := mask.Compile()

	app := wajaf.NewApplication("")
	app.AddChild(stmask)

	fmt.Println("Structure app compiled WAJAF:", app)

	json, err := json.Marshal(app)
	if err != nil {
		fmt.Println("error json marshal", err)
		return
	}
	fmt.Println("JSON CODE", string(json))
}
