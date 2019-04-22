FROM frolvlad/alpine-glibc

COPY ./main ./

ENTRYPOINT ./main
