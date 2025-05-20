FROM golang:1.22

WORKDIR /app
ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ✅ Set env var to use mock data
# ENV USE_MOCK_PRODUCTS=true

# 👇 Keep test env for running tests only
RUN echo "🧪 Running unit tests..." && \
    USE_MOCK_PRODUCTS=true go test ./... -v

# 🛠️ Build the binary
RUN go build -o main ./cmd/main.go

# 👇 Remove test env before running in production
ENV USE_MOCK_PRODUCTS=

EXPOSE 8080

CMD ["./main"]
