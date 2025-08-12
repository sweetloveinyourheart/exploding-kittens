#!/bin/bash
. ./scripts/util.sh

tag=$1

if [ -z "$tag" ]; then
    echo "Usage: $0 <tag>"
    exit 1
fi

git tag -d $tag &> /dev/null
git push origin :refs/tags/$tag &> /dev/null
kittens-echo "DELETED tag: $tag"
