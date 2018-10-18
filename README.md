# TLSWebServer

Is a very simple TLS webserver written in golang.

TLSWebServer is meant to deliver content only via HTTPS connections.

## Features

- Secure TLS settings, by default
- redirects http requests to https
- protection against sloworis kind of attacks
- fast and well scaleing

## Getting started

```TLSWebServer -http=:80 -https=:443 -staticDir=/usr/var/www -cert=./tls/cert.pem -key=./tls/key.pem```

The above command will: 
- start a http server, that just redirects every request to https on port 80 on all interfaces.
- start a https server on port 443 on all interfaces, serving files from _/usr/var/www_ and 
  useing the given cert anf key for TLS.


