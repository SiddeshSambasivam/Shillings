phony: protobuf

protobuf:
	export PATH="$(PATH):$(go env GOPATH)/bin" ; cd ./proto/ ; protoc --go_out=./ main.proto

web:
	cd cmd/WebServer; go build -o ../../bin/web .; cd ..; cd ..;\
	./bin/web;

app: 
	cd cmd/ApplicationServer; go build -o ../../bin/app .; cd ..; cd ..;\
	./bin/app;