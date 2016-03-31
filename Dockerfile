FROM ubuntu:latest

ADD amonagent_0.5-1-g8b8581a_all.deb var/agent.deb

# Install dependecy
RUN apt-get install -y sysstat
RUN mkdir -p /etc/opt/amonagent
RUN echo '{"api_key": "test", "amon_instance": "https://demo.amon.cx"}' >> /etc/opt/amonagent/amonagent.conf
RUN dpkg -i /var/agent.deb

RUN /opt/amonagent/amonagent -test
RUN apt-get remove -y amonagent
