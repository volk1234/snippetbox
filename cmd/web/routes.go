package main

import (
	"net/http"

	"snippetbox.sandx.link/ui"

	"github.com/justinas/alice"
	"github.com/julienschmidt/httprouter"
	
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		app.notFound(w)
	})
	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)
	router.HandlerFunc(http.MethodGet, "/ping", ping)
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.userLogoutPost))
	router.Handler(http.MethodGet, "/about", dynamic.ThenFunc(app.about))

	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodGet, "/account/view", protected.ThenFunc(app.accountView))
	router.Handler(http.MethodGet, "/account/password/update", protected.ThenFunc(app.accountPasswordUpdate))
	router.Handler(http.MethodPost, "/account/password/update", protected.ThenFunc(app.accountPasswordUpdatePost))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/snippet/logout", protected.ThenFunc(app.userLogoutPost))

	standartChain := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standartChain.Then(router)
}
