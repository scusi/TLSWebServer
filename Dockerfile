# create an image to build
FROM golang:alpine as builder
ARG VERSION="development"
ARG COMMIT="unknown"
ARG BRANCH="unknown"
ARG BUILDTIME="unknown" 
RUN mkdir /build	# directory where we build the binary	
RUN mkdir /webroot		# directory where data will be stored, aka webroot
ADD ./www /webroot
#COPY ./www /		# copy webroot to builder
RUN mkdir /config	# directory where configuration files are stored
ADD . /build/
WORKDIR /build
# build a static version of TLSWebserver to run from a scratch image
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
	-ldflags "-extldflags \"-static\" -X main.version=$VERSION -X main.commit=$COMMIT -X main.buildtime=$BUILDTIME -X main.branch=$BRANCH" \
	-o TLSWebserver .

# generate a self-signed cert to start with
RUN go run cmd/gencert/generate_cert.go --rsa-bits=2048 --host=localhost

# create a fitsrv docker image
FROM scratch
COPY --from=builder /build/TLSWebserver /app/
COPY --from=builder /config /app/config
COPY --from=builder /build/docker.config.yml /app/config/config.yml
COPY --from=builder /build/cert.pem /app/config/
COPY --from=builder /build/key.pem /app/config/
#COPY --from=builder /build/tokendb.cache /app/config/
COPY --from=builder /webroot /www
WORKDIR /app
CMD ["./TLSWebserver"]
