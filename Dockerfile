# ---------- Build Stage ----------
FROM golang:1.25.8-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd
# ---------- Run Stage ----------
FROM alpine:latest
WORKDIR /root/
# Only copy the compiled binary over
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]