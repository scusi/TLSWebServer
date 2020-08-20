package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	//"os/user"
	//"path/filepath"
)

/* // obsolete variables with configuratino details, now delivered via config file
var httpAddr string
var httpsAddr string
var tlsCertPath string
var tlsKeyPath string
var staticDir string
*/

var (
	// Version set during compile time, e.g. v0.1.42
	Version   string
	Commit    string // git commit, set during compiletime
	Branch    string // git branch, set during compile time
	Buildtime string // compile timestamp
	//usr                 *user.User // variable that holds the user environment
	// ConfigFilePath holds path to config file
	ConfigFilePath string
	// ConfigFileLocations holds default locations to look for a config file
	ConfigFileLocations []string
	// cfg pointer to configuration slice
	cfg *Config
	// err is a global error variable
	err error
	// showVersion, will show version info and exit if true
	showVersion bool
)

func init() {
	// determine current user
	/*usr, err := user.Current()
	if err != nil {
		log.Printf("Can not determine current user: ", err)
	}
	*/
	// fill default config locations
	ConfigFileLocations = []string{
		"/app/config/config.yml"
		"/etc/TLSWebserver/config.yml",
		"/usr/local/etc/TLSWebServer/config.yml",
		//filepath.Join(usr.HomeDir, ".config/TLSWebServer/config.yml"),
		"./config.yml",
	}
	// create a new config with default values
	cfg = NewConfig("")

	flag.StringVar(&ConfigFilePath, "conf", os.Getenv("CONFIG"), "path to config file")
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
