.DEFAULT_GOAL := start
image_name = the-drink-almanac-api-image
container_name = the-drink-almanac-api-container

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