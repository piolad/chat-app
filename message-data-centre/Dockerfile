FROM golang:latest

WORKDIR /app

COPY go.mod .
COPY main.go .

RUN go build -o bin .

ENTRYPOINT [ "/app/bin" ]