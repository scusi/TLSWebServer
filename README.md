# TLSWebServer with multiple domain support

## Setup
```
// create directory structure
mkdir -p test1/tls test1/www
mkdir -p test2/tls test2/www

// create certificates
cd test1/tls && go run $GOROOT/src/crypto/tls/generate_cert.go -rsa-bits 4096 -host test1.test && cd ../../
cd test2/tls && go run $GOROOT/src/crypto/tls/generate_cert.go -rsa-bits 4096 -host test2.test && cd ../../

// create an index.html in the webroot
echo "test1.test speaking" >> test1/www/index.html
echo "test2.test speaking" >> test2/www/index.html

// run the proof of concept
go run main.go
```

In another terminal you can check if it worked.

```
openssl s_client -servername test1.test -connect test1.test:8443 < /dev/null 2>/dev/null | openssl x509 -fingerprint -noout -in /dev/stdin

openssl s_client -servername test2.test -connect test2.test:8443 < /dev/null 2>/dev/null | openssl x509 -fingerprint -noout -in /dev/stdin
```
