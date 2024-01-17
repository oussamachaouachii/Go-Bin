package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)
	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/new", protected.ThenFunc(app.snippetNew))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePOST))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))

	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.HandlerFunc(http.MethodGet, "/ping", app.ping)

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	return alice.New(app.recoverPanic, app.logRequest, secureHeaders).Then(router)
}
