FROM --platform=linux/x86_64 golang:1.21.0-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && \
    go build -o main

FROM --platform=linux/x86_64 golang:1.21.0-alpine AS production
COPY --from=builder /app/main /main

CMD ["/main"]
