#!/usr/bin/env bash

go get github.com/mailru/go-clickhouse
go get github.com/dgrijalva/jwt-go
go build main.go
sudo docker build -t memefy-server .
docker push registry.gitlab.com/memerecommendersystemteam/memefy-server
