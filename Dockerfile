FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY service/go.mod service/go.sum ./
RUN go mod download
COPY service/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o bot .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /build/bot .
ENTRYPOINT ["./bot"]
CMD ["/config"]
