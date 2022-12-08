MOCO_VERSION := 0.14.0
IMAGE_TAG=controller:dev

.PHONY: build
build:
	git clone  -b v$(MOCO_VERSION) --depth 1 https://github.com/cybozu-go/moco.git
	patch -d moco -p1 < remove-sleep.patch
	docker build --target controller -t $(IMAGE_TAG) moco
