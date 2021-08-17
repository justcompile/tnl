FROM golang:1.16

WORKDIR $GOPATH/src/github.com/justcompile/tnl

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

CMD ["go", "run", "cmd/server/main.go"]