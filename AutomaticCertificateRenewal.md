# Automatic Certificate renewal

In order to have automatic certificate renewal you can use acme.sh

### Install acme.sh

I recommend acme.sh as a tool to create and renew TLS certificates wirh _Let's encrypt_.

Install acme.sh on the same machine you want to run TLSWebServer like this:
```
curl https://get.acme.sh | sh
```

### Get a TLS certificate

Please note the following conditions:
- make sure TLSWebServer is running.
- you need to be on a host reachable from the internet on port 80 and 443.
  On a internet server this should be no problem, if not firewalled.
  On a home network behind NAT you need to configure port forwardings in your router first.
- your domain names DNS A record has to point to the IP of your host.

Issue an initial certificate for your domain.
Exchange _yourexample.org_ with your domain in all following examples.

```
acme.sh --issue -d yourexample.org -w /var/TLSWebServer/yourexample.org/gwww
```

### Install your certificate

```
acme.sh --install-cert -d yourexample.org \
 --key-file /var/TLSWebServer/yourexample.org/gtls/key.pem \
 --fullchain-file /var/TLSWebServer/yourexample.org/gtls/cert.pem \
 --reloadcmd "sudo service tlswebserver force-reload"
```

After this your certificates will be renewed automattically by acme.sh.
Cool, isn't it?

### Renew Your Certificate - forcefully

If you ever need to renew your certifcate before automatic renewal, you can do:
```
acme.sh --renew -d yourexample.org -w /var/TLSWebServer/yourexample.org/gwww/ --force
```
