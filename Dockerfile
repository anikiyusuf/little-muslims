FROM golang:1.25 AS builder

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go binary
RUN go build -o /app/bin/main cmd/api/main.go

# -----------------------------
# Final Runtime Stage
# -----------------------------
FROM alpine:3.19

# Copy binary from build stage
COPY --from=builder /app/bin/main /main
# Copy environment file & startup script
COPY .env.docker .env
COPY scripts/start.sh start.sh

# Make script executable
RUN chmod +x start.sh

RUN go install github.com/air-verse/air@latest
# ENV PATH=$PATH:/go/bin

# Expose port
EXPOSE 8080

# Entry point
ENTRYPOINT ["./start.sh"]

# Default command
CMD ["/main"]


