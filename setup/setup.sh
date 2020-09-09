#!/bin/bash

tmp_dir="/tmp/test"
[ ! -d ${tmp_dir} ] && mkdir ${tmp_dir} || rm -fr ${tmp_dir} && mkdir ${tmp_dir}
#清空远程服务器的目录
/usr/bin/rsync -avpz --delete ${tmp_dir}/ ${remote_host}:{dest}