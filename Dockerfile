FROM golang:alpine as builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s"

FROM scratch
COPY --from=builder /build/github-app-token /github-app-token
ENTRYPOINT ["/github-app-token"]
