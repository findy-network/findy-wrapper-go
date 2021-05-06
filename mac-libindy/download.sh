#!/bin/bash

# We don't want to check error return values by our self here.
set -e

download_and_extract_to() {
	local extract_to="$1"
	local libindy_location="https://repo.sovrin.org/macos/libindy/stable/1.16.0/libindy_1.16.0.zip"
	tmpfile=$(mktemp /tmp/findy-download.XXXXXX)

	curl -o "$tmpfile" "$libindy_location"
	unzip "$tmpfile" -d "$extract_to"
	rm -f "$tmpfile"
}

download_and_extract_to "./libindy"
