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
		echo "Warning: does not exits: ""$actualpath"
		echo "We can try to install ""$libname"" for you."
		if [[ $(prompt_default_no) == "yes" ]]; then
			brew install "$libname"
		else
			echo "Terminating. Please install the missing library."
			exit 1
		fi
	fi
}

cellar_check() {
	local pattern="/usr/local/Cellar/openssl/1.0.2?"
	local files=( $pattern )
	echo "${files[0]}"
}

prompt_default_yes() {
	read -r -p "Do you want that? [Y/n] " response
	if [[ "$response" =~ ^([nN][oO]|[nN])$ ]]
	then
	    echo no
	else
	    echo yes
	fi
}

prompt_default_no() {
	read -r -p "Are you sure? [y/N] " response
	if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]
	then
	    echo yes
	else
	    echo no
	fi
}

