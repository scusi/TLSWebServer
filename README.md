# TLSWebServer

A very simple TLS webserver written in golang.

## Features

- Secure TLS settings, by default
- redirects http requests to https
- protection against sloworis kind of attacks
- fast and well scaleing

## Getting started

### Install acme.sh

I recommend acme.sh as a tool to create and renew TLS certificates wirh _Let's encrypt_.

Install acme.sh on the same machine you want to run TLSWebServer like this:
```
curl https://get.acme.sh | sh
```

### make sure DNS records point to your IP

Time to make sure that the DNS A record for your host is pointting to the IP where you want to run TLSWebServer.
How to do this depends on your setup, please find out yourself.


### create a selfsigned certificate

```
mkdir -p tls
cd tls
go run $GOROOT/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

### Get a TLS certificate

Issue an initial certificate for _example.com_

```
acme.sh --issue -d example.com --standalone
```

### Start TLSWebserver

```TLSWebServer -http=:80 -https=:443 -staticDir=/usr/var/www -cert=/home/you/.acme.sh/example.com/fullchain.cer -key=/home/you/.acme.sh/example.com/example.com.key```

The above command will:
- start a http server, that just redirects every request to https on port 80
  on all interfaces.
- start a https server on port 443 on all interfaces, serving files from
  `/usr/var/www` and useing the given cert anf key for TLS.

Note: If the _-http_ flag is omitted no http server will be started.

### Renew Your Certificate - once manual

In order to make sure that acme.sh knows where your webroot is for all future renewals,
once renew your certificate with the `-w` flag and your webroot directory.

```
acme.sh --renew -d example.com -w /usr/var/www
```

NOTE: if you change your webroot directory you need to redo this step.

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
