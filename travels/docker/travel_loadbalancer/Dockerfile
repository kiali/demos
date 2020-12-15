FROM docker.io/fortio/fortio:1.11.4 as fortio
FROM centos:centos7

LABEL maintainer="ponce.ballesteros@gmail.com"

ENV TRAVEL_HOME=/opt/travel_loadbalancer \
    PATH=$TRAVEL_HOME:$PATH

WORKDIR $TRAVEL_HOME

COPY travel_loadbalancer.sh $TRAVEL_HOME/
COPY --from=fortio /usr/bin/fortio /usr/bin/fortio

RUN chmod +x $TRAVEL_HOME/travel_loadbalancer.sh
ENTRYPOINT ["/opt/travel_loadbalancer/travel_loadbalancer.sh"]