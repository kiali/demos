FROM python:3.7.6-stretch

RUN pip install pip --upgrade
RUN pip install ansible
RUN pip install openshift jmespath 'kubernetes==11.0.0'
RUN ansible-galaxy collection install community.kubernetes

RUN apt-get update -y && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    sshpass

COPY playbooks /ansible/playbooks

WORKDIR /ansible/playbooks

ENTRYPOINT ["ansible-playbook", "-v", "scale-mesh.yml"]
