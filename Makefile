.PHONY: setup-test
setup-test:
	@go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install golang.org/x/lint/golint@latest
	@go install github.com/nishanths/exhaustive/...@latest
	@go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

.PHONY: test
test: setup-test
	@go vet ./...
	###@go vet -vettool=$(HOME)/go/bin/shadow ./...
	@$(HOME)/go/bin/goimports -d ./
	@$(HOME)/go/bin/staticcheck ./...
	@$(HOME)/go/bin/golint -min_confidence 1 ./...
	@$(HOME)/go/bin/exhaustive -default-signifies-exhaustive ./...
	@$(HOME)/go/bin/gocyclo -over 19 $(shell find . -type f -name "*.go")
	@gofmt -s -l $(shell find . -type f -name "*.go")