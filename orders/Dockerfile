FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN apk add build-base && apk cache clean
ENV CGO_ENABLED=1
RUN go build -o ./orders ./cmd/orders/main.go


FROM alpine
WORKDIR /app
COPY --from=builder /app/orders ./orders
EXPOSE 8082
CMD ["./orders"] 