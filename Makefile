USER:=$(shell id -un)

PROD_VERSION=$(shell sed 1q .version)
DEV_IMGNAME=dtr.dev.cray.com/$(USER)/dws-operator

all: src image

src:
	operator-sdk generate k8s
	operator-sdk generate openapi
	sed -i "" 's|REPLACE_IMAGE|dtr.dev.cray.com/$(USER)/dws-operator:0.0.1|g' deploy/operator.yaml

image:
	docker build -f build/Dockerfile 
		--label jhendricks/dtr.dev.cray.com/jhendricks/dws-operator:0.0.1 \
		-t dtr.dev.cray.com/jhendricks/dws-operator:0.0.1 .
	docker push $(DEV_IMGNAME):$(PROD_VERSION)

clean:
	docker container prune --force
	docker image prune --force
	docker rmi $(DEV_IMGNAME):$(PROD_VERSION)
	go clean all
