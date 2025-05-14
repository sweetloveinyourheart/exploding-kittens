#!/bin/bash
. ./scripts/util.sh

# This script is used to print the coverage stats for each module
# It is intended to be run from the Makefile

set -e

allCovTargets=$(sed -n -e '/^cov-/p' scripts/_makefiles/tests_unit.mk | awk -F ":" '{print $1}' | sort)

#remove cov-other from allCovTargets
allCovTargets=$(echo $allCovTargets | sed 's/cov-other//g')

if [[ -d tests/logs ]]; then\
    pushd . > /dev/null
    cd tests/logs
fi

kittens-echo "COVERAGE PER MODULE (percentage)"

missingCoverageData=""
for target in $allCovTargets; do
    if [[ ! -f $target/percentage.txt ]]; then
        missingCoverageData="$missingCoverageData $target"
    fi
    kittens-echo "$target: $([ -f $target/percentage.txt ] && cat $target/percentage.txt || ( echo -e "${red}no-data${nc}" ) )"
done

# Get all the cov-* directories that weren't identified in the tests_unit.mk file
otherCovOutput=$(ls -d cov-* 2>/dev/null || :)
if [[ -n $otherCovOutput ]]; then
    kittens-echo "COVERAGE PER MODULE (percentage) - other"
    for dir in $otherCovOutput; do
        if [[ ! $allCovTargets =~ $dir ]]; then
            kittens-echo "$dir: $([ -f $dir/percentage.txt ] && cat $dir/percentage.txt || ( echo -e "${red}no-data${nc}" ) )"
        fi
    done
fi

popd > /dev/null

if [[ -n $missingCoverageData ]]; then
    kittens-echo-red "Missing some coverage data!"
    statusCode=1
else
    kittens-echo-green "All coverage data is available"
    statusCode=0
fi

exit $statusCode
