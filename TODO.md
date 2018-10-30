## logging per TLSHost

TLSHost should have a logfile attribute where you can specify where to write the logfile for this host.

## add a routine to install acme.sh

Basically the routin needs to do the following:

- `curl https://get.acme.sh | sh`
- `source ~/.profile`

## add TLSHost

The idea is to have a command to add a TLSHost on the fly to the server config.

- add TLSHost to config `TLSWebServer -conf /path/to/config.json -add -host foobar.org -cert /path/to/cert.pem -key /path/to/key.pem -w /path/to/webroot`
- request a TLS certificate from Let's encrypt
  `acme.sh --issue -d foobar.org -w /path/to/webroot`
  `acme.sh --install-cert -d foobar.org --fullchain-file /path/to/cert.pem --key-file /path/to/keyfile.pem --reload-cmd "service tlswebserver restart"`

