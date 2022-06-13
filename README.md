[![Go Report Card](https://goreportcard.com/badge/github.com/webability-go/xdommask)](https://goreportcard.com/report/github.com/webability-go/xdommask)
[![GoDoc](https://godoc.org/github.com/webability-go/xdommask?status.png)](https://godoc.org/github.com/webability-go/xdommask)
[![GolangCI](https://golangci.com/badges/github.com/webability-go/xdommask.svg)](https://golangci.com)

XDOMMASK for GO v1
=============================

XDomMask is the Mask Generation (aka HTML Forms) for GO-WAJAF.

Manuals are available on godoc.org [![GoDoc](https://godoc.org/github.com/webability-go/wajaf?status.png)](https://godoc.org/github.com/webability-go/xdommask)

TO DO:
======
- Tests

Version Changes Control
=======================

v0.1.0 - 2022-06-13
------------------------
- Added colorField, maskedField (replacing pwField)
- All the mask function now pass the context to build and log messages
- Better conversion of primary keys to be compatible with numeric and string values
- Messages of information and errors are better managed
- field main interface changed to add convertValue prototype and some other minor changes
- GetName() funcion added in buttonField
- integerField, floatField now converts correctly all the possible values from the field string value
- textField now support automessages and values, password type renamed to masked type
- lovField now works and can take a query of values from a XDominion filtered object
- examples changed to meet new changes of main code

v0.0.7 - 2021-01-04
------------------------
- Implementation of type of field with constants (control, field, hidden, info)
- Add NewInfoField to the list of known fields

v0.0.6 - 2020-04-13
------------------------
- mail, pw, textarea implemented

v0.0.5 - 2020-04-13
------------------------
- Modes on text adjusted. Node Name is same as ID by default. Do not send empty fields to the wajaf structure

v0.0.3 - 2020-04-13
------------------------
- some errors corrected on fields

v0.0.2 - 2020-04-13
------------------------
- mask, text, button with all attributes and messages

v0.0.1 - 2020-04-12
------------------------
- Added Compile and zones without ID
- Added some field for button and text

v0.0.0 - 2020-04-08
------------------------
- First realease



Reference
====================================

Una mascara tiene 4 modos de funcionamiento:
- Insert
- Update
- Delete
- View

Cada modo de edicion funciona en 3 etapas:
1. Despliego la forma con todos sus campos
  pasan varias cosas:
  1.1 Se llama el Compile de XDomMask para fabricar el JSON de definicion de la forma (con (dentro del COMPILE) o sin la data del renglon)

  1.2 Se envia el JSON al cliente que corre el groupContainer.js

  1.3 El groupContainer.js despliegue de la forma (100% en en cliente: fabrica el <form y <fields )

  1.4 Si estoy en update o delete, rellena los campos con los valores de la data del renglon
     1.4.1 Si la data ya viene con el codigo de la forma, lo usa y lo rellena
     1.4.2 Sino, va a pedir la data al server para rellenar y regresa a 1.2.1 con la data (con el EXECUTE->getRecord)

  Punto y aparte:
  1.5 Para contruir la data, si se necesita durante el COMPILE, va a llamar el EXECUTE->GetRecord(key)
    Si se necesita diferido, entonces el cliente llama el punto de entrada EXECUTE->GetRecord

2. El usurio tiene interaccion con la mascara (rellena los campos)

  Todos los JS de validacion

3. Realiza la DoAction (insertar , modificar, o borrar la data en la fuente de datos)

  Cuando se pica el boton DoAction (DoInsert, DoUpdate o DoDelete), llama el Execute del server

  2.1 Si es DoUpate
    2.1.1 El Execute va a contruir el renglon a modificar
    2.1.2 Llama PreUpdate para que el usuario tenga la libertad de modificar el renglon y validar lo necesario
    2.1.3 Llama Update para fisicamente modificar el renglon   (key, record)
        En caso de NO tener una tabla, entonces el usuario TENDRA QUE definir su propia funcion Update
    2.1.4 Llama PostUpdate para cualquier post tratamiento

El modo view solo funciona en 1 sola etapa
1. Despliego la forma con todos sus campos
   relleno los campos con los valores de la data del renglon
