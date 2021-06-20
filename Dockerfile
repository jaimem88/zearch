FROM golang:1.16.4-alpine3.13
RUN apk --no-cache add make

COPY . /zearch
WORKDIR /zearch

RUN make build

ENTRYPOINT ["/zearch/out/bin/zearch"]
