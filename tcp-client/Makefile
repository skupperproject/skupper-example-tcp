IMAGE := quay.io/skupper/tcp-example-client

.PHONY: build
build:
	podman build -t ${IMAGE} .

.PHONY: run
run: build
	podman run -p 9090:9090 ${IMAGE}

.PHONY: debug
debug: build
	podman run -p 9090:9090 -it ${IMAGE} /bin/bash

# Prerequisite: podman login quay.io
.PHONY: push
push: build
	podman push -q ${IMAGE}
