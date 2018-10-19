package main

import (
	"flag"
	"log"
	"net/http"
)

var httpAddr string
var httpsAddr string
var tlsCertPath string
var tlsKeyPath string
var staticDir string

func init() {
	flag.StringVar(&httpAddr, "http", "", "http to https redirector listen address")
	flag.StringVar(&httpsAddr, "https", ":443", "https listen address")
	flag.StringVar(&tlsCertPath, "cert", "./tls/cert.pem", "tls certificate PEM file")
	flag.StringVar(&tlsKeyPath, "key", "./tls/key.pem", "tls key PEM file")
	flag.StringVar(&staticDir, "staticDir", "", "directory with static webcontent")
}

func main() {
	flag.Parse()

	if staticDir == "" {
		log.Fatal("no staticDir given but required")
	}

	if tlsCertPath == "" {
		log.Fatal("no certificate file given, but required")
	}

	if tlsKeyPath == "" {
		log.Fatal("no key file given but required")
	}

	app := &App{
		Addr:      httpsAddr,
		StaticDir: staticDir,
		TLSCert:   tlsCertPath,
		TLSKey:    tlsKeyPath,
	}

	// start a http server on httpAddr that just redirects to https
	if httpAddr != "" {
		go func() {
			httpMux := http.NewServeMux()
			httpMux.HandleFunc("/", app.TLSRedirect)
			log.Println("Starting http server on ", httpAddr)
			err := http.ListenAndServe(httpAddr, httpMux)
			log.Fatal(err)
		}()
	}

	// start a https server on httpsAddr
	app.RunTLSServer()
}
