FROM golang:1.8

WORKDIR ./
COPY ./memefy-server ./
RUN go get github.com/mailru/go-clickhouse
RUN go get github.com/dgrijalva/jwt-go
RUN go build main.go

ENTRYPOINT ./main
#COPY ./main ./
#
#EXPOSE 8085
#
#ENTRYPOINT ./main