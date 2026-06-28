FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY main.go .
COPY go.mod .
RUN go build -o spiraldb main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/spiraldb .
EXPOSE 8094
HEALTHCHECK --interval=10s --timeout=3s CMD wget -qO- http://localhost:8094/ -d '{"action":"health"}' || exit 1
ENTRYPOINT ["./spiraldb"]
