.PHONY: connect

connect:
	docker run -v $(shell pwd):/home/martin/src/github.com/martintrojer/mtorrent-go -p 1337:1337 --rm -i -t mtorrent-dev /bin/bash

dev:
	docker build --rm -f Dockerfile.dev -t mtorrent-dev .

dist:	dev
	docker run -v $(shell pwd):/home/martin/src/github.com/martintrojer/mtorrent-go -p 1337:1337 --rm -t mtorrent-dev bash -c "go build"
	docker build --rm -f Dockerfile.dist -t mtorrent-go .

# /var/lib/boot2docker/profile
# EXTRA_ARGS "--insecure-registry bubba:5000"

push:
	docker tag -f mtorrent-go:latest bubba:5000/mtorrent-go
	docker push bubba:5000/mtorrent-go
