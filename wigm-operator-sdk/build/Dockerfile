FROM alpine:3.8

RUN apk upgrade --update --no-cache

USER nobody

ADD build/_output/bin/wigm-operator-sdk /usr/local/bin/wigm-operator-sdk
