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

package-favorites-lambda:
	bash scripts/package_lambda.sh "favorites"
.PHONY:package-favorite-lambdas

package-users-lambda:
	bash scripts/package_lambda.sh "users"
.PHONY:package-users-lambdas

package-lambdas:
	make package-favorites-lambda
	make package-users-lambda
.PHONY:package-lambdas

publish-favorites-lambda:
	bash scripts/publish_lambda.sh "favorites"
.PHONY: publish-favorites-lambda

publish-users-lambda:
	bash scripts/publish_lambda.sh "users"
.PHONY: publish-users-lambda

publish-lambdas:
	make publish-favorites-lambda
	make publish-users-lambda
.PHONY:publish-lambdas

package-publish-lambdas:
	make package-lambdas
	make publish-lambdas
.PHONY: package-publish-lambdas