protobuf:
	export PATH="$(PATH):$(go env GOPATH)/bin" ; cd ./proto/ ; protoc --go_out=./ main.proto

build:
	docker compose build --no-cache

start:
	docker compose up