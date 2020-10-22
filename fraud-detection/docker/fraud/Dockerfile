FROM centos:centos7

LABEL maintainer="ponce.ballesteros@gmail.com"

ENV FRAUD_HOME=/opt/fraud \
    PATH=$FRAUD_HOME:$PATH
    
WORKDIR $FRAUD_HOME

COPY fraud $FRAUD_HOME/

ENTRYPOINT ["/opt/fraud/fraud"]
