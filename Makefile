.PHONY: lint
deploy:
	helm upgrade driving-assistant-dev charts --install --create-namespace --wait --namespace=driving-assistant-system --set=image.tag=v0.0.8 --values=charts/values.yaml  --atomic
