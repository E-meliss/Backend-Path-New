DB_URL?=$(DATABASE_URL)

.PHONY: migrate-up migrate-down migrate-status

migrate-up:
	goose -dir migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DB_URL)" down

migrate-status:
	goose -dir migrations postgres "$(DB_URL)" status
