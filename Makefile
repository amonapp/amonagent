BUILD=packaging/build
PACKAGING=packaging
AGENT=$(BUILD)/amonagent
VERSION := $(shell sh -c 'git describe --always --tags')

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
	go build -o amonagent -ldflags \
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

	mkdir -p $(BUILD)/etc/init.d
	cp -r $(PACKAGING)/amonagent.init.sh $(BUILD)/etc/init.d/amonagent
	chmod +x $(BUILD)/etc/init.d/amonagent
	chmod 755 $(BUILD)/etc/init.d/amonagent

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
--post-install $(PACKAGING)/debian/postinst \
--post-uninstall $(PACKAGING)/debian/postrm \
--pre-uninstall  $(PACKAGING)/debian/prerm \
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
--post-install	    $(PACKAGING)/rpm/postinst \
--pre-uninstall	    $(PACKAGING)/rpm/prerm \
--post-uninstall    $(PACKAGING)/rpm/postrm \
.

test_debian: build_deb
	cp $(PACKAGING)/debian/Dockerfile Dockerfile
	sed -i s/AMON_DEB_FILE/"$(DEBIAN_PACKAGE_NAME)"/g Dockerfile
	docker build --rm=true --no-cache .
	rm Dockerfile
