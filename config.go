package main

import (
	"github.com/tsuru/config"
	"log"
)

type Config struct {
	HttpAddr    string
	HttpsAddr   string
	TLSCertPath string
	TLSKeyPath  string
	StaticDir   string
}

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
		cfg.HttpAddr = ":80"
	}

	cfg.HttpsAddr, err = config.GetString("HttpsAddr")
	if err != nil {
		log.Println("Config Error: " + err.Error())
		cfg.HttpsAddr = ":443"
	}

	cfg.TLSCertPath, err = config.GetString("TLSCertPath")
	if err != nil {
		log.Println("Config Error: " + err.Error())
		cfg.TLSCertPath = "tls/cert.pem"
	}

	cfg.TLSKeyPath, err = config.GetString("TLSKeyPath")
	if err != nil {
		log.Println("Config Error: " + err.Error())
		cfg.TLSKeyPath = "tls/key.pem"
	}

	cfg.StaticDir, err = config.GetString("StaticDir")
	if err != nil {
		log.Println("Config Error: " + err.Error())
		cfg.StaticDir = "www"
	}
	return cfg
}
