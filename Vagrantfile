$deb = <<SCRIPT
rm *.deb
wget https://s3-eu-west-1.amazonaws.com/amonagent-test/amonagent.deb
apt-get update
apt-get install -y sysstat
echo '{"api_key": "test", "amon_instance": "https://demo.amon.cx"}' >> /etc/opt/amonagent/amonagent.conf
dpkg -i amonagent.deb

/opt/amonagent/amonagent -test
SCRIPT


$rpm = <<SCRIPT
rm *.rpm
wget https://s3-eu-west-1.amazonaws.com/amonagent-test/amonagent.rpm
echo '{"api_key": "test", "amon_instance": "https://demo.amon.cx"}' >> /etc/opt/amonagent/amonagent.conf
yum install amonagent.rpm -y

/opt/amonagent/amonagent -test
SCRIPT

# vagrant plugin install vagrant-vbguest
Vagrant.configure("2") do |config|

  config.vm.synced_folder ".", "/vagrant", disabled: true

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