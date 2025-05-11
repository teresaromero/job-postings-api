FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on \
    GOPROXY=https://proxy.golang.org,direct

RUN go build -o /bin/api ./cmd/api/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /bin/api ./api
CMD ["./api"]