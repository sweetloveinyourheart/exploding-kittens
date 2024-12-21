#!/bin/bash
. ./scripts/util.sh

set -e

package=$1

go test -v -p 1 -count=1 -timeout 1800s --tags development ./$package/... || exit 1