set dotenv-load

# Test the Go sources (Units).
test:
    @go test -v ./internal/app/models/...
