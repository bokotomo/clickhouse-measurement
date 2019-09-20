FROM golang:1.12.7-stretch

LABEL maintainer fukumoto

SHELL ["/bin/bash", "-c"]
ENV GO111MODULE on
WORKDIR /clickhouse

# install mysql
RUN apt update
RUN apt install -y mysql-client vim zip unzip
RUN apt-get install -y default-libmysqlclient-dev