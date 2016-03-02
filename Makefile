DEPS=go list -f '{{range .TestImports}}{{.}} {{end}}' ./...

export GO15VENDOREXPERIMENT=1

update-deps:
	rm -rf Godeps
	rm -rf vendor
	go get github.com/tools/godep
	godep save ./...

install-deps:
	go get github.com/tools/godep
	godep restore
	$(DEPS) | xargs -n1 go get -d

fmt:
	bash -c 'go list ./... | grep -v vendor | xargs -n1 go fmt'

test-oauth:
	bash -c 'go test -timeout=30s github.com/RichardKnop/pinglist-api/oauth'

test-accounts:
	bash -c 'go test -timeout=30s github.com/RichardKnop/pinglist-api/accounts'

test-alarms:
	bash -c 'go test -timeout=30s github.com/RichardKnop/pinglist-api/alarms'

test-subscriptions:
	bash -c 'go test -timeout=60s github.com/RichardKnop/pinglist-api/subscriptions'

test:
	bash -c 'go list ./... | grep -v vendor | xargs -n1 go test -timeout=60s'
