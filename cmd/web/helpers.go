package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) renderTemplate(w http.ResponseWriter, data templateData, page string) {
	buf := new(bytes.Buffer)
	err := app.template[page].ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}
