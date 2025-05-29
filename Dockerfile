FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app .
FROM scratch
COPY --from=builder /out/app /app
EXPOSE 3001
CMD ["/app"]