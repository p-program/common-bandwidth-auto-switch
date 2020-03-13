now 		  := $(shell date)
PREFIX		  ?= zeusro
APP_NAME      ?= common-bandwidth-auto-switch:latest
IMAGE		  ?= $(PREFIX)/$(APP_NAME)
MIRROR_IMAGE  ?= mirror/common-bandwidth-auto-switch:latest
auto_commit:   
	git add .
	git commit -am "$(now)"
	git push

buildAndRun:
	go build
	./common-bandwidth-auto-switch

rebuild:
	git pull
	docker build -t $(IMAGE) -f deploy/docker/Dockerfile .

mirror:
	docker tag $(IMAGE) $(MIRROR_IMAGE)
	docker push $(MIRROR_IMAGE)
	docker push $(MIRROR_IMAGE)

up:
	docker-compose build --force-rm --no-cache
	docker-compose up

