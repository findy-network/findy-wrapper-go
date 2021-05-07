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


