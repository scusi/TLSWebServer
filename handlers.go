package main

import (
	"net/http"
	"fmt"
	"os"
	"strings"
)

// Redirect all traffic to HTTPS
func (app *App) redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+app.domain+":"+app.httpsport+r.RequestURI, http.StatusMovedPermanently)
}

// Serve your index file
func (app *App) indexHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	http.ServeFile(w, r, app.StaticDir+"/index.html")
}

// Handle Static Files (files containing a .extension)
// Send a not found error for any requests containing a .extension where the file does not exist
// Redirect any requests that were not for files to your to your SPA
func (app *App) staticHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	requestPath := r.URL.Path
	fileSystemPath := app.StaticDir + r.URL.Path
	endURIPath := strings.Split(requestPath, "/")[len(strings.Split(requestPath, "/"))-1]
	splitPath := strings.Split(endURIPath, ".")
	if len(splitPath) > 1 {
	  if f, err := os.Stat(fileSystemPath); err == nil && !f.IsDir() {
		http.ServeFile(w, r, fileSystemPath)
		return
	  }
	  //If not found just send back to index.html
	  //http.NotFound(w, r)
	  http.ServeFile(w, r, app.StaticDir+"/index.html")
	  return
	}
	http.ServeFile(w, r, app.StaticDir+"/index.html")
  }

// Handle a simple healthcheck url for loadbalancers or whatnot
func (app *App) healthCheckHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	fmt.Fprintf(w, "OK!\r")
}

// Handler to redirect all NoFound to Angular index.html
func (app *App) angularHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, app.StaticDir+"/index.html")
}