FROM golang:alpine AS builder

WORKDIR /app

ADD go.mod .

COPY . .

RUN go build -o medods ./cmd/medods/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/medods /app/medods

EXPOSE 8080

CMD ["./medods"]
