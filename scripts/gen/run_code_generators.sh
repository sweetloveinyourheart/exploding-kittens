#!/bin/bash
. ./scripts/util.sh

goGenerateCmd="go generate --tags generate ./..."
goImportsCmd="go run golang.org/x/tools/cmd/goimports --local "github.com/sweetloveinyourheart/planning-poker" -w ./"

pocker-echo "Running goimports..."
$goImportsCmd

pocker-echo "Running go generate..."
$goGenerateCmd || (pocker-echo "go generate failed, retrying after goimports..." && $goImportsCmd && $goGenerateCmd)