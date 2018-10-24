# TLSWebServer with multiple domain support

## Install

### Install the .deb package

The easiest way to install TLSWebserver is to use the released .deb or .rpm package.
You can find the packages in the [releases](https://github.com/scusi/TLSWebServer/releases) folder.
Please choose a package that corresponds to your architecture and operating system.

## Configure TLSWebServer

The configuration file is located in _/etc/TLSWebServer/config.json_.
You can edit that file to your needs. 
Please take care that you end up with a valid JSON file.

Basically there are just two sections.

- `ListenAddr` is the IP:Port combination TLSWebServer should listen for incoming connections.
- `TLSHosts` is a list of hostnames TLSWebServer serves.
   For each TLSHost you need to configure:
   - a `Hostname`
   - a path to the TLS certificate and key, namely `TLSCertPath` and `TLSKeyPath`
   - a `Webroot` directory, where the content is you want to serve under the given hostname.

### add another host to your config

If you want to add another hostname to your config, simply clone the TLSHosts section 
for local host and ajust it to your needs.

Your config file for example should look like this afterwards:
```
{
  "ListenAddr": ":8443",
  "TLSHosts": [
    {
      "Hostname": "localhost",
      "TLSCertPath": "/var/TLSWebServer/localhost/tls/cert.pem",
      "TLSKeyPath": "/var/TLSWebServer/localhost/tls/key.pem",
      "Webroot": "/var/TLSWebServer/localhost/www/"
    },
    {
      "Hostname": "yourexample.org",
      "TLSCertPath": "/var/TLSWebServer/yourexample.org/tls/cert.pem",
      "TLSKeyPath": "/var/TLSWebServer/yourexample.org/tls/key.pem",
      "Webroot": "/var/TLSWebServer/yourexample.org/www/"
    },
  ]
}

```

After adding a new TLSHost section to the config file make sure: 
- the certificates are in place and readable by the server process.
- the webroot does exist and have some content

When all of the above is done you can simply restart your TLSWebServer with the following command in order to apply changes.

```
sudo service tlswebserver restart
```

Point your browser to the URL for your configured hostname, e.g. `https://yourexample.org:8443/`

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

Please see [AutomaticCertificateRenewal.md](AutomaticCertificateRenewal.md)
