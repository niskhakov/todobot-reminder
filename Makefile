.PHONY:
.SILENT:

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

build-image:
	docker build -t todobot-reminder:v0.1 .

start-container:
	docker run --name todo-bot -p 8080:8080 --env-file .env todobot-reminder:v0.1

compose:
	docker-compose up -d 

compose-down:
	docker-compose down

compose-build:
	docker-compose build

compose-logs:
	docker-compose logs -f

compose-cli:
	docker-compose exec bot sh

view-db:
	echo "DB Viewer is available at http://localhost:8089"
	boltdbweb --db-name=bot.db --port=8089