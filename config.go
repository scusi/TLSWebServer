package main

import (
	"github.com/tsuru/config"
	"log"
)

// Config holds the configuration settings
type Config struct {
	HttpAddr         string // listening address of the http -> https redirector
	ExposedHttpAddr  string // the extrernaly exposed HTTP Address (Host:Port)
	ExposedHttpsAddr string // the externally exposed TLS/HTTPS Address (Host:Port)
	HttpsAddr        string // listening address of the tls web server
	TLSCertPath      string // path to certificate file, PEM encoded
	TLSKeyPath       string // path to key file, PEM encoded, without passphrase
	StaticDir        string // path to webroot directory
}

// NewConfig - load a config file from a given path,
// if a empty path is given it will search a config file from the default config locations,
// if no config file could be found it returns a config with default values.
func NewConfig(path string) (cfg *Config) {
	cfg = new(Config)
	if path == "" {
		for _, cfl := range ConfigFileLocations {
			err := config.ReadAndWatchConfigFile(cfl)
			if err != nil {
				log.Printf("Could not read config file from '%s'\n", cfl)
				continue
			}
			log.Println("Found config file at '" + cfl + "' useing it.")
			break
		}
	} else {
		err := config.ReadAndWatchConfigFile(path)
		if err != nil {
			log.Println("Read custom config file error: " + err.Error())
			log.Fatal(err)
		}
	}
	cfg.HttpAddr, err = config.GetString("HttpAddr")
	if err != nil {
		log.Println("Config Error: " + err.Error())
		cfg.HttpAddr = ":8080" // set default value 'HttpAddr=":8080"'
	}
	cfg.ExposedHttpAddr, err = config.GetString("ExposedHttpAddr")
	if err != nil {
		log.Printf("Config Error: " + err.Error())
		cfg.ExposedHttpAddr = ":80" // set default value ":80"
	}

	cfg.HttpsAddr, err = config.GetString("HttpsAddr")
	if err != nil {
		log.Println("Config Error: " + err.Error())
		cfg.HttpsAddr = ":8443" // set default value 'HttpsAddr=":8443"'
	}
	cfg.ExposedHttpsAddr, err = config.GetString("ExposedHttpsAddr")
	if err != nil {
		log.Printf("Config Error: " + err.Error())
		cfg.ExposedHttpsAddr = ":443" // set default value ":80"
	}

	cfg.TLSCertPath, err = config.GetString("TLSCertPath")
	if err != nil {
		log.Println("Config Error: " + err.Error())
		cfg.TLSCertPath = "tls/cert.pem" // set default certificate path
	}

	cfg.TLSKeyPath, err = config.GetString("TLSKeyPath")
	if err != nil {
		log.Println("Config Error: " + err.Error())
		cfg.TLSKeyPath = "tls/key.pem" // set default key path
	}

	cfg.StaticDir, err = config.GetString("StaticDir")
	if err != nil {
		log.Println("Config Error: " + err.Error())
		cfg.StaticDir = "www" // set default path to webroot
	}
	return cfg
}
