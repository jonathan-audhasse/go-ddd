define log
	printf "\033[34m*** %s\033[0m\n" $(1)
endef

define warn
	printf "\033[33mWarning: %s\033[0m\n" $(1)
endef

define err
	printf "\033[31mError: %s\033[0m\n" $(1)
endef

help: ## Show this help message
	@echo "Valid script operations:"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | awk '{split($$0, a, ":.*##*"); printf "\033[0;32m%-10s\033[m #%s\n", a[1], a[2]}'

#-------------#
#  Variables  #
#-------------#

#------------------#
#  Common targets  #
#------------------#

build: ## Build
	docker compose build

clean:
	$(call info, "TODO script to clean")

#------------------#
#  Mount Services  #
#------------------#

up: ## Mount services
	docker compose up -d

down: ## Stop and Remove all mounted services
	docker compose down --remove-orphans

#---------#
#  Tests  #
#---------#

test: up ## Run tests
	$(call info, "TODO set tests")
	docker compose run --rm --no-deps --entrypoint=go server test ./...

#-----------#
#  Helpers  #
#-----------#

logs: ## Visualize the last 100 docker logs
	docker compose logs server | tail -100

fmt: ## Formatter
	go fmt ./...

lint: ## Run linter (ref: https://golangci-lint.run)
	golangci-lint run

#-----------------#
#  Local dev      #
#-----------------#

zsh: up ## Mount services for local dev
	docker compose exec server /bin/zsh

#---------#
#  Mocks  #
#---------#

mock-gen: # Generate mocks
	mockery --log-level=""

mock-clean: # Clean generated mock files
	@find -name .mock -print -exec rm -r {} +