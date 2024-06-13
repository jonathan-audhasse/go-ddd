define log
	printf "\033[34m*** %s\033[0m\n" $(1)
endef

define warn
	printf "\033[33mWarning: %s\033[0m\n" $(1)
endef

define err
	printf "\033[31mError: %s\033[0m\n" $(1)
endef

help: ## show this help message
	@echo "Valid script operations:"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | awk '{split($$0, a, ":.*##*"); printf "... %-20s :%s\n", a[1], a[2]}'

#-------------#
#  Variables  #
#-------------#

SERVER_NAME=server

#------------------#
#  Common targets  #
#------------------#

build: ## build
	docker compose build

clean:
	$(call info, "TODO script to clean")

#------------------#
#  Mount Services  #
#------------------#

up: ## mount services
	docker compose up -d ${SERVER_NAME}

down: ## stop and Remove all mounted services
	docker compose down --remove-orphans

#---------#
#  Tests  #
#---------#

test: up ## run tests
	$(call info, "TODO set tests")
	docker compose run --rm --no-deps --entrypoint=go ${SERVER_NAME} test ./...

#-----------#
#  Helpers  #
#-----------#

pkg := $(shell go list ./...)

logs: ## visualize the last 100 docker logs
	docker compose logs ${SERVER_NAME} | tail -100

fmt: ## formatter
	go fmt $(pkg)

#-----------------#
#  Local dev      #
#-----------------#

zsh: up ## mount services for local dev
	docker compose exec ${SERVER_NAME} /bin/zsh
