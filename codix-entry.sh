#!/usr/bin/env bash
# CODIX DevOps
#

cd /go/oci-api
export PATH=$PATH:/usr/local/go/bin
make run

if [ $# -gt 0 ];then
    exec "$@"
else
    sh
fi
