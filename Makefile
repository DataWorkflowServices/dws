USER:=$(shell id -un)

PROD_VERSION=$(shell sed 1q .version)
DEV_IMGNAME=dws-operator
DTR_IMGPATH=dtr.dev.cray.com/$(USER)/$(DEV_IMGNAME)

all: image

.PHONY: src
src:
	operator-sdk generate k8s
	operator-sdk generate openapi
	bash build/bin/genclient

image:
	docker build -f build/Dockerfile --label $(DTR_IMGPATH):$(PROD_VERSION) -t $(DTR_IMGPATH):$(PROD_VERSION) .
	docker push $(DTR_IMGPATH):$(PROD_VERSION)
	docker image prune --force
	sed -i "" 's|REPLACE_IMAGE|$(DTR_IMGPATH):$(PROD_VERSION)|g' deploy/operator.yaml

clean:
	docker container prune --force
	docker image prune --force
	docker rmi $(DTR_IMGPATH):$(PROD_VERSION)
	go clean all
