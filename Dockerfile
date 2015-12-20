FROM ubuntu:latest

ADD amonagent_0.3_all.deb var/agent.deb

# Install dependecy
RUN apt-get install -y sysstat dbus
RUN mkdir -p /etc/opt/amonagent
RUN echo '{"api_key": "test", "amon_instance": "https://demo.amon.cx"}' >> /etc/opt/amonagent/amonagent.conf
RUN dpkg -i /var/agent.deb

RUN /opt/amonagent/amonagent -test
RUN apt-get remove -y amonagent
