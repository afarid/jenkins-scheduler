FROM debian:10-slim

COPY jenkins-scheduler /jenkins-scheduler

ENTRYPOINT [ "/jenkins-scheduler" ]