phony: protobuf

protobuf:
	export PATH="$(PATH):$(go env GOPATH)/bin" ; cd ./proto/ ; protoc --go_out=./ main.proto

brun:
	docker-compose build;
	docker-compose up