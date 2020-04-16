package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"crypto/tls"
	"time"	
	"github.com/dimfeld/httptreemux"
	"strings"
)

/* // obsolete variables with configuratino details, now delivered via config file
var httpAddr string
var httpsAddr string
var tlsCertPath string
var tlsKeyPath string
var staticDir string
*/

var (
	Version             string     // Version set during compile time, e.g. v0.1.42
	Commit              string     // git commit, set during compiletime
	Branch              string     // git branch, set during compile time
	Buildtime           string     // compile timestamp
	usr                 *user.User // variable that holds the user environment
	ConfigFilePath      string     // holds path to config file
	ConfigFileLocations []string   // holds default locations to look for a config file
	cfg                 *Config    // pointer to configuration slice
	err                 error      // global error variable
	showVersion         bool       // if true programm prints version, commit, branch and buildtime, then exit
)

func init() {
	// determine current user
	usr, err := user.Current()
	if err != nil {
		log.Fatal("Can not determine current user: ", err)
	}
	// fill default config locations
	ConfigFileLocations = []string{
		"/etc/TLSWebserver/config.yml",
		"/usr/local/etc/TLSWebServer/config.yml",
		filepath.Join(usr.HomeDir, ".config/TLSWebServer/config.yml"),
		"./config.yml",
	}
	// create a new config with default values
	cfg = NewConfig("")

	flag.StringVar(&ConfigFilePath, "conf", "./config.yml", "path to config file")
	flag.BoolVar(&showVersion, "version", false, "shows version information and exists")
	//flag.StringVar(&cfg.HttpAddr, "http", "", "http to https redirector listen address")
	//flag.StringVar(&cfg.HttpsAddr, "https", ":443", "https listen address")
	//flag.StringVar(&cfg.TLSCertPath, "cert", "./tls/cert.pem", "tls certificate PEM file")
	//flag.StringVar(&cfg.TLSKeyPath, "key", "./tls/key.pem", "tls key PEM file")
	//flag.StringVar(&cfg.StaticDir, "staticDir", "", "directory with static webcontent")
}

// Redirect all traffic to HTTPS
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+cfg.domain+":"+cfg.httpsport+r.RequestURI, http.StatusMovedPermanently)
}
  
// Serve your index file
func indexHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	http.ServeFile(w, r, cfg.StaticDir+"/index.html")
}
  
  // Handle Static Files (files containing a .extension)
  // Send a not found error for any requests containing a .extension where the file does not exist
  // Redirect any requests that were not for files to your to your SPA
func staticHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	requestPath := r.URL.Path
	fileSystemPath := cfg.StaticDir + r.URL.Path
	endURIPath := strings.Split(requestPath, "/")[len(strings.Split(requestPath, "/"))-1]
	splitPath := strings.Split(endURIPath, ".")
	if len(splitPath) > 1 {
	  if f, err := os.Stat(fileSystemPath); err == nil && !f.IsDir() {
		http.ServeFile(w, r, fileSystemPath)
		return
	  }
	  //If not found just send back to index.html
	  //http.NotFound(w, r)
	  http.ServeFile(w, r, cfg.StaticDir+"/index.html")
	  return
	}
	http.ServeFile(w, r, cfg.StaticDir+"/index.html")
  }
  
  // Handle a simple healthcheck url for loadbalancers or whatnot
func healthCheckHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	fmt.Fprintf(w, "OK!\r")
}
  
func angularHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, cfg.StaticDir+"/index.html")
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
	if cfg.httpsport == "" {
		log.Fatal("no HTTPS given but required")
	}
	if cfg.httpport == "" {
		log.Fatal("no HTTP given but required")
	}
	if cfg.domain == "" {
		log.Fatal("no Domain given but required")
	}


	app := &App{
		domain: cfg.domain,
		httpport: cfg.httpport,
		httpsport: cfg.httpsport,
		StaticDir: cfg.StaticDir,
		TLSCert:   cfg.TLSCertPath,
		TLSKey:    cfg.TLSKeyPath,
	}

 
log.Printf("Starting Prime Frontend Server")
log.Printf("Listening for HTTP connections at: http://%v:%v", cfg.domain, cfg.httpport)
log.Printf("Listening for HTTPS connections at: https://%v:%v", cfg.domain, cfg.httpsport)

router := httptreemux.New()

tlsConfig := &tls.Config{
	MinVersion:               tls.VersionTLS12,
	MaxVersion:               tls.VersionTLS12,
	PreferServerCipherSuites: true,
	CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	CipherSuites: []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	},
}

srv := &http.Server{
	Addr:         app.domain+":"+app.httpsport,
	Handler:      router,
	TLSConfig:    tlsConfig,
	IdleTimeout:  time.Minute,
	ReadTimeout:  5 * time.Second,
	WriteTimeout: 10 * time.Second,
}
// make sure we can reload certificates and keys during runtime
kpr, err := NewKeypairReloader(cfg.TLSCertPath, cfg.TLSKeyPath)
if err != nil {
	log.Fatal(err)
}
srv.TLSConfig.GetCertificate = kpr.GetCertificateFunc()

// handle angular
router.NotFoundHandler = angularHandler
router.GET("/", indexHandler)
router.GET("/*", staticHandler)
router.GET("/healthcheck", healthCheckHandler)

go func() {
	if err :=  srv.ListenAndServeTLS("",""); err != nil {
		log.Fatalf("ListenAndServeTLS error: %v", err)
  	}
	}()
	if err := http.ListenAndServe(cfg.domain+":"+cfg.httpport, http.HandlerFunc(redirectHandler)); err != nil {
  		log.Fatalf("ListenAndServe error: %v", err)
	}
}
