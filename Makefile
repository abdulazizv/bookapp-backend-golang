CURRENT_DIR=$(shell pwd)

-include .env

DB_URL="postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

run:
	go run cmd/main.go

migrate_file:
	migrate create -ext sql -dir migrations -seq alter_author

migrate_up:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migrate_down:
	migrate -path migrations -database "$(DB_URL)" -verbose down 

migrate_forse:
	migrate -path migrations -database "$(DB_URL)" -verbose forse 5

swag:
	swag init -g api/router.go  -o api/docs
