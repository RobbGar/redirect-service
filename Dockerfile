FROM golang:1.17 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:latest

ENV TZ=UTC
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .

CMD ["./main"]
