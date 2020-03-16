now 		  := $(shell date)
PREFIX		  ?= zeusro
APP_NAME      ?= common-bandwidth-auto-switch:latest
IMAGE		  ?= $(PREFIX)/$(APP_NAME)
MIRROR_IMAGE  ?= registry.cn-shenzhen.aliyuncs.com/amiba/common-bandwidth-auto-switch:latest

auto_commit:   
	git add .
	git commit -am "$(now)"
	git push

buildAndRun:
	go build
	./common-bandwidth-auto-switch

mirror:
	docker build -t $(MIRROR_IMAGE) -f deploy/docker/Dockerfile .

release-mirror:
	docker push $(MIRROR_IMAGE)

rebuild:
	git pull
	docker build -t $(IMAGE) -f deploy/docker/Dockerfile .

test:
	mkdir -p artifacts/report/coverage
	go test -v -cover -coverprofile c.out.tmp ./...
	cat c.out.tmp | grep -v "_mock.go" > c.out
	go tool cover -html=c.out -o artifacts/report/coverage/index.html	

up:
	docker-compose build --force-rm --no-cache
	docker-compose up


