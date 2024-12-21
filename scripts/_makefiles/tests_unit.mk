# Targets for running unit and integration tests under `go test`

template-unit-test:
	@go clean -testcache
	@./scripts/unit_test/test_ci.sh "$(package)" || exit 1

###
### SERVICES
###

unit-test-userserver:
	@make template-unit-test package=services/user