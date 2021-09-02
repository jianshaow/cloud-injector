FROM golang:1.16 AS build

WORKDIR /workspace
COPY cmd ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -v -o build/pod-injector cmd/main.go

FROM alpine:3.14.2
COPY --from=build  /workspace/build/pod-injector /usr/local/bin/pod-injector

EXPOSE 8443

CMD ["pod-injector"]
