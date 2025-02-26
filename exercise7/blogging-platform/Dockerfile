FROM golang:1.22.3 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go get github.com/go-sql-driver/mysql

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o blogging-platform .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/blogging-platform .

EXPOSE 8080

CMD ["./blogging-platform"]
