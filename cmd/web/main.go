package main

import (
	"flag"
	"log"
	"net/http"
)

var httpAddr string
var httpsAddr string
var cert string
var key string
var staticDir string

func init() {
	flag.StringVar(&httpAddr, "http", ":80", "http listen address")
	flag.StringVar(&httpsAddr, "https", ":443", "https listen address")
	flag.StringVar(&cert, "cert", "./tls/cert.pem", "tls certificate PEM file")
	flag.StringVar(&key, "key", "./tls/key.pem", "tls key PEM file")
	flag.StringVar(&staticDir, "staticDir", "./ui/static", "directory with static webcontent")
}

func main() {
	flag.Parse()

	app := &App{
		Addr:      httpsAddr,
		StaticDir: staticDir,
		TLSCert:   cert,
		TLSKey:    key,
	}

	// start a http server on httpAddr that just redirects to https
	go func() {
		httpMux := http.NewServeMux()
		httpMux.HandleFunc("/", app.TLSRedirect)
		log.Println("Starting http server on ", httpAddr)
		err := http.ListenAndServe(httpAddr, httpMux)
		log.Fatal(err)
	}()

	// start a https server on httpsAddr
	app.RunTLSServer()
}
