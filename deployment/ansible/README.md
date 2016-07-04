## Installing the amon-agent with Ansible

- Copy all the files from deployment/ansible to a local path
- In api_vars.yml, replace the following variables with the appropriate values:


```
api_key - a valid API key 
amon_instance - the IP or domain pointing to your Amon instance
```

## Running the playbook 



```
ansible-playbook amonagent.yml
```

Tested on Debian/Ubuntu and CentOS


## Test the playbook

Install Docker, in api_vars.yml, uncomment machine_id for Docker then: 

```
docker pull williamyeh/ansible:centos7
docker pull williamyeh/ansible:ubuntu14.04

make test distro=ubuntu
make test distro=centos

```
