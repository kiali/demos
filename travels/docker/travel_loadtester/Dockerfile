FROM docker.io/fortio/fortio:1.11.4 as fortio
FROM centos:centos7

LABEL maintainer="ponce.ballesteros@gmail.com"

ENV TRAVEL_HOME=/opt/travel_loadtester \
    PATH=$TRAVEL_HOME:$PATH

WORKDIR $TRAVEL_HOME

COPY travel_loadtester.sh $TRAVEL_HOME/
COPY --from=fortio /usr/bin/fortio /usr/bin/fortio

RUN chmod +x $TRAVEL_HOME/travel_loadtester.sh
ENTRYPOINT ["/opt/travel_loadtester/travel_loadtester.sh"]