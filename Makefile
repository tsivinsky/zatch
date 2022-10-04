dev:
	docker-compose up --build -d db

start:
	docker-compose up

build:
	docker-compose build

stop:
	docker-compose down
