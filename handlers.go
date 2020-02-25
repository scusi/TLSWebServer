package main

import (
	"net"
	"net/http"
)

func (app *App) TLSRedirect(w http.ResponseWriter, req *http.Request) {
	_, appPort, _ := net.SplitHostPort(app.Addr)
	reqHost, _, _ := net.SplitHostPort(req.Host)
	localAddr := net.JoinHostPort(reqHost, appPort)
	http.Redirect(w, req,
		"https://"+localAddr+req.URL.String(),
		http.StatusMovedPermanently)
}
