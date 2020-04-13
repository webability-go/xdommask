package xdommask

import (
	"fmt"
	"testing"

	"github.com/webability-go/xdominion"
)

type MyMask struct {
	*Mask
}

func GetRecord(m *Mask) *xdominion.XRecord {
	fmt.Println("GET RECORD DE MYMASK")
	rec := xdominion.NewXRecord()
	rec.Set("Field1", "Value1")
	fmt.Println(rec)
	return rec
}

func TestNewMask(t *testing.T) {

	mask := &MyMask{Mask: NewMask("")}
	mask.GetRecord = GetRecord
	mask.Mode = INSERT
	mask.Variables["COUNTRY"] = "MX"
	mask.Variables["LANGUAGE"] = "es"

	mask.AlertMessage = "Alert: there was an error"
	mask.ServerMessage = "Error on the server side"
	mask.InsertTitle = "Insert"
	mask.UpdateTitle = "Update"
	mask.ViewTitle = "View"

	mask.SuccessJS = `
function(params)
{
  alert("Success");
}
`
	f1 := NewTextField("locale")
	f1.Title = "Locale:"
	f1.HelpDescription = "Type the desired Locale"
	f1.NotNullModes = INSERT | UPDATE
	f1.StatusNotNull = "Locale cannot be empty"
	f1.MaxLength = 200
	f1.URLVariable = "locale"
	//	f1.DefaultValue = "es_MX"
	mask.AddField(f1)

	f2 := NewButtonField("submit")
	f2.HelpDescription = "Type the desired Locale"
	f2.NotNullModes = INSERT | UPDATE
	mask.AddField(f2)

	stmask := mask.Compile()

	fmt.Println("Structure mask:", mask)

	fmt.Println("Structure mask compile WAJAF:", stmask)
}
