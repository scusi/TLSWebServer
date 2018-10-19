# TLSWebServer

A very simple TLS webserver written in golang.

## Features

- Secure TLS settings, by default
- redirects http requests to https
- protection against sloworis kind of attacks
- fast and well scaleing
- can reload certificates on the fly, no downtime

## Getting started

### make sure DNS records point to your IP

Time to make sure that the DNS A record for your host is pointting to the IP where you want to run TLSWebServer.
How to do this depends on your setup, please find out yourself.

### create a selfsigned certificate

Before you can start a TLSWebServer you need a certificate. In the following you learn how to create a selfsigned certificate, to get started.
For productive use a certificate signed by a acknowledged root CA will be more suitable. Please see [AutomaticCertRenewal.md] document how you can archive that.

```
mkdir -p tls
cd tls
go run $GOROOT/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

### Start TLSWebserver

```TLSWebServer -http=:80 -https=:443 -staticDir=/usr/var/www -cert=/home/you/.acme.sh/example.com/fullchain.cer -key=/home/you/.acme.sh/example.com/example.com.key```

The above command will:
- start a http server, that just redirects every request to https on port 80
  on all interfaces.
- start a https server on port 443 on all interfaces, serving files from
  `/usr/var/www` and useing the given cert anf key for TLS.

Note: If the _-http_ flag is omitted no http server will be started.

## Create a service file for TLSWebServer

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

## The follwoing sites use TLSWebServer

- [scusiblog.org](https://scusiblog.org)
