REG=docker.astuart.co:5000/k8s
BINARY?=$(shell basename $(PWD))
IMAGE:=$(BINARY)

.PHONY: build push deploy

TAG=$(REG)/$(IMAGE)

$(IMAGE): *.go
	go build -o $(IMAGE)
	# upx $(IMAGE)
	
build: $(IMAGE)
	docker build -t $(TAG) .

push: build
	docker push $(TAG)

deploy: push
	kubectl delete po --namespace kube-system -l app=$(IMAGE)
