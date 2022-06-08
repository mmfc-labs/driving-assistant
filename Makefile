.PHONY: test
test:
	go test -race -v -coverprofile=cover.out  ./...

.PHONY: test-ui
test-ui: test
	go tool cover -html=cover.out -o cover.html
	open cover.html

