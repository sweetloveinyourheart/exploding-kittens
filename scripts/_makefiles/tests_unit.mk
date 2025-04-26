# Targets for running unit and integration tests under `go test`

test: # Run all unit tests (see more options in Makefile)
	@./scripts/unit_test/runAllUnitTests.sh
test-verbose:
	@./scripts/unit_test/runAllUnitTests.sh verbose
test-coverage:
	@./scripts/unit_test/runAllUnitTests.sh cov
	@./scripts/unit_test/printCoverageStats.sh
test-coverage-verbose:
	@./scripts/unit_test/runAllUnitTests.sh cov verbose
	@./scripts/unit_test/printCoverageStats.sh
print-coverage:
	@./scripts/unit_test/printCoverageStats.sh


# CI Automation Conventions:
# Any makefile target that starts with ut- will run the unit tests for that package
# Any makefile target that starts with cov- will run the unit tests for that package and generate a coverage report
#
# Any unit tests that are not covered by an explicit ut-/cov- target will be covered by 'ut-other' and 'cov-other.'
# If you create a new pair of ut- and cov- targets, remember to exclude that package from 'template-other', and add a new job to .github/workflows/tests.yaml

template-ut:
	@go clean -testcache
	@./scripts/unit_test/ciTestWrapper.sh "$(optionalArg)" "$(package)" "$(verbose)" || exit 1

template-cov:
	@rm -rf tests/logs/cov-$(packageName)*
	@mkdir -p tests/logs/cov-$(packageName)
	@(make template-ut package=$(package) packageName=$(packageName) verbose=$(verbose) optionalArg="-coverprofile=tests/logs/cov-$(packageName)/cov.tmp") || exit 1
	@exclusions=$$(grep --include=\*.go -Ril "DO NOT EDIT" . | cut -c 3- | xargs | tr -s '[:blank:]' ',' | sed -E 's!,!|github.com/SandsB2B/ldx/!g'); \
	cat tests/logs/cov-$(packageName)/cov.tmp | grep -vE "github.com/SandsB2B/ldx/$${exclusions}" > tests/logs/cov-$(packageName)/cov;
	@rm -f tests/logs/cov-$(packageName)/cov.tmp
	@go tool cover -func tests/logs/cov-$(packageName)/cov                                                  >> tests/logs/cov-$(packageName)/low-level.txt  || exit 1
	@go tool cover -func tests/logs/cov-$(packageName)/cov | grep total: | awk '{print $$3}' | sed 's/.$$//' > tests/logs/cov-$(packageName)/percentage.txt || exit 1
	@go tool cover -html=tests/logs/cov-$(packageName)/cov -o                                                  tests/logs/cov-$(packageName)/visual.html    || exit 1
	@echo Reports are available in the logs directory.

###
### SERVICES
###

ut-userserver:
	@make template-ut package=services/user

cov-userserver:
	@make template-cov package=services/user

ut-gameengineserver:
	@make template-ut package=services/game_engine

cov-gameengineserver:
	@make template-cov package=services/game_engine