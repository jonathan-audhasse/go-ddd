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

#------------#
#  Database  #
#------------#

LBASE := liquibase --defaults-file=migration/liquibase/liquibase.properties --search-path=migration/liquibase

db-update: ## migrate database up
	$(LBASE) update

db-downgrade: check-liquibase-properties ## migrate database down
	$(LBASE) rollback-count 99

db-populate: ## populate database
	./set_up.py -d --health-value=${HEALTH_VALUE} datas/chaos_monkey.json
