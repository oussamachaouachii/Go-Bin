package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
	"snippetbox.oussama.com/internal/dto"
	"snippetbox.oussama.com/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.renderTemplate(w, *data, "home.tmpl.html")

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		} else {
			app.serverError(w, err)
			return
		}
	}

	rendData := app.newTemplateData(r)
	rendData.Snippet = snippet
	app.renderTemplate(w, *rendData, "view.tmpl.html")

}

func (app *application) snippetCreatePOST(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("ERROR HERE")
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	validate := validator.New()
	snippetData := &dto.SnippetCreateDTO{
		Title:   title,
		Content: content,
		Expires: expires,
	}
	err = validate.Struct(snippetData)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		fmt.Fprint(w, err)
	}

	// fieldErrors := make(map[string]string)
	// if strings.TrimSpace(title) == "" {
	// 	fieldErrors["title"] = "This field cannot be blank or filled with spaces"
	// } else if utf8.RuneCountInString(title) > 100 {
	// 	fieldErrors["title"] = "This field cannot be more than 100 characterslong"
	// }

	// if strings.TrimSpace(content) == "" {
	// 	fieldErrors["title"] = "This field cannot be blank or filled with spaces"
	// }

	// if expires != 1 && expires != 7 && expires != 365 {
	// 	fieldErrors["expires"] = "This field must equal 1, 7 or 365"
	// }

	// if len(fieldErrors) > 1 {
	// 	fmt.Fprint(w, fieldErrors)
	// 	return
	// }

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id),
		http.StatusSeeOther)
}

func (app *application) snippetNew(w http.ResponseWriter, r *http.Request) {
	rendData := app.newTemplateData(r)
	app.renderTemplate(w, *rendData, "create.tmpl.html")
}
