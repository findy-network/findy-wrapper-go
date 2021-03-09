#!/bin/bash

if [[ "$1" != "" ]]; then
	echo freeing variables
	unset ImmuUrl
	unset ImmuPort
	unset ImmuUsrName
	unset ImmuPasswd
else
	export ImmuUrl="mock"
	export ImmuPort=3322
	export ImmuUsrName="immudb"
	export ImmuPasswd="immudb"
	echo variables set
fi
