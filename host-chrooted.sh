#!/bin/bash

target=$(basename $0)
if [ -n "$TARGET" ]; then
	target="$TARGET"
fi

chroot /host /usr/bin/env -i PATH="/bin:/sbin:/usr/bin:/lib/udev" $target $@
