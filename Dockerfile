FROM golang:1.22

WORKDIR /app
ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 🧪 Run tests and print coverage percentage
RUN echo "🧪 Running tests..." && \
    go test ./... -v -cover

# 🛠️ Build the binary
RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]
