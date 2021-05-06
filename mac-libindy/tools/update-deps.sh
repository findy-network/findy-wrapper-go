#!/bin/bash

lib_name="$1""/libindy.dylib"
echo "$lib_name"

if [[ ! -e "$lib_name" ]]; then 
	echo Error: cannot find "$lib_name"
	printf "\nusage:\t$0 location_of_libindy.dylib\n"
	exit 1
fi

get_abs_filename() {
	# $1 : relative filename
	echo "$(cd "$(dirname "$1")" && pwd)/$(basename "$1")"
}

brew_location() {
	echo $(brew --prefix "$1")
}

change_lib_location() {
	local oldpath="$1"
	local newpath="$2"
	printf "old: %s\n\tnew: %s\n" $oldpath $newpath

	if [[ ! -e "$oldpath" ]]; then 
		echo does not exits: "$oldpath"
		printf "install_name_tool call: %s %s %s\n" $oldpath $newpath $lib_name
		#install_name_tool -change $oldpath $newpath $lib_name
	else
		echo "OLD PATH COULD USE"
	fi
}

update_lib_location() {
	local pathbase="$1"
	local libname="$2"
	local curpath="$(otool -L "$lib_name" | tail -n +2 | egrep "$libname" | awk '{print $1}')"
	local actualpath="$3"

	if [[ "$actualpath" == "" ]]; then
		actualpath="$(brew_location $pathbase)""/lib/""$libname"".dylib"
	fi

	#printf "base: %s lib: %s current: %s\n\tshould: %s\n" $pathbase $libname $curpath $actualpath

	if [[ ! -e "$actualpath" ]]; then 
		echo Error: does not exits: "$actualpath"
		exit 1
	fi
	change_lib_location $curpath $actualpath
}

cellar_check() {
	local pattern="/usr/local/Cellar/openssl/1.0.2?"
	local files=( $pattern )
	echo "${files[0]}"
}

# main ----

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
	exit 1
fi

abs_lib_path=$(get_abs_filename "$lib_name")
update_lib_location "libindy" "libindy" "$abs_lib_path"
echo "-----"
update_lib_location "zeromq" "libzmq" 
update_lib_location "$brew_name" "libssl" "$ssl_location"
update_lib_location "$brew_name" "libcrypto" "$crypto_location"
update_lib_location "libsodium" "libsodium"

