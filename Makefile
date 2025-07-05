SCLAER_NAME := stockprice \
	ingress-nginx \
	rocketmq
OUTPUT_DIR := bin
TARGET := keda-$(SCLAER_NAME)-scaler \
	keda-$(SCLAER_NAME)-nginx-scaler \
	keda-$(SCLAER_NAME)-rocketmq-scaler

.PHONY: build
build:
	go build -o $(OUTPUT_DIR)/$(TARGET) ./cmd/$(SCLAER_NAME)/main.go  # ✔ Tab (fixed)

.PHONY: docker
docker:
	docker build -t keda-ingress-nginx-scaler:v0.0.1 .

.PHONY: generate
generate:
	protoc --proto_path=proto \
  --go_out=./pkg/api \
  --go-grpc_out=./pkg/api \
  proto/externalscaler.proto

redeploy: docker
	kind load docker-image keda-ingress-nginx-scaler:v0.0.1
	kubectl delete pod -n keda -l app=ingress-nginx-external-scaler

log:
	kubectl logs -n keda -l app=ingress-nginx-external-scaler --tail=2000