package main

import (
	"log"
	"net/http"
	"strings"
)

// TLSRedirect is a function that will redirect all incoming requests to its https equivilent.
func (app *App) TLSRedirect(w http.ResponseWriter, req *http.Request) {
	host := strings.Split(req.Host, ":")[0]
	if len(cfg.ExposedHttpsAddr) > 0 {
		httpsPort := strings.Split(req.Host, ":")[1]
		host += ":" + httpsPort
	}
	targetURL := "https://" + host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		targetURL += "?" + req.URL.RawQuery
	}
	log.Printf("redirect %s to: %s", req.RemoteAddr, targetURL)
	http.Redirect(w, req, targetURL, http.StatusPermanentRedirect)
}
