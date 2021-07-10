#!/bin/bash

rsync -va --delete -e "ssh"  ./ root@server:/utile/api/ || exit 1
ssh server "
    docker build -t api-utile-space:latest /utile/api
    docker rm -f api-utile-space
    docker run -d -p 3000:3000 --name api-utile-space api-utile-space:latest
"