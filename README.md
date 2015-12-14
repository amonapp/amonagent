# Amon Monitoring Agent

[Amon](https://amon.cx) is an easy-to-use hosted server monitoring service.
The Amon agent reports metrics to our service.


# [Change log](https://github.com/amonapp/amonagent-go/blob/master/CHANGELOG.md)

## Manual Installation


### Install on Ubuntu/Debian

1. **Import the public key used by the package management system.** <br>
	The Ubuntu package management tools (i.e. dpkg and apt) ensure package consistency and authenticity by requiring that distributors sign packages with GPG keys. Issue the following command to import the Amon Agent key

		sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv AD53961F


2. **Create a list file for the Agent.** <br>
Create the /etc/apt/sources.list.d/amonagent.list list file using the following command:


		echo 'deb http://packages.amon.cx/repo amon contrib' | sudo tee /etc/apt/sources.list.d/amonagent.list


3. **Reload local package database and install** <br>
Issue the following commands:


		sudo apt-get update
		sudo apt-get install amonagent


### Install on CentOS/Amazon Linux

1.**Configure the package management system (YUM).** <br>
Create a /etc/yum.repos.d/amonagent.repo file to hold the following configuration information for the Amon Agent repository:

	[amonagent]
	name=Amonagent Repository
	baseurl=http://packages.amon.cx/rpm/
	gpgcheck=0
	enabled=1
	priority=1


2.**Install the Agent package**. <br>
To install the latest stable version of the Agent, issue the following command:

	yum install -y amonagent




## Credits / Contact

Contact martin@amon.cx with questions.

Primary maintainer: Martin Rusev (martin@amon.cx)
