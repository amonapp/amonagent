FROM centos:centos6

ADD amonagent_0.2.8.4-1-gfdd1e0d-1.noarch.rpm var/agent.rpm

RUN yum -t -y install /var/agent.rpm
