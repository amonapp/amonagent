$deb = <<SCRIPT
rm *.deb
apt-get update
apt-get install -y sysstat
apt-get remove -y amonagent
rm -rf /etc/opt/amonagent
rm -rf /var/log/amonagent/amonagent.log



dpkg -i /vagrant/amonagent.deb
echo '{"api_key": "test", "amon_instance": "https://demo.amon.cx"}' >> /etc/opt/amonagent/amonagent.conf
/opt/amonagent/amonagent -test

service amonagent start
service amonagent status
cat /var/log/amonagent/amonagent.log
service amonagent stop
SCRIPT


$rpm = <<SCRIPT
rm *.rpm
yum install /vagrant/amonagent.rpm -y
echo '{"api_key": "test", "amon_instance": "https://demo.amon.cx"}' >> /etc/opt/amonagent/amonagent.conf


/opt/amonagent/amonagent -test

service amonagent start
service amonagent status
cat /var/log/amonagent/amonagent.log
service amonagent stop
SCRIPT

# vagrant plugin install vagrant-vbguest
Vagrant.configure("2") do |config|

  config.vm.synced_folder "~/go/src/github.com/amonapp/amonagent/",
		"/vagrant/",
		:mount_options => [ "dmode=777", "fmode=777" ]

  config.vm.define "ubuntu1404" do |ubuntu1404|
    ubuntu1404.vm.box = "ubuntu/trusty64"

    ubuntu1404.vm.provision "shell", inline: $deb
  end

  config.vm.define "debian8" do |debian8|
    debian8.vm.box = "debian/jessie64"
    debian8.vm.provision "shell", inline: $deb
  end

  config.vm.define "debian7" do |debian7|
    debian7.vm.box = "debian/wheezy64"
    debian7.vm.provision "shell", inline: $deb
  end

  config.vm.define "centos6" do |centos6|
    centos6.vm.box = "puphpet/centos65-x64"
    centos6.vm.provision "shell", inline: $rpm
  end

  config.vm.define "centos7" do |centos7|
    centos7.vm.box = "bento/centos-7.1"
    centos7.vm.provision "shell", inline: $rpm
  end

end