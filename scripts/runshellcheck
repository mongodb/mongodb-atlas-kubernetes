#!/bin/sh

# shellcheck disable=SC1004

find . \
	-type d -name '.git' -prune \
	-o \
	-type f \( \
		-exec sh -c 'file --brief "$1" \
			  | grep -qE "((POSIX|Korn|Bourne-Again) shell|/usr/bin/env k?sh) script"' _ {} \; \
		-o \
		-name '*.sh' -o -name '*.bash' -o -name '*.ksh' \
	\) \
	-exec shellcheck --color=always {} +