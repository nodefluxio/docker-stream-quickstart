PKG := "gitlab.com/nodefluxio/vanilla-dashboard"
PKG_LIST := $(shell go list ${PKG}... | grep -v /client/)

.PHONY: test
test:
	@go test ${PKG_LIST} -coverprofile .testCoverage.txt -timeout 30s \
		&& go tool cover -func=.testCoverage.txt
	@rm .testCoverage.txt

.PHONY: lint
lint:  ## Run linter
	$(eval export PATH=${GOPATH}/bin:$(PATH))
	@golint -set_exit_status ${PKG_LIST}
