FROM ubuntu:latest

RUN apt-get update
RUN apt-get install -y gdebi-core
ADD amonagent_0.1.1-12-g22b86de_all.deb var/agent.deb

RUN gdebi -n /var/agent.deb

RUN /etc/init.d/amon-agent status


CMD ["/bin/bash"]
