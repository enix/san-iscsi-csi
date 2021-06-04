#!/bin/bash

# if [ ! -d /host/proc ]; then
# 	# echo "setup /host2 fake root"
# 	# cd /host
# 	# for f in $(find -maxdepth 1 -type d ! -path .); do mkdir -p /host2/$(basename $f) && /bin/mount --bind /host/$f /host2/$(basename $f); done
# 	# /bin/umount /host2/proc
# 	mkdir /host/proc
# 	/bin/mount --bind /proc /host/proc
# fi

target=$(basename $0)
if [ -n "$TARGET" ]; then
	target="$TARGET"
fi

# if [[ "$target" == "multipath" ]]; then
# 	# target="strace -f $target"
# 	LD_LIBRARY_PATH=/root/multipath-tools-0.7.4/libmultipath/:/root/multipath-tools-0.7.4/libmpathcmd/
# 	target="LD_LIBRARY_PATH=$LD_LIBRARY_PATH gdbserver :9942 /root/multipath-tools-0.7.4/$target/$target"
# fi

chroot /host /usr/bin/env -i PATH="/bin:/sbin:/usr/bin:/lib/udev" $target $@
