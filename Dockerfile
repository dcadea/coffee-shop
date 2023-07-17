FROM golang:1.20-alpine

WORKDIR /app

COPY ./src ./src
COPY go.mod go.mod
COPY go.sum go.sum

RUN go build -o main ./src/main/

EXPOSE 8080

CMD ["./main"]

