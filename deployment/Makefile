build_puppet_container:
	cp puppet/Dockerfile.base Dockerfile
	docker build --force-rm=true --rm=true --no-cache --tag="martinrusev/puppet-master" .    
	rm Dockerfile
	docker rmi $$(docker images -q --filter dangling=true)

# =====================
# Test
# =====================

# Ansible test_ansible distro = ubuntu|centos
test_ansible:
	cp ansible/Dockerfile.$(distro) Dockerfile
	docker build --force-rm=true --rm=true --no-cache .  
	rm Dockerfile
	docker rmi $$(docker images -q --filter dangling=true)




test_puppet:
	cp puppet/Dockerfile Dockerfile
	docker build --force-rm=true --rm=true --no-cache .  
	rm Dockerfile
	docker rmi $$(docker images -q --filter dangling=true)
