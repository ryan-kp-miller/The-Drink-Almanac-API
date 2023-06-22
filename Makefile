.DEFAULT_GOAL := start
image_name = the-drink-almanac-api-image
container_name = the-drink-almanac-api-container
test_image_name = the-drink-almanac-api-test-image
test_container_name = the-drink-almanac-api-test-container

build:
	docker build -t $(image_name) ./go_api
.PHONY:build

run:
	ls -a
	docker run -it --rm -p 8080:8080 --env-file=.env --name $(container_name) $(image_name)
.PHONE:run

start:
	make build
	make run
.PHONY:start

up:
	docker compose up --build -d
.PHONY:up

down:
	docker compose down
.PHONY:down

test:
	docker build -t $(test_image_name) -f ./go_api/Dockerfile.test ./go_api
	docker run -it --rm --name $(test_container_name) $(test_image_name)
.PHONY:test

package-lambdas:
	bash scripts/package_favorites_lambda.sh
.PHONY:package-lambdas

publish-lambdas:
	bash scripts/publish_favorites_lambda.sh
.PHONY: publish-lambdas

package-publish-lambdas:
	make package-lambdas
	make publish-lambdas
.PHONY: package-publish-lambdas