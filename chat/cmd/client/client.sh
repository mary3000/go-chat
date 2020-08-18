#!/bin/bash

cd ../../../

force=0

while [[ "$1" =~ ^- && ! "$1" == "--" ]]; do case $1 in
  -f | --force )
    force=1
    ;;
esac; shift; done
if [[ "$1" == '--' ]]; then shift; fi

if [ "$force" == 1 ] || [[ "$(docker images -q chat_client 2> /dev/null)" == "" ]]; then
  docker build -t chat_client -f chat/cmd/client/Dockerfile .
fi

docker run --network host -it chat_client
