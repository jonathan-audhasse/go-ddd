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

cover:
	# go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
	# go test -coverprofile=coverage.out ./...
	go test -cover ./...

run: # run the server
	go run cmd/api/main.go

migrate-up: ## migrate database up
	go run cmd/migrate/main.go -cmd=up