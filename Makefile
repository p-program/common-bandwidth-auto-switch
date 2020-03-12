now 		  := $(shell date)
PREFIX		  ?= zeusro

auto_commit:   
	git add .
	git commit -am "$(now)"
	git push

buildAndRun:
	go build
	./common-bandwidth-auto-switch

buildDockerImage:
	docker build -t $(PREFIX)/common-bandwidth-auto-switch:1.0 -f deploy/docker/Dockerfile .