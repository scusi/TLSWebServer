package main

import (
	"net/http"
	"fmt"
	"os"
	"strings"
)

// Redirect all traffic to HTTPS
func (app *App) redirectHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(os.Stdout,"Inside redirectHandler  redirecting to %s ","https://"+app.domain+":"+app.httpsport+"/index.html")
	http.Redirect(w, r, "https://"+app.domain+":"+app.httpsport+"/index.html", http.StatusMovedPermanently)
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
	//fmt.Fprintf(os.Stdout,"Inside staticHandler 1.  requestPath:  %s fileSystemPath : %s  endURIPath: %s  splitPath: %s ",requestPath, fileSystemPath, endURIPath, splitPath )
	if len(splitPath) > 1 {
	  if f, err := os.Stat(fileSystemPath); err == nil && !f.IsDir() {
		fmt.Fprintf(os.Stdout,"Inside staticHandler 2. requestPath:  %s ",requestPath)
		http.ServeFile(w, r, fileSystemPath)
		return
	  }
	  //If not found just send back to index.html
	  //http.NotFound(w, r)
	  //fmt.Fprintf(os.Stdout,"Inside staticHandler 3. serving from   %s ",app.StaticDir+"/index.html")
	  http.ServeFile(w, r, app.StaticDir+"/index.html")
	  return
	}
	//fmt.Fprintf(os.Stdout,"Inside staticHandler 4. serving from   %s ",app.StaticDir+"/index.html")
	http.ServeFile(w, r, app.StaticDir+"/index.html")
  }

// Handle a simple healthcheck url for loadbalancers or whatnot
func (app *App) healthCheckHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	//fmt.Fprintf(os.Stdout,"Inside healthCheckHandler ")
	fmt.Fprintf(w, "OK!\r")
}

// Handler to redirect all NoFound to Angular index.html
func (app *App) angularHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(os.Stdout,"Inside angularHandler ")
	http.ServeFile(w, r, app.StaticDir+"/index.html")
}