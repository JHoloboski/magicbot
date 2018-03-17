FROM golang:1.10 AS build

ARG DEP_VERSION=v0.4.1
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    wget -O /usr/local/bin/dep https://github.com/golang/dep/releases/download/$DEP_VERSION/dep-linux-amd64 && \
    chmod +x /usr/local/bin/dep

WORKDIR $GOPATH/src/github.com/jholoboski/magicbot
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -v -vendor-only

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -installsuffix static -ldflags "-s -w" .

FROM scratch AS run

COPY --from=build /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=build /go/bin/magicbot .

ENTRYPOINT ["./magicbot"]
