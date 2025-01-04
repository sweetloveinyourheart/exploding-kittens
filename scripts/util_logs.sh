#!/bin/bash

# Console colors
red='\033[0;31m'
yellow='\033[0;33m'
green='\033[0;32m'
gray='\033[0;90m'
cyan='\033[0;36m'
nc='\033[0m' # No Color, use this to terminate color sequences

kittens-timestamp() {
    date +"%H:%M:%S"
    return 0
}

kittens-echo() {
    local timestamp=$(kittens-timestamp)
    local callLocation=$(caller)
    local callLocationFile=$(echo $callLocation | cut -d' ' -f2)
    local callLocationLine=$(echo $callLocation | cut -d' ' -f1)
    echo -e "${gray}[$timestamp][$callLocationFile:$callLocationLine]:${nc} $@"
    return 0
}

kittens-echo-red() {
    local timestamp=$(kittens-timestamp)
    local callLocation=$(caller)
    local callLocationFile=$(echo $callLocation | cut -d' ' -f2)
    local callLocationLine=$(echo $callLocation | cut -d' ' -f1)
    echo -e "${gray}[$timestamp][$callLocationFile:$callLocationLine]:${nc} ${red}$@${nc}"
    return 0
}

kittens-echo-yellow() {
    local timestamp=$(kittens-timestamp)
    local callLocation=$(caller)
    local callLocationFile=$(echo $callLocation | cut -d' ' -f2)
    local callLocationLine=$(echo $callLocation | cut -d' ' -f1)
    echo -e "${gray}[$timestamp][$callLocationFile:$callLocationLine]:${nc} ${yellow}$@${nc}"
    return 0
}

kittens-echo-green() {
    local timestamp=$(kittens-timestamp)
    local callLocation=$(caller)
    local callLocationFile=$(echo $callLocation | cut -d' ' -f2)
    local callLocationLine=$(echo $callLocation | cut -d' ' -f1)
    echo -e "${gray}[$timestamp][$callLocationFile:$callLocationLine]:${nc} ${green}$@${nc}"
    return 0
}

kittens-echo-blue() {
    local timestamp=$(kittens-timestamp)
    local callLocation=$(caller)
    local callLocationFile=$(echo $callLocation | cut -d' ' -f2)
    local callLocationLine=$(echo $callLocation | cut -d' ' -f1)
    echo -e "${gray}[$timestamp][$callLocationFile:$callLocationLine]:${nc} ${cyan}$@${nc}"
    return 0
}
