package main

import (
	"net/http"
)

// Routes basically routes requests through the LogRequest function
func (app *App) Routes() http.Handler {
	fileServer := http.FileServer(http.Dir(app.StaticDir))
	// Pass the router as the 'next' parameter to the LogRequest middleware.
	// Because LogRequest() is just a function, and the function returns a
	// http.Handler we don't need to do anything else.
	return LogRequest(SecureHeaders(fileServer))
}
