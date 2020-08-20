FROM centos:centos7

LABEL maintainer="ponce.ballesteros@gmail.com"

ENV SERVER_HOME=/opt/server \
    PATH=$SERVER_HOME:$PATH
    
WORKDIR $SERVER_HOME

COPY server $SERVER_HOME/

ENTRYPOINT ["/opt/server/server"]
