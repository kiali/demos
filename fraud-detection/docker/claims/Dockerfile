FROM centos:centos7

LABEL maintainer="ponce.ballesteros@gmail.com"

ENV CLAIMS_HOME=/opt/claims \
    PATH=$CLAIMS_HOME:$PATH
    
WORKDIR $CLAIMS_HOME

COPY claims $CLAIMS_HOME/

ENTRYPOINT ["/opt/claims/claims"]
