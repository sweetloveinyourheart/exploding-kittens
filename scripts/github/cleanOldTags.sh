#!/bin/bash
. ./scripts/util.sh

# Remove any tags that don't match the patterns we want to keep, or are older than 3 months
# This script is intended to be run as a cron job

# delete all of the local tags
git tag | xargs git tag -d

# fetch all of the remote tags
git fetch --tags

# get all of the tags in this repo
tags=$(git tag)

# get the current date
now=$(date +%s)

# loop through the tags
for tag in $tags; do

    # keep auto-generated main tags for 2 months
    if [[ "$tag" == "20"*"_main_"* ]]; then
        tagDate=$(git log -1 --format=%at $tag)
        if (( $now - $tagDate > 5184000 )); then
            ./scripts/github/deleteTag.sh $tag
        fi

    # keep auto-generated tags from other branches for 1 month
    elif [[ "$tag" == "20"* ]]; then
        tagDate=$(git log -1 --format=%at $tag)
        # if the tag is older than 1 month
        if (( $now - $tagDate > 2592000 )); then
            ./scripts/github/deleteTag.sh $tag
        fi
    else
        # keep it :
        kittens-echo "KEEP $tag"
    fi
done
