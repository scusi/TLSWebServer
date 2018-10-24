# TLSWebServer with multiple domain support

## Install

### Install the .deb package

The easiest way to install TLSWebserver is to use the released .deb or .rpm package.
You can find the packages in the [dist](dist/) folder.
Please choose a package that corresponds to your architecture and operating system.

## Configure TLSWebServer

The configuration file is located in _/etc/TLSWebServer/config.json_.
You can edit that file to your needs. Please take care that you end up with a valif JSON file.

## Certificates

Certificates and corresponding key are based in _/var/TLSWebServer/{hostname}/tls/_.
Where _{hostname}_ is your configured hostname, e.g. _localhost_.

During Setup a self signed certificate for localhost will be generated and moved into the approriate location.

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
