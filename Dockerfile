FROM golang:1.8

#WORKDIR ./memefy-server
COPY ./memefy-server .

RUN go get github.com/mailru/go-clickhouse
RUN go get github.com/dgrijalva/jwt-go

RUN go build ./main.go
EXPOSE 8085

ENTRYPOINT ./main