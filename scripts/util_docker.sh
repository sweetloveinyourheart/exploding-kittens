#!/bin/bash

pocker-docker-cleanup() {
    pkill compose &> /dev/null || : # compose logs processes are sometimes left running, kill those too
    docker network rm -f fs_test_net &> /dev/null || :
    docker network prune -f &> /dev/null || :
    docker container prune -f &> /dev/null || :
    docker system prune -f &> /dev/null || :
    return 0
}

pocker-compose-up() {
    local composeFile=$1

    if [ -z "$composeFile" ]; then
        pocker-echo-red "call requires a composeFile"
        return 1
    fi

    composeCommand="docker compose -f $composeFile up --wait-timeout 300"
    pocker-echo "Running: $composeCommand"
    $composeCommand || (pocker-echo "Failed to up the compose stack" && return 1)
    return $?
}

pocker-compose-down() {
    local composeFile=$1

    if [ -z "$composeFile" ]; then
        pocker-echo-red "call requires a composeFile"
        return 1
    fi

    composeCommand="docker compose -f $composeFile down --timeout=0 --volumes --remove-orphans"
    pocker-echo "Running: $composeCommand"
    $composeCommand || ( pocker-echo "Failed to down the compose stack" && return 1 )
    return $?
}
