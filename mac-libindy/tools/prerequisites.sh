#!/bin/bash

brew_location() {
	echo $(brew --prefix "$1")
}

check_lib_exists() {
	local pathbase="$1"
	local libname="$2"
	local actualpath="$3"

	if [[ "$actualpath" == "" ]]; then
		actualpath="$(brew_location $pathbase)""/lib/""$libname"".dylib"
	fi

	if [[ ! -e "$actualpath" ]]; then 
		echo "Error: does not exits: ""$actualpath"
		echo "do you want us to install it?"
		exit 1
	fi
}

cellar_check() {
	local pattern="/usr/local/Cellar/openssl/1.0.2?"
	local files=( $pattern )
	echo "${files[0]}"
}

# =========== main =================

if ! command -v brew &> /dev/null
then
	echo "home brew is not installed for the current user"
	exit 1
fi


if [[ "$openssl_path" == "" ]]; then
	cellar_path=$(cellar_check)
	if [[ -e "$cellar_path" ]]; then
		openssl_path="$cellar_path"
		ssl_location="$openssl_path/lib/libssl.dylib"
		crypto_location="$openssl_path/lib/libcrypto.dylib"
	else
		brew_name="openssl@1.0" 
		openssl_path=$(brew_location "openssl@1.0")
	fi
fi

if [[ ! -e "$openssl_path" ]]; then
	echo "open ssl 1.0 not installed..."
	echo "do you want us to try to install it?"
	exit 1
fi

check_lib_exists "zeromq" "libzmq" 
check_lib_exists "$brew_name" "libssl" "$ssl_location"
check_lib_exists "$brew_name" "libcrypto" "$crypto_location"
check_lib_exists "libsodium" "libsodium"

