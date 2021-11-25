FROM alpine:3.15.0

ADD ./azure-scheduled-events /bin/azure-scheduled-events

RUN apk add --update ca-certificates \
    && rm -rf /var/cache/apk/*

ENTRYPOINT ["/bin/azure-scheduled-events"]
