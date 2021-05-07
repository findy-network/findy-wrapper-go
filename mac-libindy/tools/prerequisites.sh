#!/bin/bash

if ! command -v brew &> /dev/null
then
	echo "home brew is not installed for the current user"
	exit 1
fi

