set dotenv-load

# Benchmark the Go sources.
benchmark:
    @cd ./consistency && go test -bench .

# Create a local CA and sign a server certificate.
# This will only be used if domains = ["localhost"].
cert-dir := "./security/testdata"
make-certs:
    @brew install mkcert
    @rm -rf {{cert-dir}} ; mkdir {{cert-dir}}
    @mkcert -install
    @mkcert -cert-file {{cert-dir}}/server.crt \
        -key-file {{cert-dir}}/server.key \
        localhost 127.0.0.1 ::1

# Test the Go sources (Units).
test:
    @go test -v -coverprofile=.coverprofile.out github.com/andygeiss/cloud-native-utils/...

# Test module integration like the Server.
test-integration:
    @go test -v --tags=integration ./...
