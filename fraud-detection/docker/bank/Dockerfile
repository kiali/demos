FROM centos:centos7

LABEL maintainer="ponce.ballesteros@gmail.com"

ENV BANK_HOME=/opt/bank \
    PATH=$BANK_HOME:$PATH
    
WORKDIR $BANK_HOME

COPY bank $BANK_HOME/

ENTRYPOINT ["/opt/bank/bank"]
