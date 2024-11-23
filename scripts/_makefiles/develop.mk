# Targets for local development and testing
FULL_SERVER_STACK_COMPOSE_FILE := ./dockerfiles/docker-compose.yml

base-compose-up:
	@source ./scripts/util.sh && pocker-compose-up "$(COMPOSE_FILE)"

base-compose-down:
	@source ./scripts/util.sh && pocker-compose-down "$(COMPOSE_FILE)"

compose-up: # Start the full-server stack
	@make base-compose-up COMPOSE_FILE=$(FULL_SERVER_STACK_COMPOSE_FILE)

compose-down: 
	@make base-compose-down COMPOSE_FILE=$(FULL_SERVER_STACK_COMPOSE_FILE)
