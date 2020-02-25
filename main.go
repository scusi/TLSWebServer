package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

/*
var httpAddr string
var httpsAddr string
var tlsCertPath string
var tlsKeyPath string
var staticDir string
*/

// Variables for version information
var (
	Version   string
	Commit    string
	Branch    string
	Buildtime string

	usr                 *user.User
	ConfigFilePath      string
	ConfigFileLocations []string
	cfg                 *Config
	err                 error
	showVersion         bool
)

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("Can not determine current user: ", err)
	}
	ConfigFileLocations = []string{
		"/etc/TLSWebserver/config.yml",
		"/usr/local/etc/TLSWebServer/config.yml",
		filepath.Join(usr.HomeDir, ".config/TLSWebServer/config.yml"),
		"./config.yml",
	}
	cfg = NewConfig("")

	flag.StringVar(&ConfigFilePath, "conf", "./config.yml", "path to config file")
	flag.BoolVar(&showVersion, "version", false, "shows version information and exists")
	//flag.StringVar(&cfg.HttpAddr, "http", "", "http to https redirector listen address")
	//flag.StringVar(&cfg.HttpsAddr, "https", ":443", "https listen address")
	//flag.StringVar(&cfg.TLSCertPath, "cert", "./tls/cert.pem", "tls certificate PEM file")
	//flag.StringVar(&cfg.TLSKeyPath, "key", "./tls/key.pem", "tls key PEM file")
	//flag.StringVar(&cfg.StaticDir, "staticDir", "", "directory with static webcontent")
}

func main() {
	flag.Parse()
	if showVersion {
		fmt.Printf("Version: %s, based on branch: %s (commit: %s), Buildtime: %s\n", Version, Branch, Commit, Buildtime)
		os.Exit(1)
	}
	if ConfigFilePath != "" {
		cfg = NewConfig(ConfigFilePath)
	}
	log.Printf("config:\n%+v\n", cfg)
	if cfg.StaticDir == "" {
		log.Fatal("no staticDir given but required")
	}

	if cfg.TLSCertPath == "" {
		log.Fatal("no certificate file given, but required")
	}

	if cfg.TLSKeyPath == "" {
		log.Fatal("no key file given but required")
	}

	app := &App{
		Addr:      cfg.HttpsAddr,
		StaticDir: cfg.StaticDir,
		TLSCert:   cfg.TLSCertPath,
		TLSKey:    cfg.TLSKeyPath,
	}

	// start a http server on httpAddr that just redirects to https
	if cfg.HttpAddr != "" {
		log.Printf("Starting a HTTP redirector\n")
		go func() {
			httpMux := http.NewServeMux()
			httpMux.HandleFunc("/", app.TLSRedirect)
			log.Println("Starting http server on ", cfg.HttpAddr)
			err := http.ListenAndServe(cfg.HttpAddr, httpMux)
			log.Fatal(err)
		}()
	}

	// start a https server on httpsAddr
	app.RunTLSServer()
}
