BUILD=packaging/build
PACKAGING=packaging
AGENT=$(BUILD)/amonagent
VERSION = $(shell sh -c 'git describe --always --tags')
# VERSION = "0.2.1"
INITD_SCRIPT=packaging/init.sh
SYSTEMD_SCRIPT=packaging/amonagent.service

PACKAGES_PATH=/home/martin/amon-packages
DEBIAN_REPO_PATH=$(PACKAGES_PATH)/debian/
RPM_REPO_PATH=$(PACKAGES_PATH)/centos/

PACKAGE="amonagent"
DEBIAN_PACKAGE_NAME="${PACKAGE}_${VERSION}_all.deb"
CENTOS_PACKAGE_NAME="${PACKAGE}-${VERSION}-1.noarch.rpm"

FPM_BUILD=fpm --epoch 1 -s dir -e -C $(BUILD) \
-a all -m "Amon Packages <packages@amon.cx>" \
--url "http://amon.cx/" \
--description "`cat $(PACKAGING)/desc`"\
-v $(VERSION) \
--vendor Amon

clean:
	rm -rf $(BUILD)

install_repo_base:
	sudo apt-get install reprepro createrepo -y --force-yes


update_dep:
	godep save
	godep update github.com/shirou/...

build:
	godep go build -o amonagent -ldflags \
		"-X main.Version=$(VERSION)" \
		./cmd/amonagent.go


# Layout all of the files common to both versions of the Agent in
# the build directory.
install_base: build
	mkdir -p $(BUILD)
	mkdir -p $(BUILD)/etc/opt/amonagent
	mkdir -p $(BUILD)/etc/opt/amonagent/plugins-enabled
	mkdir -p $(BUILD)/usr/bin

	cp amonagent $(BUILD)/usr/bin/amonagent

	mkdir -p $(BUILD)/var/log/amonagent
	chmod 755 $(BUILD)/var/log/amonagent

	mkdir -p $(BUILD)/opt/amonagent/scripts
	cp $(INITD_SCRIPT) $(BUILD)/opt/amonagent/scripts/init.sh
	cp $(SYSTEMD_SCRIPT) $(BUILD)/opt/amonagent/scripts/amonagent.service

	@echo $(VERSION)

# =====================
# Ubuntu/Debian
# =====================
build_deb: clean install_base
	rm -f *.deb
	FPM_EDITOR="echo 'Replaces: amonagent (<= $(VERSION))' >>" \
$(FPM_BUILD) -t deb \
-n amonagent \
-d "adduser" \
-d "sysstat" \
-d "dbus" \
--post-install $(PACKAGING)/postinst.sh \
--post-uninstall $(PACKAGING)/postrm.sh \
--pre-uninstall  $(PACKAGING)/prerm.sh \
.

# =====================
# CentOS/Fedora
# =====================
build_rpm: clean install_base
	rm -f *.rpm
	FPM_EDITOR="echo ''>>"  \
$(FPM_BUILD) -t rpm \
-n "amonagent" \
-d "sysstat" \
-d "dbus" \
--conflicts "amonagent < $(VERSION)" \
--post-install	    $(PACKAGING)/postinst.sh \
--post-uninstall    $(PACKAGING)/postrm.sh \
--pre-uninstall  $(PACKAGING)/prerm.sh \
.

build_all: build_deb build_rpm

update_debian_repo:
	cp "amonagent_$(VERSION)_all.deb" $(DEBIAN_REPO_PATH)
	find $(DEBIAN_REPO_PATH)  -name \*.deb -exec reprepro --ask-passphrase -Vb $(DEBIAN_REPO_PATH)repo includedeb amon {} \;

update_rpm_repo:
	cp "amonagent-$(VERSION)-1.noarch.rpm" $(RPM_REPO_PATH)
	createrepo --update $(RPM_REPO_PATH)


deploy: update_debian_repo update_rpm_repo
	sudo ntpdate -u pool.ntp.org
	aws s3 sync $(PACKAGES_PATH)/debian/repo s3://packages.amon.cx/repo --region=eu-west-1
	aws s3 sync $(PACKAGES_PATH)/centos s3://packages.amon.cx/rpm --region=eu-west-1



upload:
	sudo ntpdate -u pool.ntp.org
	aws s3 sync $(PACKAGES_PATH)/debian s3://packages.amon.cx/repo/ --region=eu-west-1
	aws s3 sync $(PACKAGES_PATH)/centos s3://packages.amon.cx/rpm/ --region=eu-west-1


build_and_deploy: build_all deploy

upload_packages: build_all 
	sudo ntpdate -u pool.ntp.org
	
	find . -iname "*.deb*" -execdir mv {} amonagent.deb \;
	find . -iname "*.rpm*" -execdir mv {} amonagent.rpm \;
	aws s3 cp amonagent.deb s3://amonagent-test --region=eu-west-1
	aws s3 cp amonagent.rpm s3://amonagent-test --region=eu-west-1
