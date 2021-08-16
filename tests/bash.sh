#!/bin/bash
set -e

cd "$(dirname "$0")"

for foo in *.sh; do
	if [ $foo = "bash.sh" ]; then
		# do not run self
		continue
	fi
	echo -ne "\rTest: $foo...\033[K"
	bash "$foo"
done

echo -e "\rAll tests OK\033[K"
