FROM golang

WORKDIR /app

COPY . .

RUN go build -o main src/main.go

CMD ["./main"]