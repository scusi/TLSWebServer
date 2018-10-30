# TLSWebServer with multiple domain support

TLSWebserver is a two daemon package that combines a tls only webserver with
a scp (secure copy protocol) server.

It's purpose is to run a secure webserver easily.

## Design Goals

- TLS certificates can be updated automatically.
- Webcontent of hosts can easily be updated over the network, 
  also by 3rd party users.
- relative secure solution.
- sane TLS settings

## Install

Please see [INSTALL.md](INSTALL.md)

## Certificates

Certificates and corresponding key are based in _/var/TLSWebServer/{hostname}/tls/_.
Where _{hostname}_ is your configured hostname, e.g. _localhost_.

During Setup a self signed certificate for localhost will be generated and moved 
into the appropriate location.

## Starting and Stoping TLSWebServer

During Install a systemd service file was installed unter _/etc/systemd/system/tlswebserver.service_.
Hence you can start your TLSWebServer like any other installed service.

To start TLSWebServer use:
```
sudo service tlswebserver start
```

To see the current status of the service use:
```
sudo service tlswebserver status
```

You can stop the service like this:

```
sudo service tlswebserver stop
```

## Automatic Certificate Renewal

Please see [AutomaticCertificateRenewal.md](AutomaticCertificateRenewal.md) for information how you can use _acme.sh_ and _Let's Encrypt_ in order to have valid signed TLS certificates which renew automatically.
