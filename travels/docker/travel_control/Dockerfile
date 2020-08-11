FROM centos:centos7

LABEL maintainer="ponce.ballesteros@gmail.com"

ENV TRAVEL_HOME=/opt/travel_control \
    PATH=$TRAVEL_HOME:$PATH
    
WORKDIR $TRAVEL_HOME

COPY travel_control $TRAVEL_HOME/
COPY html $TRAVEL_HOME/html

ENTRYPOINT ["/opt/travel_control/travel_control"]
