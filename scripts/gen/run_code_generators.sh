#!/bin/bash
. ./scripts/util.sh

set -e

function resetFiles() {
    pocker-echo "Resetting files matching $1"
    local filePattern=$1
    local files=$(find proto -name "$filePattern" -type f)
    for file in $files; do
        git checkout --ours $file &> /dev/null || ( git checkout --theirs $file &> /dev/null || : ) # :)
    done
}

resetFiles "*.pb.go"
resetFiles "*.connect.go"

goGenerateCmd="go generate --tags generate ./..."
goImportsCmd="go run golang.org/x/tools/cmd/goimports --local "github.com/sweetloveinyourheart/planning-poker" -w ./"

pocker-echo "Running goimports..."
$goImportsCmd

pocker-echo "Running go generate..."
$goGenerateCmd || (pocker-echo "go generate failed, retrying after goimports..." && $goImportsCmd && $goGenerateCmd)