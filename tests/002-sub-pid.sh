#!/bin/sh

a() {
	echo $$
}

PID="$$"
PID2="$(a)"

if [ $PID != $PID2 ]; then
	echo "invalid, $PID should be equal to $PID2"
	exit 1
fi
