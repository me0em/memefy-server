FROM frolvlad/alpine-glibc

ENV SERVERPORT = ":8085"
ENV DB = "http://default@memefy.fun:8123/memefy"
ENV SECRETKAY = "wer6YTIFpojneEfe34fr4go8ukcyyjr45y8867"
ENV MLMODELHOST = "http://127.0.0.1:8228/hero"
ENV ERRORHOST = "http://127.0.0.1:8228/error"
COPY ./main ./

ENTRYPOINT ./main
