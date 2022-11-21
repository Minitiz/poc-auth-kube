build-docker: build-api build-sh 

build-sh:
	cp go.mod demo/services/storage-hub
	$(MAKE) -C demo/services/storage-hub build-docker
	rm demo/services/storage-hub/go.mod

build-api:
	$(MAKE) -C demo/services/api build-docker

start-minikube:
	minikube start --extra-config=apiserver.service-account-signing-key-file=/var/lib/minikube/certs/sa.key --extra-config=apiserver.service-account-issuer=kubernetes/serviceaccount --extra-config=apiserver.service-account-api-audiences=api

apply-api:
	$(MAKE) -C demo/services/api apply
apply-api2:
	$(MAKE) -C demo/services/api apply-api2
apply-api3:
	$(MAKE) -C demo/services/api apply-api3
apply-sh:
	$(MAKE) -C demo/services/storage-hub apply

apply-all: apply-api apply-api2 apply-api3 apply-sh

clean: clean-sh clean-api

clean-sh:
	-$(MAKE) -C demo/services/storage-hub clean

clean-api:
	-$(MAKE) -C demo/services/api clean

all: clean build-docker apply-all

pf-api:
	$(MAKE) -C demo/services/api pf

pf-api2:
	$(MAKE) -C demo/services/api pf-api2

pf-api3:
	$(MAKE) -C demo/services/api pf-api3

pf-sh:
	$(MAKE) -C demo/services/storage-hub pf

# no sense
pf-all:
	$(MAKE) -C demo/services/api pf &
	$(MAKE) -C demo/services/api pf-api2 &
	$(MAKE) -C demo/services/storage-hub pf &

logs-api:
	$(MAKE) -C demo/services/api log

logs-api2:
	$(MAKE) -C demo/services/api2 log-api2

logs-api3:
	$(MAKE) -C demo/services/api2 log-api2

logs-sh:
	$(MAKE) -C demo/services/storage-hub log

deploy-sh: clean-sh build-sh apply-sh

deploy-api: clean-api build-api apply-api apply-api2 apply-api3

