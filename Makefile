
IMAGE_NAME := quotaguard
IMAGE_TAG := v1.0.5
K8S_NAMESPACE := webhook-system
CERTS_DIR := ./certs

all: apply docker-build docker-push apply-namespace generate-certs apply-secret clean-and-deploy

apply:
	kubectl apply -f resources/quotapolicy.yaml

docker-build:
	@echo "Building Docker image:  ${IMAGE_NAME}:${IMAGE_TAG}"
	docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .
	@echo "Image built"

docker-push:
	@docker tag ${IMAGE_NAME}:${IMAGE_TAG} ccr.ccs.tencentyun.com/malyue/${IMAGE_NAME}:${IMAGE_TAG}
	@docker push ccr.ccs.tencentyun.com/malyue/${IMAGE_NAME}:${IMAGE_TAG}

generate-certs:
	@echo "Generating TLS certificates..."
	rm -rf $(CERTS_DIR)
	mkdir -p $(CERTS_DIR)
	openssl req -x509 -newkey rsa:2048 -keyout $(CERTS_DIR)/tls.key -out $(CERTS_DIR)/tls.crt -days 365 -nodes \
		-subj "/CN=quotaguard.$(K8S_NAMESPACE).svc" \
		-addext "subjectAltName=DNS:quotaguard.$(K8S_NAMESPACE).svc"
	@echo "Certificates generated in $(CERTS_DIR)"

apply-namespace:
	kubectl apply -f resources/namespace.yaml

apply-secret:
	kubectl create secret tls webhook-certs \
	--cert=certs/tls.crt --key=certs/tls.key  \
	--namespace=$(K8S_NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -

clean-and-deploy:
	@echo "Performing clean redeploy..."
	@kubectl delete -f resources/deploy.yaml --ignore-not-found
	# 在创建 deploy 之前，不能存在 webhook 配置，否则会构成死锁（ 创建 pod 依赖 webhook，而 webhook 服务需要存在 pod）
	@kubectl delete ValidatingWebhookConfiguration quotaguard --ignore-not-found
	@sleep 2  # 等待资源完全删除
	@kubectl apply -f resources/rbac.yaml
	@kubectl apply -f resources/deploy.yaml
	@kubectl apply -f resources/service.yaml
	@echo "Waiting for webhook pod to be ready"
	@kubectl wait --for=condition=ready pod \
		-n $(K8S_NAMESPACE) \
		-l app=quotaguard \
		--timeout=120s
	@cat resources/validwebhookconf.template.yaml | \
		sed "s/{{CA_BUNDLE}}/$$(cat $(CERTS_DIR)/tls.crt | base64 | tr -d '\n')/" | \
		kubectl apply -f -
	@#kubectl apply -f resources/validwebhookconf.yaml

deploy-webhook:
	@cat resources/validwebhookconf.template.yaml | \
     		sed "s/{{CA_BUNDLE}}/$$(cat $(CERTS_DIR)/tls.crt | base64 | tr -d '\n')/" | \
     		kubectl apply -f -