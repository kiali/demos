FROM centos:centos7

LABEL maintainer="ponce.ballesteros@gmail.com"

ENV TRAVEL_HOME=/opt/travel_portal \
    PATH=$TRAVEL_HOME:$PATH
    
WORKDIR $TRAVEL_HOME

COPY travel_portal $TRAVEL_HOME/

ENTRYPOINT ["/opt/travel_portal/travel_portal"]
