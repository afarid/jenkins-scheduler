FROM debian:10-slim

RUN apt-get update && apt-get install ca-certificates -y \
 && rm -rf /var/lib/apt/lists/*

COPY jenkins-scheduler /jenkins-scheduler

ENTRYPOINT [ "/jenkins-scheduler" ]