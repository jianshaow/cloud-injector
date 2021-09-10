FROM golang:1.16 AS build

WORKDIR /workspace
COPY go.mod go.sum ./
COPY cmd ./cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -v -o pod-injector cmd/main.go cmd/config.go

FROM alpine:3.14.2
COPY --from=build /workspace/pod-injector /usr/local/bin/pod-injector

EXPOSE 8443

CMD ["pod-injector"]
