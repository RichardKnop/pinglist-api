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

test-facebook:
	bash -c 'go test -timeout=30s github.com/RichardKnop/pinglist-api/facebook'

test-subscriptions:
	bash -c 'go test -timeout=120s github.com/RichardKnop/pinglist-api/subscriptions'

test-alarms:
	bash -c 'go test -timeout=30s github.com/RichardKnop/pinglist-api/alarms'

test-metrics:
	bash -c 'go test -timeout=30s github.com/RichardKnop/pinglist-api/metrics'

test-teams:
	bash -c 'go test -timeout=30s github.com/RichardKnop/pinglist-api/teams'

test-notifications:
	bash -c 'go test -timeout=30s github.com/RichardKnop/pinglist-api/notifications'

test:
	bash -c 'go list ./... | grep -v vendor | xargs -n1 go test -timeout=120s'
