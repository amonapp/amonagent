FROM ubuntu:latest

RUN apt-get update
RUN apt-get install -y gdebi-core
ADD amonagent_0.1.1-13-g22c4b4c_all.deb var/agent.deb

RUN gdebi -n /var/agent.deb

RUN /etc/init.d/amonagent start


CMD ["/bin/bash"]
