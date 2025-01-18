#!/usr/bin/env bash
# @author zouxiaoliang
# @date 2025/01/17

SOURCE="$0"

while [ -h "$SOURCE" ]; do
    DIR="$(cd -P "$(dirname "$SOURCE")" && pwd)"
    SOURCE="$(readlink "$SOURCE")"
    # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
    [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$(cd -P "$(dirname "$SOURCE")" && pwd)"

rm -rf jump-client jump-server

go build -o jump-client github.com/zouxiaoliang/jump/client
go build -o jump-server github.com/zouxiaoliang/jump/server
