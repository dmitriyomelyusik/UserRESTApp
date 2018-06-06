run: test
	bin/main

lint:
	gometalinter .
	gometalinter postgres/.
	gometalinter entity/.
	gometalinter handlers/. --disable=gas
	gometalinter controller/.
	gometalinter errors/.

test:
	go test postgres/postgres_test.go
	go test handlers/handlers_test.go
