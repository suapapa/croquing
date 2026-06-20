.PHONY: all build push lint

TAG=latest
CR=icn.vultrcr.com/homincr1
IMAGE_TAG_BE=${CR}/croquis-king-backend:${TAG}
IMAGE_TAG_FE=${CR}/croquis-king-frontend:${TAG}

all: lint push

build_be:
	docker buildx build --platform linux/amd64 -t ${IMAGE_TAG_BE} .

build_fe:
	docker buildx build --platform linux/amd64 \
	    --target prod \
		--build-arg VITE_API_BASE= \
	    -t ${IMAGE_TAG_FE} ./frontend

push: build_be build_fe
	docker push ${IMAGE_TAG_BE}
	docker push ${IMAGE_TAG_FE}

lint:
	golangci-lint run ./...
