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
	app.renderTemplate(w, http.StatusOK, data, "home.tmpl.html")

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

	app.renderTemplate(w, http.StatusOK, rendData, "view.tmpl.html")

}

func (app *application) snippetCreatePOST(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
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
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusFound)
}

func (app *application) snippetNew(w http.ResponseWriter, r *http.Request) {
	rendData := app.newTemplateData(r)
	app.renderTemplate(w, http.StatusOK, rendData, "create.tmpl.html")
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	rendData := app.newTemplateData(r)
	app.renderTemplate(w, http.StatusOK, rendData, "signup.tmpl.html")
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	validate := validator.New()
	userData := &dto.UserCreateDTO{
		Name:     name,
		Email:    email,
		Password: password,
	}
	err = validate.Struct(userData)
	if err != nil {
		rendData := app.newTemplateData(r)
		app.renderTemplate(w, http.StatusUnprocessableEntity, rendData, "signup.tmpl.html")
		return
	}

	err = app.users.Insert(name, email, password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			app.sessionManager.Put(r.Context(), "flash", "There is an account associated with this email.")
			rendData := app.newTemplateData(r)
			app.renderTemplate(w, http.StatusUnprocessableEntity, rendData, "signup.tmpl.html")
			fmt.Println("heeeereeee")
			return
		} else {
			app.serverError(w, err)
			return
		}
	}
	app.sessionManager.Put(r.Context(), "flash", "User successfully created!")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	validate := validator.New()
	userData := &dto.UserLoginDTO{
		Email:    email,
		Password: password,
	}

	err = validate.Struct(userData)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	id, username, err := app.users.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			app.sessionManager.Put(r.Context(), "flash", "Invalid Credentials.")
			rendData := app.newTemplateData(r)
			app.renderTemplate(w, http.StatusOK, rendData, "login.tmpl.html")
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	app.sessionManager.Put(r.Context(), "username", username)
	app.sessionManager.Put(r.Context(), "flash", "Logged in successfully, WELCOME :)")

	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/snippet/new", http.StatusSeeOther)

}
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	rendData := app.newTemplateData(r)
	app.renderTemplate(w, http.StatusOK, rendData, "login.tmpl.html")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
