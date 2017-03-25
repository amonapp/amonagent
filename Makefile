BUILD=packaging/build
PACKAGING=packaging
AGENT=$(BUILD)/amonagent
VERSION = $(shell sh -c 'git describe --always --tags')
INITD_SCRIPT=packaging/init.sh
SYSTEMD_SCRIPT=packaging/amonagent.service
TEMPFILE_CONF=packaging/tmpfilesd_amonagent.conf

PACKAGES_PATH=/home/martin/amon/amon-packages
DEBIAN_REPO_PATH=$(PACKAGES_PATH)/debian/
RPM_REPO_PATH=$(PACKAGES_PATH)/centos/

PACKAGE="amonagent"
DEBIAN_PACKAGE_NAME="${PACKAGE}_${VERSION}_all.deb"
CENTOS_PACKAGE_NAME="${PACKAGE}-${VERSION}-1.noarch.rpm"
GO_SKIP_VENDOR=$(shell sh -c 'go list ./... | grep -v /vendor/')

FPM_BUILD=fpm --epoch 1 -s dir -e -C $(BUILD) \
-a all -m "Amon Packages <packages@amon.cx>" \
--url "http://amon.cx/" \
--description "`cat $(PACKAGING)/desc`"\
-v $(VERSION) \
--vendor Amon


setup_test_env:
	sudo apt-get install ruby-dev
	sudo gem install sensu-plugins-disk-checks --no-ri --no-rdoc

setup_dev_env:
	sudo apt-get install ruby-dev gcc make reprepro createrepo -y --force-yes
	sudo gem install fpm
	vagrant up

clean:
	rm -rf $(BUILD)


build:
	CGO_ENABLED=0 go build -o amonagent -ldflags \
		"-X main.Version=$(VERSION)" \
		./cmd/amonagent.go


build_32bit:
	CGO_ENABLED=0 GOARCH=386 go build -o amonagent32 -ldflags \
		"-X main.Version=$(VERSION)" \
		./cmd/amonagent.go

build_arm:
	CGO_ENABLED=0 GOARCH=arm go build -o amonagent-arm -ldflags \
		"-X main.Version=$(VERSION)" \
		./cmd/amonagent.go

# Layout all of the files common to both versions of the Agent in
# the build directory.
install_base: build build_32bit
	mkdir -p $(BUILD)
	mkdir -p $(BUILD)/etc/opt/amonagent
	mkdir -p $(BUILD)/etc/opt/amonagent/plugins-enabled
	mkdir -p $(BUILD)/opt/amonagent
	mkdir -p $(BUILD)/usr/bin

	cp amonagent $(BUILD)/opt/amonagent/amonagent
	cp amonagent $(BUILD)/usr/bin/amonagent

	# Add the 32bit binary
	cp amonagent32 $(BUILD)/usr/bin/amonagent32

	mkdir -p $(BUILD)/var/log/amonagent
	chmod 755 $(BUILD)/var/log/amonagent

	# /var/run permissions for RPM distros
	mkdir -p $(BUILD)/usr/lib/tmpfiles.d
	cp $(TEMPFILE_CONF) $(BUILD)/usr/lib/tmpfiles.d/amonagent.conf

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
	aws s3 sync $(PACKAGES_PATH)/debian/repo s3://packages.amon.cx/repo --region=eu-west-1 --profile=personal
	aws s3 sync $(PACKAGES_PATH)/centos s3://packages.amon.cx/rpm --region=eu-west-1 --profile=personal



upload:
	sudo ntpdate -u pool.ntp.org
	aws s3 sync $(PACKAGES_PATH)/debian s3://packages.amon.cx/repo/ --region=eu-west-1 --profile=personal
	aws s3 sync $(PACKAGES_PATH)/centos s3://packages.amon.cx/rpm/ --region=eu-west-1 --profile=personal


build_and_deploy: build_all deploy

upload_packages: build_all
	sudo ntpdate -u pool.ntp.org

	find . -iname "*.deb*" -execdir mv {} amonagent.deb \;
	find . -iname "*.rpm*" -execdir mv {} amonagent.rpm \;
	aws s3 cp amonagent.deb s3://amonagent-test --region=eu-west-1
	aws s3 cp amonagent.rpm s3://amonagent-test --region=eu-west-1

test_deb: build_deb
	find . -iname "*.deb*" -execdir mv {} amonagent.deb \;
	# vagrant reload ubuntu1404 --provision
	# vagrant reload debian8 --provision
	vagrant reload debian7 --provision

test_rpm: build_rpm
	find . -iname "*.rpm*" -execdir mv {} amonagent.rpm \;
	vagrant reload centos6 --provision
	# vagrant reload centos7 --provision

test_output:
	go build
	./amonagent -test


setup_travis:
	sudo add-apt-repository ppa:masterminds/glide && sudo apt-get update
	sudo apt-get install glide
	glide i


# Run full unit tests using docker containers (includes setup and teardown)
test: vet
	go test -race $(GO_SKIP_VENDOR)

# Run "short" unit tests
test-short: vet
	go test -short $(GO_SKIP_VENDOR)

race:
	go run -race cmd/amonagent.go

vet:
	go vet $(GO_SKIP_VENDOR)
