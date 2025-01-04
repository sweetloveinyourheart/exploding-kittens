#!/bin/bash

kittens-docker-cleanup() {
    pkill compose &> /dev/null || : # compose logs processes are sometimes left running, kill those too
    docker network rm -f fs_test_net &> /dev/null || :
    docker network prune -f &> /dev/null || :
    docker container prune -f &> /dev/null || :
    docker system prune -f &> /dev/null || :
    return 0
}

kittens-compose-up() {
    local composeFile=$1

    if [ -z "$composeFile" ]; then
        kittens-echo-red "call requires a composeFile"
        return 1
    fi

    composeCommand="docker compose -f $composeFile up --wait-timeout 300"
    kittens-echo "Running: $composeCommand"
    $composeCommand || (kittens-echo "Failed to up the compose stack" && return 1)
    return $?
}

kittens-compose-down() {
    local composeFile=$1

    if [ -z "$composeFile" ]; then
        kittens-echo-red "call requires a composeFile"
        return 1
    fi

    composeCommand="docker compose -f $composeFile down --timeout=0 --volumes --remove-orphans"
    kittens-echo "Running: $composeCommand"
    $composeCommand || ( kittens-echo "Failed to down the compose stack" && return 1 )
    return $?
}
