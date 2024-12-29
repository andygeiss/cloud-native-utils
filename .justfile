set dotenv-load

# Benchmark the Go sources.
benchmark:
    @cd ./consistency && go test -bench .

# Create a local CA and sign a server certificate.
# This will only be used if domains = ["localhost"].
cert-dir := "./security/testdata"
make-certs:
    @rm -rf {{cert-dir}} ; mkdir {{cert-dir}}
    @brew install mkcert
    @mkcert -install
    @mkcert -key-file {{cert-dir}}/server.key -cert-file {{cert-dir}}/server.crt localhost 127.0.0.1 ::1

# Create the plugins.
plugin:
    @go build -buildmode=plugin -o ./extensibility/testdata/adapter.so ./extensibility/testdata/adapter.go

# Test the Go sources (Units).
test: plugin
    @go test -v -coverprofile=.coverprofile.out github.com/andygeiss/cloud-native-utils/...

# Test module integration like the Server.
test-integration:
    @go test -v --tags=integration ./...
