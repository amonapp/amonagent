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
	mkdir -p $(BUILD)/opt/amonagent

	cp amonagent $(BUILD)/opt/amonagent/amonagent

	mkdir -p $(BUILD)/var/log/amonagent
	mkdir -p $(BUILD)/var/run/amonagent

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
--post-install $(PACKAGING)/postinst \
--post-uninstall $(PACKAGING)/postrm \
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
--conflicts "amonagent < $(VERSION)" \
--post-install	    $(PACKAGING)/postinst \
--post-uninstall    $(PACKAGING)/postrm \
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
	aws s3 sync $(PACKAGES_PATH)/debian/repo s3://beta.packages.amon.cx/repo --region=eu-west-1
	aws s3 sync $(PACKAGES_PATH)/centos s3://beta.packages.amon.cx/rpm --region=eu-west-1

upload:
	sudo ntpdate -u pool.ntp.org
	aws s3 sync $(PACKAGES_PATH)/debian s3://beta.packages.amon.cx/repo/ --region=eu-west-1
	aws s3 sync $(PACKAGES_PATH)/centos s3://beta.packages.amon.cx/rpm/ --region=eu-west-1


build_all_and_deploy: build_all deploy

build_test_debian_container:
	cp $(PACKAGING)/debian/Dockerfile.base Dockerfile
	docker build --force-rm=true --rm=true --no-cache -t=amonagent/ubuntu-base .
	rm Dockerfile
	docker rmi $$(docker images -q --filter dangling=true)

test_debian: build_deb
	cp $(PACKAGING)/debian/Dockerfile Dockerfile
	sed -i s/AMON_DEB_FILE/"$(DEBIAN_PACKAGE_NAME)"/g Dockerfile
	docker build --rm=true --no-cache -t=amonagent-$(VERSION) .
	rm Dockerfile
	docker rmi $$(docker images -q --filter dangling=true)
