FROM centos:centos7

LABEL maintainer="ponce.ballesteros@gmail.com"

ENV CARDS_HOME=/opt/cards \
    PATH=$CARDS_HOME:$PATH
    
WORKDIR $CARDS_HOME

COPY cards $CARDS_HOME/

ENTRYPOINT ["/opt/cards/cards"]
