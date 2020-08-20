FROM centos:centos7

LABEL maintainer="ponce.ballesteros@gmail.com"

ENV CLIENT_HOME=/opt/client \
    PATH=$CLIENT_HOME:$PATH
    
WORKDIR $CLIENT_HOME

COPY client $CLIENT_HOME/

ENTRYPOINT ["/opt/client/client"]
