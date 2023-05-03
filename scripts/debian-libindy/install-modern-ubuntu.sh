#!/bin/bash

set -e

install_libssl() {
	wget http://archive.ubuntu.com/ubuntu/pool/main/o/openssl/libssl1.1_1.1.0g-2ubuntu4_amd64.deb
	sudo dpkg -i ./libssl1.1_1.1.0g-2ubuntu4_amd64.deb
	rm libssl1.1_1.1.0g-2ubuntu4_amd64.deb
}

install_indy() {
	INDY_LIB_VERSION="1.16.0"
	UBUNTU_VERSION="bionic"

	sudo apt-get update && \
	    sudo apt-get install -y software-properties-common apt-transport-https && \
	    sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 68DB5E88 && \
	    sudo add-apt-repository -y "deb https://repo.sovrin.org/sdk/deb $UBUNTU_VERSION stable" && \
	    sudo add-apt-repository -y "deb https://repo.sovrin.org/sdk/deb xenial stable" && \
	    sudo apt-get update

	sudo apt-get install -y libindy-dev="$INDY_LIB_VERSION-xenial" \
	    libindy="$INDY_LIB_VERSION-$UBUNTU_VERSION"
}

install_libssl
install_indy

echo "libindy-dev is installed and ready to serve!"

