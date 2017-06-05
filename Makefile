
deploy:
	CGO_ENABLED=0 go build .
	docker build -t quay.io/chronojam/vault-audit-bridge:latest .
	docker push quay.io/chronojam/vault-audit-bridge
