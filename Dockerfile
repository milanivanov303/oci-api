FROM golang:alpine3.18

RUN apk update
RUN apk add git make bash
RUN git config --global http.sslVerify false
RUN git clone 'https://gdancheva:glpat-zQsukkvz8y72GjxJp2hg@gitlab.codixfr.private/enterpriseapps/oci-api.git'

COPY codix-entry.sh /

ENTRYPOINT ["/codix-entry.sh"]
