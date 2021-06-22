FROM golang:alpine as builder
WORKDIR /build
COPY . .
RUN apk --update add ca-certificates && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s"

FROM scratch
COPY --from=builder /build/github-app-token /github-app-token
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/github-app-token"]
