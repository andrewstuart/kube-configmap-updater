REG=andrewstuart
BINARY?=$(shell basename $(PWD))
IMAGE:=$(BINARY)

.PHONY: build push deploy

TAG=$(REG)/$(IMAGE)

$(IMAGE): *.go
	go build -o $(IMAGE)

clean:
	-rm $(IMAGE)
	
build: $(IMAGE)
	-upx $(IMAGE)
	docker build -t $(TAG) .

push: build
	docker push $(TAG)

deploy: push
	kubectl delete po --namespace kube-system -l app=$(IMAGE) --grace-period=0
