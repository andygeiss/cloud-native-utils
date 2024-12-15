set dotenv-load

# Benchmark the Go sources.
benchmark:
    @cd ./utils/consistency && go test -v -bench .

# Test the Go sources.
test:
    @go test -v ./utils/...
