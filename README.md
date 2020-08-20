# TLSWebServer (master branch)

A very simple TLS webserver written in golang.

## Features

- Secure TLS settings, by default
  - TLSv1.2 only
  - uses only strong ciphers with forward secrecy
- redirects http requests to https
- protection against slowloris kind of attacks
- can reload certificates on the fly, no downtime

## Docker

There is a docker image (scusi/tls-webserver:latest) that can be used to run TLSWebserver.

You can test run it like this:
```
docker pull scusi/tls-webserver:latest
docker run -e CONFIG=/app/config/config.yml -p 80:8080 -p 443:8443 tls-webserver:latest
```

In practice you probably want to keep the webroot and the config files persistant.
This can be done useing volumes.
The following docker command starts TSLWebserver with two volumes, 
one for the config files and one for the webroot.

The local folder `./www` will be mounted into the docker container at `/www`.
The local folder `./config` will be mounted into the container at `/app/config`.
In the config folder there are usually thre files, a config.yml, the TLS certificate in PEM format and the TLS key, also in PEM format and without a passphrase.

```
docker run -e CONFIG=/app/config/config.yml \
	-p 80:8080 -p 443:8443 \
	-v `pwd`/www:/www -v `pwd`/config:/app/config \
	scusi/tls-webserver:latest
```

## Build instructions

This software comes with a ```Makefile```, so you can simply say:

```make build``` to actually biuild the binary on your platform.

If you just execute ```make``` you build the binary and directly start it.
Befor you do that read the rest of this document.

## Getting started

### make sure DNS records point to your IP

Time to make sure that the DNS A record for your host is pointting to the IP where you want to run TLSWebServer.
How to do this depends on your setup, please find out yourself.

### create a selfsigned certificate

Before you can start a TLSWebServer you need a certificate. In the following you learn how to create a selfsigned certificate, to get started.
For productive use a certificate signed by a acknowledged root CA will be more suitable. Please see [AutomaticCertRenewal.md](AutomaticCertRenewal.md) document how you can archive that.

```
mkdir -p tls
cd tls
go run $GOROOT/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

### Allow TLSWebServer to bind to privileged ports
When you want to use the default ports for HTTP (port 80) and HTTPS (port 443) 
on most unice operating systems you need to have root privileges. In order to 
avoid that - on linux - you can do:
```
sudo setcap cap_net_bind_service=ep /path/to/TLSWebServer
```
This will allow your TLSWebServer binary to bind to ports below 1024.
This needs to be redone when the binary was updated.

Otherwise you have to run TLSWebServer with root privileges or use non default ports above 1024, like 8080 (HTTP-ALT) and 8443 (HTTPS-ALT).

### Prepare a TLSWebserver config file

Create a file called _config.yml_ with a content like in the following listing.

```
HttpAddr: 127.0.0.1:8080
HttpsAddr: 127.0.0.1:8443
TLSCertPath: tls/cert.pem
TLSKeyPath: tls/key.pem
StaticDir: /path/to/your/webroot
```

### Start TLSWebserver

You can start a TLSWebServer now with your config, like this:

```TLSWebServer -conf /path/to/your/config.yml```

Note: If the _-conf_ flag is omitted TLSWebServer will search for a config file in the following loctions:
- /etc/TLSWebServer/config.yml
- /usr/local/etc/TLSWebServer/config.yml
- ~/.config/TLSWebServer/config.yml
- ./config.yaml 

If no config file could be found on the above default locations and there was none given on the command line it will use a default configuration.

### Create a service file for TLSWebServer

Edit the follwoing content to your needs and save it under your systemd service directory as `tlswebserver.service`. The systemd directory on most linux systems is usually `/etc/systemd/system/`.

Make sure it is executable.
You can make it executable with the following command:
```
chmod 755 /etc/systemd/system/tlswebserver.service
```

### systemd service file _example

The following listing contains an example service file.
Please make sure you change the pathes and filenames in the _**ExecStart**_ variable to match your environment.

```
[Unit]
Description=tlswebserver
After=network.target

[Service]
Type=simple
User=root
ExecStart=/home/you/TLSWebServer-linux-x86.bin -http 89.238.79.203:80 -https 89.238.79.203:443 -staticDir /home/you/testWeb/www/ -cert /home/you/.acme.sh/example.com/fullchain.cer -key /home/you/.acme.sh/example.com/example.com.key
Restart=on-failure

[Install]
WantedBy=multi-user.target
```
After that you should be able to start and stop your tlswebserver like this:

```
sudo service tlswebserver start
```

```
sudo service tlswebserver stop
```

You can also seee the current status with:

```
sudo service tlswebserver status
```
