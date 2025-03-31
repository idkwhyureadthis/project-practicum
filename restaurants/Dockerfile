FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN apk add build-base && apk cache clean
ENV CGO_ENABLED=1
RUN go build -o ./restaurants ./cmd/restaurants/main.go


FROM alpine
WORKDIR /app
COPY --from=builder /app/restaurants ./restaurants
EXPOSE 8082
CMD ["./restaurants"] 