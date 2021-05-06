#!/bin/bash

if ! command -v brew &> /dev/null
then
	echo "how brew is not installed for this user"
	exit 1
fi

