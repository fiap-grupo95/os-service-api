APP_SERVICE_NAME=app
APP_CONTAINER_NAME=os-service-api
DB_SERVICE_NAME=db
DB_CONTAINER_NAME=db
APP_BINARY_PATH=/app/os-service-api

.PHONY: init up down logs swag-generate-docker swag-run-docker test coverage coverage-html

init:
	cp .env-example .env
	docker-compose up -d --build

	@echo "Aguardando banco ficar pronto..."
	docker-compose exec $(DB_SERVICE_NAME) sh -c 'until pg_isready -U $$POSTGRES_USER -d $$POSTGRES_DB; do sleep 1; done'

	@echo "Aguardando container $(APP_CONTAINER_NAME) estar rodando..."
	@while [ -z "$$(docker-compose ps -q $(APP_SERVICE_NAME))" ] || \
		[ "$$(docker inspect -f '{{.State.Running}}' $$(docker-compose ps -q $(APP_SERVICE_NAME)))" != "true" ]; do \
		echo "Aguardando container..."; sleep 2; \
	done

	@echo "Container $(APP_CONTAINER_NAME) está rodando!"

	@echo "Aguardando app iniciar (sleep 10s dentro do container)..."
	docker-compose exec $(APP_SERVICE_NAME) sh -c 'sleep 10'

up:
	docker-compose up -d --build

down:
	docker-compose down

logs:
	docker-compose logs -f $(APP_SERVICE_NAME)

dev-up:
	cp .env-example .env
	docker-compose up -d dev

swag-generate: dev-up
	docker-compose exec dev sh -c "go install github.com/swaggo/swag/cmd/swag@latest && swag init -g ./internal/adapter/http/routes/routes.go --output ./docs --parseDependency --parseInternal"

test: dev-up
	docker-compose exec dev go test ./... -v

coverage: dev-up
	docker-compose exec dev go test ./... -coverprofile=coverage.out
	docker-compose exec dev go tool cover -func=coverage.out

coverage-html: dev-up
	docker-compose exec dev go test ./... -coverprofile=coverage.out
	docker cp $$(docker-compose ps -q dev):/app/coverage.out .
	go tool cover -html=coverage.out

get-all:
	kubectl get all -n os-service-api

deploy:
	kubectl apply -f ./internal/infrastructure/k8s/deploy/namespace.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/secret.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/configmap.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/deployment-api.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/service-api.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/hpa-api.yml

deploy-local:
	kubectl apply -f ./internal/infrastructure/k8s/deploy/namespace.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/secret-local.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/configmap-local.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/deployment-db.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/deployment-api.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/service-api.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/hpa-api.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/mockoon.yml

run-jobs:
	kubectl apply -f ./internal/infrastructure/k8s/job-process

delete:
	kubectl delete -f ./internal/infrastructure/k8s/deploy/namespace.yml

local-api:
	@echo "Expondo o database em host - localhost"
	kubectl port-forward service/os-service-api-service 8080:8080 -n os-service-api

local-db:
	@echo "Expondo a db em http://localhost:8080"
	kubectl port-forward service/mongodb 27017:27017 -n os-service-api

deploy-newrelic:
	kubectl apply -f ./internal/infrastructure/k8s/deploy/newrelic-namespace.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/newrelic-secret.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/newrelic-configmap.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/newrelic-fluent-bit-configmap.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/newrelic-events-configmap.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/newrelic-rbac.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/newrelic-daemonset-infra.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/newrelic-daemonset-fluent-bit.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/newrelic-kube-state-metrics.yml
	kubectl apply -f ./internal/infrastructure/k8s/deploy/newrelic-deployment-events.yml

delete-newrelic:
	kubectl delete -f ./internal/infrastructure/k8s/deploy/newrelic-namespace.yml

get-newrelic:
	kubectl get all -n newrelic