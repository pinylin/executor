FROM ubuntu:14.04.1

MAINTAINER  pinylin

RUN apt-get update
RUN apt-get install -y ca-certificates
RUN apt-get install dialog
RUN mkdir -p /usr/share/cmd

COPY tasksService /usr/share/cmd/wechatService

ENTRYPOINT ["/usr/share/cmd/wechatService"]