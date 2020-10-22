FROM centos:centos7

LABEL maintainer="ponce.ballesteros@gmail.com"

ENV POLICIES_HOME=/opt/policies \
    PATH=$POLICIES_HOME:$PATH
    
WORKDIR $POLICIES_HOME

COPY policies $POLICIES_HOME/

ENTRYPOINT ["/opt/policies/policies"]
