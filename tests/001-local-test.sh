#!/bin/bash

b() {
	X=b
}

a() {
	local X=a
	b
	echo $X
}

X=z
Y="$(a)"

if [ $X != z ] || [ $Y != b ]; then
	echo "FAIL, expected $X=z and $Y=b"
	exit 1
fi
