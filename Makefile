# Copied from the baseline (stack/makefile.md). Adjust per its rule 5; record
# any other deviation in the README.

# Where make certs writes the mTLS material the integration tests read.
CERTS = ./security/testdata

# Targets are alphabetical, so the default is named rather than first.
.DEFAULT_GOAL = check
.PHONY: benchmark certs check ci fmt test test-integration

# Rule-3 target: the allocation-sensitive packages, benchmarked by hand.
benchmark:
	go test -bench=. -benchmem ./consistency/... ./efficiency/... ./resource/...

# Rule-3 target: a local CA and an mTLS pair for the integration tests.
# $(CERTS) is gitignored; mkcert is a prerequisite, not something this
# installs (rule 4).
certs:
	rm -rf $(CERTS)
	mkdir -p $(CERTS)
	mkcert -install
	mkcert -key-file $(CERTS)/server.key -cert-file $(CERTS)/server.crt localhost 127.0.0.1 ::1
	mkcert -client -key-file $(CERTS)/client.key -cert-file $(CERTS)/client.crt localhost
	cp "$$(mkcert -CAROOT)/rootCA.pem" $(CERTS)/ca.crt

# Default. Every gate, in this order (operations/ci.md), against the working
# tree. Run before every commit.
check:
	test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)
	go vet ./...
	go fix -diff ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go mod tidy -diff
	go test -race -shuffle=on ./...
	CGO_ENABLED=0 go build -trimpath ./...

# The same gates against the commit: a file never added, or a .env, cannot
# make it green. Run before every push. go version runs first, inside the
# copy, so the run records which toolchain ran. The archive goes through a
# file so git's exit status stops the run; one shell line so the trap cleans
# up however check ends.
ci:
	t=$$(mktemp); d=$$(mktemp -d); trap 'rm -rf "$$t" "$$d"' EXIT; git archive -o "$$t" HEAD && tar -xf "$$t" -C "$$d" && go -C "$$d" version && $(MAKE) -C "$$d" check

# goimports first: go fix type-checks, so a missing import would stop the
# recipe before goimports could add it. go fix manages the imports its own
# rewrites need.
fmt:
	go run golang.org/x/tools/cmd/goimports@latest -w .
	go fix ./...

# The inner loop.
test:
	go test -race -shuffle=on ./...

# Rule-3 target: the tagged tests, which need make certs first.
test-integration:
	go test -race -shuffle=on -tags=integration ./...
