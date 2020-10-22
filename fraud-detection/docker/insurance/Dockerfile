FROM centos:centos7

LABEL maintainer="ponce.ballesteros@gmail.com"

ENV INSURANCE_HOME=/opt/insurance \
    PATH=$INSURANCE_HOME:$PATH
    
WORKDIR $INSURANCE_HOME

COPY insurance $INSURANCE_HOME/

ENTRYPOINT ["/opt/insurance/insurance"]
