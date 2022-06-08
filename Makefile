.PHONY: test
test:
	go test -race -v -coverprofile=cover.out  ./...

.PHONY: test-ui
test-ui: test
	go tool cover -html=cover.out -o cover.html
	open cover.html

update-configmap:
	kubectl create configmap driving-assistant-dev-conf --from-file=config.yaml -n driving-assistant-system -o yaml --dry-run | kubectl apply -f -
