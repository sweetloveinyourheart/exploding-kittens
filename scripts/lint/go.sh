#!/bin/bash
. ./scripts/util.sh

# This script is used to run the linter on the codebase.
# It is intended to be run from the Makefile

arg1=$1  # 'fix-goimports' = Fix the malformed goimports

set -e

runGoImports="go run golang.org/x/tools/cmd/goimports"

function checkForMalformedFile() {
    if [ -n "$1" ]; then
        return 1
    fi
    return 0
}

function isGeneratedFile() {
    grep 'DO NOT EDIT' "$1" &> /dev/null
}

function trimNonExistingFiles() { # If a file no longer exists in the latest commit, we shouldn't check it against goimports
    local filesToCheck=$1
    local existingFiles=""
    for file in $filesToCheck; do
        if [ -f "$file" ] && ! $(isGeneratedFile "$file"); then
            existingFiles="$existingFiles $file"
        fi
    done
    echo $existingFiles
}


function fixBadImports() {
    local badImports=$1
    for file in $badImports; do
        if ! $(isGeneratedFile "$file"); then
            $runGoImports --local "github.com/sweetloveinyourheart/exploding-kittens" -w $file
        fi
    done
}

kittens-echo "Running goimports..."

currentBranch=$(git rev-parse --abbrev-ref HEAD)

commitedFiles=$(git diff --name-only $currentBranch $(git merge-base $currentBranch origin/main))
filesToCheck=$(trimNonExistingFiles "$(echo "$commitedFiles" | grep '.go$' || echo "")")
if [ ! -z "$filesToCheck" ]; then
    #kittens-echo "goimports: Checking committed files: \n$filesToCheck"
    badImports1=$(echo "$filesToCheck" | xargs $runGoImports --local "github.com/sweetloveinyourheart/exploding-kittens" -l)
fi

stagedFiles=$(git diff --name-only --cached)
filesToCheck=$(trimNonExistingFiles "$(echo "$stagedFiles" | grep '.go$' || echo "")")
if [ ! -z "$filesToCheck" ]; then
    #kittens-echo "goimports: Checking staged files: \n$filesToCheck"
    badImports2=$(echo "$filesToCheck" | xargs $runGoImports --local "github.com/sweetloveinyourheart/exploding-kittens" -l)
fi

unstagedFiles=$(git diff --name-only)
filesToCheck=$(trimNonExistingFiles "$(echo "$unstagedFiles" | grep '.go$' || echo "")")
if [ ! -z "$filesToCheck" ]; then
    #kittens-echo "goimports: Checking unstaged files: \n$filesToCheck"
    badImports3=$(echo "$filesToCheck" | xargs $runGoImports --local "github.com/sweetloveinyourheart/exploding-kittens" -l)
fi

untrackedFiles=$(git ls-files --others --exclude-standard)
filesToCheck=$(trimNonExistingFiles "$(echo "$untrackedFiles" | grep '.go$' || echo "")")
if [ ! -z "$filesToCheck" ]; then
    #kittens-echo "goimports: Checking untracked files: \n$filesToCheck"
    badImports4=$(echo "$filesToCheck" | xargs $runGoImports --local "github.com/sweetloveinyourheart/exploding-kittens" -l)
fi

if [[ "$arg1" == "fix-goimports" ]]; then
    fixBadImports $badImports1
    fixBadImports $badImports2
    fixBadImports $badImports3
    fixBadImports $badImports4
    [ -z "$badImports1" ] || kittens-echo-red "$badImports1"
    [ -z "$badImports2" ] || kittens-echo-red "$badImports2"
    [ -z "$badImports3" ] || kittens-echo-red "$badImports3"
    [ -z "$badImports4" ] || kittens-echo-red "$badImports4"
    kittens-echo-green "goimports fixed."
    exit 0
fi

goImportsExitCode=0

checkForMalformedFile $badImports1 || goImportsExitCode=1
checkForMalformedFile $badImports2 || goImportsExitCode=1
checkForMalformedFile $badImports3 || goImportsExitCode=1
checkForMalformedFile $badImports4 || goImportsExitCode=1

if [ $goImportsExitCode -eq 1 ]; then
    kittens-echo-red "goimports failed for the following files:"
    [ -z "$badImports1" ] || kittens-echo-red "$badImports1"
    [ -z "$badImports2" ] || kittens-echo-red "$badImports2"
    [ -z "$badImports3" ] || kittens-echo-red "$badImports3"
    [ -z "$badImports4" ] || kittens-echo-red "$badImports4"
    kittens-echo-red "Please run 'goimports --local "github.com/sweetloveinyourheart/exploding-kittens" -w <file>' on the above files (or run 'make goimports') and commit the changes."
    exit 1
fi

kittens-echo "Running golangci-lint..."
go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2 run --timeout 10m0s ./... || ( kittens-echo-red "Linting failed." && exit 1 )

kittens-echo-green "Linting passed."
