#!/bin/bash

# in case of using cache dir, we need to initialize it
/usr/sbin/squid -d 1 --foreground -f /etc/squid/squid.conf -z

# now start the squid primary process with supplied options
/usr/sbin/squid -d 1 --foreground -f /etc/squid/squid.conf $@
