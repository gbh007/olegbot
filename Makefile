BUILD_ENV = GOOS=linux GOARCH=amd64 CGO_ENABLED=0

build:
	$(BUILD_ENV) go build -o _build/bot ./cmd/bot
	cd react-ui && npm run build
	docker build -f Dockerfile -t olegbot:latest .
	docker image save -o olegbot.tar olegbot:latest


up: build
	docker compose up --build --remove-orphans

down:
	docker compose down --remove-orphans