#!/bin/bash

chroot /host /usr/bin/env -i PATH="/bin:/sbin:/usr/bin" $(basename $0) $@
