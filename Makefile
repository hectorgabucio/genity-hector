
## Runs a local environment using docker-compose
.PHONY: start
start:
	docker-compose -f deployments/docker-compose.yml up --build

## Autogenerates mocks
.PHONY: mocks
mocks: 
	mockery -all -recursive -output ./test/mocks

## Removes volumes and all containers
.PHONY: clean
clean: 
	docker-compose -f deployments/docker-compose.yml down --rmi local -v

# Help documentation à la https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@cat Makefile | grep -v '\.PHONY' |  grep -v '\help:' | grep -B1 -E '^[a-zA-Z0-9_.-]+:.*' | sed -e "s/:.*//" | sed -e "s/^## //" |  grep -v '\-\-' | sed '1!G;h;$$!d' | awk 'NR%2{printf "\033[36m%-30s\033[0m",$$0;next;}1' | sort