protobuf:
	export PATH="$(PATH):$(go env GOPATH)/bin" ; cd ./proto/ ; protoc --go_out=./ main.proto

perfTest:
	wrk -t200 -c1000 -d1 --latency http://0.0.0.0:8080/v1/account -s ./tests/performance/profile.lua --timeout 10

prod:
	docker compose up --scale web=5

build:
	docker compose build --no-cache

start:
	docker compose up