FROM golang:1.14.0-alpine3.11 as build-env
ADD . /src/github.com/zeusro/common-bandwidth-auto-switch
WORKDIR /src/github.com/zeusro/common-bandwidth-auto-switch
RUN apk update
# 分层构建,复用缓存
RUN GOARCH=amd64 CGO_ENABLED=0 go build -mod vendor -o common-bandwidth-auto-switch  github.com/zeusro/common-bandwidth-auto-switch 

FROM alpine:3.10
COPY --from=build-env /src/github.com/zeusro/common-bandwidth-auto-switch /app/common-bandwidth-auto-switch
ENTRYPOINT ["/app/common-bandwidth-auto-switch"]
