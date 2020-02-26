FROM scratch
COPY TLSWebServer /
COPY config.yml /
COPY www /
COPY www/index.html /www/
COPY tls /
COPY tls/cert.pem /tls/
COPY tls/key.pem /tls/
ENTRYPOINT ["/TLSWebServer"]
EXPOSE 8080
EXPOSE 8443
