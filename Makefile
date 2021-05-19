USER:=$(shell id -un)

PROD_VERSION=$(shell sed 1q .version)
DEV_REPONAME=dws-operator
DEV_IMGNAME=cray-dws-operator
DTR_IMGPATH=arti.dev.cray.com/$(DEV_REPONAME)/$(DEV_IMGNAME)
OPERATOR_SDK_IMGPATH=arti.dev.cray.com/kj-docker-unstable-local/cray-operator-sdk-build:0.19.2-20210318151410_80def2a

all: codestyle image

.PHONY: test

code-generation:
	docker run --rm -v $(PWD)/kubernetes:/go/src/stash.us.cray.com/dpm/$(DEV_REPONAME)/kubernetes -v $(PWD)/vendor:/go/src/stash.us.cray.com/dpm/$(DEV_REPONAME)/vendor -v $(PWD)/pkg:/go/src/stash.us.cray.com/dpm/$(DEV_REPONAME)/pkg -v $(PWD)/cmd:/go/src/stash.us.cray.com/dpm/$(DEV_REPONAME)/cmd -v $(PWD)/build:/go/src/stash.us.cray.com/dpm/$(DEV_REPONAME)/build $(OPERATOR_SDK_IMGPATH) stash.us.cray.com/dpm/$(DEV_REPONAME)/build/codeGenerationOperatorSdk.sh $(DEV_REPONAME)

vendor: code-generation
	GOPRIVATE=stash.us.cray.com go mod vendor

fmt: code-generation
	go fmt `go list -f {{.Dir}} ./...`

image: code-generation
	docker build -f build/Dockerfile --label $(DTR_IMGPATH):$(PROD_VERSION) -t $(DTR_IMGPATH):$(PROD_VERSION) .

test:
	docker build -f build/Dockerfile --label $(DTR_IMGPATH)-$@:$(PROD_VERSION)-$@ -t $(DTR_IMGPATH)-$@:$(PROD_VERSION) --target $@ .
	docker run --rm -t --name $@  -v $(PWD):$(PWD):rw,z $(DTR_IMGPATH)-$@:$(PROD_VERSION)

lint:
	docker build -f build/Dockerfile --label $(DTR_IMGPATH)-$@:$(PROD_VERSION)-$@ -t $(DTR_IMGPATH)-$@:$(PROD_VERSION) --target $@ .
	docker run --rm -t --name $@  -v $(PWD):$(PWD):rw,z $(DTR_IMGPATH)-$@:$(PROD_VERSION)

codestyle:
	docker build -f build/Dockerfile --label $(DTR_IMGPATH)-$@:$(PROD_VERSION) -t $(DTR_IMGPATH)-$@:$(PROD_VERSION) --target $@ .
	docker run --rm -t --name $@  -v $(PWD):$(PWD):rw,z $(DTR_IMGPATH)-$@:$(PROD_VERSION)

clean-lint:
	docker rmi $(DTR_IMGPATH)-lint:$(PROD_VERSION) || true

clean-codestyle:
	docker rmi $(DTR_IMGPATH)-codestyle:$(PROD_VERSION) || true

push:
	docker push $(DTR_IMGPATH):$(PROD_VERSION)

clean:
	docker container prune --force
	docker image prune --force
	docker rmi $(DTR_IMGPATH):$(PROD_VERSION)
	go clean all
