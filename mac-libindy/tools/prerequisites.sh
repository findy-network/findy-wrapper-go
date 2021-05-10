#!/bin/bash

set -e

. tools/functions.sh

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
	echo "open ssl 1.0 not installed. We can install it for you..."
	if [[ $(prompt_default_yes) == "no" ]]; then
		exit 1
	else
		echo "installing openssl 1.0..."
		brew install rbenv/tap/openssl@1.0
	fi
fi

check_lib_exists "zeromq" "libzmq" 
check_lib_exists "$brew_name" "libssl" "$ssl_location"
check_lib_exists "$brew_name" "libcrypto" "$crypto_location"
check_lib_exists "libsodium" "libsodium"

