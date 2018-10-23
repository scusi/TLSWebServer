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
Exchange _scusiblog.org_ with your domain in all following examples.

```
acme.sh --issue -d scusiblog.org -w /home/jupp/tlsWeb/www
```

### Install your certificate

```
acme.sh --install-cert -d scusiblog.org \
 --key-file /home/jupp/testWeb/tls/key.pem \
 --fullchain-file /home/jupp/testWeb/tls/cert.pem \
 --reloadcmd "sudo service tlswebserver force-reload"
```

After this your certificates will be renewed automattically by acme.sh.
Cool, isn't it?

### Renew Your Certificate - forcefully

If you ever need to renew your certifcate before automatic renewal, you can do:
```
acme.sh --renew -d scusiblog.org -w /home/jupp/tlsWeb/www/ --force
```

### Manually force TLSWebServer to reload it's certificate and key

In case you want manually to force your TLSWebServer process to reload the 
certificate you can do as described in the this section.

#### Step 1

Find the PID of your running TLSWebServer process.
```
$> ps ax | grep TLSWebServer
14169 pts/19   Sl+    0:00 TLSWebServer -staticDir /var/www/public/ -https :8443
```
The PID is that number in the first column, 14169 in this example.

#### Step 2

Send a HUP signal to that PID, in order to force reloading certificate and key.

```
kill -s HUP 14169
```

There is also a small command shipping with TLSWebServer that does Step 1 and Step 2 for you.
Just start it without any arguments, like in the example below.
```
$ ./reloadCerts
2018/10/19 21:37:57 found a 'TLSWebServer' process, PID = 2694
2018/10/19 21:37:57 sent HUP signal to PID: 2694
```
#### Verification

Your TLSWebServer process will output a line similar to the following one, on standard error (stderr):
```
2018/10/19 12:15:43 Received SIGHUP, reloading TLS certificate and key from "./tls/cert.pem" and "./tls/key.pem"
```

When you connect to your TLSWebServer process with a browser you should get 
the new certificate.
You can also do it from the commandline, like this:
```
openssl s_client -connect localhost:8443 < /dev/null 2>/dev/null | openssl x509 -fingerprint -noout -in /dev/stdin
```
In the above example you have to adjust your hostname and port (localhost:8443) 
to your needs.
