package main

import (
	"net"
	"net/http"
)

// TLSRedirect is a function that will redirect all incoming requests to its https equivilent.
func (app *App) TLSRedirect(w http.ResponseWriter, req *http.Request) {
	_, appPort, _ := net.SplitHostPort(app.httpsport)
	reqHost, _, _ := net.SplitHostPort(req.Host)
	localAddr := net.JoinHostPort(reqHost, appPort)
	http.Redirect(w, req,
		"https://"+localAddr+req.URL.String(),
		http.StatusMovedPermanently)
}
