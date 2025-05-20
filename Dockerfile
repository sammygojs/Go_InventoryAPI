FROM golang:1.22

WORKDIR /app
ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN echo "Unit tests..." && \
    USE_MOCK_PRODUCTS=true go test ./... -v

RUN go build -o main ./cmd/main.go

ENV USE_MOCK_PRODUCTS=

EXPOSE 8080

CMD ["./main"]
