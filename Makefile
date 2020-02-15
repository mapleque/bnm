.PHONY: test

test:
	cd server && \
		DB_DSN='root@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=true&loc=Local' \
		go test -coverprofile cover.out

cover:
	cd server && \
		go tool cover -html=cover.out -o cover.html && \
		open cover.html
