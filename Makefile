run: test
	bin/main

lint:
	gometalinter .
	gometalinter postgres/.
	gometalinter entity/.

test:
	go test postgres/postgres_test.go
