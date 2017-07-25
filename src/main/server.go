package main

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/secure"
	"log"
	"net/http"
)

// StartServer Wraps the mux Router and uses the Negroni Middleware
// like flask app.run
func StartServer(ctx AppContext) {
	//create a new routerController
	router := mux.NewRouter().StrictSlash(true)
	//circle routes from global array
	for _, route := range routes {
		//define a new val by http.Handler while cycle every time
		var handler http.Handler
		//makeHandler
		handler = makeHandler(ctx, route.HandlerFunc)
		// add method header Path,name for route array
		router.
		Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	//get current env
	// security
	var isDevelopment = false
	if ctx.Env == local {
		isDevelopment = true
	}
	secureMiddleware := secure.New(secure.Options{
		IsDevelopment:      isDevelopment, // This will cause the AllowedHosts, SSLRedirect, and STSSeconds/STSIncludeSubdomains options to be ignored during development. When deploying to production, be sure to set this to false.
		AllowedHosts:       []string{},    // AllowedHosts is a list of fully qualified domain names that are allowed (CORS)
		ContentTypeNosniff: true,          // If ContentTypeNosniff is true, adds the X-Content-Type-Options header with the value `nosniff`. Default is false.
		BrowserXssFilter:   true,          // If BrowserXssFilter is true, adds the X-XSS-Protection header with the value `1; mode=block`. Default is false.
	})
	// start now
	//create a new negroni,register mux route in it
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))
	n.UseHandler(router)
	log.Println("===> Starting app (v" + ctx.Version + ") on port " + ctx.Port + " in " + ctx.Env + " mode.")
	if ctx.Env == local {
		n.Run("localhost:" + ctx.Port)
	} else {
		n.Run(":" + ctx.Port)
	}
}
