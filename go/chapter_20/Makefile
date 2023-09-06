.PHONY: build build-no-cache up down ps app-bash watch exec

# docker-compose関連のコマンド

build:
	docker compose build

build-no-cache:
	docker compose build --no-cache

up:
	docker compose up

down:
	docker compose down

ps:
	docker compose ps -a

app-bash:
	docker compose exec app bash

db-bash:
	docker compose exec db bash

redis-bash:
	docker compose exec redis bash



