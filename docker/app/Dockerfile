FROM golang:1.18-alpine

RUN mkdir /go/src/cTCP
WORKDIR /go/src/cTCP

RUN apk update && apk add libc-dev && apk add gcc && apk add make

COPY go.mod go.sum ./
RUN go mod download && go mod verify

RUN go install github.com/githubnemo/CompileDaemon@latest
RUN apk add --no-cache git
# RUN go get github.com/githubnemo/CompileDaemon

COPY . .
COPY ./entrypoint_app.sh /entrypoint_app.sh

ADD https://raw.githubusercontent.com/eficode/wait-for/v2.1.0/wait-for /usr/local/bin/wait-for
RUN chmod +rx /usr/local/bin/wait-for /entrypoint_app.sh

ENTRYPOINT [ "sh", "/entrypoint_app.sh" ]


# RUN  go build -o bin/app ./cmd/ApplicationServer

# CMD [ "./bin/app" ]