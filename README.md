# TLSWebServer with multiple domain support

This is a proof of concept TLSWebserver that supports multiple domains with it's own certificate each.
This variant of TLSWebServer does not support reloading certificates on the fly. You have to restart the server as a whole.

## Config file
save the following listing as _config.json_ in the local directory
```
{
  "ListenAddr": ":8443",
  "TLSHosts": [
    {
      "Hostname": "localhost",
      "TLSCertPath": "localhost/tls/cert.pem",
      "TLSKeyPath": "localhost/tls/key.pem",
      "Webroot": "localhost/www/"
    },
    {
      "Hostname": "test1.test",
      "TLSCertPath": "test1/tls/cert.pem",
      "TLSKeyPath": "test1/tls/key.pem",
      "Webroot": "test1/www"
    },
    {
      "Hostname": "test2.test",
      "TLSCertPath": "test2/tls/cert.pem",
      "TLSKeyPath": "test2/tls/key.pem",
      "Webroot": "test2/www"
    }
  ]
}

```
## Setup
```
// create directory structure
mkdir -p localhost/tls localhost/www
mkdir -p test1/tls test1/www
mkdir -p test2/tls test2/www

// create certificates
cd localhost/tls && go run $GOROOT/src/crypto/tls/generate_cert.go -rsa-bits 4096 -host localhost && cd ../../
cd test1/tls && go run $GOROOT/src/crypto/tls/generate_cert.go -rsa-bits 4096 -host test1.test && cd ../../
cd test2/tls && go run $GOROOT/src/crypto/tls/generate_cert.go -rsa-bits 4096 -host test2.test && cd ../../

// create an index.html in the webroot for each host
echo "localhost speaking" >> localhost/www/index.html
echo "test1.test speaking" >> test1/www/index.html
echo "test2.test speaking" >> test2/www/index.html

// run the proof of concept
go run main.go
```

## Verify

In another terminal you can check if it worked.

```
openssl s_client -servername localhost -connect localhost:8443 < /dev/null 2>/dev/null | openssl x509 -fingerprint -noout -in /dev/stdin
openssl s_client -servername test1.test -connect test1.test:8443 < /dev/null 2>/dev/null | openssl x509 -fingerprint -noout -in /dev/stdin
openssl s_client -servername test2.test -connect test2.test:8443 < /dev/null 2>/dev/null | openssl x509 -fingerprint -noout -in /dev/stdin
```

You should see three different fingerprints.

Note: The TLSHost which comes first in the config ('localhost' in above example) is the default host. In case of confusion this certificate will be presented.

