FROM golang:1.18

RUN mkdir /go/src/cTCP
WORKDIR /go/src/cTCP

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN  go build -o bin/app ./cmd/ApplicationServer

CMD [ "./bin/app" ]