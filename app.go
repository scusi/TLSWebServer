package main

// App holds the settings for tls web server
type App struct {
	Addr      string
	StaticDir string
	TLSCert   string
	TLSKey    string
}
