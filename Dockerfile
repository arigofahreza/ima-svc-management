FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY go.* .
RUN go mod download
COPY . .
RUN go build -o server server.go


FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/ima-svc-management .
EXPOSE 61123
ENTRYPOINT ["./ima-svc-management"]