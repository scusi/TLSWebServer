// example TLS webserver that supports multiple domains
package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type TLSHost struct {
	Hostname    string
	TLSCertPath string
	TLSKeyPath  string
	Webroot     string
}

type Config struct {
	ListenAddr string
	TLSHosts   []TLSHost
}

var configFile string

func init() {
	flag.StringVar(&configFile, "conf", "./config", "config file to load")
}

func WriteConfig(filename string, conf Config) (err error) {
	jBytes, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return
	}
	err = ioutil.WriteFile("config.json", jBytes, 0700)
	if err != nil {
		return
	}
	log.Printf("config written to: %s\n", "config.json")
	return
}

func ReadConfig(filename string) (conf Config, err error) {
	cfg := &Config{}
	jBytes, err := ioutil.ReadFile("config.json")
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
		log.Printf(pattern, r.RemoteAddr, r.Host, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func main() {
	flag.Parse()
	t := log.Logger{}
	var err error
	//conf, err := config.ReadConfigFile(configFile)
	conf, err := ReadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("ReadConf:\n%+v\n", conf)

	if err != nil {
		t.Fatal(err)
	}
	//log.Printf("read conf:\n%+v\n", conf)
	// setup routing for the TLSHosts
	mux := http.NewServeMux()
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
		//log.Printf("host:\n%+v\n", host)
		tlsConfig.Certificates[i], err = tls.LoadX509KeyPair(host.TLSCertPath, host.TLSKeyPath)
		if err != nil {
			t.Fatal(err)
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
	listener, err := tls.Listen("tcp", conf.ListenAddr, tlsConfig)
	if err != nil {
		t.Fatal(err)

	}
	log.Fatal(server.Serve(listener))
}
