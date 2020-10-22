FROM centos:centos7

LABEL maintainer="ponce.ballesteros@gmail.com"

ENV ACCOUNTS_HOME=/opt/accounts \
    PATH=$ACCOUNTS_HOME:$PATH
    
WORKDIR $ACCOUNTS_HOME

COPY accounts $ACCOUNTS_HOME/

ENTRYPOINT ["/opt/accounts/accounts"]
