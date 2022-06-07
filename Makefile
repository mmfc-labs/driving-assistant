.PHONY: lint
deploy:
	helm upgrade driving-assistant-dev charts --install --create-namespace --wait --namespace=driving-assistant-system --set=image.tag=v0.0.8 --values=charts/values.yaml  --atomic

.PHONY: test
test:
	go test -race -v -coverprofile=cover.out  ./...

.PHONY: test-ui
test-ui: test
	go tool cover -html=cover.out -o cover.html
	open cover.html

