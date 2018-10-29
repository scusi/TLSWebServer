// example TLS webserver that supports multiple domains
package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type TLSHost struct {
	Hostname    string
	TLSCertPath string
	TLSKeyPath  string
	Webroot     string
}

type Config struct {
	Logfile      string
	Loglevel     string // Trace, Info, Warning, Error
	ScpAddr      string
	TLSAddr      string
	RedirectHttp bool
	TLSHosts     []TLSHost
}

var configFile string
var showVersion bool
var addHost bool
var delHost bool
var newHost *TLSHost
var err error
var logger *log.Logger

func init() {
	flag.StringVar(&configFile, "conf", "/etc/TLSWebServer/config.json", "config file to load")
	flag.BoolVar(&showVersion, "version", false, "shows version info and exits")
	flag.BoolVar(&addHost, "add", false, "add a TLSHost")
	flag.BoolVar(&delHost, "del", false, "delete a TLSHost")
	newHost = new(TLSHost)
	flag.StringVar(&newHost.Hostname, "host", "", "hostname to add or delete")
	flag.StringVar(&newHost.TLSCertPath, "cert", "", "path to cert")
	flag.StringVar(&newHost.TLSKeyPath, "key", "", "path to key file")
	flag.StringVar(&newHost.Webroot, "w", "", "path to webroot")
}

func (conf *Config) TLSRedirect(w http.ResponseWriter, r *http.Request) {
	for _, h := range conf.TLSHosts {
		if r.Host == h.Hostname {
			http.Redirect(w, r,
				"https://"+r.Host+r.URL.String(),
				http.StatusMovedPermanently)
			pattern := `%s - "%s %s %s %s"`
			Info.Printf(pattern, r.RemoteAddr, r.Host, r.Proto, r.Method, r.URL.RequestURI())
			return
		}
	}
	http.NotFound(w, r)
	return
}

func main() {
	flag.Parse()
	//t := log.Logger{}
	// configure logfile to stdout, for now
	logger, err = AllLogger("-")
	if err != nil {
		Error.Println(err)
	}
	if showVersion {
		fmt.Printf("TLSWebServer Version: %s, Commit: %s, Builddate: %s\n", version, commit, date)
		return
	}
	logger.Printf("TLSWebServer Version: %s, Commit: %s, Builddate: %s\n", version, commit, date)
	conf, err := ReadConfig(configFile)
	if err != nil {
		logger.Fatal(err)
	}
	// reconfig logfile from config value
	logger, err = AllLogger(conf.Logfile)
	if err != nil {
		Error.Println(err)
	}

	if addHost {
		logger.Printf("going to add host %s\n", newHost.Hostname)
		conf.TLSHosts = append(conf.TLSHosts, *newHost)
		WriteConfig(configFile, conf)
		return
	}

	if delHost {
		logger.Printf("going to delete host %s\n", newHost.Hostname)
		var tlsHosts = conf.TLSHosts
		for i, h := range tlsHosts {
			if h.Hostname == newHost.Hostname {
				log.Printf("found host to delete at position %d in conf.TLSHosts\n", i)
				if len(tlsHosts) <= i+1 {
					tlsHosts = tlsHosts[:i]
				} else {
					tlsHosts = append(tlsHosts[:i], tlsHosts[i+1:]...)
				}
				logger.Printf("after delete conf.TLSHosts looks like: %+v\n", tlsHosts)
			}
		}
		//log.Printf("after range tlsHosts looks like: %+v\n", tlsHosts)
		conf.TLSHosts = tlsHosts
		WriteConfig(configFile, conf)
		return
	}
	// setup a http to https redirector on port 80
	if conf.RedirectHttp == true {
		logger.Printf("Starting a HTTP redirector\n")
		go func() {
			httpMux := http.NewServeMux()
			httpMux.HandleFunc("/", conf.TLSRedirect)
			host, _, err := net.SplitHostPort(conf.TLSAddr)
			httpAddr := net.JoinHostPort(host, "80")
			logger.Printf("Starting http server on %s\n", httpAddr)
			err = http.ListenAndServe(httpAddr, httpMux)
			logger.Fatal(err)
		}()
	}
	// setup routing for the TLSHosts
	mux := http.NewServeMux()
	// you can see a version info if you request it like with the below curl command
	// curl -k --resolve tlswebserver:8443:127.0.0.1 https://tlswebserver:8443/server-version/
	mux.HandleFunc("tlswebserver/server-version/", Version)
	for _, host := range conf.TLSHosts {
		mux.Handle(host.Hostname+"/", LogRequest(http.FileServer(http.Dir(host.Webroot))))
	}

	// setup a TLS server
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
	tlsConfig.Certificates = make([]tls.Certificate, len(conf.TLSHosts))
	// go http server treats the 0'th key as a default fallback key
	for i, host := range conf.TLSHosts {
		tlsConfig.Certificates[i], err = tls.LoadX509KeyPair(host.TLSCertPath, host.TLSKeyPath)
		if err != nil {
			logger.Fatal(err)
		}
	}
	tlsConfig.BuildNameToCertificate()

	server := &http.Server{
		IdleTimeout:    time.Minute,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      tlsConfig,
		Handler:        mux,
	}
	listener, err := tls.Listen("tcp", conf.TLSAddr, tlsConfig)
	if err != nil {
		logger.Fatal(err)

	}
	logger.Printf("TLS listener started at '%s'\n", conf.TLSAddr)
	logger.Fatal(server.Serve(listener))
}

func WriteConfig(filename string, conf Config) (err error) {
	jBytes, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return
	}
	err = ioutil.WriteFile(filename, jBytes, 0700)
	if err != nil {
		return
	}
	logger.Printf("config written to: %s\n", filename)
	return
}

func ReadConfig(filename string) (conf Config, err error) {
	cfg := &Config{}
	jBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = json.Unmarshal(jBytes, cfg)
	if err != nil {
		return
	}
	return *cfg, nil
}

func LogRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		pattern := `%s - "%s %s %s %s"`
		logger.Printf(pattern, r.RemoteAddr, r.Host, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func Version(w http.ResponseWriter, r *http.Request) {
	pattern := `%s - "%s %s %s %s"`
	logger.Printf(pattern, r.RemoteAddr, r.Host, r.Proto, r.Method, r.URL.RequestURI())
	fmt.Fprintf(w, "TLSWebServer (multiDomain)\nVersion: %v,\nCommit %v,\nbuilt at %v\n", version, commit, date)
}
