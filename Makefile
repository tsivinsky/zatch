dev:
	docker-compose up --build -d db redis

start:
	docker-compose up

build:
	docker-compose build

stop:
	docker-compose down
