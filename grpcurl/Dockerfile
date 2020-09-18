FROM golang:1.13.3-alpine3.10 as builder

RUN apk add --update --no-cache git bash

RUN go get github.com/fullstorydev/grpcurl/...
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl

FROM curlimages/curl as runner

COPY --from=builder /go/bin/grpcurl /usr/bin