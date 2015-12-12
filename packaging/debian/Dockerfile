FROM ubuntu:latest

RUN apt-get update
RUN apt-get install -y gdebi-core
ADD AMON_DEB_FILE var/agent.deb

RUN gdebi -n /var/agent.deb

RUN /etc/init.d/amon-agent status


CMD ["/bin/bash"]
