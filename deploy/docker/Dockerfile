FROM golang:1.14.0-alpine3.11
ADD . /src/github.com/zeusro/common-bandwidth-auto-switch
WORKDIR /src/github.com/zeusro/common-bandwidth-auto-switch
RUN apk update -y  && \
	GOARCH=amd64 CGO_ENABLED=0 go build  github.com/zeusro/common-bandwidth-auto-switch
ENTRYPOINT ["/common-bandwidth-auto-switch"]