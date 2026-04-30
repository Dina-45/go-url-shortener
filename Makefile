DB_URL=postgres://postgres:123456@localhost:5432/mydb?sslmode=disable

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)