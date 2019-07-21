TAG ?= latest

.PHONY: all
all: lint publish

.PHONY: lint
lint:
	@echo "Linting"
	docker build . --target deps -t "justcompile/ecs2cw:${TAG}"
	docker run --rm -it "justcompile/ecs2cw:${TAG}" ./bin/golangci-lint run ./...

.PHONY: publish
publish:
	@echo "--> Publishing Image"
	docker build . -t "justcompile/ecs2cw:${TAG}"
	docker push justcompile/ecs2cw:${TAG}
