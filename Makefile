api:
	go run . http

migrate-up:
	go run . migrate up

migrate-down:
	go run . migrate down

migrate-seed:
	go run . migrate seed

migrate-refresh:
	go run . migrate refresh

cmd-hello:
	go run . cmd hello