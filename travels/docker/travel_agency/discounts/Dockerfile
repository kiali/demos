FROM centos:centos7

LABEL maintainer="ponce.ballesteros@gmail.com"

ENV TRAVEL_HOME=/opt/travel_agency \
    PATH=$TRAVEL_HOME:$PATH
    
WORKDIR $TRAVEL_HOME

COPY discounts $TRAVEL_HOME/

ENTRYPOINT ["/opt/travel_agency/discounts"]
