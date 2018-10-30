#!/bin/sh

echo "[*] generating a self-signed tls certificate for TLSWebServer"
/usr/bin/openssl req -x509 -subj '/CN=localhost' -nodes -days 365 -newkey rsa:2048 -keyout /var/TLSWebServer/localhost/tls/key.pem -out /var/TLSWebServer/localhost/tls/cert.pem

echo "[*] generating a ssh hostkey for scpdrop"
/usr/bin/ssh-keygen -t rsa -b 4096 -C scpDrop -f /scpdrop/id_rsa

echo "[*] start tlswebserver"
/bin/chmod +x /etc/systemd/system/tlswebserver.service
/usr/sbin/service tlswebserver start

echo "[*] start scpdrop"
/bin/chmod +x /etc/systemd/system/scpdrop.service
/usr/sbin/service scpdrop start
