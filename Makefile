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

dockerrun:
	docker run --rm --name restapp -p 8080:8080 --net=host restapp

dockerbuild:
	docker build -t restapp .

build:
	go build -o bin/main main.go