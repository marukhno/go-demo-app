FROM golang:1.16.0
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY db.go ./db.go
COPY models.go ./models.go
COPY main.go ./main.go

RUN go build -o main .
EXPOSE 8080

CMD ["./main"]
