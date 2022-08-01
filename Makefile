phony: protobuf

protobuf:
	cd ./proto/ ; protoc --go_out=./ main.proto
