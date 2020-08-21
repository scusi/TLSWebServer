# TLSWebServer (master branch)

A very simple TLS webserver written in golang.

## Features

- Secure TLS settings, by default
  - TLSv1.2 only
  - uses only strong ciphers with forward secrecy
- redirects http requests to https
- protection against slowloris kind of attacks
- can reload certificates on the fly, no downtime

## Build instructions

This software comes with a ```Makefile```, so you can simply say:

```make build``` to actually biuild the binary on your platform.

If you just execute ```make``` you build the binary and directly start it.
Before you do that read the rest of this document.

### build a docker image

Use the provided [Dockerfile](Dockerfile) to build a docker image.

```
git clone git@github.com:scusi/TLSWebServer.git
cd TLSWebserver
docker build -t tls-webserver:my_build .
```

## run TLSWebserver in docker

There is a docker image ([scusi/tls-webserver:latest](https://hub.docker.com/r/scusi/tls-webserver)) that can be used to run TLSWebserver.

If you have docker installed you can test run it like this:

```
docker run -e CONFIG=/app/config/config.yml -p 80:8080 -p 443:8443 scusi/tls-webserver:latest
```

In practice you probably want to keep the webroot and the config files persistant.
This can be done useing volumes.
The following docker command starts TSLWebserver with two volumes, 
one for the config files and one for the webroot.

The local folder `./www` will be mounted into the docker container at `/www`.
The local folder `./config` will be mounted into the container at `/app/config`.
In the config folder there are usually three files, a config.yml, the TLS certificate in PEM format and the TLS key, also in PEM format and without a passphrase.

```
docker run -e CONFIG=/app/config/config.yml \
	-p 80:8080 -p 443:8443 \
	-v `pwd`/www:/www -v `pwd`/config:/app/config \
	scusi/tls-webserver:latest
```

## Getting started with TLSWebserver 

The following describes the steps to install and run TLSWebserver on a unix like operating system.

### make sure DNS records point to your IP

Time to make sure that the DNS A record for your host is pointting to the IP where you want to run TLSWebServer.
How to do this depends on your setup, please find out yourself.

### create a selfsigned certificate

Before you can start a TLSWebServer you need a certificate. In the following you learn how to create a selfsigned certificate, to get started.
For productive use a certificate signed by a acknowledged root CA will be more suitable. 
Please see [AutomaticCertRenewal.md](AutomaticCertRenewal.md) document how you can archive that.

```
mkdir -p tls
cd tls
go run $GOROOT/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

**NOTE**: The docker image already comes with self-signed certificate. 
However, this certificate is the same for everyone useing the docker image and
*NOT* meant for production use. For production use you should have your custom
certificates in the config directory and mount it as a volume into the docker
container. As shown in the docker example above.

### Allow TLSWebServer to bind to privileged ports

When you want to use the default ports for HTTP (port 80) and HTTPS (port 443) 
on most unice operating systems you need to have root privileges. In order to 
avoid that - on linux - you can do:
```
sudo setcap cap_net_bind_service=ep /path/to/TLSWebServer
```
This will allow your TLSWebServer binary to bind to ports below 1024.
This needs to be redone when the binary was updated.

Otherwise you have to run TLSWebServer with root privileges (*NOT RECOMMENDED*)
or use non default ports above 1024, like 8080 (HTTP-ALT) and 8443 (HTTPS-ALT).

### Prepare a TLSWebserver config file

Create a file called _config.yml_ with a content like in the following listing.

```
# HttpAddr is the listening address for http (usually :80).
# On http TLSWebserver does only redirects to https.
# HttpAddr is optional and can be empty
# If HttpAddr is empty, no http->https redirector will be started.
HttpAddr: 127.0.0.1:8080

# ExposedHttpAddr defines the external http endpoint (host:port).
# ExposedHttpAddr is optional and can be empty.
# ExposedHttpAddr is only needed in scenarios where all following conditions are met.
# * the http redirector feature is enabled (aka, HttpAddr is not empty)
# * the server is listening on an internal address (not directly reachable by clients) 
# * the exposed (external) address is not on the default port (80)
ExposedHttpAddr: :8080

# HttpsAddr defines the listening address for https (usually :433)
# HttpsAddr is required. 
HttpsAddr: 127.0.0.1:8443

# ExposedHttpsAddr is the same as ExposedHttpAddr, just for TLS / HTTPS.
# ExposedHttpAddr is optional and can be empty.
# ExposedHttpsAddr is only needed in scenarios where all following conditions are met.
# * HttpsAddr is on an internal address (not directly reachable by clients) 
# * the exposed (external) address is not on the default port (443)
ExposedHttpsAddr: :8433	

# TLSCertPath defines the path where the tls certificate the server uses is found.
# The certificate needs to be in PEM encoded format.
# TLSCertPath must be an absolute path
# TLSCertPath is required
TLSCertPath: /path/to/tls/cert.pem

# TLSKeyPath defines the path where the tls key the server uses is found.
# The key must not have a passphrase and needs to be in PEM encoded format.
# TLSKeyPath is required
TLSKeyPath: /path/to/tls/key.pem

# StaticDir defines the webroot, where all the files to be served are located.
# StaticDir is required
StaticDir: /path/to/your/webroot	# path to the webroot directory
```

### Start TLSWebserver

You can start a TLSWebServer now with your config, like this:

```TLSWebServer -conf /path/to/your/config.yml```

Note: the default value of the `-conf` flag can be set with an environment
variable called `CONFIG`. So you can start TLSWebserver with no argument, but
still with a custom config, as long as `CONFIG` environment variable contains a
path to a valid config file.

```
$>: export CONFIG="/app/config-files/config.yml"
$>: ./TLSWebserver
```

Note: If the _-conf_ flag is omitted and no `CONFIG` variable was set, TLSWebServer will search for a config file in the following loctions, in that order.

- /app/config/config.yml
- /etc/TLSWebServer/config.yml
- /usr/local/etc/TLSWebServer/config.yml
- ./config.yaml 

The first file beeing found is taken. If no config file could be found on the
above default locations and there was none given on the command line it will
use a default configuration.

### Create a service file for TLSWebServer

In order to start TLSWebserver at system startup you need a startup file. This section shows an example of how to do that on a linux / systemd machine.
Edit the following content to your needs and save it under your systemd service directory as `tlswebserver.service`. The systemd directory on most linux systems is usually `/etc/systemd/system/`.

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
ExecStart=/path/to/binary/of/TLSWebServer -c /etc/TLSWebserver/config.yml 
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

If you want it to start automatically on boot-up, run:

```
systemctl enable tlswebserver
```
