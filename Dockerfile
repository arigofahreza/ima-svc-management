FROM golang:1.18-alpine AS builder
WORKDIR /app
ENV CGO_ENABLED=0
RUN go install -v -a std
COPY vendor ./vendor/
RUN grep -v '#' vendor/modules.txt  | xargs -I{} sh -c "go build -mod=vendor -a -v -i {} || true"
COPY . .
RUN go build -mod=vendor -v -o ima-svc-management


FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder . ima-svc-management
EXPOSE 76331
ENTRYPOINT ["./ima-svc-management"]