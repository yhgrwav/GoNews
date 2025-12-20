FROM golang:alpine AS builder

WORKDIR /app

# Install git and other dependencies if needed (alpine is minimal)
# RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
# Copy config.json as fallback
COPY config.json . 

EXPOSE 8080
CMD ["./main"]