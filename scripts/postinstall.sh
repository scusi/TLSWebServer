#!/bin/sh

/usr/bin/openssl req -x509 -subj '/CN=localhost' -nodes -days 365 -newkey rsa:2048 -keyout /var/TLSWebServer/localhost/tls/key.pem -out /var/TLSWebServer/localhost/tls/cert.pem
