include ./blackbox.conf
export

DATA_DIR?=~/data

CHAINS?=
chains-compose-files := $(foreach service,$(CHAINS),-f ./services/$(service)/docker-compose.yml)

docker-compose = DATA_DIR=$(DATA_DIR) docker-compose -p blackbox -f ./docker-compose.yml -f ./docker-compose-deps.yml $(chains-compose-files)
# Load only the supporting containers
# not the API container
docker-compose-dev = DATA_DIR=$(DATA_DIR) docker-compose -p blackbox -f ./docker-compose-deps.yml $(chains-compose-files)

build:
	./scripts/build-docker.sh

configuration: setup
	$(docker-compose) config

devup:
	$(docker-compose-dev) pull && \
	$(docker-compose-dev) up -t 60
devdown:
	$(docker-compose-dev) down --remove-orphans


pull: setup
	$(docker-compose) pull

start: pull
	$(docker-compose) up -t 60

update: pull
	$(docker-compose) up -d -t 60

stop:
	$(docker-compose) down --remove-orphans


install-services: install-blackbox-service
	systemctl daemon-reload

uninstall-services: uninstall-blackbox-service
	systemctl daemon-reload

# Installs the systemd service, enables it and starts it
install-blackbox-service:
	cp services/blackbox.service /etc/systemd/system/
	systemctl enable /etc/systemd/system/blackbox.service
	systemctl start blackbox.service

# Uninstalls the service
uninstall-blackbox-service:
	systemctl stop blackbox.service
	systemctl disable blackbox.service
	rm /etc/systemd/system/blackbox.service

# PIVX (and maybe other CHAINs) require a swap file to function properly.
# * This must be run as root and is NOT idempotent
# * Use $ swapon --show before running this command!
# * See: https://linuxize.com/post/how-to-add-swap-space-on-ubuntu-18-04/
install-swapfile:
	fallocate -l 2G /swapfile && \
	chmod 600 /swapfile && \
	mkswap /swapfile && \
	swapon /swapfile && \
	echo "/swapfile swap swap defaults 0 0" >> /etc/fstab && \
	swapon --show

install-docker:
	apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common
	curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
	apt-key fingerprint 0EBFCD88
	add-apt-repository "deb [arch=arm64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
	apt-get update && apt-get install -y docker-ce docker-ce-cli containerd.io

setup: check-chains encryption-key

# Generates an ecryption key to be used in the API
encryption-key:
	@./scripts/generate-enc-key.sh

# DATA_DIR=/path/to/pivxdata
check-chains:
ifndef CHAINS
	$(error 'CHAINS' is undefined)
else
	@echo "configured for ${CHAINS}"
endif
