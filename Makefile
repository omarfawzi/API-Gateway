# ========================
# Configuration
# ========================
CONFIG_FILE                ?= config/$(basename $(CONFIG_TEMPLATE_FILE)).json
PORT                       ?= 8080

# ========================
# Build and Run
# ========================
.PHONY: build run clean wire install-wire render-templates

run: wire render-templates
	[ -f .env ] || cp .env.dist .env
	go build -o ./bin/gateway ./cmd/main.go
	./bin/gateway -p $(PORT) -c config/$(basename $(CONFIG_TEMPLATE_FILE)).json

clean:
	rm -rf bin/

wire:
	wire ./internal

install-wire:
	go install github.com/google/wire/cmd/wire@latest

# ========================
# Docker Compose
# ========================
.PHONY: docker docker-up docker-stop docker-destroy

docker: wire
	[ -f .env ] || cp .env.dist .env
	docker network create gateway || true
	docker-compose up

docker-up: ## Start docker stack
	docker-compose up -d --build

docker-stop: ## Stop docker stack
	docker-compose stop

docker-destroy: ## Destroy docker stack
	docker-compose down

# ========================
# Quality Assurance
# ========================
.PHONY: qa lint lint-golangci vet test

qa: lint vet test ## Run quality assurance checks

lint: lint-golangci
	@fmt_issues=$$(gofmt -l . | grep -v vendor/); \
	if [ -n "$$fmt_issues" ]; then \
		echo "Formatting issues found in:" $$fmt_issues; \
		exit 1; \
	fi

vet: ## Vet the code
	go vet ./...
	golangci-lint run

test: ## Run unit tests with coverage
	go test ./... -failfast -coverpkg=./... -coverprofile .testCoverage.txt
	@go tool cover -func .testCoverage.txt | grep total | awk '{print "Total coverage: " $$3}'

# ========================
# Krakend Visualization
# ========================
.PHONY: install-krakend-utils verify-krakend draw-krakend

install-krakend-utils:
	brew install krakend graphviz
	go install github.com/krakendio/krakend-config2dot/v2/cmd/krakend-config2dot@latest

verify-krakend:
	krakend check --config ${CONFIG_FILE}

draw-krakend: render-templates verify-krakend
	krakend-config2dot -c ${CONFIG_FILE} | dot -Tpng -o docs/$(basename $(CONFIG_TEMPLATE_FILE))-gateway.png

# ========================
# Vulnerability Check
# ========================
.PHONY: check-vulnerabilities

check-vulnerabilities:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

# ========================
# Template Rendering
# ========================
.PHONY: install-gomplate

install-gomplate:
	brew install gomplate

render-templates:
	@if [ -z "$${CONFIG_TEMPLATE_FILE}" ]; then \
		echo "Error: CONFIG_TEMPLATE_FILE is not set."; \
		exit 1; \
	fi
	gomplate -f config/$${CONFIG_TEMPLATE_FILE} -o config/rendered_config_$(basename $(CONFIG_TEMPLATE_FILE)).yaml
	yq eval -o=json config/rendered_config_$(basename $(CONFIG_TEMPLATE_FILE)).yaml > config/$(basename $(CONFIG_TEMPLATE_FILE)).json
	rm config/rendered_config_$(basename $(CONFIG_TEMPLATE_FILE)).yaml
