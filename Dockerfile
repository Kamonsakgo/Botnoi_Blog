FROM golang:alpine as build-env
ENV CGO_ENABLED=0
ENV GO111MODULE=on
WORKDIR $GOPATH/src/bn-crud-ads

COPY . .

RUN go mod download

RUN go build


FROM scratch

# Copy our static executable.
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /go/src/bn-crud-ads/bn-crud-ads /go/src/bn-crud-ads/bn-crud-ads

# Run the hello binary.
ENTRYPOINT ["/go/src/bn-crud-ads/bn-crud-ads"]

