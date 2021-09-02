FROM alpine:3.14.2

COPY build/pod-injector /usr/local/bin/pod-injector

EXPOSE 8443

CMD ["pod-injector"]
