package main

// Route is the model for the router setup
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc HandlerFunc
}

// Routes are the main setup for our Router
type Routes []Route

var routes = Routes{
	Route{"Healthcheck", "GET", "/healthcheck", HealthcheckHandler},
	Route{"restfulTest", "GET", "/restfulTest/{id}", GetFuncHandler},
	Route{"restfulTest", "POST", "/restfulTest/{id}", PostFuncHandler},
	Route{"restfulTest", "PUT", "/restfulTest/{id}", PutFuncHandler},
	Route{"restfulTest", "DELETE", "/restfulTest/{id}", DeleteFuncHandler},

}
