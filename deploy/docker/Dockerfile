FROM golang:1.14.0-alpine3.11
ENV GOPATH /:/src/github.com/zeusro/common-bandwidth-auto-switch
ADD . /src/github.com/zeusro/common-bandwidth-auto-switch
WORKDIR /src/github.com/zeusro/common-bandwidth-auto-switch
RUN apk update
# 分层构建,复用缓存
RUN GOARCH=amd64 CGO_ENABLED=0 go build -o common-bandwidth-auto-switch  github.com/zeusro/common-bandwidth-auto-switch && chmod +x common-bandwidth-auto-switch
ENTRYPOINT ["/src/github.com/zeusro/common-bandwidth-auto-switch/common-bandwidth-auto-switch"]