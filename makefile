include .env


.PHONY: network
network:
	docker network create test_canal

.PHONY: up
up:
	docker-compose -p $(PROJECT) -f docker-compose.yaml up -d


.PHONY: down
down:
	docker-compose -p $(PROJECT) -f docker-compose.yaml down

MySQLContainer = $(PROJECT)_mysql_1

.PHONY: db
db:
	docker cp ./db.sql $(MySQLContainer):/
	docker exec $(MySQLContainer) sh -c 'mysql -u root -p$(MYSQL_ROOT_PASSWORD) < db.sql'

.PHONY: run
run:
	go run main.go