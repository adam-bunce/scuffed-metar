# Trying to build this in the script kills the instance
# can't use package manager version for some resaon (mayb skill issue)
FROM golang:1.17 as build-step

ENV GO111MODULE off
ENV CGO_ENABLED 0
ENV REPO github.com/awslabs/amazon-ecr-credential-helper/ecr-login/cli/docker-credential-ecr-login

RUN go get -u $REPO

RUN rm /go/bin/docker-credential-ecr-login

RUN go build \
 -o /go/bin/docker-credential-ecr-login \
 /go/src/$REPO

FROM alpine:latest

COPY --from=build-step /go/bin /go/bin

WORKDIR /go/bin/